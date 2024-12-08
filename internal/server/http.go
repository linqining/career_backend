package server

import (
	companyv1 "career_backend/api/company/v1"
	userv1 "career_backend/api/user/v1"
	"career_backend/internal/routers"
	"career_backend/internal/service"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/http"
)

// NewHTTPServer creates a HTTP server
func NewHTTPServer(c *app.Config, userSvc *service.UserServiceServer, companyService *service.CompanyServiceServer) *http.Server {
	router := routers.NewRouter()

	srv := http.NewServer(
		http.WithAddress(c.HTTP.Addr),
		http.WithReadTimeout(c.HTTP.ReadTimeout),
		http.WithWriteTimeout(c.HTTP.WriteTimeout),
	)

	srv.Handler = router

	userv1.RegisterUserServiceHTTPServer(router, userSvc)
	companyv1.RegisterCompanyServiceHTTPServer(router, companyService)

	return srv
}
