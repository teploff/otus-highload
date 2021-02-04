package http

import (
	"messenger/internal/domain"
	"net/http"

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

		user, err := authSvc.Authenticate(c, request.AccessToken)
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

		wsSvc.EstablishConn(c, user.ID, conn)
	}
}
