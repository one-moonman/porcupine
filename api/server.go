package api

import "github.com/gin-gonic/gin"

type Server struct {
	router   *gin.Engine
	handlers Handler
}

func (s Server) Run() {
	s.router = gin.Default()
	s.Routes()
	s.router.Run()
}
