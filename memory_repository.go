package main

import (
	"encoding/hex"
	"errors"
	"os"
	"sync"
)

const TestKey = "test"

type MemoryRepository struct {
	data map[string]string
	mu   *sync.Mutex
}

func (repo *MemoryRepository) Get(key string) (string, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	defer delete(repo.data, key)
	encryptedLetter, ok := repo.data[key]

	if !ok {
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

func (repo *MemoryRepository) Check(key string) bool {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	_, ok := repo.data[key]
	return ok
}

func (repo *MemoryRepository) Set(key, letter string, ttl int) error {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	secretKeyStr := hex.EncodeToString(secretKey)
	encryptedLetter, err := encrypt(secretKeyStr, letter)
	if err != nil {
		return errors.New("failed to encrypt secret")
	}
	repo.data[key] = encryptedLetter
	return nil
}
