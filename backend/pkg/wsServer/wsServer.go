package ws

import "social-network/pkg/models"

// represent websocket server
type Server struct {
	Clients map[*Client]bool
	Repos   *models.Repositories
}

func StartServer(repos *models.Repositories) *Server {
	server := &Server{
		make(map[*Client]bool),
		repos,
	}
	return server
}

// register client
func (s *Server) RegisterNewClient(client *Client) {
	s.Clients[client] = true //update client list
}

// register and unregister clients
func (s *Server) UnregisterClient(client *Client) {
	if _, ok := s.Clients[client]; ok {
		delete(s.Clients, client)
	}
}