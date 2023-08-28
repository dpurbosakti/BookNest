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

type Server func(db *gorm.DB) error

func Serve(db *gorm.DB) error {
	// newrelic, err := logger.SetupLogger()
	// if err != nil {
	// 	logrus.Fatal("failed to setup logger: ", err)
	// }

	file, err := os.Open("./panic.log")
	if err != nil {
		logrus.Info(err)
		file, err = os.Create("./panic.log")
		if err != nil {
			logrus.Error(err)
		}
	}

	gin.ForceConsoleColor()
	r := gin.New()
	r.Use(gin.Recovery())
	// r.Use(nrgin.Middleware(newrelic))
	r.Use(gin.RecoveryWithWriter(file))
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/metrics"}}))

	// settingCors := cors.DefaultConfig() // TODO: need to check
	// settingCors.AllowOrigins = config.GetCorsAllowedUrl()
	// settingCors.AllowCredentials = true // TODO: need to check
	// settingCors.AddAllowHeaders(config.GetCorsAllowedHeaders()...)
	// r.Use(cors.New(settingCors)) // TODO: need to check
	// //? Prometheus Middleware
	// monitor.RegisterPromMetrics() // Kalau mau nambahin metric baru, tambahin di sini
	// r.Use(monitor.PrometheusMiddleware())
	// r.GET("/metrics", gin.WrapH(promhttp.Handler())) //TODO: kalau bisa ini nanti di hide atau di encrypt dsb
	//? End Prometheus Middleware

	InitRouter(r, db)
	port := fmt.Sprintf(":%s", "8080")

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

	return err
}
