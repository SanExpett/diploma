package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/SanExpett/diploma/internal/domain"
	myerrors "github.com/SanExpett/diploma/internal/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxIface interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Close()
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type storage struct {
	pool PgxIface
}

func newStorage(pool PgxIface) (*storage, error) {
	return &storage{
		pool: pool,
	}, nil
}

const insertSubscription = `INSERT INTO subscription (title, description, amount, duration) VALUES ($1, $2, $3, $4);`

func (storage *storage) CreateSubscription(sub domain.Subscription) error {
	_, err := storage.pool.Exec(context.Background(), insertSubscription, sub.Title, sub.Description, sub.Amount, sub.Duration)
	if err != nil {
		return fmt.Errorf("failed to create subscripion: %w: %w", err,
			myerrors.ErrFailInExec)
	}

	return nil
}

func (filmsPageHandlers *FilmsPageHandlers) AddSubscriptions(w http.ResponseWriter, r *http.Request) {
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"postgres",
		"5432",
		"postgres",
		"postgres",
		"nimbus",
	))
	if err != nil {
		log.Fatal(err)
	}

	s, err := newStorage(pool)
	if err != nil {
		log.Fatalf("failed to create storage: %v\n", err)

		return
	}

	subs := []domain.Subscription{
		{
			Title:       "Ежемесячный платеж",
			Description: "Наслаждайтесь обширной библиотекой фильмов и сериалов с разнообразным контентом.",
			Amount:      299,
			Duration:    1,
		},
		{
			Title:       "Ежегодный платеж",
			Description: "Покупка на 12 месяцев без продления. Выгоднее на 30%: 208₽ в месяц вместо 299₽ в месяц за ежемесячную подписку",
			Amount:      2490,
			Duration:    12,
		},
	}

	for _, sub := range subs {
		err = s.CreateSubscription(sub)
		if err != nil {
			log.Fatalf("failed to create subscription: %v \n", err)
			return
		}
	}

	return
}
