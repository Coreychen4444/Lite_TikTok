package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Coreychen4444/Lite_TikTok/model"
	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{rdb: rdb}
}

var ctx = context.Background()

// 在redis 中获取好友列表id
func (r *RedisRepository) GetFriendList(userId int64) ([]int64, error) {
	cacheKey := fmt.Sprintf("user:%d:friends_list", userId)
	ids, err := r.rdb.SMembers(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(ids) > 0 {
		var idList []int64
		for _, id := range ids {
			userId, err := strconv.ParseInt(id, 10, 64)
			if err == nil {
				idList = append(idList, userId)
			}
		}
		return idList, nil
	}
	return []int64{}, nil
}

// 将好友列表id存入redis
func (r *RedisRepository) AddFriendList(userId int64, friends []model.User) error {
	cacheKey := fmt.Sprintf("user:%d:friends_list", userId)
	var ids []interface{}
	for _, friend := range friends {
		ids = append(ids, friend.ID)
	}
	_, err := r.rdb.SAdd(ctx, cacheKey, ids...).Result()
	if err != nil {
		return err
	}
	r.rdb.Expire(ctx, cacheKey, time.Hour)
	return nil
}
