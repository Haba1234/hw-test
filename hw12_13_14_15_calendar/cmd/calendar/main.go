package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/server/http"
	initdb "github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage/initdb"
	"github.com/joho/godotenv"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/calendar_config.toml", "Path to configuration file")
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
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

	// Connect to base.
	dsn := fmt.Sprintf("user=%s password=%s sslmode=%s", config.Storage.User, config.Storage.Password, config.Storage.SSLMode)
	storage, err := initdb.New(context.Background(), config.Storage.Type, dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	calendar := app.New(logg, storage)

	http := internalhttp.NewServer(logg, calendar)
	grpc := internalgrpc.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := http.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if err := grpc.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
		if err := calendar.Close(ctx); err != nil {
			logg.Error("failed close storage: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		addrServer := net.JoinHostPort(config.Server.Address, config.Server.GRPCPort)
		if err := grpc.Start(ctx, addrServer); err != nil {
			logg.Error("failed to start gRPC server: " + err.Error())
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		addrServer := net.JoinHostPort(config.Server.Address, config.Server.HTTPPort)
		if err := http.Start(ctx, addrServer); err != nil {
			logg.Error("failed to start HTTP server: " + err.Error())
			cancel()
		}
	}()

	<-ctx.Done()
	wg.Wait()

	logg.Warn("calendar stoped.")
}
