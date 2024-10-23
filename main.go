package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// 服务器配置
const (
	port          = ":8777"          // 服务器监听端口
	authPassword   = "1234"          // 硬编码的密码，用于身份验证
	maxTimestamps  = 20               // 允许的最大时间戳数量
)

// 记录使用过的时间戳及其顺序
var usedTimestamps []string

// setClipboard 将数据设置到剪贴板
func setClipboard(data string) error {
	if data == "" {
		return fmt.Errorf("输入数据为空，无法设置剪贴板")
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
		return "", fmt.Errorf("无法获取剪贴板内容: %v", err)
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
	// 解析时间戳
	ts, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return false // 无法解析时间戳
	}

	// 检查与当前时间的差异
	if time.Since(ts) > 10*time.Second {
		return false // 超过10秒
	}

	// 检查是否已经使用过
	for _, usedTs := range usedTimestamps {
		if usedTs == timestamp {
			return false // 已使用过
		}
	}

	// 记录时间戳
	if len(usedTimestamps) >= maxTimestamps {
		// 超过最大记录数，移除最旧的时间戳
		usedTimestamps = usedTimestamps[1:] // 移除最旧的
	}
	usedTimestamps = append(usedTimestamps, timestamp) // 记录当前时间戳

	return true
}

// handleRequest 处理客户端请求
func handleRequest(w http.ResponseWriter, r *http.Request) {
	timestamp := r.Header.Get("Timestamp")
	md5Hash := r.Header.Get("MD5")

	// 检查时间戳有效性
	if !isTimestampValid(timestamp) {
		http.Error(w, "403 Forbidden: Invalid or expired timestamp", http.StatusForbidden)
		return
	}

	if !authenticate(timestamp, md5Hash) {
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
		w.Write([]byte(clipboardContent))
	case http.MethodPost:
		body, err := io.ReadAll(r.Body) // 去掉长度限制
		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}
		if err := setClipboard(string(body)); err != nil {
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
	fmt.Println("服务器在端口" + port + "上监听...")
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}

// main 函数
func main() {
	// 捕捉中断信号，安全关闭服务器
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go startServer()

	// 等待中断信号
	<-sigs
	fmt.Println("\n服务器已关闭")
}
