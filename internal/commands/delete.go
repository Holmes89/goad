package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"smail/internal/mail"
)

func NewDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete Message",
		Long:  "delete a message given an id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			service := mail.NewMailService(":8080")
			defer service.Close()

			//TODO fetch username from config
			err := service.DeleteMessage(args[0], "test")
			if err != nil {
				fmt.Println("Unable to delete message")
			}
		},
	}
}

