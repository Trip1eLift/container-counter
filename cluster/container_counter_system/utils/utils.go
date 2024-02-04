package utils

import (
	"math/rand"
	"time"
)

func RandTimeS(timeRangeSecond int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(1000 * timeRangeSecond)
	//fmt.Printf("rand int: %d\n", num)
	return time.Duration(num) * time.Millisecond
}

func RandSleepS(timeRangeSecond int) {
	time.Sleep(RandTimeS(timeRangeSecond))
}
