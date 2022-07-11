package api

import "github.com/gin-gonic/gin"

func (s *Server) Routes() *gin.Engine {
	router := s.router

	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", s.handlers.Register())
		v1.POST("/login", s.handlers.Login(), s.handlers.GenerateTokens())
		v1.POST("/refresh", s.handlers.VerifyRefreshToken(), s.handlers.GenerateTokens())
		v1.POST("/logout", s.handlers.VerifyAccessToken(), s.handlers.Logout())
		v1.GET("/me", s.handlers.VerifyAccessToken(), s.handlers.Me())
	}
	return router
}
