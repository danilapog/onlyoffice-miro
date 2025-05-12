package service

import "context"

type Storage[K comparable, V any] interface {
	Find(context.Context, K) (V, error)
	Insert(context.Context, K, V) (V, error)
	Update(context.Context, K, V) (V, error)
	Delete(context.Context, K) error
}
