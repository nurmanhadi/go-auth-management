package service

import (
	"auth-management/internal/cache"
	"auth-management/internal/entity"
	"auth-management/internal/event/publisher"
	"auth-management/internal/repository"
	"auth-management/pkg"
	"auth-management/pkg/dto"
	"auth-management/pkg/enum"
	"auth-management/pkg/response"
	"auth-management/pkg/security"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	logger         zerolog.Logger
	validator      *validator.Validate
	userRepository *repository.UserRepository
	tokenCache     *cache.TokenCache
	userPublisher  *publisher.UserPublisher
}

func NewUserService(logger zerolog.Logger, validator *validator.Validate, userRepository *repository.UserRepository, tokenCache *cache.TokenCache, userPublisher *publisher.UserPublisher) *UserService {
	return &UserService{
		logger:         logger,
		validator:      validator,
		userRepository: userRepository,
		tokenCache:     tokenCache,
		userPublisher:  userPublisher,
	}
}
func (s *UserService) UserRegister(request *dto.UserRequest) error {
	if err := s.validator.Struct(request); err != nil {
		s.logger.Warn().Err(err).Msg("failed to validate request")
		return err
	}
	newUsername := strings.ToLower(request.Username)
	totalUser, err := s.userRepository.CountByUsername(newUsername)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed count by username to database")
		return err
	}
	if totalUser > 0 {
		s.logger.Warn().Err(nil).Msgf("username %s already exists", newUsername)
		return response.Except(http.StatusConflict, "username already exists")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed hash password to bcrypt")
		return err
	}
	id := uuid.NewString()
	createdAt := time.Now()
	user := &entity.User{
		Id:        id,
		Username:  newUsername,
		Password:  string(hashPassword),
		Role:      enum.ROLE_USER,
		CreatedAt: createdAt,
	}
	if err := s.userRepository.Create(user); err != nil {
		s.logger.Error().Err(err).Msg("failed create user to database")
		return err
	}
	go func() {
		data := &dto.EventUserPayload{
			Event:     "user.registered",
			Timestamp: time.Now(),
			Data: dto.EventUserData{
				UserId:        id,
				Username:      newUsername,
				Registered_at: createdAt,
			},
		}
		if err := s.userPublisher.PublishUserRegistered(data); err != nil {
			s.logger.Error().Err(err).Msg("failed publish user registered to publisher")
			return
		}
	}()
	s.logger.Info().Str("username", newUsername).Msg("user register success")
	return nil
}
func (s *UserService) UserLogin(request *dto.UserRequest) (*dto.TokenResponse, error) {
	if err := s.validator.Struct(request); err != nil {
		s.logger.Warn().Err(err).Msg("failed to validate request")
		return nil, err
	}
	newUsername := strings.ToLower(request.Username)
	user, err := s.userRepository.FindByUsername(newUsername)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.Warn().Err(err).Msg("username or password wrong")
			return nil, response.Except(http.StatusBadRequest, "username or password wrong")
		} else {
			s.logger.Error().Err(err).Msg("failed find by username to database")
			return nil, err
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		s.logger.Warn().Err(err).Msg("username or password wrong")
		return nil, response.Except(http.StatusBadRequest, "username or password wrong")
	}
	secret := []byte(os.Getenv("JWT_SECRET"))
	accessToken, err := security.JwtGenerateAccessToken(user.Id, user.Role, secret)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed generate access token to jwt")
		return nil, err
	}
	refreshToken, expUnix, err := security.JwtGenerateRefreshToken(user.Id, user.Role, secret)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed generate refresh token to jwt")
		return nil, err
	}
	hashRefreshToken := pkg.HashToSha256(refreshToken)
	if err := s.tokenCache.SetRefreshToken(hashRefreshToken, refreshToken, expUnix); err != nil {
		s.logger.Error().Err(err).Msg("failed set refresh token to cache")
		return nil, err
	}
	resp := &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: hashRefreshToken,
	}
	s.logger.Info().Str("username", user.Username).Msg("user login success")
	return resp, nil
}
func (s *UserService) UserGenerateToken(request *dto.UserGenerateToken) (*dto.TokenResponse, error) {
	if err := s.validator.Struct(request); err != nil {
		s.logger.Warn().Err(err).Msg("failed to validate request")
		return nil, err
	}
	token, err := s.tokenCache.GetRefreshToken(request.RefreshToken)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			s.logger.Warn().Err(err).Msg("token missing")
			return nil, response.Except(http.StatusUnauthorized, "token missing")
		}
		s.logger.Error().Err(err).Msg("failed get refresh token to cache")
		return nil, err
	}
	claim, err := security.JwtVerify(token)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed verify token to jwt")
		return nil, response.Except(http.StatusUnauthorized, "cannot processed token")
	}
	secret := []byte(os.Getenv("JWT_SECRET"))
	accessToken, err := security.JwtGenerateAccessToken(claim.Subject, claim.Role, secret)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed generate access token to jwt")
		return nil, err
	}
	refreshToken, expUnix, err := security.JwtGenerateRefreshToken(claim.Subject, claim.Role, secret)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed generate refresh token to jwt")
		return nil, err
	}
	hashRefreshToken := pkg.HashToSha256(refreshToken)
	if err := s.tokenCache.RemoveRefreshToken(request.RefreshToken); err != nil {
		s.logger.Error().Err(err).Msg("failed remove refresh token to cache")
		return nil, err
	}
	if err := s.tokenCache.SetRefreshToken(hashRefreshToken, refreshToken, expUnix); err != nil {
		s.logger.Error().Err(err).Msg("failed set refresh token to cache")
		return nil, err
	}

	resp := &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: hashRefreshToken,
	}
	s.logger.Info().Str("user_id", claim.Subject).Msg("user generate token success")
	return resp, nil
}
