package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"smail/internal/mail"
)

func NewReadCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "read [id]",
		Short: "Read Message",
		Long:  "get all content for message by id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			service := mail.NewMailService(":8080")
			defer service.Close()

			username := viper.Get("username").(string)
			m, err := service.GetSingleMessage(username, args[0])
			if err != nil {
				fmt.Println("Unable to read message")
				return
			}

			if m == nil{
				fmt.Println("Invalid ID")
			} else {
				fmt.Printf("\n\nFrom: %s\n", m.From)
				fmt.Printf("Subject: %s\n", m.Subject)
				fmt.Printf("Message:\n\n%s\n\n", m.Body)
			}
		},
	}
}

