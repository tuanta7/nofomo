package lightspeed

const maxStreams = 200

type Resolution string

const (
	Resolution1Minute   Resolution = "1"
	Resolution3Minutes  Resolution = "3"
	Resolution5Minutes  Resolution = "5"
	Resolution15Minutes Resolution = "15"
	Resolution30Minutes Resolution = "30"
	Resolution1Hour     Resolution = "1H"
	Resolution1Day      Resolution = "1D"
	Resolution1Week     Resolution = "1W"
)

func (r Resolution) Valid() bool {
	switch r {
	case
		Resolution1Minute,
		Resolution3Minutes,
		Resolution5Minutes,
		Resolution15Minutes,
		Resolution30Minutes,
		Resolution1Hour,
		Resolution1Day,
		Resolution1Week:
		return true
	default:
		return false
	}
}
