package repository

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"twitch_chat_analysis/internal/model"
)

type Repository interface {
	Save(ctx context.Context, msg string) error
	GetReports(ctx context.Context) []model.Message
}

type repository struct {
	db *redis.Client
}

func NewRepository(db *redis.Client) Repository {
	return &repository{
		db: db,
	}
}

func (r repository) Save(ctx context.Context, msg string) error {
	key := "report_" + uuid.New().String()
	return r.db.Set(ctx, key, msg, 0).Err()
}

func (r repository) GetReports(ctx context.Context) []model.Message {
	var res []model.Message
	iter := r.db.Scan(ctx, 0, "report_*", 0).Iterator()

	for iter.Next(ctx) {
		var m model.Message
		key := iter.Val()
		val := r.db.Get(ctx, key)
		b, err := val.Bytes()
		if err != nil {
			log.Println(err)
			continue
		}

		err = json.Unmarshal(b, &m)

		if err != nil {
			log.Println(err)
			continue
		}

		res = append(res, m)
	}

	if err := iter.Err(); err != nil {
		return res
	}

	return res
}
