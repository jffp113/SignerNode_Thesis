package client

import (
	"github.com/go-ping/ping"
	"math/rand"
	"strings"
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
		num := r.Intn(end-start) + start
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

func GetSubsetMembership(hosts []string, n int) []string {
	positionsToSelect := generateRandomNumber(0, len(hosts), n)

	result := make([]string, 0, n)

	for _, v := range positionsToSelect {
		result = append(result, hosts[v])
	}

	return result
}

func GetRandomGroupMember(GroupMembership []string) string {
	seed := time.Now().UTC().UnixNano()
	rnd := rand.New(rand.NewSource(seed))

	pos := rnd.Intn(len(GroupMembership))
	return GroupMembership[pos]
}

func GetNNodesRandomly(signernodes []string, n int) []string {
	var result []string
	alreadyChoosen := make(map[int]bool)
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	var i int
	for i < n {
		pos := rnd.Intn(len(signernodes))

		if alreadyChoosen[pos] {
			continue
		}
		alreadyChoosen[pos] = true
		elm := signernodes[pos]
		result = append(result, elm)
		i++
	}

	return result
}

type RTT struct {
	rtt     time.Duration
	address string
}

func GetNNearestNodes(signernodes []string, n int, parallel bool) []string {
	if parallel {
		return GetNNearestNodesParallel(signernodes,n)
	}else{
		return GetNNearestNodesSerial(signernodes,n)
	}
}

func GetNNearestNodesSerial(signernodes []string, n int) []string{
	var rtts []RTT

	if len(signernodes) < n {
		return []string{}
	}

	for i := range signernodes {
		addr := strings.Split(signernodes[i], ":")
		pinger, err := ping.NewPinger(addr[0])

		if err != nil {
			continue
		}

		pinger.Count = 3
		pinger.Run()
		stats := pinger.Statistics()

		rtts = insertOrder(rtts, RTT{
			rtt:     stats.AvgRtt,
			address: signernodes[i],
		})
	}

	if len(rtts) < n {
		return []string{}
	}

	result := make([]string, 0, n)
	for i := range rtts {
		result = append(result, rtts[i].address)
	}

	return result[:n]
}

func GetNNearestNodesParallel(signernodes []string, n int) []string {
	var rtts []RTT

	if len(signernodes) < n {
		return []string{}
	}

	responseCh := make(chan RTT,1)

	for i := range signernodes {
		go func(i int) {
			addr := strings.Split(signernodes[i], ":")

			if pinger, err := ping.NewPinger(addr[0]); err == nil {
				pinger.Count = 3
				pinger.Run()
				stats := pinger.Statistics()

				responseCh <- RTT{
					rtt:     stats.AvgRtt,
					address: signernodes[i],
				}
			}
		}(i)
	}

	for i := 0 ; i < n ; i++ {
		rtts = insertOrder(rtts, <-responseCh)
	}

	if len(rtts) < n {
		return []string{}
	}

	result := make([]string, 0, n)
	for i := range rtts {
		result = append(result, rtts[i].address)
	}

	return result[:n]
}



func insertOrder(rtts []RTT, v RTT) []RTT {
	for i := range rtts {
		if rtts[i].rtt > v.rtt {
			rtts = append(rtts[:i+1], rtts[i:]...)
			rtts[i] = v
			return rtts
		}
	}
	return append(rtts, v)
}
