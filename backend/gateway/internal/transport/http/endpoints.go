package http

import (
	"gateway/internal/config"
	"gateway/internal/domain"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

type Endpoints struct {
	cfg       *config.Config
	Auth      *AuthEndpoints
	Social    *SocialEndpoints
	Messenger *MessengerEndpoints
}

func MakeEndpoints(cfg *config.Config, service domain.GRPCMessengerProxyService) *Endpoints {
	return &Endpoints{
		cfg: cfg,
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
			WS:          makeHTTPProxyEndpoint(cfg.Messenger.Addr),
			CreateChat:  makeCreateChatEndpoint(service),
			GetChats:    makeGetChatsEndpoint(service),
			GetMessages: makeGetMessagesEndpoint(service),
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

type NewsEndpoints struct {
	GetNews gin.HandlerFunc
}

type MessengerEndpoints struct {
	WS          gin.HandlerFunc
	CreateChat  gin.HandlerFunc
	GetChats    gin.HandlerFunc
	GetMessages gin.HandlerFunc
}

func makeCreateChatEndpoint(svc domain.GRPCMessengerProxyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request CreateChatRequest
		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		chatID, err := svc.CreateChat(c, header.AccessToken, request.CompanionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, CreateChatResponse{ChatID: chatID})
	}
}

func makeGetChatsEndpoint(svc domain.GRPCMessengerProxyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request GetChatsRequest
		if err := c.BindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		response, err := svc.GetChats(c, header.AccessToken, request.Offset, request.Limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, GetChatsResponse{
			Total:  response.Total,
			Limit:  response.Limit,
			Offset: response.Offset,
			Chats:  response.Chats,
		})
	}
}

func makeGetMessagesEndpoint(svc domain.GRPCMessengerProxyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request GetMessagesRequest
		if err := c.BindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		response, err := svc.GetMessages(c, header.AccessToken, request.ChatID, request.Offset, request.Limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, GetMessagesResponse{
			Total:    response.Total,
			Limit:    response.Limit,
			Offset:   response.Offset,
			Messages: response.Messages,
		})
	}
}
