package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SSEHandler struct {
	model *ModelGenerator
}

func (h *SSEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	for token := range h.model.Stream(ctx, r.URL.Query().Get("prompt")) {
		fmt.Fprintf(w, "data: %s\n\n", token)
		flusher.Flush() // 关键：立即发送到客户端
	}
}

func main() {
	model := &ModelGenerator{bufferSize: 10}
	handler := &SSEHandler{model: model}

	mux := http.NewServeMux()
	mux.Handle("/stream", SafeStream(handler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		// 调优参数
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// 启动监控
	go monitorConnections()

	log.Fatal(server.ListenAndServe())
}

func monitorConnections() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		// 获取当前连接数（需与前面metrics集成）
		log.Printf("活跃连接数: %d", getActiveConnections())
	}
}

// 获取活跃连接数
func getActiveConnections() int {
	// 这里实现获取连接数的逻辑
	// 由于是示例代码，暂时返回固定值
	return 0
}

func HijackStream(w http.ResponseWriter) (net.Conn, *bufio.ReadWriter) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		panic("hijack not supported")
	}

	conn, rw, err := hj.Hijack()
	if err != nil {
		panic(err)
	}

	// 设置TCP参数
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true) // 禁用Nagle算法
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)

		// 设置内核级缓冲区
		if err := tcpConn.SetWriteBuffer(128 * 1024); err != nil {
			log.Printf("设置写缓冲失败: %v", err)
		}
	}

	return conn, rw
}

type StreamPipeline struct {
	inputChan  chan *StreamRequest
	outputChan chan *StreamResponse
	model      *ModelGenerator
}

// 流式请求结构体
type StreamRequest struct {
	ID     string
	Ctx    context.Context
	Prompt string
}

// 流式响应结构体
type StreamResponse struct {
	ID     string
	Tokens []string
}

func (p *StreamPipeline) StartWorkers(num int) {
	for i := 0; i < num; i++ {
		go func() {
			for req := range p.inputChan {
				ctx, cancel := context.WithCancel(req.Ctx)

				// 二级缓冲管道
				intermediate := make(chan string, 10)
				go func() {
					defer close(intermediate)
					for token := range p.model.Stream(ctx, req.Prompt) {
						intermediate <- token
					}
				}()

				// 组装最终响应
				res := &StreamResponse{ID: req.ID}
				for token := range intermediate {
					res.Tokens = append(res.Tokens, token)
					if len(res.Tokens)%5 == 0 { // 每5个token发送一次
						p.outputChan <- res
						res = &StreamResponse{ID: req.ID}
					}
				}
				cancel()
			}
		}()
	}
}

func SafeStream(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("流式异常: %v", r)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		// 连接状态检测
		if _, err := w.Write(nil); err != nil {
			log.Printf("连接已断开: %v", err)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func StartServer() {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 流式处理器
		}),
	}

	// 优雅关闭
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("关闭错误: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// 模拟大模型生成器
type ModelGenerator struct {
	bufferSize int
}

func (m *ModelGenerator) Stream(ctx context.Context, prompt string) <-chan string {
	out := make(chan string, m.bufferSize)
	go func() {
		defer close(out)
		for i := 0; i < 50; i++ { // 模拟50个token生成
			select {
			case <-ctx.Done():
				log.Println("生成中断")
				return
			case <-time.After(100 * time.Millisecond): // 模拟计算延迟
				out <- fmt.Sprintf("token-%d", i)
			}
		}
	}()
	return out
}
