package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"social/internal/domain"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gobwas/ws"
	"github.com/opentracing/opentracing-go"
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

		users, count, err := profileSvc.SearchByAnthroponym(c.Request.Context(), header.AccessToken,
			request.Anthroponym, *request.Offset, *request.Limit)
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

		user, err := authSvc.GetUserByToken(c.Request.Context(), header.AccessToken)
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

		user, err := authSvc.GetUserByToken(c.Request.Context(), header.AccessToken)
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

		user, err := authSvc.GetUserByToken(c.Request.Context(), header.AccessToken)
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

		user, err := authSvc.GetUserByToken(c.Request.Context(), header.AccessToken)
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

		user, err := authSvc.GetUserByToken(c.Request.Context(), header.AccessToken)
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

		user, err := authSvc.GetUserByToken(c.Request.Context(), header.AccessToken)
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

		user, err := authSvc.GetUserByToken(c.Request.Context(), header.AccessToken)
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

		news, count, err := socialSvc.RetrieveNews(c, user.ID, *request.Offset, *request.Limit)
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

type AuthProxyEndpoints struct {
	GetUserIDByAccessToken endpoint.Endpoint
	GetUsersByAnthroponym  endpoint.Endpoint
	GetUsersByIDs          endpoint.Endpoint
}

func NewAuthProxyEndpoints(authAddr string) *AuthProxyEndpoints {
	return &AuthProxyEndpoints{
		GetUserIDByAccessToken: kitopentracing.TraceClient(opentracing.GlobalTracer(), "get-user-by-token")(makeGetUserIDByAccessTokenEndpoint("http://" + authAddr)),
		GetUsersByAnthroponym:  kitopentracing.TraceClient(opentracing.GlobalTracer(), "get-users-by-anthroponym")(makeGetUsersByAnthroponymEndpoint("http://" + authAddr)),
		GetUsersByIDs:          kitopentracing.TraceClient(opentracing.GlobalTracer(), "get-users-by-ids")(makeGetUsersByIDsEndpoint("http://" + authAddr)),
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

func makeGetUsersByAnthroponymEndpoint(proxyURL string) endpoint.Endpoint {
	tgt, _ := url.Parse(proxyURL)

	return httptransport.NewClient(
		"GET",
		tgt,
		encodeGetUsersByAnthroponymRequest,
		decodeGetUsersByAnthroponymResponse,
		httptransport.ClientBefore(kitopentracing.ContextToHTTP(opentracing.GlobalTracer(), log.NewNopLogger())),
	).Endpoint()
}

func makeGetUsersByIDsEndpoint(proxyURL string) endpoint.Endpoint {
	tgt, _ := url.Parse(proxyURL)

	return httptransport.NewClient(
		"POST",
		tgt,
		encodeGetUsersByIDsRequest,
		decodeGetUsersByIDsResponse,
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

func encodeGetUsersByAnthroponymRequest(_ context.Context, req *http.Request, request interface{}) error {
	r := request.(GetUsersByAnthroponymRequest)

	req.URL.Path = "/auth/user/get-by-anthroponym"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Token)

	q := req.URL.Query()
	q.Add("anthroponym", r.Anthroponym)
	q.Add("offset", strconv.Itoa(r.Offset))
	q.Add("limit", strconv.Itoa(r.Limit))

	req.URL.RawQuery = q.Encode()

	return nil
}

func encodeGetUsersByIDsRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(GetUsersByIDsRequest)

	req.URL.Path = "/auth/user/get-by-ids"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Token)

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

func decodeGetUserIDByAccessTokenResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response GetUserIDByAccessTokenResponse
	err := json.NewDecoder(resp.Body).Decode(&response)

	return response, err
}

func decodeGetUsersByAnthroponymResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response GetUsersByAnthroponymResponse
	err := json.NewDecoder(resp.Body).Decode(&response)

	return response, err
}

func decodeGetUsersByIDsResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response GetUsersByIDsResponse
	err := json.NewDecoder(resp.Body).Decode(&response)

	return response, err
}
