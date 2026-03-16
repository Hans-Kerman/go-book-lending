package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Hans-Kerman/go-book-lending/backend/config"
	"github.com/Hans-Kerman/go-book-lending/backend/pkg"
	"github.com/Hans-Kerman/go-book-lending/backend/routers"
	"github.com/lmittmann/tint"
)

// @title Go Book Lending API
// @version 1.0
// @description This is a book lending system API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description "Type 'Bearer YOUR_JWT_TOKEN' to authenticate."
func main() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.Kitchen,
			AddSource:  true,
			NoColor:    false,
		}),
	))

	if err := pkg.InitStaticDir(); err != nil {
		slog.Error("error when touch pictures folder", "error", err)
		os.Exit(1)
	}
	if err := config.InitConfig(); err != nil {
		slog.Error("error when init config", "error", err)
		os.Exit(1)
	}
	if err := config.InitDataBase(); err != nil {
		slog.Error("error when connect database", "error", err)
		os.Exit(1)
	}

	r := routers.SetupRouter()
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(config.AppConfig.Server.Port),
		Handler: r,
	}

	go func() {
		fmt.Printf("Starting server on :%v\n", config.AppConfig.Server.Port) // 提前打印
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error when start the server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit //阻塞住主程序
	slog.Info("closing the server...")

	//收到信号则开始关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 10)

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("error when shutdown", "error", err)
		os.Exit(1)
	}
	slog.Info("Server exiting...")
}
