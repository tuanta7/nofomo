package strategy

type Signal uint8

const (
	Hold Signal = iota
	Buy
	Sell
)

func (s Signal) String() string {
	switch s {
	case Buy:
		return "buy"
	case Sell:
		return "sell"
	default:
		return "hold"
	}
}
