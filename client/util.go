package client

import (
	"math/rand"
	"time"
)

//Generate count [start, end] non repeated random numbers
func generateRandomNumber(start int, end int, count int) []int {
	// scope check
	if end < start || (end-start) < count {
		return nil
	}
	// slice for storing results
	nums := make([]int, 0)

	// random number generator, add time stamp to ensure that the random number generated each time is different
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		// generate random number
		num := r.Intn(end - start) + start
		// duplicate check
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}

func getLocalSubsetMembership(hosts []string,n int) []string{
	positionsToSelect := generateRandomNumber(0,len(hosts),n)

	result := make([]string,0,n)

	for _,v := range positionsToSelect {
		result = append(result,hosts[v])
	}

	return result
}
