package database

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"smail/internal/mail"
)

const collection_name = "mail"

type mailRepo struct {
	session *mgo.Session
	database string
}

func NewMailRepo(session *mgo.Session, database string) mail.MailRepo{
	return &mailRepo{
		session,
		database,
	}
}

func (r *mailRepo) GetMail(id, username string) (*mail.Mail, error){
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.database).C(collection_name)
	querier := bson.M{"uuid": id, "to": username}
	m := &mail.Mail{}

	if err := c.Find(querier).One(&m); err != nil {
		log.Error("Unable to fetch single message: ", err.Error())
		return nil, err
	}

	return m, nil

}

func (r *mailRepo) GetAllMail(username string, all bool) ([]*mail.Mail, error){
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.database).C(collection_name)
	var m []*mail.Mail

	var querier bson.M
	if all {
		querier = bson.M{"to": username}
	} else {
		querier = bson.M{"to": username, "read": false}
	}

	if err := c.Find(querier).All(&m); err != nil {
		log.Error("Unable to update message", err.Error())
		return nil, err
	}

	return m, nil
}

func (r *mailRepo) CreateMail(m *mail.Mail) error{
	sess := r.session.Copy()
	defer sess.Close()


	c := sess.DB(r.database).C(collection_name)

	if err := c.Insert(m); err != nil {
		log.Error("Unable to insert message", err.Error())
		return err
	}
	return nil
}

func (r *mailRepo) UpdateMail(m *mail.Mail) error{
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.database).C(collection_name)

	querier := bson.M{"uuid": m.Uuid}
	change := bson.M{"$set": bson.M{"read": m.Read}}
	if err := c.Update(querier, change); err != nil {
		log.Error("Unable to update message", err.Error())
		return err
	}
	return nil
}

func (r *mailRepo) DeleteMail(id string, username string) error{
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.database).C(collection_name)
	querier := bson.M{"uuid": id, "to": username}

	if err := c.Remove(querier); err != nil {
		log.Error("Unable to delete message: ", err.Error())
		return err
	}
	return nil
}