package main

import (
	"github.com/Phanile/go-exchange-auth/internal/app"
	"github.com/Phanile/go-exchange-auth/internal/config"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(envLocal)
	setupEnv(envLocal)

	application := app.NewApp(log, cfg)
	go application.GRPCServer.MustRun()
	go handleMetrics()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupEnv(env string) {
	switch env {
	case envLocal:
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
	}
}

func handleMetrics() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":2111", mux)

	if err != nil {
		panic("error starting metrics server: " + err.Error())
	}
}
