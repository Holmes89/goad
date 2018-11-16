package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

			username := viper.Get("username").(string)
			err := service.DeleteMessage(args[0], username)
			if err != nil {
				fmt.Println("Unable to delete message")
			}
		},
	}
}

