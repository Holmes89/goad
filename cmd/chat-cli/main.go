package main

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"io"
	"smail/internal/message"
	"time"
)

func main() {
	conn, _ := grpc.Dial(":8080", grpc.WithInsecure())
	defer conn.Close()
	client := message.NewMessengerClient(conn)

	clientDeadline := time.Now().Add(10 * time.Hour)
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	messenger, _ := client.SendMessage(ctx)

	app := tview.NewApplication()
	inputField := tview.NewInputField().
		SetLabel("Message: ").
		SetFieldWidth(255)
	inputField.SetDoneFunc(func(key tcell.Key) {
			m := &message.Message{
				Uuid: "ABC",
				Room: "Unknown",
				From: "Client",
				Body: inputField.GetText(),
			}
		messenger.Send(m)
			inputField.SetText("")
		})

	messageArea := tview.NewTextView()
	messageArea.SetTitle("Messages").SetBorder(true)
	go func() {
		for {
			resp, err := messenger.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				app.Stop()
				continue
			}
			fmt.Fprintf(messageArea, "%s: %s\n", resp.From, resp.Body)
			app.Draw()
		}
	}()


	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(messageArea, 0, 10, false).
			AddItem(inputField, 1, 1, true), 0, 4, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Commands"), 20, 1, false)
	if err := app.SetRoot(flex, true).SetFocus(inputField).Run(); err != nil {
		conn.Close()
		panic(err)
	}

}

