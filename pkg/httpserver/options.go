package httpserver

import (
	"net"
	"time"
)

// Option allows for managing server options.
type Option func(*Server)

// Port is a server option for setting the port.
func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

// ReadTimeout sets the read timeout for the server.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

// WriteTimeout sets the write timeout for the server.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

// ShutdownTimeout sets the shutdown timeout for the server.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
