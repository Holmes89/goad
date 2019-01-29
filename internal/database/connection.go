package database

import (
	"crypto/tls"
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	"net"
)

func MongoConnect(url string) (*mgo.Session, error) {

	log.Info("connecting to database")
	dialInfo, err := mgo.ParseURL(url)
	if err != nil {
		log.Fatal("unable to get info for database: ", err)
	}

	//Below part is similar to above.
	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		log.Warn("unable to connect to database ", err)
		return nil, err
	}
	return session, nil
}

func MongoSimpleConnect(url string) (*mgo.Session, error)  {

	log.Info("connecting to database")
	c, err := mgo.Dial(url)
	if err != nil {
		log.WithField("error", err.Error()).Warn("Unable to connect to mongo database")
		return nil, err
	}
	return c, nil
}