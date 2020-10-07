package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social-network/internal/domain"
	"time"
)

type AuthEndpoints struct {
	SignUp       gin.HandlerFunc
	SignIn       gin.HandlerFunc
	RefreshToken gin.HandlerFunc
}

func makeSignUpEndpoint(svc domain.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request SignUpRequest
		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		credentials, err := domain.NewCredentials(request.Email, request.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		if err = svc.SignUp(c, &domain.User{
			Credentials: *credentials,
			Name:        request.Name,
			Surname:     request.Surname,
			Birthday:    time.Unix(request.Birthday, 0),
			Sex:         request.Sex,
			City:        request.City,
			Interests:   request.Interests,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, EmptyResponse{})
	}
}

func makeSignInEndpoint(svc domain.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request SignInRequest
		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		tokenPair, err := svc.SignIn(c, &domain.Credentials{
			Login:    request.Login,
			Password: request.Password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, SignInResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		})
	}
}

func makeRefreshTokenEndpoint(svc domain.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request RefreshTokenRequest
		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		tokenPair, err := svc.RefreshToken(c, request.RefreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, SignInResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		})
	}
}

func MakeAuthEndpoints(svc domain.AuthService) *AuthEndpoints {
	return &AuthEndpoints{
		SignUp:       makeSignUpEndpoint(svc),
		SignIn:       makeSignInEndpoint(svc),
		RefreshToken: makeRefreshTokenEndpoint(svc),
	}
}
