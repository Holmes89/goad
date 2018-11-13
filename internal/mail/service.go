package mail

import (
	"context"
	"errors"
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

	if req.Id == "" {
		mail, err := s.repo.GetAllMail(req.Username, req.All)
		if err != nil {
			log.WithFields(log.Fields{
				"username": req.Username,
			}).Error("Unable to fetching emails")
			return nil, err
		}
		resp.Mail = mail
	} else {
		mail, err := s.fetchMail(req.Id, req.Username)
		if err != nil {
			return nil, err
		}
		resp.Mail = []*Mail{mail}
	}

	return resp, nil
}

func (s *mailerServer) SendMail(ctx context.Context, mail *Mail) (*Mail, error) {
	if err := s.repo.CreateMail(mail); err != nil {
		log.WithFields(log.Fields{
			"to": mail.To,
			"from": mail.From,
		}).Error("Unable to creat email")
		return nil, err
	}
	return mail, nil
}

func (s *mailerServer) DeleteMail(ctx context.Context, mail *MailID) (*Empty, error) {
	if mail.Id == "" || mail.Username == "" {
		err := errors.New("Id and Username required")
		log.Error("Missing username or id")
		return nil, err
	}
	if err := s.repo.DeleteMail(mail.Id, mail.Username); err != nil {
		log.Error("Error deleting mail", err)
		return nil, err
	}
	return &Empty{}, nil

}

func (s *mailerServer) MarkUnread(ctx context.Context, mail *MailID) (*Mail, error){
	return s.updateMailReadStatus(mail, false)
}


func (s *mailerServer) MarkRead(ctx context.Context, mail *MailID) (*Mail, error){
	return s.updateMailReadStatus(mail, true)
}

func (s *mailerServer) updateMailReadStatus(mail *MailID, status bool) (*Mail, error) {
	m, err := s.fetchMail(mail.Id, mail.Username)
	if err != nil {
		return nil, err
	}

	m.Read = status
	if err := s.repo.UpdateMail(m); err != nil {
		log.WithFields(log.Fields{
			"username": mail.Username,
			"id": mail.Id,
		}).Error("Unable to fetching emails")

		return nil, err
	}
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
		}).Error("Unable to fetching emails")
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