package data

import "github.com/DaanV2/f1-game-dashboards/server/sessions"

type (
	Database interface {
		Chairs() Storage[sessions.Chair]
		Config() RawStorage
	}

	Storage[T any] interface {
		Get(id string) (T, error)
		Set(id string, value T) error
		Keys() []string
	}

	RawStorage interface {
		Get(id string) ([]byte, error)
		Set(id string, value []byte) error
		Keys() []string
	}
)