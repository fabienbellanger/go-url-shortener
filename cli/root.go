package cli

import (
	"time"

	server "github.com/fabienbellanger/go-url-shortener"
	"github.com/fabienbellanger/go-url-shortener/db"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "go-url-shortener",
	Short:   "A URL shortener application",
	Long:    "A URL shortener application written with Fiber and GORM",
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

// initConfig initializes configuration from config file.
func initConfig() error {
	viper.SetConfigFile(".env")
	return viper.ReadInConfig()
}

// TODO: Choose what to initialize
func initConfigLoggerDatabase() (*zap.Logger, *db.DB, error) {
	// Configuration initialization
	// ----------------------------
	if err := initConfig(); err != nil {
		return nil, nil, err
	}

	// Logger initialization
	// ---------------------
	logger, err := server.InitLogger()
	if err != nil {
		return nil, nil, err
	}
	defer logger.Sync()

	// Database connection
	// -------------------
	db, err := db.New(&db.DatabaseConfig{
		Driver:          viper.GetString("DB_DRIVER"),
		Host:            viper.GetString("DB_HOST"),
		Username:        viper.GetString("DB_USERNAME"),
		Password:        viper.GetString("DB_PASSWORD"),
		Port:            viper.GetInt("DB_PORT"),
		Database:        viper.GetString("DB_DATABASE"),
		Charset:         viper.GetString("DB_CHARSET"),
		Collation:       viper.GetString("DB_COLLATION"),
		Location:        viper.GetString("DB_LOCATION"),
		MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
		MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
		ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME") * time.Hour,
	})
	if err != nil {
		return nil, nil, err
	}

	return logger, db, err
}

func displayLevel(l string) aurora.Value {
	switch l {
	case "DEBUG":
		return aurora.Cyan(l)
	case "INFO":
		return aurora.Green(l)
	case "WARN":
		return aurora.Brown(l)
	default:
		return aurora.Red(l)
	}
}
