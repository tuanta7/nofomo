package accounts

import (
	"github.com/tuanta7/nofomo/pkg/lightspeed/trading"
)

type Client struct {
	restClient *trading.Client
}
