package cmd

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SawitProRecruitment/UserService/core"
	"github.com/SawitProRecruitment/UserService/core/service/authsvc"
	"github.com/SawitProRecruitment/UserService/generated"
	sawithttp "github.com/SawitProRecruitment/UserService/handler/http"
	"github.com/SawitProRecruitment/UserService/repository/postgres"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Env    string       `json:"env"`
	Auth   AuthConfig   `json:"auth"`
	Server ServerConfig `json:"http"`
	DB     PsqlConfig   `json:"postgresql"`
}

type ServerConfig struct {
	Address           string        `json:"address"`
	PrefixPath        string        `json:"prefixPath"`
	ReadTimeout       time.Duration `json:"readTimeout"`
	WriteTimeout      time.Duration `json:"writeTimeout"`
	ReadHeaderTimeout time.Duration `json:"readHeaderTimeout"`
}

type PsqlConfig struct {
	Host        string        `json:"host"`
	Port        int           `json:"port"`
	Database    string        `json:"db"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	SSLMode     bool          `json:"sslMode"`
	MaxOpenConn int           `json:"maxOpenConn"`
	MaxIdleConn int           `json:"maxIdleConn"`
	MaxIdleTime time.Duration `json:"maxIdleTime"`
}

type AuthConfig struct {
	TokenPrivateKeyPath string        `json:"privateKeyPath"`
	TokenPublicKeyPath  string        `json:"publicKeyPath"`
	TokenExpDuration    time.Duration `json:"tokenExpDuration"`
	EncryptSecretKey    string        `json:"encryptSecretKey"`
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "service",
	Short: "SawitPro service",
}

func initConfig() Config {
	var cfg Config

	f := func() error {
		// required: base config file
		viper.AddConfigPath("config")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		// optional: config file for local development
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
		if err := godotenv.Load("config/.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Println("Config file changed:", e.Name)
		})

		return viper.Unmarshal(&cfg, func(dc *mapstructure.DecoderConfig) {
			dc.TagName = "json"
		})
	}

	err := f()
	if err != nil {
		panic(err)
	}

	return cfg
}

func initLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}

func initPostgres(ctx context.Context, cfg PsqlConfig) (*postgres.Repository, error) {
	opts := postgres.NewRepoOptions{
		Username:    cfg.Username,
		Password:    cfg.Password,
		Host:        cfg.Host,
		Port:        cfg.Port,
		Database:    cfg.Database,
		SSLMode:     cfg.SSLMode,
		MaxIdleConn: cfg.MaxIdleConn,
		MaxOpenConn: cfg.MaxOpenConn,
		MaxIdleTime: cfg.MaxIdleTime,
	}

	repo, err := postgres.New(ctx, opts)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func initServer(cfg ServerConfig, handler *sawithttp.Handler) http.Server {
	e := echo.New()
	e.Use(handler.MiddlewareLogging)
	e.Use(handler.MiddlewareError)

	generated.RegisterHandlersWithBaseURL(e, handler, cfg.PrefixPath)
	s := http.Server{
		Addr:              cfg.Address,
		Handler:           e,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	return s
}

func initAuthSvc(cfg AuthConfig, repo core.UserRepo) (*authsvc.Service, error) {
	opts := authsvc.ServiceOpts{
		PrvKeyPath:       cfg.TokenPrivateKeyPath,
		PubKeyPath:       cfg.TokenPublicKeyPath,
		TokenExpDuration: cfg.TokenExpDuration,
		EncryptSecretKey: cfg.EncryptSecretKey,
	}

	return authsvc.New(opts, repo)
}
