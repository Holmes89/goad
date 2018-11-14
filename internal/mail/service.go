package mail

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

type mailService struct {
	conn *grpc.ClientConn
}

type MailService interface{
	Send(to, from, sub, message string) error
	GetMessages(username string, all bool) ([]*Mail, error)
	GetSingleMessage(username, id string) (*Mail, error)
	UpdateMessageStatus(id string, unread bool) error
	DeleteMessage(id, username string) error
	Close()
}

func NewMailService(endpoint string) MailService{

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		log.Panic("did not connect: %s", err)
	}

	return &mailService{
		conn,
	}
}

func (s *mailService) Close() {
	s.conn.Close()
}

func (s *mailService) Send(to, from, sub, message string) error {
	c := NewMailerClient(s.conn)

	m := &Mail{
		To:      to,
		Subject: sub,
		From:    "placeholder",
		Body:    message,
	}

	_, err := c.SendMail(context.Background(), m)
	return err
}

func (s *mailService) GetMessages(username string, all bool) ([]*Mail, error) {
	c := NewMailerClient(s.conn)

	r, err := c.GetMail(context.Background(), &GetMailRequest{
		Username: username,
		All: all,
	})

	return r.Mail, err
}


func (s *mailService) GetSingleMessage(username, id string) (*Mail, error) {
	c := NewMailerClient(s.conn)

	r, err := c.GetMail(context.Background(), &GetMailRequest{
		Username: username,
		Uuid: id,
	})

	if r == nil || len(r.Mail) == 0 {
		return nil, err
	}

	s.UpdateMessageStatus(id, false)
	return r.Mail[0], err
}

func (s *mailService) UpdateMessageStatus(id string, unread bool) error {
	c := NewMailerClient(s.conn)

	//TODO fix this
	m := &MailID{
		Username: "test",
		Uuid: id,
	}

	if unread {
		if _, err := c.MarkUnread(context.Background(), m); err != nil {
			return err
		}
	} else {
		if _, err := c.MarkRead(context.Background(), m); err != nil {
			return err
		}
	}

	return nil

}

func (s *mailService) DeleteMessage(id, username string) error {
	c := NewMailerClient(s.conn)

	//TODO fix this
	m := &MailID{
		Username: "test",
		Uuid: id,
	}

	_, err := c.DeleteMail(context.Background(), m)
	return err
}