package message

type grpcServer struct {
	hub Hub
}

func NewGRPCServer(hub Hub) MessengerServer{
	return &grpcServer{
		hub,
	}
}

func (s *grpcServer)  SendMessage(messenger Messenger_SendMessageServer) error {
	c := NewGRPCClient(s.hub, messenger)
	s.hub.Register(c)
	c.Run()
	return nil
}
