package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var serverLogsFlag bool
var dbLogsFlag bool

func init() {
	logReaderCmd.Flags().BoolVarP(&serverLogsFlag, "server", "s", false, "server error logs")
	logReaderCmd.Flags().BoolVarP(&dbLogsFlag, "db", "d", false, "database logs")

	rootCmd.AddCommand(logReaderCmd)
}

var logReaderCmd = &cobra.Command{
	Use:   "log-reader",
	Short: "Log reader",
	Long:  `Log reader`,
	Run: func(cmd *cobra.Command, args []string) {
		if serverLogsFlag == dbLogsFlag {
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if line, err := parseLine(scanner.Bytes(), serverLogsFlag, dbLogsFlag); err == nil {
				fmt.Println(line)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	},
}

type errorLog struct {
	Level   string    `json:"level"`
	Time    time.Time `json:"time"`
	Caller  string    `json:"caller"`
	Message string    `json:"message"`
}

func parseLine(line []byte, serverLogs, dbLogs bool) (string, error) {
	if serverLogs {
		return parseLineServer(line)
	} else if dbLogs {
		return string(line), nil
	}
	return "", errors.New("invalid flag")
}

func parseLineServer(line []byte) (string, error) {
	var errLog errorLog
	err := json.Unmarshal(line, &errLog)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("%s %7s | %s | %s",
		errLog.Time.Format(time.RFC3339),
		displayLevel(errLog.Level),
		errLog.Caller,
		errLog.Message)

	return result, nil
}
