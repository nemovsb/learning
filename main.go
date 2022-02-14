package main

import (
	"errors"
	"fmt"
	errors2 "github.com/pkg/errors"
	"go.uber.org/zap"
	"learning/internal/app"
	"learning/internal/configuration/di"
	"learning/internal/http_server"
	"learning/internal/storage"
	"learning/pkg/zaplogger"
	"log"
	"os"
	"os/signal"
	"syscall"

	group "github.com/oklog/run"
)

var ErrOsSignal = errors.New("got os signal")

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /

// @securityDefinitions.basic  BasicAuth
func main() {

	config, err := di.ViperConfigurationProvider(os.Getenv("GOLANG_ENVIRONMENT"), false)
	if err != nil {
		log.Fatal("Read config error: ", err)
	}

	logger, zapLoggerCleanup, err := zaplogger.Provider(config.ZapLoggerMode)
	if err != nil {
		log.Fatal(errors2.WithMessage(err, "zap logger provider"))
	}

	logger.Info("application", zap.String("event", "initializing"))
	logger.Info("application", zap.Any("resolved_configuration", config))

	redisConnection, err := storage.NewRedisClient(di.GetRedisConfig(config), logger)
	if err != nil {
		logger.Fatal("redis can't start")
		return
	}

	//TCPConnection, err := app.NewTCPConnect(di.GetTCPConfig(config))

	redisHandler := app.NewIncValueHandler(redisConnection, logger)
	hmacHandler := app.NewHMACHandler(logger)
	tcpHandler, err := app.NewTCPHandler(di.GetTCPConfig(config), app.NewStringConvertTCP(logger, di.GetTCPConfig(config)), logger)
	fmt.Printf("\ntcp handler: %+v\n", tcpHandler)
	if err != nil {
		log.Println("tcp service can't start:", err)
	}
	handlerSet := http_server.NewHandlerSet(redisHandler, hmacHandler, tcpHandler)

	handler := http_server.NewRouter(handlerSet)

	server := http_server.NewServer(di.GetHTTPServerConfig(config), handler)

	var (
		serviceGroup        group.Group
		interruptionChannel = make(chan os.Signal, 1)
	)

	serviceGroup.Add(func() error {
		signal.Notify(interruptionChannel, syscall.SIGINT, syscall.SIGTERM)
		osSignal := <-interruptionChannel

		return fmt.Errorf("%w: %s", ErrOsSignal, osSignal)
	}, func(error) {
		interruptionChannel <- syscall.SIGINT
	})

	serviceGroup.Add(func() error {
		logger.Info("server", zap.String("event", "HTTP API started"))

		return server.Run()
	}, func(error) {
		// Graceful http server shutdown
		err = server.Shutdown()
		logger.Info("shutdown Http Server error", zap.Error(err))
	})

	err = serviceGroup.Run()
	logger.Info("services stopped", zap.Error(err))

	err = redisConnection.Close()
	logger.Info("close Redis connection error", zap.Error(err))

	zapLoggerCleanup()

	//err = server.Run()
}
