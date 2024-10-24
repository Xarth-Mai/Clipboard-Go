package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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

// 服务器配置
const (
	port          = ":8777"
	authPassword  = "1234"
	maxTimestamps = 20
)

// 记录使用过的时间戳及其顺序
var usedTimestamps []string

type RequestData struct {
	Data string `json:"data"`
}

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
	return md5Hash(timestamp+authPassword) == hash
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

// startServer 启动 HTTP 服务器
func startServer() {
	http.HandleFunc("/", handleRequest)
	log.Println(i18n.Translate("Server is starting, listening on port"), port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(i18n.Translate("Server startup failed:"), err)
	}
}

func main() {
	// 捕捉中断信号，安全关闭服务器
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 自动设置翻译
	i18n.SetCustomTranslations(EasyI18nTranslations)
	i18n.InitLanguage()

	go startServer()

	// 等待中断信号
	<-sigs
	fmt.Println("\n\b", i18n.Translate("Server is down"))
}
