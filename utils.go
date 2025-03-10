package go_ethernet_ip

import (
	"math/rand"
	"time"

	"github.com/teldio-operations/go-ethernet-ip/typedef"
)

func CtxGenerator() typedef.Ulint {
	rand.Seed(time.Now().UnixNano())
	return typedef.Ulint(rand.Uint64())
}
