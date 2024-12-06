package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	port            string
	router          *http.ServeMux
	customerHandler *CustomerHandler
	orderHandler    *OrderHandler
	robotHandler    *RobotHandler
}

// NewServer returns http router to work with.
func NewServer(port string, router *http.ServeMux, customerHandler *CustomerHandler, orderHandler *OrderHandler, robotHandler *RobotHandler) *Server {
	return &Server{
		port:            port,
		router:          router,
		customerHandler: customerHandler,
		orderHandler:    orderHandler,
		robotHandler:    robotHandler}
}

// setRoutes activating handlers and sets routes for http router.
func (s *Server) setRoutes() {
	//customer routes
	s.router.HandleFunc("POST /customers", s.customerHandler.CreateCustomer)

	//robot routes
	s.router.HandleFunc("POST /robots", s.robotHandler.CreateRobot)
	s.router.HandleFunc("GET /robots/week", s.robotHandler.RobotsCreatedThisWeek)

	//order routes
	s.router.HandleFunc("POST /orders", s.orderHandler.OrderRobot)
}

func (s *Server) Start() error {
	s.setRoutes()

	fmt.Println("Server is listening... at port:", s.port)

	return http.ListenAndServe(":"+s.port, nil)

}
