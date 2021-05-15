package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestGetNNodesRandomly(t *testing.T) {
//	cli := NewPermissionlessClient()
//	test := []string{"1","2","3","4","5","6","7","8"}
//	s := cli.GetNNodesRandomly(test,5)
//	fmt.Println(s)
//}
//

//func TestGetNNearestNodes(t *testing.T) {
//	cli := NewPermissionlessClient()
//
//	cli.GetNNearestNodes([]string{"149.210.219.5:1234","51.83.75.29:1234","95.217.57.18:1234","192.199.248.75:1234"},2)
//
//}

func TestGetInsertSort(t *testing.T) {
	v1 := RTT{
		rtt:     10,
		address: "1",
	}
	v2 := RTT{
		rtt:     9,
		address: "1",
	}
	v3 := RTT{
		rtt:     11,
		address: "1",
	}
	v4 := RTT{
		rtt:     100,
		address: "1",
	}

	rtts := insertOrder([]RTT{}, v1)
	assert.Equal(t, v1, rtts[0])

	rtts = insertOrder(rtts, v2)
	assert.Equal(t, v2, rtts[0])
	assert.Equal(t, v1, rtts[1])

	rtts = insertOrder(rtts, v3)
	assert.Equal(t, v2, rtts[0])
	assert.Equal(t, v1, rtts[1])
	assert.Equal(t, v3, rtts[2])

	rtts = insertOrder(rtts, v4)
	assert.Equal(t, v2, rtts[0])
	assert.Equal(t, v1, rtts[1])
	assert.Equal(t, v3, rtts[2])
	assert.Equal(t, v4, rtts[3])

}
