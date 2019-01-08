package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"net"
	"net/http"
	"os"
	"smail/internal/database"
	"smail/internal/mail"
	"smail/internal/message"
	"sync"
)

const (
	defaultMongoURL = "localhost:27017"
	defaultDb    = "smail"
	defaultGRPCPort     = ":8080"
	defaultWSPort     = ":8081"
)

func main() {
	var (
		mongoURL        = envString("SMAIL_DB_URL", defaultMongoURL)
		grpcPort            = envString("GRPC_PORT", defaultGRPCPort)
		wsPort            = envString("WS_PORT", defaultWSPort)
	)

	var mongoSession *mgo.Session
	if mongoURL == defaultMongoURL {
		mongoSession = database.MongoSimpleConnect(mongoURL)
	} else{
		mongoSession = database.MongoConnect(mongoURL)
	}

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

	var wg sync.WaitGroup
	wg.Add(2)

	chatServer := message.NewGRPCServer(hub)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Panic("Failed to listen: ", err.Error())
	}

	grpcServer := grpc.NewServer()

	mail.RegisterMailerServer(grpcServer, mailServer)
	message.RegisterMessengerServer(grpcServer, chatServer)

	go func() {
		log.WithFields(log.Fields{"port": grpcPort, "type": "grpc"}).Info("listening")
		if err := grpcServer.Serve(lis); err != nil {
			log.WithField("err", err.Error()).Fatal("unable to serve grpc")
		}
		wg.Done()
	}()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsClient := message.NewWebSocketClient(hub, w, r)
		hub.Register(wsClient)
		wsClient.Run()
	})

	go func() {
		log.WithFields(log.Fields{"port": wsPort, "type": "websocket"}).Info("listening")
		err = http.ListenAndServe(wsPort, nil)
		if err != nil {
			log.WithField("err", err.Error()).Fatal("unable to serve websocket")
		}
		wg.Done()
	}()

	wg.Wait()
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}