package xrand

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func RandCode() string {
	var str string
	randint := rand.Perm(9)[0:4]
	for _, v := range randint {
		str = fmt.Sprintf("%s%d", str, v)
	}
	return str
}
