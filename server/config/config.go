package config

type (
	ItemNotFoundError struct {
		Key string
	}

	ItemNotType struct {
		Key string
	}
)

func (e *ItemNotFoundError) Error() string {
	return "item not found: " + e.Key
}

func (e *ItemNotType) Error() string {
	return "item not of expected type for key: " + e.Key
}

func IsNotFound(err error) bool {
	_, ok := err.(*ItemNotFoundError)
	return ok
}

func IsNotType(err error) bool {
	_, ok := err.(*ItemNotType)
	return ok
}

func Get[T any](c map[string]interface{}, key string) (T, error) {
	var result T
	item, ok := c[key]
	if !ok {
		return result, &ItemNotFoundError{Key: key}
	}

	if v, ok := item.(T); ok {
		return v, nil
	}

	return result, &ItemNotType{Key: key}
}
