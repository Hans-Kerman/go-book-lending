package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/Hans-Kerman/go-book-lending/backend/config"
	"github.com/Hans-Kerman/go-book-lending/backend/routers"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("error when init config: %s", err.Error())
	}
	if err := config.InitDataBase(); err != nil {
		log.Fatalf("error when init database: %s", err.Error())
	}

	r := routers.SetupRouter()
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(config.AppConfig.Server.Port),
		Handler: r,
	}

	go func() {
		fmt.Printf("Starting server on :%v\n", config.AppConfig.Server.Port) // 提前打印
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error when start the server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit //阻塞住主程序
	log.Println("closing the server...")

	//收到信号则开始关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 10)

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error when shutdown: %s\n", err.Error())
	}
	log.Println("Server exiting...")
}
