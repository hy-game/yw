package logic

import (
	"math"
	"sync"
)

var gateCliCnt sync.Map
var gateEndPoint string

func updateGateCli(endPoint string, cliCnt int32){
	if cliCnt == -1{
		gateCliCnt.Delete(endPoint)
	}else {
		gateCliCnt.Store(endPoint, cliCnt)
	}

	var minCnt = int32(math.MaxInt32)
	gateCliCnt.Range(func(k, v interface{})bool{
		if v.(int32) < minCnt{
			gateEndPoint = endPoint
		}
		return  true
	})
}

func getGateEndPoint()string{
	return gateEndPoint
}

