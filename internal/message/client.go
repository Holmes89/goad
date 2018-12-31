package message

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"sync"
)

type Client interface{
	Send(message *Message)
	Run()
	Close()
}

type grpcClient struct {
	hub Hub
	messenger Messenger_SendMessageServer
	wg sync.WaitGroup
}

func NewGRPCClient(hub Hub, messenger Messenger_SendMessageServer) Client {
	c := &grpcClient{
		hub: hub,
		messenger: messenger,
	}
	logrus.Info("client created")
	return c
}

func (c *grpcClient) Send(message *Message) {
	c.messenger.Send(message)
}

func (c *grpcClient) Close() {
	c.hub.Unregister(c)
}

func (c *grpcClient) Run() {
	for {
		resp, err := c.messenger.Recv()
		if err == io.EOF {
			logrus.Error("eof")
			c.Close()
			break
		}
		if err != nil {
			handleErrorMessage(err)
			c.Close()
			break
		}

		//TODO Add uuid, timestamp
		message := &Message{
			From: resp.From,
			Body: resp.Body,
		}
		c.hub.Broadcast(message)
	}
}

func handleErrorMessage(err error) {
	status, _ := status.FromError(err)
	if codes.Canceled == status.Code() {
		logrus.Info("connection closed")
	} else {
		logrus.WithField("error", status.Message()).Error("received error from client")
	}
}