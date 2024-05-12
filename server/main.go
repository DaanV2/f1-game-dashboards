/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"
	"time"

	"github.com/DaanV2/f1-game-dashboards/server/cmd"
	"github.com/charmbracelet/log"

	_ "go.uber.org/automaxprocs"
)

func init() {
	// Initialize the default logger.
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller: true,
		ReportTimestamp: true,
		TimeFormat: time.DateTime,
	})
	log.SetDefault(logger)
}

func main() {
	cmd.Execute()
}
