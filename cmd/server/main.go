package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"smail/internal/database"
	"smail/internal/mail"
	"smail/internal/message"
)

const (
	defaultMongoURL = "localhost:27017"
	defaultDb    = "smail"
	defaultPort     = ":8080"
)

func main() {
	var (
		mongoURL        = envString("SMAIL_DB_URL", defaultMongoURL)
		port            = envString("PORT", defaultPort)
	)

	mongoSession := database.MongoConnect(mongoURL)
	defer mongoSession.Close()

	err := mongoSession.Ping()

	if err != nil {
		log.Panic("Error connecting to database ", err.Error())
	}

	log.Info("connected to database")

	mailRepo := database.NewMailRepo(mongoSession, defaultDb)
	mailServer := mail.NewMailServer(mailRepo)

	hub := message.NewHub()
	hub.Run()
	defer hub.Close()

	chatServer := message.NewGRPCServer(hub)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panic("Failed to listen: ", err.Error())
	}

	grpcServer := grpc.NewServer()

	mail.RegisterMailerServer(grpcServer, mailServer)
	message.RegisterMessengerServer(grpcServer, chatServer)

	log.WithField("port", port).Info("listening")
	if err := grpcServer.Serve(lis); err != nil {
		log.Panic("failed to serve: ", err.Error())
	}

}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}