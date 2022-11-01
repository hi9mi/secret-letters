package main

import "github.com/google/uuid"

type KeyGen interface {
	Get() string
}

type TestKeyGen struct{}

func (t *TestKeyGen) Get() string {
	return "test"
}

type UUIDKeyGen struct{}

func (u *UUIDKeyGen) Get() string {
	return uuid.New().String()
}

func getTestKeyGen() KeyGen {
	return &TestKeyGen{}
}

func getUUIDKeyGen() KeyGen {
	return &UUIDKeyGen{}
}
