package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/0x726f6f6b6965/web3-auth/internal/api"
	"github.com/0x726f6f6b6965/web3-auth/internal/api/router"
	appCfg "github.com/0x726f6f6b6965/web3-auth/internal/config"
	"github.com/0x726f6f6b6965/web3-auth/internal/utils"
	"github.com/0x726f6f6b6965/web3-auth/pkg/dynamo"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func main() {
	godotenv.Load()
	path := os.Getenv("CONFIG")
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("read yaml error", err)
		return
	}
	var cfg appCfg.AppConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("unmarshal yaml error", err)
		return
	}

	if cfg.IsDevEnv() {
		dynamo.NewDevLocalClient(cfg.DynamoDB.Host, cfg.DynamoDB.Table)
	} else {
		awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(cfg.DynamoDB.Region))
		if err != nil {
			log.Fatalf(fmt.Sprintf("Failed to create dynamo client: %s", err))
		}
		dynamo.NewDynamoClient(context.Background(), awsCfg, cfg.DynamoDB.Table)
	}

	secret, err := os.ReadFile(cfg.JwtSecretKey)
	if err != nil {
		log.Fatal("read secret key error", err)
		return
	}

	if secret == nil {
		log.Fatal("secret key is nil")
		return
	}
	utils.JwtSecretKey = secret

	api.NewAuthAPI()
	api.NewUserAPI()

	if err := startServer(&cfg); err != nil {
		log.Fatalf(fmt.Sprintf("Failed to start server: %s", err))
	}
}

func startServer(cfg *appCfg.AppConfig) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpPort),
		Handler: initEngine(cfg),
	}
	ctx, cancel := context.WithCancel(context.Background())

	go listenToSystemSignals(cancel)

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf(fmt.Sprintf("Failed to shutdown server: %s", err))
		}
	}()
	log.Println("Server started success")
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server was shutdown gracefully")
		return nil
	}
	return err
}

func initEngine(cfg *appCfg.AppConfig) *gin.Engine {
	gin.SetMode(func() string {
		if cfg.IsDevEnv() {
			return gin.DebugMode
		}
		return gin.ReleaseMode
	}())
	engine := gin.New()
	engine.Use(cors.Default())
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "Service internal exception!",
		})
	}))
	router.RegisterRoutes(engine)
	return engine
}

func listenToSystemSignals(cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	cancel()
}
