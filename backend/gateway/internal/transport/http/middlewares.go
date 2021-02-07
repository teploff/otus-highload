package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticateMiddleware(endpoints *AuthProxyEndpoints) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := AuthenticateRequest{
			Header:   findToken(c),
			Resource: c.Request.URL.Path,
		}

		resp, err := endpoints.Authenticate(c.Request.Context(), request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		response := resp.(AuthenticateResponse)

		if !response.IsAuthenticated {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Message: "User isn't authenticated",
			})

			return
		}

		c.Next()
	}
}

func findToken(c *gin.Context) string {
	if token := c.Request.Header.Get("Authorization"); token != "" {
		return token
	}

	v := c.Request.URL.Query()
	if token, exist := v["token"]; exist {
		return token[0]
	}

	return ""
}
