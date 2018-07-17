package storage

import (
	"math/rand"
)

//RandInt64 get the random number in [min, max]
func RandInt64(min, max int64) int64 {
	if min >= max || max == 0 {
		return max
	}
	x := rand.Int63n(max-min) + min
	//fmt.Println(x)
	return x
}

//RandInt get the random numer in [min, max]
func RandInt(min, max int) int {
	if min >= max || max == 0 {
		return max
	}

	x := rand.Intn(max-min) + min

	//fmt.Println("RandInt: = ",x)
	return x
}
