package main

import (
    "context"
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "os/exec"
    "os/signal"
    "strconv"
    "strings"
    "syscall"
    "time"

    i18n "github.com/Xarth-Mai/EasyI18n-Go"
)

// 程序配置
const (
    configFilePath = "/etc/clipboard-go/config.json"
    maxTimestamps  = 20
)

// 记录使用过的时间戳及其顺序
var usedTimestamps []string

// RequestData 期望请求体
type RequestData struct {
    Data string `json:"data"`
}

// Config HTTP服务器配置
type Config struct {
    Port         string `json:"port"`
    AuthPassword string `json:"auth_password"`
}

// config 加载配置
var config Config

// loadConfig 加载配置
func loadConfig() error {
    file, err := os.Open(configFilePath)
    if err != nil {
        return err
    }
    defer func() {
        if err := file.Close(); err != nil {
            log.Printf("Error closing file: %v", err) // 记录关闭文件的错误
        }
    }()

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return err
    }
    return nil
}

// setClipboard 设置剪贴板
func setClipboard(data string) error {
    if data == "" {
        return fmt.Errorf(i18n.Translate("Input data is empty"))
    }
    cmd := exec.Command("xclip", "-selection", "clipboard")
    cmd.Stdin = strings.NewReader(data)
    return cmd.Run()
}

// getClipboard 从剪贴板获取数据
func getClipboard() (string, error) {
    cmd := exec.Command("xclip", "-selection", "clipboard", "-o")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return string(output), nil
}

// md5Hash 计算输入字符串的 MD5 哈希值
func md5Hash(input string) string {
    hash := md5.Sum([]byte(input))
    return hex.EncodeToString(hash[:])
}

// authenticate 进行身份验证
func authenticate(timestamp, hash string) bool {
    return md5Hash(timestamp+config.AuthPassword) == hash
}

// isTimestampValid 检查时间戳是否有效（不超过10秒，并且未使用过）
func isTimestampValid(timestamp string) bool {
    tsInt, err := strconv.ParseInt(timestamp, 10, 64)
    if err != nil {
        return false
    }
    ts := time.Unix(tsInt, 0)
    if time.Since(ts) > 10*time.Second {
        return false
    }

    // 检查是否已经使用过
    for _, usedTs := range usedTimestamps {
        if usedTs == timestamp {
            return false
        }
    }

    // 记录时间戳
    if len(usedTimestamps) >= maxTimestamps {
        usedTimestamps = usedTimestamps[1:]
    }
    usedTimestamps = append(usedTimestamps, timestamp)
    return true
}

// handleRequest 处理客户端请求
func handleRequest(w http.ResponseWriter, r *http.Request) {
    timestamp := r.Header.Get("Timestamp")
    md5Hash := r.Header.Get("MD5")

    if !isTimestampValid(timestamp) || !authenticate(timestamp, md5Hash) {
        http.Error(w, "403 Forbidden", http.StatusForbidden)
        return
    }

    switch r.Method {
    case http.MethodGet:
        clipboardContent, err := getClipboard()
        if err != nil {
            http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
            return
        }
        _, err = w.Write([]byte(clipboardContent))
        if err != nil {
            return
        }

    case http.MethodPost:
        var requestData RequestData
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "400 Bad Request", http.StatusBadRequest)
            return
        }
        if err := json.Unmarshal(body, &requestData); err != nil {
            http.Error(w, "400 Bad Request", http.StatusBadRequest)
            return
        }
        if err := setClipboard(requestData.Data); err != nil {
            http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusNoContent)

    default:
        http.Error(w, "400 Bad Request", http.StatusBadRequest)
    }
}

// initServer 初始化 HTTP 服务器
func initServer() *http.Server {
    server := &http.Server{
        Addr:    config.Port,
        Handler: http.DefaultServeMux,
    }

    http.HandleFunc("/", handleRequest)
    log.Println(i18n.Translate("Server is starting, listening on port"), config.Port)

    return server
}

func main() {
    // 捕捉中断信号，安全关闭服务器
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

    // 自动设置翻译
    i18n.SetCustomTranslations(EasyI18nTranslations)
    i18n.InitLanguage()

    if err := loadConfig(); err != nil {
        log.Fatal(i18n.Translate("Failed to load config:"), err)
    }

    server := initServer()

    // 在 goroutine 中启动 server
    go func() {
        if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatal(i18n.Translate("Server startup failed:"), err)
        }
    }()

    // 等待中断信号
    <-sigs
    log.Println(i18n.Translate("Shutting down server..."))

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal(i18n.Translate("Server forced to shutdown:"), err)
    }
}
