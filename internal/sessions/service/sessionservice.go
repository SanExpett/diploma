package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/SanExpett/diploma/internal/metrics"
	"github.com/SanExpett/diploma/internal/requestId"
)

type sessionStorage interface {
	Add(login string, token string, version uint32) (err error)
	DeleteSession(login string, token string) (err error)
	Update(login string, token string) (err error)
	CheckVersion(login string, token string, usersVersion uint32) (hasSession bool, err error)
	GetVersion(login string, token string) (version uint32, err error)
	HasSession(login string, token string) error
	CheckAllUserSessionTokens(login string) error
}

type SessionService struct {
	sessionStorage sessionStorage
	metrics        *metrics.GrpcMetrics
	logger         *zap.SugaredLogger
}

func NewSessionService(sessionStorage sessionStorage, metrics *metrics.GrpcMetrics,
	logger *zap.SugaredLogger) *SessionService {
	return &SessionService{
		sessionStorage: sessionStorage,
		metrics:        metrics,
		logger:         logger,
	}
}

func (service *SessionService) Add(ctx context.Context, login, token string, version uint32) (err error) {
	service.metrics.IncRequestsTotal("Add")
	err = service.sessionStorage.Add(login, token, version)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to add session: %v", ctx.Value(requestId.ReqIDKey), err)
		return err
	}
	return nil
}

func (service *SessionService) DeleteSession(ctx context.Context, login, token string) (err error) {
	service.metrics.IncRequestsTotal("DeleteSession")
	err = service.sessionStorage.DeleteSession(login, token)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to delete session: %v",
			ctx.Value(requestId.ReqIDKey), err)
		return err
	}
	return nil
}

func (service *SessionService) Update(ctx context.Context, login, token string) (err error) {
	service.metrics.IncRequestsTotal("Update")
	err = service.sessionStorage.Update(login, token)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to update session: %v", ctx.Value(requestId.ReqIDKey), err)
		return err
	}
	return nil
}

func (service *SessionService) CheckVersion(ctx context.Context, login, token string,
	usersVersion uint32) (hasSession bool, err error) {
	service.metrics.IncRequestsTotal("CheckVersion")
	hasSession, err = service.sessionStorage.CheckVersion(login, token, usersVersion)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to check version: %v", ctx.Value(requestId.ReqIDKey), err)
		return hasSession, err
	}
	return hasSession, nil
}

func (service *SessionService) GetVersion(ctx context.Context, login, token string) (version uint32, err error) {
	service.metrics.IncRequestsTotal("GetVersion")
	version, err = service.sessionStorage.GetVersion(login, token)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to get version: %v", ctx.Value(requestId.ReqIDKey), err)
		return version, err
	}
	return version, nil
}

func (service *SessionService) HasSession(ctx context.Context, login, token string) (err error) {
	service.metrics.IncRequestsTotal("HasSession")
	err = service.sessionStorage.HasSession(login, token)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to has session: %v", ctx.Value(requestId.ReqIDKey), err)
		return err
	}
	return nil
}

func (service *SessionService) CheckAllUserSessionTokens(ctx context.Context, login string) error {
	service.metrics.IncRequestsTotal("CheckAllUserSessionTokens")
	err := service.sessionStorage.CheckAllUserSessionTokens(login)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to check all user's session tokens: %v",
			ctx.Value(requestId.ReqIDKey), err)
		return err
	}
	return nil
}
