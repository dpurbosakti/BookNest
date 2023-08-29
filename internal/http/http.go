package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server func(db *gorm.DB, port string)

func Serve(db *gorm.DB, port string) {
	gin.ForceConsoleColor()
	r := gin.New()
	r.Use(gin.Logger())

	InitRouter(r, db)
	port = fmt.Sprintf(":%s", port)

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatalf("server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	q := <-quit
	logrus.Info("Shutdown server ...")
	logrus.Info("Exit signal: ", q)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server shutdown: ", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()

}
