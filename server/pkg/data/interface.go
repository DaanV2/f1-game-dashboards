package data

import "github.com/DaanV2/f1-game-dashboards/server/sessions"

type (
	IDatabase interface {
		Chairs() IStorage[sessions.Chair]
	}

	IStorage[T any] interface {
		Get(id string) (T, error)
		Set(id string, value T) error
		Keys() []string
	}
)