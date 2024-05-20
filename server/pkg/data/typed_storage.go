package data

import "encoding/json"

type TypedStorage[T any] struct {
	base RawStorage
}

func NewTypedStorage[T any](base RawStorage) *TypedStorage[T] {
	return &TypedStorage[T]{
		base: base,
	}
}

func (ds *TypedStorage[T]) Get(id string) (T, error) {
	var result T
	data, err := ds.base.Get(id)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)

	return result, err
}

func (ds *TypedStorage[T]) Set(id string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ds.base.Set(id, data)
}

func (ds *TypedStorage[T]) Delete(id string) error {
	return ds.base.Delete(id)
}

func (ds *TypedStorage[T]) Keys() []string {
	return ds.base.Keys()
}