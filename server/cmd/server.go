/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DaanV2/f1-game-dashboards/server/api"
	"github.com/DaanV2/f1-game-dashboards/server/pkg/data"
	"github.com/DaanV2/f1-game-dashboards/server/sessions"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "TODO",
	Long:  `TODO`,
	Run:   ServerCmd,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ServerCmd(cmd *cobra.Command, args []string) {
	chairs := sessions.NewChairManager()
	server := api.NewApiServer(chairs)
	database, err := data.NewStorage(cmd.Flags())
	if err != nil {
		log.Fatal("could not create storage", "error", err)
	}

	data.DatabaseHooks(database, chairs)

	if err := server.Start(); err != nil {
		log.Fatal("could not start server", "error", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-signals

	if err := server.Stop(); err != nil {
		log.Error("could not stop server", "error", err)
	}
}
