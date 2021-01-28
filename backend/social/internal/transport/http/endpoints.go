package http

import (
	"net/http"
	"social/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
)

type Endpoints struct {
	WS         gin.HandlerFunc
	Profile    *ProfileEndpoints
	Friendship *FriendshipEndpoints
	News       *NewsEndpoints
}

func MakeEndpoints(auth domain.AuthService, profile domain.ProfileService, social domain.SocialService, ws domain.WSService) *Endpoints {
	return &Endpoints{
		WS: makeWSEndpoint(auth, ws),

		Profile: &ProfileEndpoints{
			SearchByAnthroponym: makeSearchProfileByAnthroponym(auth, profile),
		},

		Friendship: &FriendshipEndpoints{
			Create:       makeCreateFriendshipEndpoint(auth, social),
			Confirm:      makeConfirmFriendshipEndpoint(auth, social),
			Reject:       makeRejectFriendshipEndpoint(auth, social),
			SplitUp:      makeSplitUpFriendshipEndpoint(auth, social),
			GetFriends:   makeGetFriendsEndpoint(auth, social),
			GetFollowers: makeGetFollowersEndpoint(auth, social),
		},

		News: &NewsEndpoints{
			GetNews: makeGetNewsEndpoint(auth, social),
		},
	}
}

// WS godoc
// @Summary Websocket handshake endpoint.
// @Description Websocket handshake endpoint.
// @Accept  json
// @Produce json
// @Param token query string true "User's Access-JWT"
// @Success 200 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Router /social/ws [get].
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

		wsSvc.EstablishConn(c, user, conn)
	}
}

type ProfileEndpoints struct {
	SearchByAnthroponym gin.HandlerFunc
}

// SearchProfileByAnthroponym godoc
// @Summary Search user's by some variation of name and surname.
// @Description Search user's by some variation of name and surname.
// @Tags Profile
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param payload body SearchProfileByAnthroponymRequest true "Search payload"
// @Success 200 {object} SearchProfileByAnthroponymResponse
// @Failure 400 {object} ErrorResponse
// @Router /social/profile/search-by-anthroponym [post].
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

		users, count, err := profileSvc.SearchByAnthroponym(c, request.Anthroponym, header.AccessToken, *request.Offset,
			*request.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, SearchProfileByAnthroponymResponse{
			Users: users,
			Count: count,
		})
	}
}

type FriendshipEndpoints struct {
	Create             gin.HandlerFunc
	Confirm            gin.HandlerFunc
	Reject             gin.HandlerFunc
	SplitUp            gin.HandlerFunc
	GetFriends         gin.HandlerFunc
	GetFollowers       gin.HandlerFunc
	GetUserFriendships gin.HandlerFunc
}

// CreateFriendship godoc
// @Summary Create friendship between two users.
// @Description Create friendship between two users.
// @Tags Friendship
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param payload body FriendshipRequest true "Friends' id"
// @Success 200 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Router /friendship/create [post].
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

		user, err := authSvc.Authenticate(c, header.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		if err = socialSvc.CreateFriendship(c, user.ID, request.FriendsID); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, EmptyResponse{})
	}
}

// ConfirmFriendship godoc
// @Summary Confirming friendship between two users.
// @Description Confirming friendship between two users.
// @Tags Friendship
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param payload body FriendshipRequest true "Friends' id"
// @Success 200 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Router /friendship/confirm [post].
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

		user, err := authSvc.Authenticate(c, header.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		if err = socialSvc.ConfirmFriendship(c, user.ID, request.FriendsID); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, EmptyResponse{})
	}
}

// RejectFriendship godoc
// @Summary Rejecting friendship between two users.
// @Description Rejecting friendship between two users.
// @Tags Friendship
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param payload body FriendshipRequest true "Friends' id"
// @Success 200 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Router /friendship/reject [post].
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

		user, err := authSvc.Authenticate(c, header.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		if err = socialSvc.RejectFriendship(c, user.ID, request.FriendsID); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, EmptyResponse{})
	}
}

// SplitUpFriendship godoc
// @Summary Splitting up friendship between two users.
// @Description Splitting up friendship between two users.
// @Tags Friendship
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param payload body FriendshipRequest true "Friends' id"
// @Success 200 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Router /friendship/split-up [post].
func makeSplitUpFriendshipEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
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

		user, err := authSvc.Authenticate(c, header.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		if err = socialSvc.BreakFriendship(c, user.ID, request.FriendsID); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, EmptyResponse{})
	}
}

// GetFriends godoc
// @Summary Retrieving user's friends.
// @Description Retrieving user's friends.
// @Tags Friendship
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Success 200 {object} GetFriendsResponse
// @Failure 400 {object} ErrorResponse
// @Router /friendship/get-friends [get].
func makeGetFriendsEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		user, err := authSvc.Authenticate(c, header.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		friends, err := socialSvc.GetFriends(c, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, GetFriendsResponse{Friends: friends})
	}
}

// GetFollowers godoc
// @Summary Retrieving user's followers.
// @Description Retrieving user's followers.
// @Tags Friendship
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Success 200 {object} GetFollowersResponse
// @Failure 400 {object} ErrorResponse
// @Router /friendship/get-followers [get].
func makeGetFollowersEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		user, err := authSvc.Authenticate(c, header.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		followers, err := socialSvc.GetFollowers(c, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, GetFollowersResponse{Followers: followers})
	}
}

type NewsEndpoints struct {
	GetNews gin.HandlerFunc
}

// GetNews godoc
// @Summary Retrieving user's news.
// @Description Retrieving user's news.
// @Tags News
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Success 200 {object} GetNewsResponse
// @Failure 400 {object} ErrorResponse
// @Router /news/get-news [get].
func makeGetNewsEndpoint(authSvc domain.AuthService, socialSvc domain.SocialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthorizationHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		var request GetNewsRequest
		if err := c.BindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		user, err := authSvc.Authenticate(c, header.AccessToken)
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

		news, count, err := socialSvc.RetrieveNews(c, user.ID, *request.Limit, *request.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, GetNewsResponse{
			News:  news,
			Count: count,
		})
	}
}
