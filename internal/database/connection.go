package database

import (
	"crypto/tls"
	"gopkg.in/mgo.v2"
	log "github.com/sirupsen/logrus"
	"net"
)

func MongoConnect(url string) *mgo.Session {

	log.Info("Connecting to database")
	dialInfo, err := mgo.ParseURL(url)
	if err != nil {
		log.Fatal("Unable to get info for database: ", err)
	}

	//Below part is similar to above.
	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		log.Fatal("Unable to connect to database ", err)
	}
	return session
}
