package main

import (
	"context"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	rdb redis.Client
	ctx context.Context
}

func (repo *RedisRepository) Get(key string) (string, error) {
	encryptedLetter, err := repo.rdb.GetDel(repo.ctx, key).Result()


	if err == redis.Nil || err != nil {
		return "", errors.New("letter not found")
	}
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	secretKeyStr := hex.EncodeToString(secretKey)
	letter, err := decrypt(secretKeyStr, encryptedLetter)

	if err != nil {
		return "", errors.New("failed to decrypt letter")
	}

	return letter, nil
}

func (repo *RedisRepository) Check(key string) bool {
	_, err := repo.rdb.Get(repo.ctx, key).Result()
	return err != redis.Nil || err != nil 
}

func (repo *RedisRepository) Set(key, letter string, ttl int) error {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	secretKeyStr := hex.EncodeToString(secretKey)
	encryptedLetter, err := encrypt(secretKeyStr, letter)
	if err != nil {
		return errors.New("failed to encrypt secret")
	}
	err = repo.rdb.Set(repo.ctx, key, encryptedLetter, time.Duration(ttl) * time.Second).Err()

	if err != nil || err == redis.Nil {
		return errors.New("unable to set letter")
	}

	return nil
}
