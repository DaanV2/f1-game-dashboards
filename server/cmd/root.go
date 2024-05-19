/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "", // TODO: Add a short description here
	Long:  ``, // TODO: Add a short description here
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logOptions := log.Options{
			TimeFormat: time.DateTime,
			ReportCaller: cmd.Flag("log-report-caller").Value.String() == "true",
		}

		// log-level
		level, err := log.ParseLevel(cmd.Flag("log-level").Value.String())
		if err != nil {
			log.Fatal("invalid log level", "error", err)
		}
		logOptions.Level = level

		// log-format
		switch cmd.Flag("log-format").Value.String() {
		default:
			logOptions.Formatter = log.TextFormatter
		case "json":
			logOptions.Formatter = log.JSONFormatter
		case "logfmt":
			logOptions.Formatter = log.LogfmtFormatter
		}

		// Initialize the default logger.
		logger := log.NewWithOptions(os.Stderr, logOptions)
		log.SetDefault(logger)

		maxprocs.Set(maxprocs.Logger(func(s string, i ...interface{}) {
			msg := fmt.Sprintf(s, i...)
			log.Info(msg)
		}))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	pFlags := rootCmd.PersistentFlags()
	pFlags.String("log-level", "info", "The log level to use (debug, info, warn, error, fatal)")
	pFlags.String("log-format", "text", "The log format to use (text, json, logfmt)")
	pFlags.Bool("log-report-caller", true, "Whether to report the caller location")

	pFlags.String("storage-type", "files", "Storage type to use (files)")
	pFlags.String("files-storage-directory", "", "The directory to store files in (default: ./data/files)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
