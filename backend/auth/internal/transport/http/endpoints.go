package http

import (
	"auth/internal/domain"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Endpoints struct {
	Auth *AuthEndpoints
}

func MakeEndpoints(auth domain.AuthService) *Endpoints {
	return &Endpoints{
		Auth: &AuthEndpoints{
			SignUp:           makeSignUpEndpoint(auth),
			SignIn:           makeSignInEndpoint(auth),
			RefreshToken:     makeRefreshTokenEndpoint(auth),
			GetUserIDByEmail: makeGetUserIDByEmail(auth),
		},
	}
}

type AuthEndpoints struct {
	SignUp                     gin.HandlerFunc
	SignIn                     gin.HandlerFunc
	RefreshToken               gin.HandlerFunc
	GetUserIDByEmail           gin.HandlerFunc
	SearchProfileByAnthroponym gin.HandlerFunc
}

// SignUp godoc
// @Summary Sign up user by credentials.
// @Description Sign up user by credentials.
// @Tags auth
// @Accept  json
// @Produce json
// @Param payload body SignUpRequest true "Sign up payload"
// @Success 200 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/sign-up [post].
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
			Birthday:    request.Birthday,
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

// SignIn godoc
// @Summary Sign in user by credentials.
// @Description Sign in user by credentials.
// @Tags auth
// @Accept  json
// @Produce json
// @Param payload body SignInRequest true "Sign in payload"
// @Success 200 {object} SignInResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/sign-in [post].
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
			Email:    request.Email,
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

// RefreshToken godoc
// @Summary Refresh expired token on new JWT pair.
// @Description Refresh expired token on new JWT pair.
// @Tags auth
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param payload body RefreshTokenRequest true "JWT refresh token"
// @Success 200 {object} SignInResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/token [put].
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

		switch err {
		case nil:
			c.JSON(http.StatusOK, SignInResponse{
				AccessToken:  tokenPair.AccessToken,
				RefreshToken: tokenPair.RefreshToken,
			})
		case sql.ErrNoRows:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: "token is expired",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})
		}
	}
}

// GetUserIDByEmail godoc
// @Summary Retrieving user's id by email.
// @Description Retrieving user's id by email.
// @Tags auth
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param email query string true "User's Email"
// @Success 200 {object} GetUserIDByEmailResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/user/get-by-email [get].
func makeGetUserIDByEmail(svc domain.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		_, err := svc.Authenticate(c, header.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request GetUserIDByEmailRequest
		if err = c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		userID, err := svc.GetUserIDByEmail(c, request.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, GetUserIDByEmailResponse{
			UserID: userID,
		})
	}
}
