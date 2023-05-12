package server

import (
	"server/internal/executor"
	"server/internal/handler"
)

type Server struct {
	H *handler.Handler
}

func NewServer() (*Server, error) {

	execcpp := executor.NewExecutorCPP()

	handler := handler.NewHandler(execcpp)

	return &Server{
		H: handler,
	}, nil
}
