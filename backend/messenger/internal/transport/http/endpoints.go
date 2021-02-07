package http

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/opentracing/opentracing-go"
	"messenger/internal/domain"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
)

type Endpoints struct {
	Messenger *MessengerEndpoints
}

func MakeEndpoints(auth domain.AuthService, ws domain.WSService) *Endpoints {
	return &Endpoints{
		Messenger: &MessengerEndpoints{
			WS: makeWSEndpoint(auth, ws),
		},
	}
}

type MessengerEndpoints struct {
	WS gin.HandlerFunc
}

// WS godoc
// @Summary Websocket handshake endpoint.
// @Description Websocket handshake endpoint.
// @Accept  json
// @Produce json
// @Param token query string true "User's Access-JWT"
// @Success 200 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Router /messenger/ws [get].
func makeWSEndpoint(authSvc domain.AuthService, wsSvc domain.WSService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request WSRequest

		if err := c.BindQuery(&request); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		user, err := authSvc.GetUserByToken(c.Request.Context(), request.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		wsSvc.EstablishConn(c, user, conn)
	}
}

type AuthProxyEndpoints struct {
	GetUserIDByAccessToken endpoint.Endpoint
}

func NewAuthProxyEndpoints(authAddr string) *AuthProxyEndpoints {
	return &AuthProxyEndpoints{
		GetUserIDByAccessToken: kitopentracing.TraceClient(opentracing.GlobalTracer(), "get-user-by-token")(makeGetUserIDByAccessTokenEndpoint("http://" + authAddr)),
	}
}

func makeGetUserIDByAccessTokenEndpoint(proxyURL string) endpoint.Endpoint {
	tgt, _ := url.Parse(proxyURL)

	return httptransport.NewClient(
		"GET",
		tgt,
		encodeGetUserIDByAccessTokenRequest,
		decodeGetUserIDByAccessTokenResponse,
		httptransport.ClientBefore(kitopentracing.ContextToHTTP(opentracing.GlobalTracer(), log.NewNopLogger())),
	).Endpoint()
}

func encodeGetUserIDByAccessTokenRequest(_ context.Context, req *http.Request, request interface{}) error {
	r := request.(GetUserIDByAccessTokenRequest)

	req.URL.Path = "/auth/user/get-by-token"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Token)

	return nil
}

func decodeGetUserIDByAccessTokenResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response GetUserIDByAccessTokenResponse
	err := json.NewDecoder(resp.Body).Decode(&response)

	return response, err
}
