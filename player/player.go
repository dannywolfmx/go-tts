package player

type Player interface {
	Play() error
	Skip()
	GetText() string
}
