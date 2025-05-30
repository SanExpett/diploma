package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/SanExpett/diploma/internal/domain"
	"github.com/SanExpett/diploma/internal/metrics"
	"github.com/SanExpett/diploma/internal/requestId"
)

type usersStorage interface {
	CreateUser(user domain.UserSignUp) error
	RemoveUser(email string) error
	HasUser(email, password string) error
	GetUser(email string) (domain.User, error)
	ChangeUserPassword(email, newPassword string) (domain.User, error)
	ChangeUserName(email, newName string) (domain.User, error)
	GetUserDataByUuid(uuid string) (domain.User, error)
	GetUserPreview(uuid string) (domain.UserPreview, error)
	ChangeUserPasswordByUuid(uuid, newPassword string) (domain.User, error)
	ChangeUserNameByUuid(uuid, newName string) (domain.User, error)
	ChangeUserAvatarByUuid(uuid, filename string) (domain.User, error)
	HasSubscription(uuid string) (bool, error)
	AddSubscription(uuid string, newDate string) error
	GetSubscriptions() ([]domain.Subscription, error)
	GetSubscription(uuid string) (domain.Subscription, error)
}

type UsersService struct {
	storage   usersStorage
	metrics   *metrics.GrpcMetrics
	secretKey string
	logger    *zap.SugaredLogger
}

func NewUsersService(storage usersStorage, metrics *metrics.GrpcMetrics,
	logger *zap.SugaredLogger) *UsersService {
	return &UsersService{
		storage:   storage,
		metrics:   metrics,
		logger:    logger,
		secretKey: os.Getenv("SECRETKEY"),
	}
}

func (service *UsersService) CreateUser(ctx context.Context, user domain.UserSignUp) error {
	service.metrics.IncRequestsTotal("CreateUser")
	err := service.storage.CreateUser(user)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to create user: %v", ctx.Value(requestId.ReqIDKey),
			err)
		return err
	}
	return nil
}

func (service *UsersService) RemoveUser(ctx context.Context, login string) error {
	service.metrics.IncRequestsTotal("RemoveUser")
	err := service.storage.RemoveUser(login)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to remove user: %v", ctx.Value(requestId.ReqIDKey),
			err)
		return err
	}
	return nil
}

func (service *UsersService) HasUser(ctx context.Context, login, password string) error {
	service.metrics.IncRequestsTotal("HasUser")
	err := service.storage.HasUser(login, password)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to has user: %v", ctx.Value(requestId.ReqIDKey), err)
		return err
	}
	return nil
}

func (service *UsersService) GetUser(ctx context.Context, login string) (domain.User, error) {
	service.metrics.IncRequestsTotal("GetUser")
	user, err := service.storage.GetUser(login)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to get user: %v", ctx.Value(requestId.ReqIDKey), err)
		return domain.User{}, err
	}
	return user, nil
}

func (service *UsersService) ChangeUserPassword(ctx context.Context, login, newPassword string) (domain.User, error) {
	service.metrics.IncRequestsTotal("ChangeUserPassword")
	user, err := service.storage.ChangeUserPassword(login, newPassword)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to change password: %v",
			ctx.Value(requestId.ReqIDKey), err)
		return domain.User{}, err
	}
	return user, nil
}

func (service *UsersService) ChangeUserName(ctx context.Context, login, newName string) (domain.User, error) {
	service.metrics.IncRequestsTotal("ChangeUserName")
	user, err := service.storage.ChangeUserName(login, newName)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to change username: %v", ctx.Value(requestId.ReqIDKey),
			err)
		return domain.User{}, err
	}
	return user, nil
}

func (service *UsersService) GetUserDataByUuid(ctx context.Context, uuid string) (domain.User, error) {
	service.metrics.IncRequestsTotal("GetUserDataByUuid")
	user, err := service.storage.GetUserDataByUuid(uuid)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to get user data: %v",
			ctx.Value(requestId.ReqIDKey), err)
		return domain.User{}, err
	}

	service.metrics.IncRequestsTotal("HasSubscription")
	stat, err := service.storage.HasSubscription(uuid)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to check subscription: %v", ctx.Value(requestId.ReqIDKey),
			err)
		return domain.User{}, err
	}

	user.HasSubscription = stat

	return user, nil
}

func (service *UsersService) GetUserPreview(ctx context.Context, uuid string) (domain.UserPreview, error) {
	service.metrics.IncRequestsTotal("GetUserPreview")
	userPreview, err := service.storage.GetUserPreview(uuid)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to get user preview: %v", ctx.Value(requestId.ReqIDKey),
			err)
		return domain.UserPreview{}, err
	}
	return userPreview, nil
}

func (service *UsersService) ChangeUserPasswordByUuid(ctx context.Context, uuid, newPassword string) (domain.User,
	error) {
	service.metrics.IncRequestsTotal("ChangeUserPasswordByUuid")
	user, err := service.storage.ChangeUserPasswordByUuid(uuid, newPassword)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to change password: %v",
			ctx.Value(requestId.ReqIDKey), err)
		return domain.User{}, err
	}
	return user, nil
}

func (service *UsersService) ChangeUserNameByUuid(ctx context.Context, uuid, newName string) (domain.User, error) {
	service.metrics.IncRequestsTotal("ChangeUserNameByUuid")
	user, err := service.storage.ChangeUserNameByUuid(uuid, newName)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to change username: %v", ctx.Value(requestId.ReqIDKey),
			err)
		return domain.User{}, err
	}
	return user, nil
}

func (service *UsersService) ChangeUserAvatarByUuid(ctx context.Context, uuid, newAvatar string) (domain.User, error) {
	service.metrics.IncRequestsTotal("ChangeUserAvatarByUuid")
	user, err := service.storage.ChangeUserAvatarByUuid(uuid, newAvatar)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to change username: %v", ctx.Value(requestId.ReqIDKey),
			err)
		return domain.User{}, err
	}
	return user, nil
}

func (service *UsersService) HasSubscription(ctx context.Context, uuid string) (bool, error) {
	service.metrics.IncRequestsTotal("HasSubscription")
	stat, err := service.storage.HasSubscription(uuid)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to check subscription: %v", ctx.Value(requestId.ReqIDKey),
			err)
		return false, err
	}
	return stat, nil
}

func (service *UsersService) PaySubscription(ctx context.Context, uuid, subId string) (string, error) {
	service.metrics.IncRequestsTotal("PaySubscription")

	sub, err := service.storage.GetSubscription(subId)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to get subscription: %v", ctx.Value(requestId.ReqIDKey),
			err)
		return "", err
	}

	//requestBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{
	//    "amount": {
	//      "value": %f,
	//      "currency": "RUB"
	//    },
	//    "payment_method_data": {
	//      "type": "bank_card"
	//    },
	//    "confirmation": {
	//      "type": "redirect",
	//      "return_url": "https://nimbus.ru/"
	//    },
	//    "description": "Подписка %s"
	//  }`, sub.Amount, uuid)))
	//req, err := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", requestBody)
	//if err != nil {
	//	return "", err
	//}
	//req.SetBasicAuth("393063", "test_qaG8b_fmJMDHP-Htdq7a_kCwhnAKTEM9ZWAOA0OgDJ0")
	//req.Header.Set("Content-Type", "application/json")
	//id := guuid.New().String()
	//req.Header.Set("Idempotence-Key", id)
	//
	//resp, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	return "", err
	//}
	//
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return "", err
	//}
	//defer resp.Body.Close()
	//
	//c := make(map[string]json.RawMessage)
	//
	//e := json.Unmarshal([]byte(string(body)), &c)
	//if e != nil {
	//	return "", e
	//}
	//
	//go func() {
	//	for {
	//		select {
	//		case <-time.After(time.Minute * 10):
	//			fmt.Println("Payment failed")
	//			return
	//		default:
	//			checkReq, err := http.NewRequest("GET", fmt.Sprintf("https://api.yookassa.ru/v3/payments/%s",
	//				string(c["id"])[1:len(c["id"])-1]), nil)
	//			if err != nil {
	//				fmt.Println(err)
	//			}
	//			checkReq.SetBasicAuth("393063", "test_qaG8b_fmJMDHP-Htdq7a_kCwhnAKTEM9ZWAOA0OgDJ0")
	//
	//			checkResp, err := http.DefaultClient.Do(checkReq)
	//			if err != nil {
	//				fmt.Println(err)
	//			}
	//
	//			body, err := io.ReadAll(checkResp.Body)
	//			if err != nil {
	//				fmt.Println(err)
	//			}
	//			defer checkResp.Body.Close()
	//
	//			resp := make(map[string]json.RawMessage)
	//
	//			respBytes := json.Unmarshal([]byte(string(body)), &resp)
	//			if respBytes != nil {
	//				fmt.Println(err)
	//			}
	//
	//			if string(resp["status"]) == `"waiting_for_capture"` {
	//				req, err := http.NewRequest("POST", fmt.Sprintf(
	//					"https://api.yookassa.ru/v3/payments/%s/capture", string(c["id"])[1:len(c["id"])-1]), nil)
	//				if err != nil {
	//					fmt.Println(err)
	//				}
	//				req.SetBasicAuth("393063", "test_qaG8b_fmJMDHP-Htdq7a_kCwhnAKTEM9ZWAOA0OgDJ0")
	//				req.Header.Set("Content-Type", "application/json")
	//				req.Header.Set("Idempotence-Key", id)
	//
	//				resp, err := http.DefaultClient.Do(req)
	//				if err != nil {
	//					fmt.Println(err)
	//				}
	//				defer resp.Body.Close()
	//
	//				err = service.storage.AddSubscription(uuid, time.Now().AddDate(0, int(sub.Duration), 0).Format("2006-01-02"))
	//				if err != nil {
	//					fmt.Println(err)
	//				}
	//				return
	//			}
	//			time.Sleep(time.Second * 10)
	//		}
	//	}
	//}()
	//
	//e = json.Unmarshal([]byte(string(c["confirmation"])), &c)
	//if e != nil {
	//	return "", e
	//}

	err = service.storage.AddSubscription(uuid, time.Now().AddDate(0, int(sub.Duration), 0).Format("2006-01-02"))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//return string(c["confirmation_url"])[1 : len(c["confirmation_url"])-1], nil
	return "", nil
}

func (service *UsersService) GetSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	service.metrics.IncRequestsTotal("AddSubscription")
	subs, err := service.storage.GetSubscriptions()
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to get subscriptions: %v", ctx.Value(requestId.ReqIDKey),
			err)
		return nil, err
	}
	return subs, nil
}

func (service *UsersService) GetSubscription(ctx context.Context, uuid string) (domain.Subscription, error) {
	service.metrics.IncRequestsTotal("GetSubscription")
	sub, err := service.storage.GetSubscription(uuid)
	if err != nil {
		service.logger.Errorf("[reqid=%s] failed to get subscription: %v", ctx.Value(requestId.ReqIDKey),
			err)
		return domain.Subscription{}, err
	}
	return sub, nil
}
