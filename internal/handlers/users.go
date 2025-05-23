package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	"github.com/SanExpett/diploma/internal/domain"
	myerrors "github.com/SanExpett/diploma/internal/errors"
	"github.com/SanExpett/diploma/internal/metrics"
	reqid "github.com/SanExpett/diploma/internal/requestId"
	session "github.com/SanExpett/diploma/internal/session/proto"
)

type UserPageHandlers struct {
	usersClient    *session.UsersClient
	sessionsClient *session.SessionsClient
	metrics        *metrics.HttpMetrics
	logger         *zap.SugaredLogger
}

func NewUserPageHandlers(usersClient *session.UsersClient, sessionsClient *session.SessionsClient,
	metrics *metrics.HttpMetrics, logger *zap.SugaredLogger) *UserPageHandlers {
	return &UserPageHandlers{
		usersClient:    usersClient,
		sessionsClient: sessionsClient,
		metrics:        metrics,
		logger:         logger,
	}
}

func (UserPageHandlers *UserPageHandlers) GetProfileData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := ctx.Value(reqid.ReqIDKey)
	uuid := mux.Vars(r)["uuid"]
	req := session.GetUserDataByUuidRequest{Uuid: uuid}
	userProto, err := (*UserPageHandlers.usersClient).GetUserDataByUuid(ctx, &req)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestID, err)
		}
		return
	}
	user := convertUserToRegular(userProto.User)

	escapeUserData(&user)
	response := domain.ProfileResponse{
		Status:   http.StatusOK,
		UserInfo: user,
	}

	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to marshal response: %v\n", requestID, err)
		}
		return
	}

	err = WriteResponse(w, r, UserPageHandlers.metrics, jsonResponse, requestID)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestID, err)
		}
		return
	}
}

func (UserPageHandlers *UserPageHandlers) GetProfilePreview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := ctx.Value(reqid.ReqIDKey)

	uuid := mux.Vars(r)["uuid"]
	req := session.GetUserPreviewRequest{Uuid: uuid}
	userPreviewProto, err := (*UserPageHandlers.usersClient).GetUserPreview(ctx, &req)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestID, err)
		}
		return
	}
	userPreview := convertUserPreviewToRegular(userPreviewProto.User)

	escapeUserPreviewData(&userPreview)

	response := domain.ProfilePreviewResponse{
		Status:      http.StatusOK,
		UserPreview: userPreview,
	}

	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to marshal response: %v\n", requestID, err)
		}
		return
	}

	err = WriteResponse(w, r, UserPageHandlers.metrics, jsonResponse, requestID)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestID, err)
		}
		return
	}
}

func (UserPageHandlers *UserPageHandlers) ProfileEditByUuid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value(reqid.ReqIDKey)

	userToken, err := r.Cookie("access")
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, myerrors.ErrNoActiveSession)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}

	_, err = IsTokenValid(userToken, os.Getenv("SECRETKEY"))
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}

	uuid := mux.Vars(r)["uuid"]
	req := session.GetUserDataByUuidRequest{Uuid: uuid}
	getUserByDataRes, err := (*UserPageHandlers.usersClient).GetUserDataByUuid(ctx, &req)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}
	currUserProto := getUserByDataRes.User

	newData := r.FormValue("newData")
	switch r.FormValue("action") {
	case "chPassword":
		err = ValidatePassword(newData)
		if err != nil {
			err = WriteError(w, r, UserPageHandlers.metrics, err)
			if err != nil {
				UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
			}
			return
		}

		reqPass := session.ChangeUserPasswordByUuidRequest{Uuid: uuid, NewPassword: newData}
		changePassRes, err := (*UserPageHandlers.usersClient).ChangeUserPasswordByUuid(ctx, &reqPass)
		if err != nil {
			err = WriteError(w, r, UserPageHandlers.metrics, err)
			if err != nil {
				UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
			}
			return
		}
		currUserProto = changePassRes.User

	case "chUsername":
		err = ValidateUsername(newData)
		if err != nil {
			err = WriteError(w, r, UserPageHandlers.metrics, err)
			if err != nil {
				UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
			}
			return
		}

		reqName := session.ChangeUserNameByUuidRequest{Uuid: uuid, NewUsername: newData}
		changeNameRes, err := (*UserPageHandlers.usersClient).ChangeUserNameByUuid(ctx, &reqName)
		if err != nil {
			err = WriteError(w, r, UserPageHandlers.metrics, err)
			if err != nil {
				UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
			}
			return
		}
		currUserProto = changeNameRes.User
	case "chAvatar":
		files := r.MultipartForm.File["avatar"]

		avatar64, err := Encode(files[0])
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to encode into base64: %v\n", requestId, err)
			err = WriteError(w, r, UserPageHandlers.metrics, err)
			if err != nil {
				UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
			}
			return
		}

		reqAvatar := session.ChangeUserAvatarByUuidRequest{Uuid: uuid, NewAvatar: avatar64}
		changeAvatarRes, err := (*UserPageHandlers.usersClient).ChangeUserAvatarByUuid(ctx, &reqAvatar)
		if err != nil {
			err = WriteError(w, r, UserPageHandlers.metrics, err)
			if err != nil {
				UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
			}
			return
		}
		currUserProto = changeAvatarRes.User
	}

	version := currUserProto.Version + 1

	tokenSigned, err := GenerateTokens(currUserProto.Email, false, version)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}

	reqAdd := session.AddRequest{Login: currUserProto.Email, Token: tokenSigned, Version: version}
	_, err = (*UserPageHandlers.sessionsClient).Add(ctx, &reqAdd)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}

	tokenCookie := &http.Cookie{
		Name:     "access",
		Value:    tokenSigned,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   tokenCookieExpirationTime,
	}

	http.SetCookie(w, tokenCookie)

	err = WriteSuccess(w, r, UserPageHandlers.metrics)
	if err != nil {
		UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
	}
}

func (UserPageHandlers *UserPageHandlers) HasSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value(reqid.ReqIDKey)

	uuid := mux.Vars(r)["uuid"]
	req := session.HasSubscriptionRequest{Uuid: uuid}
	stat, err := (*UserPageHandlers.usersClient).HasSubscription(ctx, &req)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}

	response := domain.HasSubsctiptionsResponse{
		Status:          http.StatusOK,
		HasSubscription: stat.Status,
	}

	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to marshal response: %v\n", requestId, err)
		}
		return
	}
	err = WriteResponse(w, r, UserPageHandlers.metrics, jsonResponse, requestId)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}
}

func (UserPageHandlers *UserPageHandlers) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value(reqid.ReqIDKey)

	subs, err := (*UserPageHandlers.usersClient).GetSubscriptions(ctx, &session.GetSubscriptionsRequest{})
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}
	subscriptions := convertSubsToRegular(subs.Subscriptions)

	response := domain.SubsctiptionsResponse{
		Status:        http.StatusOK,
		Subscriptions: subscriptions,
	}

	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to marshal response: %v\n", requestId, err)
		}
		return
	}
	err = WriteResponse(w, r, UserPageHandlers.metrics, jsonResponse, requestId)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}
}

type payRequest struct {
	SubId string `json:"subId"`
}

func (UserPageHandlers *UserPageHandlers) PaySubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value(reqid.ReqIDKey)

	var request payRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		UserPageHandlers.logger.Errorf("[reqid=%s] failed to decode: %v\n", requestId, myerrors.ErrFailedDecode)
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}

	uuid := mux.Vars(r)["uuid"]
	req := session.HasSubscriptionRequest{Uuid: uuid}
	stat, err := (*UserPageHandlers.usersClient).HasSubscription(ctx, &req)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}

	if stat.Status {
		err = WriteError(w, r, UserPageHandlers.metrics, myerrors.ErrAlreadyHaveSubscription)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}

	link, err := (*UserPageHandlers.usersClient).PaySubscription(ctx, &session.PaySubscriptionRequest{
		Uuid:  uuid,
		SubId: request.SubId,
	})
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}

	response := domain.PayResponse{
		Link: link.PaymentResponse,
	}

	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to marshal response: %v\n", requestId, err)
		}
		return
	}
	err = WriteResponse(w, r, UserPageHandlers.metrics, jsonResponse, requestId)
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}
}

func (UserPageHandlers *UserPageHandlers) RemoveUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value(reqid.ReqIDKey)

	var request domain.RemoveUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		UserPageHandlers.logger.Errorf("[reqid=%s] failed to decode: %v\n", requestId, myerrors.ErrFailedDecode)
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}

	_, err = (*UserPageHandlers.usersClient).RemoveUser(ctx, &session.RemoveUserRequest{
		Login: request.Login,
	})
	if err != nil {
		err = WriteError(w, r, UserPageHandlers.metrics, err)
		if err != nil {
			UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
		}
		return
	}

	err = WriteSuccess(w, r, UserPageHandlers.metrics)
	if err != nil {
		UserPageHandlers.logger.Errorf("[reqid=%s] failed to write response: %v\n", requestId, err)
	}
}
