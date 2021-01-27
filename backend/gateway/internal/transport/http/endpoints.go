package http

import (
	"gateway/internal/config"
	"gateway/internal/transport/grpc"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

type Endpoints struct {
	cfg           *config.Config
	grpcendpoints *grpc.MessengerProxyEndpoints
	Auth          *AuthEndpoints
	Social        *SocialEndpoints
	Messenger     *MessengerEndpoints
}

func MakeEndpoints(cfg *config.Config, grpcendpoints *grpc.MessengerProxyEndpoints) *Endpoints {
	return &Endpoints{
		cfg:           cfg,
		grpcendpoints: grpcendpoints,
		Auth: &AuthEndpoints{
			SignUp:           makeHTTPProxyEndpoint(cfg.Auth.Addr),
			SignIn:           makeHTTPProxyEndpoint(cfg.Auth.Addr),
			RefreshToken:     makeHTTPProxyEndpoint(cfg.Auth.Addr),
			GetUserIDByEmail: makeHTTPProxyEndpoint(cfg.Auth.Addr),
		},
		Social: &SocialEndpoints{
			WS: makeHTTPProxyEndpoint(cfg.Social.Addr),
			Profile: &SocialProfileEndpoints{
				SearchByAnthroponym: makeHTTPProxyEndpoint(cfg.Social.Addr),
			},
			Friendship: &SocialFriendshipEndpoints{
				Create:       makeHTTPProxyEndpoint(cfg.Social.Addr),
				Confirm:      makeHTTPProxyEndpoint(cfg.Social.Addr),
				Reject:       makeHTTPProxyEndpoint(cfg.Social.Addr),
				SplitUp:      makeHTTPProxyEndpoint(cfg.Social.Addr),
				GetFriends:   makeHTTPProxyEndpoint(cfg.Social.Addr),
				GetFollowers: makeHTTPProxyEndpoint(cfg.Social.Addr),
			},
			News: &NewsEndpoints{
				GetNews: makeHTTPProxyEndpoint(cfg.Social.Addr),
			},
		},
		Messenger: &MessengerEndpoints{
			WS: makeHTTPProxyEndpoint(cfg.Messenger.Addr),
		},
	}
}

type AuthEndpoints struct {
	SignUp           gin.HandlerFunc
	SignIn           gin.HandlerFunc
	RefreshToken     gin.HandlerFunc
	GetUserIDByEmail gin.HandlerFunc
}

type SocialProfileEndpoints struct {
	SearchByAnthroponym gin.HandlerFunc
}

type SocialFriendshipEndpoints struct {
	Create       gin.HandlerFunc
	Confirm      gin.HandlerFunc
	Reject       gin.HandlerFunc
	SplitUp      gin.HandlerFunc
	GetFriends   gin.HandlerFunc
	GetFollowers gin.HandlerFunc
}

type SocialEndpoints struct {
	WS         gin.HandlerFunc
	Profile    *SocialProfileEndpoints
	Friendship *SocialFriendshipEndpoints
	News       *NewsEndpoints
}

type NewsEndpoints struct {
	GetNews gin.HandlerFunc
}

type MessengerEndpoints struct {
	WS gin.HandlerFunc
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
