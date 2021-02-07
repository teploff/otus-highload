package http

import (
	"bytes"
	"context"
	"encoding/json"
	"gateway/internal/config"
	"gateway/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Endpoints struct {
	cfg       *config.Config
	Auth      *AuthEndpoints
	AuthProxy *AuthProxyEndpoints
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
		AuthProxy: &AuthProxyEndpoints{
			Authenticate: kitopentracing.TraceClient(opentracing.GlobalTracer(), "authenticate")(makeAuthenticateProxyEndpoint("http://" + cfg.Auth.Addr)),
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
			WS:          makeHTTPProxyEndpoint(cfg.Messenger.HTTPAddr),
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

		var clientSpan opentracing.Span
		tracer := opentracing.GlobalTracer()

		if parentSpan := opentracing.SpanFromContext(c.Request.Context()); parentSpan != nil {
			clientSpan = tracer.StartSpan(
				c.Request.RequestURI,
				opentracing.ChildOf(parentSpan.Context()),
			)
		} else {
			clientSpan = tracer.StartSpan(c.Request.RequestURI)
		}
		defer clientSpan.Finish()

		ext.SpanKindRPCClient.Set(clientSpan)
		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), clientSpan))

		if span := opentracing.SpanFromContext(c.Request.Context()); span != nil {
			opentracing.GlobalTracer().Inject(
				span.Context(),
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(c.Request.Header),
			)
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

type AuthProxyEndpoints struct {
	Authenticate endpoint.Endpoint
}

func makeAuthenticateProxyEndpoint(proxyURL string) endpoint.Endpoint {
	tgt, _ := url.Parse(proxyURL)

	return httptransport.NewClient(
		"POST",
		tgt,
		encodePostAddressRequest,
		decodeDeleteAddressResponse,
		httptransport.ClientBefore(kitopentracing.ContextToHTTP(opentracing.GlobalTracer(), log.NewNopLogger())),
	).Endpoint()
}

func encodePostAddressRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(AuthenticateRequest)

	req.URL.Path = "/auth/authenticate"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Header)

	return encodeRequest(ctx, req, request)
}

func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)

	return nil
}

func decodeDeleteAddressResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response AuthenticateResponse
	err := json.NewDecoder(resp.Body).Decode(&response)

	return response, err
}
