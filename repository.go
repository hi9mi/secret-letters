package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
)

type Repository interface {
	Get(key string) (string, error)
	Set(key, letter string, ttl int) error
	Check(key string) bool
}

func getMemoryRepository() Repository {
	return &MemoryRepository{
		data: make(map[string]string),
		mu:   &sync.Mutex{},
	}
}

func getRedisRepository() Repository {
	var redisOpts *redis.Options

	if strings.Contains(os.Getenv("LOCAL"), "true") {
		redisOpts = &redis.Options{
			Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_URL")),
			Password: "",
			DB:       0,
		}
	} else {
		builtOpts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
		if err != nil {
			log.Fatal(err)
		}
		redisOpts = builtOpts
	}

	return &RedisRepository{
		*redis.NewClient(redisOpts), context.Background()}
}

func encrypt(keyString string, stringToEncrypt string) (string, error) {
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), err
}

func decrypt(keyString string, stringToDecrypt string) (string, error) {
	key, _ := hex.DecodeString(keyString)
	ciphertext, _ := base64.URLEncoding.DecodeString(stringToDecrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
