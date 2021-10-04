package repositories

type Mode uint8

const (
	FreeMode   Mode = 0
	StrictMode Mode = 1
)

func (m Mode) IsStrict() bool {
	return m == StrictMode
}

func (m Mode) IsFree() bool {
	return m == FreeMode
}
