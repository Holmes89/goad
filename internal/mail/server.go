package mail

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type mailerServer struct {
	repo MailRepo
}

func NewMailServer(repo MailRepo) MailerServer {
	return &mailerServer{
		repo: repo,
	}
}

func (s *mailerServer) GetMail(ctx context.Context, req *GetMailRequest) (*GetMailResponse, error) {

	resp := &GetMailResponse{}

	if req.Uuid == "" {
		mail, err := s.repo.GetAllMail(req.Username, req.All)
		if err != nil {
			log.WithFields(log.Fields{
				"username": req.Username,
			}).Error("Unable to fetching emails")
			return nil, err
		}
		resp.Mail = mail
	} else {
		mail, err := s.fetchMail(req.Uuid, req.Username)
		if err != nil {
			return nil, err
		}
		resp.Mail = []*Mail{mail}
	}

	log.WithFields(log.Fields{
		"username": req.Username,
		"count": len(resp.Mail),
	}).Info("Get Mail")

	return resp, nil
}

func (s *mailerServer) SendMail(ctx context.Context, mail *Mail) (*Mail, error) {
	mail.Uuid = generateId()
	if err := s.repo.CreateMail(mail); err != nil {
		log.WithFields(log.Fields{
			"to": mail.To,
			"from": mail.From,
		}).Error("Unable to creat email")
		return nil, err
	}
	log.WithFields(log.Fields{
		"to": mail.To,
		"from": mail.From,
	}).Info("Sent Mail")
	return mail, nil
}

func (s *mailerServer) DeleteMail(ctx context.Context, mail *MailID) (*Empty, error) {
	if mail.Uuid == "" || mail.Username == "" {
		err := errors.New("Id and Username required")
		log.Error("Missing username or id")
		return nil, err
	}
	if err := s.repo.DeleteMail(mail.Uuid, mail.Username); err != nil {
		log.Error("Error deleting mail", err)
		return nil, err
	}

	log.WithFields(log.Fields{
		"id": mail.Uuid,
	}).Info("Deleted Mail")
	return &Empty{}, nil

}

func (s *mailerServer) MarkUnread(ctx context.Context, mail *MailID) (*Mail, error){
	return s.updateMailReadStatus(mail, false)
}


func (s *mailerServer) MarkRead(ctx context.Context, mail *MailID) (*Mail, error){
	return s.updateMailReadStatus(mail, true)
}

func (s *mailerServer) updateMailReadStatus(mail *MailID, status bool) (*Mail, error) {
	m, err := s.fetchMail(mail.Uuid, mail.Username)
	if err != nil {
		return nil, err
	}

	m.Read = status
	if err := s.repo.UpdateMail(m); err != nil {
		log.WithFields(log.Fields{
			"username": mail.Username,
			"id": mail.Uuid,
		}).Error("Unable to fetching emails")

		return nil, err
	}
	log.WithFields(log.Fields{
		"id": mail.Uuid,
		"read": status,
	}).Info("Updated Mail")
	return m, nil
}

func (s *mailerServer) fetchMail(id string, username string) (*Mail, error) {
	if id == "" || username == "" {
		err := errors.New("Id and Username required")
		log.Error("Missing username or id")
		return nil, err
	}

	m, err := s.repo.GetMail(id, username)
	if err != nil {
		log.WithFields(log.Fields{
			"username": username,
			"id": id,
		}).Error("Unable to fetch email")
		return nil, err
	}
	if m == nil {
		err := errors.New("Email doesn't exits")
		log.WithFields(log.Fields{
			"username": username,
			"id": id,
		}).Error("Invalid email")
		return nil, err
	}
	return m, nil
}

func generateId() string {
	n := 4
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}