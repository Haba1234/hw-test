package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/calendar_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	addrServer := net.JoinHostPort(config.Server.Address, config.Server.Port)
	server := internalhttp.NewServer(addrServer, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel = context.WithTimeout(ctx, time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
