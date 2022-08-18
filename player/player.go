package player

type Player interface {
	Play() error
	Stop()
	GetText() string
}
