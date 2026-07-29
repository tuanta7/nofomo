package strategy

type Context struct{}

type Signal uint8

type Strategy interface {
	Evaluate(Context) Signal
}
