package implementation

import (
	"context"
	"net"
	"social/internal/domain"
	"social/internal/infrastructure/stan"
	stantransport "social/internal/transport/stan"
	"sort"
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
	UserID string `json:"user_id"`
}

func (a *authService) Authenticate(_ context.Context, token string) (string, error) {
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": token,
	}

	// only url is required, others are optional.
	r, err := req.Get("http://"+a.authAddr+"/auth/user/get-id-by-token", header)
	if err != nil {
		return "", err
	}

	var response getUserIDByAccessTokenResponse
	if err = r.ToJSON(&response); err != nil {
		return "", err
	}

	return response.UserID, nil
}

type profileService struct {
	repository domain.UserRepository
}

func NewProfileService(repository domain.UserRepository) *profileService {
	return &profileService{repository: repository}
}

func (p *profileService) SearchByAnthroponym(ctx context.Context, anthroponym, userID string, limit, offset int) ([]*domain.Questionnaire, int, error) {
	tx, err := p.repository.GetTx(ctx)
	if err != nil {
		return nil, 0, err
	}

	users, count, err := p.repository.GetByAnthroponym(tx, anthroponym, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	questionnaires := make([]*domain.Questionnaire, 0, len(users))
	for _, user := range users {
		questionnaires = append(questionnaires, &domain.Questionnaire{
			ID:               user.ID,
			Email:            user.Email,
			Name:             user.Name,
			Surname:          user.Surname,
			Birthday:         user.Birthday,
			Sex:              user.Sex,
			City:             user.City,
			Interests:        user.Interests,
			FriendshipStatus: user.FriendshipStatus,
		})
	}

	return questionnaires, count, p.repository.CommitTx(tx)
}

type socialService struct {
	userRepository        domain.UserRepository
	socialRepository      domain.SocialRepository
	socialCacheRepository domain.SocialCacheRepository
	stanClient            *stan.Client
}

func NewSocialService(userRep domain.UserRepository, socialRep domain.SocialRepository, socialCacheRep domain.SocialCacheRepository, stanClient *stan.Client) *socialService {
	return &socialService{
		userRepository:        userRep,
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

func (s *socialService) GetFriends(ctx context.Context, userID string) ([]*domain.Questionnaire, error) {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	users, err := s.socialRepository.GetFriends(tx, userID)
	if err != nil {
		return nil, err
	}

	questionnaires := make([]*domain.Questionnaire, 0, len(users))
	for _, user := range users {
		questionnaires = append(questionnaires, &domain.Questionnaire{
			ID:               user.ID,
			Email:            user.Email,
			Name:             user.Name,
			Surname:          user.Surname,
			Birthday:         user.Birthday,
			Sex:              user.Sex,
			City:             user.City,
			Interests:        user.Interests,
			FriendshipStatus: user.FriendshipStatus,
		})
	}

	return questionnaires, s.socialRepository.CommitTx(tx)
}

func (s *socialService) GetFollowers(ctx context.Context, userID string) ([]*domain.Questionnaire, error) {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	users, err := s.socialRepository.GetFollowers(tx, userID)
	if err != nil {
		return nil, err
	}

	questionnaires := make([]*domain.Questionnaire, 0, len(users))
	for _, user := range users {
		questionnaires = append(questionnaires, &domain.Questionnaire{
			ID:               user.ID,
			Email:            user.Email,
			Name:             user.Name,
			Surname:          user.Surname,
			Birthday:         user.Birthday,
			Sex:              user.Sex,
			City:             user.City,
			Interests:        user.Interests,
			FriendshipStatus: user.FriendshipStatus,
		})
	}

	return questionnaires, s.socialRepository.CommitTx(tx)
}

func (s *socialService) RetrieveNews(ctx context.Context, userID string, limit, offset int) ([]*domain.News, int, error) {
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

func (s *socialService) PublishNews(ctx context.Context, userID string, newsContent []string) error {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return err
	}

	user, err := s.userRepository.GetByID(tx, userID)
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

	if err = s.socialRepository.PublishNews(tx, userID, news); err != nil {
		return err
	}

	if err = s.socialRepository.CommitTx(tx); err != nil {
		return err
	}

	if err = s.stanClient.Publish("news", &stantransport.NewsPersistRequest{OwnerID: userID, News: news}); err != nil {
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
	userRep        domain.UserRepository
	socialRep      domain.SocialRepository
	socialCacheRep domain.SocialCacheRepository
	wsPoolRep      domain.WSPoolRepository
	stanClient     *stan.Client
	logger         *zap.Logger
}

func NewWSService(userRep domain.UserRepository, socialRep domain.SocialRepository,
	socialCacheRep domain.SocialCacheRepository, wsPoolRep domain.WSPoolRepository,
	stanClient *stan.Client, logger *zap.Logger) *wsService {
	return &wsService{
		userRep:        userRep,
		socialRep:      socialRep,
		socialCacheRep: socialCacheRep,
		wsPoolRep:      wsPoolRep,
		stanClient:     stanClient,
		logger:         logger,
	}
}

func (w *wsService) EstablishConn(ctx context.Context, userID string, conn net.Conn) {
	w.wsPoolRep.AddConnection(userID, conn)

	go func(userID string) {
		defer conn.Close()
		defer w.wsPoolRep.RemoveConnection(userID, conn)

		tx, err := w.userRep.GetTx(ctx)
		if err != nil {
			w.logger.Error("can't get transaction", zap.Error(err))

			return
		}

		user, err := w.userRep.GetByID(tx, userID)
		if err != nil {
			w.logger.Error("fail get user by id", zap.Error(err))

			return
		}

		if err = tx.Commit(); err != nil {
			w.logger.Error("fail commit transaction", zap.Error(err))

			return
		}

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
	}(userID)
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
