package http

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net/http"
	"social-network/internal/domain"
)

type Endpoints struct {
	Auth      *AuthEndpoints
	Profile   *ProfileEndpoints
	Social    *SocialEndpoints
	Messenger *MessengerEndpoints
	Ws        *WsEndpoints
}

func MakeEndpoints(auth domain.AuthService, profile domain.ProfileService, social domain.SocialService, messenger domain.MessengerService) *Endpoints {
	return &Endpoints{
		Auth: &AuthEndpoints{
			SignUp:       makeSignUpEndpoint(auth),
			SignIn:       makeSignInEndpoint(auth),
			RefreshToken: makeRefreshTokenEndpoint(auth),
		},
		Profile: &ProfileEndpoints{
			Search: &SearchEndpoints{
				GetByAnthroponym: makeSearchProfileByAnthroponym(auth, profile),
			},
		},
		Social: &SocialEndpoints{
			CreateFriendship:                  makeCreateFriendshipEndpoint(auth, social),
			ConfirmFriendship:                 makeConfirmFriendshipEndpoint(auth, social),
			RejectFriendship:                  makeRejectFriendshipEndpoint(auth, social),
			BreakFriendship:                   makeBreakFriendshipEndpoint(auth, social),
			GetFriends:                        makeGetFriendsEndpoint(auth, social),
			GetFollowers:                      makeGetFollowersEndpoint(auth, social),
			GetAllQuestionnaires:              makeGetAllQuestionnairesEndpoint(auth, social),
			GetQuestionnairesByNameAndSurname: makeGetQuestionnairesByNameAndSurnameEndpoint(auth, social),
		},
		Messenger: &MessengerEndpoints{
			CreateChat:  makeCreateChatEndpoint(auth, messenger),
			GetChat:     makeGetChatEndpoint(auth, messenger),
			DeleteChats: makeDeleteChatsEndpoint(auth, messenger),
			GetChats:    makeGetChatsEndpoint(auth, messenger),
			GetMessages: makeGetMessagesEndpoint(auth, messenger),
			SendMessage: makeSendMessageEndpoint(auth, messenger),
		},
		Ws: &WsEndpoints{
			Connect: makeWsConnectEndpoint(auth),
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

type ProfileEndpoints struct {
	Search *SearchEndpoints
}

type SearchEndpoints struct {
	GetByAnthroponym gin.HandlerFunc
}

func makeSearchProfileByAnthroponym(authSvc domain.AuthService, profileSvc domain.ProfileService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request SearchProfileByAnthroponymRequest
		if err := c.BindQuery(&request); err != nil {
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

		const (
			defaultLimit  = 10
			defaultOffset = 0
		)

		if request.Limit == nil {
			request.Limit = new(int)
			*request.Limit = defaultLimit
		}
		if request.Offset == nil {
			request.Offset = new(int)
			*request.Offset = defaultOffset
		}

		quest, count, err := profileSvc.SearchByAnthroponym(c, request.Anthroponym, userID, *request.Limit, *request.Offset)
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

type SocialEndpoints struct {
	CreateFriendship                  gin.HandlerFunc
	ConfirmFriendship                 gin.HandlerFunc
	RejectFriendship                  gin.HandlerFunc
	BreakFriendship                   gin.HandlerFunc
	GetFriends                        gin.HandlerFunc
	GetFollowers                      gin.HandlerFunc
	GetAllQuestionnaires              gin.HandlerFunc
	GetQuestionnairesByNameAndSurname gin.HandlerFunc
}

func makeCreateFriendshipEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request FriendshipRequest
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

		if err = socialSvc.CreateFriendship(c, userID, request.FriendID); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, EmptyResponse{})
	}
}

func makeConfirmFriendshipEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request FriendshipRequest
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

		if err = socialSvc.ConfirmFriendship(c, userID, request.FriendID); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, EmptyResponse{})
	}
}

func makeRejectFriendshipEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request FriendshipRequest
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

		if err = socialSvc.RejectFriendship(c, userID, request.FriendID); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, EmptyResponse{})
	}
}

func makeBreakFriendshipEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request FriendshipRequest
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

		if err = socialSvc.BreakFriendship(c, userID, request.FriendID); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, EmptyResponse{})
	}
}

func makeGetFriendsEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
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

		friends, err := socialSvc.GetFriends(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, GetFriendsResponse{Friends: friends})
	}
}

func makeGetFollowersEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
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

		followers, err := socialSvc.GetFollowers(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, GetFollowersResponse{Followers: followers})
	}
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

type MessengerEndpoints struct {
	CreateChat  gin.HandlerFunc
	GetChat     gin.HandlerFunc
	DeleteChats gin.HandlerFunc
	GetChats    gin.HandlerFunc
	GetMessages gin.HandlerFunc
	SendMessage gin.HandlerFunc
}

func makeCreateChatEndpoint(authSvc domain.AuthService, messSvc domain.MessengerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
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

		var request CreateChatRequest
		if err = c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		chatID, err := messSvc.CreateChat(c, userID, request.CompanionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, CreateChatResponse{
			ChatID: chatID,
		})
	}
}

func makeGetChatEndpoint(authSvc domain.AuthService, messSvc domain.MessengerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
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

		var request GetChatRequest
		if err = c.BindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		chat, err := messSvc.GetChat(c, userID, request.CompanionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, GetChatResponse{ID: chat.ID, CreateTime: chat.CreateTime})
	}
}

func makeDeleteChatsEndpoint(authSvc domain.AuthService, messSvc domain.MessengerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		_, err := authSvc.Authenticate(c, header.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}
	}
}

func makeGetChatsEndpoint(authSvc domain.AuthService, messSvc domain.MessengerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
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

		var request GetChatsRequest
		if err = c.BindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		const (
			defaultLimit  = 10
			defaultOffset = 0
		)

		if request.Limit == nil {
			request.Limit = new(int)
			*request.Limit = defaultLimit
		}
		if request.Offset == nil {
			request.Offset = new(int)
			*request.Offset = defaultOffset
		}

		chats, total, err := messSvc.GetChats(c, userID, *request.Limit, *request.Offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, GetChatsResponse{
			Total:  total,
			Limit:  request.Limit,
			Offset: request.Offset,
			Chats:  chats,
		})
	}
}

func makeGetMessagesEndpoint(authSvc domain.AuthService, messSvc domain.MessengerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
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

		var request GetMessagesRequest
		if err = c.BindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		const (
			defaultLimit  = 10
			defaultOffset = 0
		)

		if request.Limit == nil {
			request.Limit = new(int)
			*request.Limit = defaultLimit
		}
		if request.Offset == nil {
			request.Offset = new(int)
			*request.Offset = defaultOffset
		}

		messages, total, err := messSvc.GetMessages(c, userID, request.ChatID, *request.Limit, *request.Offset)
		switch err {
		case nil:
			c.JSON(http.StatusOK, GetMessagesResponse{
				Total:    total,
				Limit:    request.Limit,
				Offset:   request.Offset,
				Messages: messages,
			})
		case sql.ErrNoRows:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Message: fmt.Sprintf("chat id [%s] not found", request.ChatID),
			})
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
		}
	}
}

func makeSendMessageEndpoint(authSvc domain.AuthService, messSvc domain.MessengerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
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

		var request SendMessagesRequest
		if err = c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		if err = messSvc.SendMessages(c, userID, request.ChatID, request.Messages); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, EmptyResponse{})
	}
}

type WsEndpoints struct {
	Connect gin.HandlerFunc
}

func makeWsConnectEndpoint(authSvc domain.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
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

		conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		//messenger.AddConnection(userID, conn)

		go func() {
			defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					//messenger.RemoveConnection(userID, conn)
					return
				}

				log.Println(msg)
				//messenger.CreateMessage(msg)

				err = wsutil.WriteServerMessage(conn, op, []byte(fmt.Sprintf("Pong to %s", userID)))
				if err != nil {
					//c.JSON(http.StatusUnauthorized, ErrorResponse{
					//	Message: err.Error(),
					//})
					return
				}
			}
		}()

	}
}
