package timeutils

import (
	"fmt"
	"testing"
	"time"
)

func TestUnix64TimeToUnix32Time(t *testing.T) {
	var data int64
	data = time.Now().Unix()

	d2 := Unix64TimeToUnix32Time(data)

	fmt.Println(data)
	fmt.Println(d2)
}
