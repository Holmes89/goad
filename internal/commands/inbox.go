package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"smail/internal/mail"
)

func NewInboxCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "inbox [all]",
		Short: "Get messages for user",
		Long:  "gets messages for user, can define if all messages should be returned or just read",
		Run: func(cmd *cobra.Command, args []string) {

			service := mail.NewMailService(":8080")
			defer service.Close()

			all := (len(args) > 0) && (args[0] == "all")

			//TODO fetch username from config
			ms, err := service.GetMessages("test", all)
			if err != nil {
				fmt.Println("Unable to send message")
			} else {
				fmt.Printf("ID\t\t     Unread\t\tSubject\n\n")
				for _, m := range ms {
					status := "*"
					if m.Read{
						status = " "
					}
					fmt.Printf("%s\t\t%s\t\t%s\n", m.Uuid, status, m.Subject)
				}
				fmt.Printf("\n\n\n")
			}
		},
	}
}

