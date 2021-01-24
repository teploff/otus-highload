package http

import (
	"gateway/internal/config"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

type Endpoints struct {
	cfg  *config.Config
	Auth *AuthEndpoints
	//Messenger *MessengerEndpoints
	//Social    *SocialEndpoints
}

func MakeEndpoints(cfg *config.Config) *Endpoints {
	return &Endpoints{
		cfg: cfg,
		Auth: &AuthEndpoints{
			SignUp:           makeHTTPProxyEndpoint(cfg.Auth.Addr),
			SignIn:           makeHTTPProxyEndpoint(cfg.Auth.Addr),
			RefreshToken:     makeHTTPProxyEndpoint(cfg.Auth.Addr),
			GetUserIDByEmail: makeHTTPProxyEndpoint(cfg.Auth.Addr),
		},
		//Messenger: &MessengerEndpoints{
		//	Upload: makeUploadEndpoint(cfg.Messenger.Addr),
		//	Get:    makeGetEndpoint(cfg.Messenger.Addr),
		//},
		//Social: &SocialEndpoints{
		//	Calculate: makeCalculateEndpoint(cfg.Messenger.Addr),
		//},
	}
}

type AuthEndpoints struct {
	SignUp           gin.HandlerFunc
	SignIn           gin.HandlerFunc
	RefreshToken     gin.HandlerFunc
	GetUserIDByEmail gin.HandlerFunc
}

func makeHTTPProxyEndpoint(targetHost string) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy := httputil.ReverseProxy{
			Director: func(request *http.Request) {
				request.Header.Add("X-Forwarded-Host", request.Host)
				request.Header.Add("X-Origin-Host", targetHost)
				request.URL.Scheme = "http"
				request.URL.Host = targetHost
			},
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
