package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"smail/internal/commands"
)

var rootCmd = &cobra.Command{
	Use:   "smail",
	Short: "Simple Mail",
	Long: `A Simple Mail client built as an example on how to use GRPC, Cobra, and Mongo`,
}

func main() {

	sendCmd := commands.NewSendCommand()
	inboxCmd := commands.NewInboxCommand()
	readCmd := commands.NewReadCommand()
	markCmd := commands.NewMarkCommand()
	deleteCmd := commands.NewDeleteCommand()
	rootCmd.AddCommand(sendCmd, inboxCmd, readCmd, markCmd, deleteCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}