package internal

type Content interface {
	GetKey() string
	GetValue() string
}

type Header struct {
	key   string
	value string
}

func (h *Header) GetKey() string {
	return h.key
}

func (h *Header) GetValue() string {
	return h.value
}

type Body struct {
	key   string
	value string
}

func (b *Body) GetKey() string {
	return b.key
}

func (b *Body) GetValue() string {
	return b.value
}
