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

// initConfigLoggerDatabase initializes configuration, logger and database.
// Logger and database initialization are not required.
func initConfigLoggerDatabase(initLogger, initDatabase bool) (logger *zap.Logger, database *db.DB, err error) {
	// Configuration initialization
	// ----------------------------
	if err = initConfig(); err != nil {
		return nil, nil, err
	}

	// Logger initialization
	// ---------------------
	if initLogger {
		logger, err = server.InitLogger()
		if err != nil {
			return nil, nil, err
		}
		defer logger.Sync()
	}

	// Database connection
	// -------------------
	if initDatabase {
		database, err = db.New(&db.DatabaseConfig{
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
	}

	return logger, database, err
}

func displayLogLevel(l string) aurora.Value {
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

func displayLogMethod(m string) aurora.Value {
	switch m {
	case "GET":
		return aurora.Cyan(m)
	case "POST":
		return aurora.Blue(m)
	case "PUT":
		return aurora.Brown(m)
	case "PATCH":
		return aurora.Magenta(m)
	case "DELETE":
		return aurora.Red(m)
	default:
		return aurora.Gray(12, m)
	}
}

func displayLogStatusCode(c uint) aurora.Value {
	if c < 200 {
		return aurora.Cyan(c)
	} else if c < 300 {
		return aurora.Green(c)
	} else if c < 400 {
		return aurora.Magenta(c)
	} else if c < 500 {
		return aurora.Brown(c)
	} else {
		return aurora.Red(c)
	}
}
