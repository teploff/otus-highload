package implementation

import (
	"context"
	"net"
	"social/internal/domain"
	"social/internal/infrastructure/stan"
	stantransport "social/internal/transport/stan"
	"sort"
	"strconv"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/imroc/req"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type authService struct {
	authAddr string
}

func NewAuthService(authAddr string) *authService {
	return &authService{
		authAddr: authAddr,
	}
}

type getUserIDByAccessTokenResponse struct {
	*domain.User
}

func (a *authService) Authenticate(_ context.Context, token string) (*domain.User, error) {
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": token,
	}

	r, err := req.Get("http://"+a.authAddr+"/auth/user/get-by-token", header)
	if err != nil {
		return nil, err
	}

	var response getUserIDByAccessTokenResponse
	if err = r.ToJSON(&response); err != nil {
		return nil, err
	}

	return response.User, nil
}

type getUsersByAnthroponymResponse struct {
	Count int            `json:"count"`
	Users []*domain.User `json:"users"`
}

func (a *authService) GetUsersByAnthroponym(ctx context.Context, token, anthroponym string, offset, limit int) ([]*domain.User, int, error) {
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": token,
	}

	param := req.Param{
		"anthroponym": anthroponym,
		"offset":      strconv.Itoa(offset),
		"limit":       strconv.Itoa(limit),
	}

	r, err := req.Get("http://"+a.authAddr+"/auth/user/get-by-anthroponym", header, param)
	if err != nil {
		return nil, 0, err
	}

	var response getUsersByAnthroponymResponse
	if err = r.ToJSON(&response); err != nil {
		return nil, 0, err
	}

	return response.Users, response.Count, nil
}

type getUserByIDsRequest struct {
	UserIDs []string `json:"user_ids" binding:"required"`
}

type getUserByIDsResponse struct {
	Users []*domain.User `json:"users"`
}

func (a *authService) GetUsersByIDs(ctx context.Context, ids []string) ([]*domain.User, error) {
	header := req.Header{
		"Accept": "application/json",
	}

	body := getUserByIDsRequest{
		UserIDs: ids,
	}

	r, err := req.Post("http://"+a.authAddr+"/auth/user/get-by-ids", header, req.BodyJSON(body))
	if err != nil {
		return nil, err
	}

	var response getUserByIDsResponse
	if err = r.ToJSON(&response); err != nil {
		return nil, err
	}

	return response.Users, nil
}

type profileService struct {
	authSvc   domain.AuthService
	socialRep domain.SocialRepository
}

func NewProfileService(service domain.AuthService, repository domain.SocialRepository) *profileService {
	return &profileService{
		authSvc:   service,
		socialRep: repository,
	}
}

func (p *profileService) SearchByAnthroponym(ctx context.Context, token, anthroponym string, offset, limit int) ([]*domain.User, int, error) {
	user, err := p.authSvc.Authenticate(ctx, token)
	if err != nil {
		return nil, 0, err
	}

	users, count, err := p.authSvc.GetUsersByAnthroponym(ctx, token, anthroponym, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	tx, err := p.socialRep.GetTx(ctx)
	if err != nil {
		return nil, 0, nil
	}

	friendships, err := p.socialRep.GetUserFriendships(tx, user.ID)
	if err != nil {
		return nil, 0, nil
	}

	if err = p.socialRep.CommitTx(tx); err != nil {
		return nil, 0, nil
	}

	for _, u := range users {
		u.FriendshipStatus = friendshipNonameStatus

		for _, friendship := range friendships {
			if friendship.MasterUserID == user.ID {
				switch friendship.Status {
				case friendshipExpectedStatus:
					user.FriendshipStatus = friendshipConfirmedStatus
				default:
					user.FriendshipStatus = friendshipAcceptedStatus
				}
			} else if friendship.SlaveUserID == user.ID {
				user.FriendshipStatus = friendship.Status
			}
		}
	}

	return users, count, nil
}

type socialService struct {
	authSvc               domain.AuthService
	socialRepository      domain.SocialRepository
	socialCacheRepository domain.SocialCacheRepository
	stanClient            *stan.Client
}

func NewSocialService(authSvc domain.AuthService, socialRep domain.SocialRepository, socialCacheRep domain.SocialCacheRepository, stanClient *stan.Client) *socialService {
	return &socialService{
		authSvc:               authSvc,
		socialRepository:      socialRep,
		socialCacheRepository: socialCacheRep,
		stanClient:            stanClient,
	}
}

func (s *socialService) CreateFriendship(ctx context.Context, userID string, friendsID []string) error {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return err
	}

	if err = s.socialRepository.CreateFriendship(tx, userID, friendsID); err != nil {
		return err
	}

	return s.socialRepository.CommitTx(tx)
}

func (s *socialService) ConfirmFriendship(ctx context.Context, userID string, friendsID []string) error {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return err
	}

	if err = s.socialRepository.ConfirmFriendship(tx, userID, friendsID); err != nil {
		return err
	}

	if err = s.socialRepository.CommitTx(tx); err != nil {
		return err
	}

	if err = s.stanClient.Publish("friends", stantransport.FriendsActionRequest{
		Action:    "persist",
		UserID:    userID,
		FriendsID: friendsID,
	}); err != nil {
		return err
	}

	for _, friendID := range friendsID {
		if err = s.stanClient.Publish("friends", stantransport.FriendsActionRequest{
			Action:    "persist",
			UserID:    friendID,
			FriendsID: []string{userID},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *socialService) RejectFriendship(ctx context.Context, userID string, friendsID []string) error {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return err
	}

	if err = s.socialRepository.RejectFriendship(tx, userID, friendsID); err != nil {
		return err
	}

	return s.socialRepository.CommitTx(tx)
}

func (s *socialService) BreakFriendship(ctx context.Context, userID string, friendsID []string) error {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return err
	}

	if err = s.socialRepository.BreakFriendship(tx, userID, friendsID); err != nil {
		return err
	}

	if err = s.socialRepository.CommitTx(tx); err != nil {
		return err
	}

	if err = s.stanClient.Publish("friends", stantransport.FriendsActionRequest{
		Action:    "delete",
		UserID:    userID,
		FriendsID: friendsID,
	}); err != nil {
		return err
	}

	for _, friendID := range friendsID {
		if err = s.stanClient.Publish("friends", stantransport.FriendsActionRequest{
			Action:    "delete",
			UserID:    friendID,
			FriendsID: []string{userID},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *socialService) GetFriends(ctx context.Context, userID string) ([]*domain.User, error) {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	friendsID, err := s.socialRepository.GetFriends(tx, userID)
	if err != nil {
		return nil, err
	}

	if err = s.socialRepository.CommitTx(tx); err != nil {
		return nil, err
	}

	return s.authSvc.GetUsersByIDs(ctx, friendsID)
}

func (s *socialService) GetFollowers(ctx context.Context, userID string) ([]*domain.User, error) {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	followersID, err := s.socialRepository.GetFollowers(tx, userID)
	if err != nil {
		return nil, err
	}

	if err = s.socialRepository.CommitTx(tx); err != nil {
		return nil, err
	}

	return s.authSvc.GetUsersByIDs(ctx, followersID)
}

func (s *socialService) RetrieveNews(ctx context.Context, userID string, offset, limit int) ([]*domain.News, int, error) {
	ids, err := s.socialCacheRepository.RetrieveFriendsID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	ids = append(ids, userID)

	news := make([]*domain.News, 0, 1)
	for _, id := range ids {
		n, err := s.socialCacheRepository.RetrieveNews(ctx, id)
		if err != nil {
			return nil, 0, err
		}

		news = append(news, n...)
	}

	count := len(news)

	// Sort by age, keeping original order or equal elements.
	sort.SliceStable(news, func(i, j int) bool {
		return news[i].CreateTime.Unix() > news[j].CreateTime.Unix()
	})

	return s.paginate(news, offset, limit), count, nil
}

func (s *socialService) paginate(n []*domain.News, skip int, size int) []*domain.News {
	if skip > len(n) {
		skip = len(n)
	}

	end := skip + size
	if end > len(n) {
		end = len(n)
	}

	return n[skip:end]
}

func (s *socialService) PublishNews(ctx context.Context, token string, newsContent []string) error {
	user, err := s.authSvc.Authenticate(ctx, token)
	if err != nil {
		return err
	}

	news := make([]*domain.News, 0, len(newsContent))
	for _, content := range newsContent {
		news = append(news, &domain.News{
			ID: uuid.NewV4().String(),
			Owner: struct {
				Name    string `json:"name"`
				Surname string `json:"surname"`
				Sex     string `json:"sex"`
			}{
				user.Name,
				user.Surname,
				user.Sex,
			},
			Content:    content,
			CreateTime: time.Now().UTC(),
		})
	}

	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return err
	}

	if err = s.socialRepository.PublishNews(tx, token, news); err != nil {
		return err
	}

	if err = s.socialRepository.CommitTx(tx); err != nil {
		return err
	}

	if err = s.stanClient.Publish("news", &stantransport.NewsPersistRequest{OwnerID: token, News: news}); err != nil {
		return err
	}

	return nil
}

type cacheService struct {
	repository domain.SocialCacheRepository
}

func NewCacheService(repository domain.SocialCacheRepository) *cacheService {
	return &cacheService{repository: repository}
}

func (c *cacheService) AddFriends(ctx context.Context, userID string, friendsID []string) error {
	return c.repository.PersistFriend(ctx, userID, friendsID)
}

func (c *cacheService) DeleteFriends(ctx context.Context, userID string, friendsID []string) error {
	return c.repository.DeleteFriend(ctx, userID, friendsID)
}

func (c *cacheService) AddNews(ctx context.Context, userID string, news []*domain.News) error {
	return c.repository.PersistNews(ctx, userID, news)
}

// easyjson:json
type WSRequest struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

// easyjson:json
type WSNewsRequest struct {
	Content string `json:"content"`
}

type wsService struct {
	authSvc        domain.AuthService
	socialRep      domain.SocialRepository
	socialCacheRep domain.SocialCacheRepository
	wsPoolRep      domain.WSPoolRepository
	stanClient     *stan.Client
	logger         *zap.Logger
}

func NewWSService(authSvc domain.AuthService, socialRep domain.SocialRepository,
	socialCacheRep domain.SocialCacheRepository, wsPoolRep domain.WSPoolRepository,
	stanClient *stan.Client, logger *zap.Logger) *wsService {
	return &wsService{
		authSvc:        authSvc,
		socialRep:      socialRep,
		socialCacheRep: socialCacheRep,
		wsPoolRep:      wsPoolRep,
		stanClient:     stanClient,
		logger:         logger,
	}
}

func (w *wsService) EstablishConn(ctx context.Context, user *domain.User, conn net.Conn) {
	w.wsPoolRep.AddConnection(user.ID, conn)

	go func(user *domain.User) {
		defer conn.Close()
		defer w.wsPoolRep.RemoveConnection(user.ID, conn)

		for {
			msg, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				w.logger.Error("fail read from ws", zap.Error(err))

				return
			}

			if err = w.parseRequest(ctx, user, msg); err != nil {
				w.logger.Error("", zap.Error(err))
			}
		}
	}(user)
}

func (w *wsService) parseRequest(ctx context.Context, user *domain.User, msg []byte) error {
	var request WSRequest
	if err := request.UnmarshalJSON(msg); err != nil {
		return err
	}

	switch request.Topic {
	case "news":
		var r WSNewsRequest

		if err := r.UnmarshalJSON([]byte(request.Payload)); err != nil {
			return err
		}

		n := &domain.News{
			ID: uuid.NewV4().String(),
			Owner: struct {
				Name    string `json:"name"`
				Surname string `json:"surname"`
				Sex     string `json:"sex"`
			}{
				user.Name,
				user.Surname,
				user.Sex,
			},
			Content:    r.Content,
			CreateTime: time.Now().UTC(),
		}
		news := []*domain.News{n}

		tx, err := w.socialRep.GetTx(ctx)
		if err != nil {
			return err
		}

		if err = w.socialRep.PublishNews(tx, user.ID, news); err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		if err = w.stanClient.Publish("news", &stantransport.NewsPersistRequest{OwnerID: user.ID, News: []*domain.News{n}}); err != nil {
			return err
		}
	}

	return nil
}

func (w *wsService) SendNews(ctx context.Context, ownerID string, news []*domain.News) error {
	ids, err := w.socialCacheRep.RetrieveFriendsID(ctx, ownerID)
	if err != nil {
		return err
	}
	ids = append(ids, ownerID)

	for _, id := range ids {
		conns := w.wsPoolRep.RetrieveConnByUserID(id)
		for _, conn := range conns {
			for _, n := range news {
				data, err := n.MarshalJSON()
				if err != nil {
					return err
				}

				if err = wsutil.WriteServerMessage(conn, ws.OpText, data); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (w *wsService) Close() {
	w.wsPoolRep.FlushConnections()
}
