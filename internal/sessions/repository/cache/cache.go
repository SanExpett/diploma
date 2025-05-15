package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"

	myerrors "github.com/SanExpett/diploma/internal/errors"
)

var (
	maxVersion uint32 = 255
)

type SessionStorage struct {
	redisClient *redis.Client
}

func NewSessionStorage() *SessionStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return &SessionStorage{
		redisClient: client,
	}
}

func (sessionStorage *SessionStorage) Add(login string, token string, version uint32) error {
	ctx := context.Background()

	// Получаем текущие сессии пользователя
	sessions := make(map[string]uint32)
	val, err := sessionStorage.redisClient.Get(ctx, login).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	if err != redis.Nil {
		if err := json.Unmarshal([]byte(val), &sessions); err != nil {
			return err
		}

		if _, exists := sessions[token]; exists {
			return myerrors.ErrItemsIsAlreadyInTheCache
		}
	}

	sessions[token] = version

	// Сохраняем обновленные сессии
	sessionsJSON, err := json.Marshal(sessions)
	if err != nil {
		return err
	}

	return sessionStorage.redisClient.Set(ctx, login, sessionsJSON, 0).Err()
}

func (sessionStorage *SessionStorage) DeleteSession(login string, token string) error {
	ctx := context.Background()

	sessions := make(map[string]uint32)
	val, err := sessionStorage.redisClient.Get(ctx, login).Result()
	if err == redis.Nil {
		return myerrors.ErrNoSuchUserInTheCache
	}
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), &sessions); err != nil {
		return err
	}

	if _, exists := sessions[token]; !exists {
		return myerrors.ErrNoSuchSessionInTheCache
	}

	delete(sessions, token)

	sessionsJSON, err := json.Marshal(sessions)
	if err != nil {
		return err
	}

	return sessionStorage.redisClient.Set(ctx, login, sessionsJSON, 0).Err()
}

func (sessionStorage *SessionStorage) Update(login string, token string) error {
	ctx := context.Background()

	sessions := make(map[string]uint32)
	val, err := sessionStorage.redisClient.Get(ctx, login).Result()
	if err == redis.Nil {
		return myerrors.ErrNoSuchUserInTheCache
	}
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), &sessions); err != nil {
		return err
	}

	version, exists := sessions[token]
	if !exists {
		return myerrors.ErrNoSuchSessionInTheCache
	}

	if version == maxVersion {
		return myerrors.ErrTooHighVersion
	}

	sessions[token] = version + 1

	sessionsJSON, err := json.Marshal(sessions)
	if err != nil {
		return err
	}

	return sessionStorage.redisClient.Set(ctx, login, sessionsJSON, 0).Err()
}

func (sessionStorage *SessionStorage) CheckVersion(login string, token string, usersVersion uint32) (bool, error) {
	ctx := context.Background()

	sessions := make(map[string]uint32)
	val, err := sessionStorage.redisClient.Get(ctx, login).Result()
	if err == redis.Nil {
		return false, myerrors.ErrNoSuchItemInTheCache
	}
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(val), &sessions); err != nil {
		return false, err
	}

	version, exists := sessions[token]
	if !exists {
		return false, myerrors.ErrNoSuchItemInTheCache
	}

	if version == usersVersion {
		return true, nil
	}

	return false, myerrors.ErrWrongSessionVersion
}

func (sessionStorage *SessionStorage) HasSession(login string, token string) error {
	ctx := context.Background()

	sessions := make(map[string]uint32)
	val, err := sessionStorage.redisClient.Get(ctx, login).Result()
	if err == redis.Nil {
		return myerrors.ErrNoSuchUser
	}
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), &sessions); err != nil {
		return err
	}

	if _, exists := sessions[token]; !exists {
		return myerrors.ErrNoSuchUser
	}

	return nil
}

func (sessionStorage *SessionStorage) GetVersion(login string, token string) (uint32, error) {
	ctx := context.Background()

	sessions := make(map[string]uint32)
	val, err := sessionStorage.redisClient.Get(ctx, login).Result()
	if err == redis.Nil {
		return 0, myerrors.ErrNoSuchItemInTheCache
	}
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal([]byte(val), &sessions); err != nil {
		return 0, err
	}

	version, exists := sessions[token]
	if !exists {
		return 0, myerrors.ErrNoSuchItemInTheCache
	}

	return version, nil
}

func (sessionStorage *SessionStorage) CheckAllUserSessionTokens(login string) error {
	ctx := context.Background()

	sessions := make(map[string]uint32)
	val, err := sessionStorage.redisClient.Get(ctx, login).Result()
	if err == redis.Nil {
		return myerrors.ErrNoSuchUserInTheCache
	}
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), &sessions); err != nil {
		return err
	}

	for token, version := range sessions {
		fmt.Println("token:", token, "version:", version)
	}

	return nil
}
