package util

import (
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"
)

//PrintStack 打印堆栈，并打印传入的变量
func PrintStack(vars ...interface{}) {
	i := 0
	funcName, file, line, ok := runtime.Caller(i)
	for ; ok; i++ {
		log.Errorf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
		funcName, file, line, ok = runtime.Caller(i)
	}

	for v := range vars {
		log.Errorf("param: %v\n", v)
	}
}

//FuncCaller 得到调用者
func FuncCaller(lvl int) string {
	funcName, file, line, ok := runtime.Caller(lvl)
	if ok {
		return fmt.Sprintf("func:%v,file:%v,line:%v", runtime.FuncForPC(funcName).Name(), file, line)
	}
	return ""
}
