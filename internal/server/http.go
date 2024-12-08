package server

import (
	userv1 "career_backend/api/user/v1"
	"career_backend/internal/routers"
	"career_backend/internal/service"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/http"
)

// NewHTTPServer creates a HTTP server
func NewHTTPServer(c *app.Config, userSvc *service.UserServiceServer) *http.Server {
	router := routers.NewRouter()

	srv := http.NewServer(
		http.WithAddress(c.HTTP.Addr),
		http.WithReadTimeout(c.HTTP.ReadTimeout),
		http.WithWriteTimeout(c.HTTP.WriteTimeout),
	)

	srv.Handler = router

	// v1.RegisterGreeterServiceHTTPServer(router, greeterSvc)
	userv1.RegisterUserServiceHTTPServer(router, userSvc)

	return srv
}
