package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/twotwo/go-blueprint/pkg/database"
	"github.com/twotwo/go-blueprint/server"
	"github.com/twotwo/go-blueprint/server/user"
)

func main() {
	// 初始化数据库
	db, err := database.Setup()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 设置用户模块的数据库连接
	user.DB = db

	// 创建根路由
	r := chi.NewRouter()

	// 中间件
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	// CORS 中间件,允许跨域请求
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "vscode-webview://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, // 限制可以跨域的 请求头
		AllowCredentials: true,
		MaxAge:           300, // 缓存预检请求 300 秒
	}))

	// 注册API路由
	server.RegisterRoutes(r)

	// 服务配置
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// 服务启动
	go func() {
		log.Printf("服务启动在 :%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听错误: %s\n", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭错误: %s\n", err)
	}

	log.Println("服务已关闭")
}
