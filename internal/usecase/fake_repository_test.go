package usecase_test

import (
	"context"
	"time"
)

type FakeRepository struct {
	counters map[string]int64
	blocked  map[string]bool
}

func NewFakeRepository() *FakeRepository {
	return &FakeRepository{
		counters: make(map[string]int64),
		blocked:  make(map[string]bool),
	}
}

func (f *FakeRepository) IsBlocked(
	ctx context.Context,
	key string,
) (bool, error) {
	return f.blocked[key], nil
}

func (f *FakeRepository) Increment(
	ctx context.Context,
	key string,
	window time.Duration,
) (int64, error) {

	f.counters[key]++

	return f.counters[key], nil
}

func (f *FakeRepository) Block(
	ctx context.Context,
	key string,
	duration time.Duration,
) error {

	f.blocked[key] = true

	return nil
}
