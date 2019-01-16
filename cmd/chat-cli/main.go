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

	a := &App{
		messenger: messenger,
		app: app,
	}

	a.CreateNameInput()
	if err := a.app.Run(); err != nil {
		panic(err)
	}

}

type App struct {
	messenger message.Messenger_SendMessageClient
	app *tview.Application
	name string
}

func (a *App) CreateNameInput() {
	nameInputField := tview.NewInputField().
		SetLabel("Enter a name: ").
		SetFieldWidth(150)
	nameInputField.SetDoneFunc(func(key tcell.Key) {
		a.name = nameInputField.GetText()
		if a.name == "" {
			a.name = "Ron Weasley"
		}
		a.CreateChat()
	})
	a.app.SetRoot(nameInputField, true).SetFocus(nameInputField)
}

func (a *App) CreateChat(){

	messageArea := tview.NewTextView()
	messageArea.SetTitle("Messages").SetBorder(true)

	go func() {
		for {
			resp, err := a.messenger.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				a.app.Stop()
				continue
			}
			fmt.Fprintf(messageArea, "%s: %s\n", resp.From, resp.Body)
			a.app.Draw()
		}
	}()
	inputField := tview.NewInputField().
		SetLabel("Message: ").
		SetFieldWidth(255)
	inputField.SetDoneFunc(func(key tcell.Key) {
		m := &message.Message{
			Uuid: "ABC",
			Room: "Unknown",
			From: a.name,
			Body: inputField.GetText(),
		}
		a.messenger.Send(m)
		inputField.SetText("")
	})

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(messageArea, 0, 10, false).
			AddItem(inputField, 1, 1, true), 0, 4, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Commands"), 20, 1, false)

	a.app.SetRoot(flex, true).SetFocus(inputField)
}
