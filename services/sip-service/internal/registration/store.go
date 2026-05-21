package registration

import (
	"context"
	"fmt"
	"time"
)

type Store interface {
	SaveBinding(ctx context.Context, tenantID, username, contact, remoteAddr string, expires time.Duration) error
}

type RedisStore struct {
	setFn func(ctx context.Context, key, value string, ttl time.Duration) error
}

func NewRedisStore(setFn func(ctx context.Context, key, value string, ttl time.Duration) error) *RedisStore {
	return &RedisStore{
		setFn: setFn,
	}
}

func (s *RedisStore) SaveBinding(ctx context.Context, tenantID, username, contact, remoteAddr string, expires time.Duration) error {
	key := fmt.Sprintf("sip:reg:%s:%s", tenantID, username)
	val := fmt.Sprintf("%s|%s", contact, remoteAddr)
	return s.setFn(ctx, key, val, expires)
}
