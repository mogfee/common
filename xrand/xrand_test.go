package xrand

import (
	"fmt"
	"testing"
)

func TestRandString(t *testing.T) {
	for i := 0; i < 1000; i++ {
		fmt.Println(RandCode())
	}
}
