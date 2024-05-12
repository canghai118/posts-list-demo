package service

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewPageService(rdb *redis.Client, pageKey string) *PageService {
	return &PageService{
		rdb:     rdb,
		pageKey: pageKey,
	}
}

type PageService struct {
	rdb     *redis.Client
	pageKey string
}

func (p *PageService) UpdateIndex(ctx context.Context, id int, score float64) error {
	_, err := p.rdb.ZAdd(ctx, p.pageKey, redis.Z{Member: id, Score: score}).Result()
	return err
}

func (p *PageService) GetPage(ctx context.Context, pageNum int, pageSize int) ([]int, error) {
	start := (pageNum - 1) * pageSize
	end := start + pageSize - 1

	ids, err := p.rdb.ZRevRange(ctx, p.pageKey, int64(start), int64(end)).Result()
	if err != nil {
		return nil, err
	}

	return convertStringSliceToInt(ids)
}

func (p *PageService) GetTotoal(ctx context.Context) (int64, error) {
	return p.rdb.ZCard(ctx, p.pageKey).Result()
}

func convertStringSliceToInt(ids []string) ([]int, error) {
	result := make([]int, len(ids))
	for i, id := range ids {
		var err error
		result[i], err = strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
