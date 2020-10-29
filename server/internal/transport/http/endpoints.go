package http

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"social-network/internal/domain"
)

type Endpoints struct {
	Auth   *AuthEndpoints
	Social *SocialEndpoints
}

func MakeEndpoints(auth domain.AuthService, social domain.SocialService) *Endpoints {
	return &Endpoints{
		Auth: &AuthEndpoints{
			SignUp:       makeSignUpEndpoint(auth),
			SignIn:       makeSignInEndpoint(auth),
			RefreshToken: makeRefreshTokenEndpoint(auth),
		},
		Social: &SocialEndpoints{
			GetAllQuestionnaires:              makeGetAllQuestionnairesEndpoint(auth, social),
			GetQuestionnairesByNameAndSurname: makeGetQuestionnairesByNameAndSurnameEndpoint(auth, social),
		},
	}
}

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

type SocialEndpoints struct {
	GetAllQuestionnaires              gin.HandlerFunc
	GetQuestionnairesByNameAndSurname gin.HandlerFunc
}

func makeGetAllQuestionnairesEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request GetAllQuestionnairesRequest
		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		userID, err := authSvc.Authenticate(c, header.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		quest, count, err := socialSvc.GetQuestionnaires(c, userID, *request.Limit, request.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, QuestionnairesResponse{
			Questionnaires: quest,
			Count:          count,
		})
	}
}

func makeGetQuestionnairesByNameAndSurnameEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request GetQuestionnairesByNameAndSurnameRequest
		if err := c.BindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		if _, err := authSvc.Authenticate(c, header.AccessToken); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		quest, err := socialSvc.GetQuestionnairesByNameAndSurname(c, request.Prefix)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, QuestionnairesResponse{
			Questionnaires: quest,
			Count:          len(quest),
		})
	}
}
