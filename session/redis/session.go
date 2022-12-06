package redis

import (
	"Go_Web/session"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"time"
)

var (
	errSessionNotFound = errors.New("session: id 对应的 session 不存在")
)

type StoreOption func(store *Store)

// 在redis中以 hset 这种数据结构形式存在
//    sess id     key    value
// map[string]map[string]string
type Store struct {
	prefix     string
	client     redis.Cmdable
	expiration time.Duration
}

func NewStore(client redis.Cmdable, opts ...StoreOption) *Store {
	res := &Store{
		expiration: time.Millisecond,
		client:     client,
		prefix:     "sessid",
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func StoreWithPrefix(prefix string) StoreOption {
	return func(store *Store) {
		store.prefix = prefix
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	key := redisKey(s.prefix, id)
	// 在redis上创建一个map然后和id绑定在一起
	_, err := s.client.HSet(ctx, key, id, id).Result()
	if err != nil {
		return nil, err
	}
	_, err = s.client.Expire(ctx, key, s.expiration).Result()
	if err != nil {
		return nil, err
	}
	return &Session{
		key:    key,
		client: s.client,
		id:     id,
	}, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	key := redisKey(s.prefix, id)
	ok, err := s.client.Expire(ctx, key, s.expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errSessionNotFound
	}
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	key := redisKey(s.prefix, id)
	_, err := s.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	// 代表的是 id 对应的 session 不存在, 没有删除掉任何东西
	//if cnt == 0 {
	//
	//}
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	// 自由决策要不要提前把 session 存储的用户数据都拿过来
	// 1. 都不要
	// 2. 只拿热点数据
	// 3. 都拿
	key := redisKey(s.prefix, id)
	cnt, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if cnt != 1 {
		return nil, errSessionNotFound
	}
	return &Session{
		key:    key,
		id:     id,
		client: s.client,
	}, nil
}

type Session struct {
	key    string
	prefix string
	id     string
	client redis.Cmdable
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	val, err := s.client.HGet(ctx, s.key, key).Result()
	return val, err
}

func (s *Session) Set(ctx context.Context, key string, val any) error {
	const lua = `
if redis.call("exists", KEYS[1])
then
	return redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
else
	return -1
end
`
	res, err := s.client.Eval(ctx, lua, []string{s.key}, key, val).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return errSessionNotFound
	}
	return nil
}

func (s *Session) ID() string {
	return s.id
}

func redisKey(prefix, id string) string {
	return fmt.Sprintf("%s-%s", prefix, id)
}
