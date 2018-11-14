package commands

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"smail/internal/mail"
)

func NewSendCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "send [to] [subject]",
		Short: "Send a message to someone",
		Long:  "sends a message to a user given a user",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			service := mail.NewMailService(":8080")
			defer service.Close()

			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Please type message:")
			text, _ := reader.ReadString('\n')

			to := args[0]

			sub := "No Subject"
			if len(args) == 2 {
				sub = args[1]
			}

			err := service.Send(to, "not implemented", sub, text)
			if err != nil {
				fmt.Println("Unable to send message")
			} else {
				fmt.Println("Message Sent!")
			}
		},
	}
}

