package auxiliary

import (
	"math/rand"
	"strconv"
	"time"
)

// PortRandom generates pseudo random port in range 20000-65000
func PortRandom() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return strconv.Itoa(rand.Intn(65000-20000) + 20000)
}
