package hooks

type Hook[T any] struct {
	// handlers is a slice of functions that take a value of type T, use get or set to access it.
	handlers []func(T)
}

func NewHook[T any]() Hook[T] {
	return Hook[T]{}
}

func (h *Hook[T]) get() []func(T) {
	return h.handlers
}

func (h *Hook[T]) set(handlers []func(T)) {
	h.handlers = handlers
}

func (h *Hook[T]) Active() bool {
	return len(h.get()) > 0
}

func (h *Hook[T]) Add(handler func(T)) {
	handlers := h.get()
	handlers = append(handlers, handler)
	h.set(handlers)
}

func (h *Hook[T]) Call(value T) {
	handlers := h.get()

	for _, handler := range handlers {
		go handler(value)
	}
}
