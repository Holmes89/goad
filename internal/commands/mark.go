package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"smail/internal/mail"
)

func NewMarkCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "mark [read unread] [id]",
		Short: "Mark message as read or unread",
		Long:  "update message status",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			service := mail.NewMailService(":8080")
			defer service.Close()

			id := args[1]

			username := viper.Get("username").(string)
			switch args[0] {
			case "read":
				service.UpdateMessageStatus(id, username, false)
			case "unread":
				service.UpdateMessageStatus(id, username, true)
			default:
				fmt.Println("Invalid param")
			}
		},
	}
}
