package main

import (
	"fmt"
	"github.com/CodisLabs/codis/pkg/utils/log"
	"localhostTest/go-Call/call"
	"localhostTest/go-Call/test"
	"localhostTest/go-worker-base/worker"
	"net/http"
	"os"
	"runtime/pprof"
	"time"
)

//声明一号池子
var poolOne worker.WorkPool

//声明一号池子
var poolTwo worker.WorkPool

//声明回调变量
var funcs call.CallMap

//以结构体方式调用
type runWorker struct{}

//初始化协程池 和回调参数
func init() {
	poolOne = worker.GetPool("one")
	poolOne.Start(50)
	funcs = call.CreateCall()

}

//通用回调
func (f *runWorker) Run(param []worker.ParamType) {
	name := param[0].(string)
	var callParam []call.Params
	//将结果转换成 回调需要的数据
	for k, v := range param {
		if k == 0 {
			continue
		}
		callParam = append(callParam, v)
	}
	//调用回调并拿回结果
	funcs.Call(name, callParam...)
	//var re string = result[0].String()
	//fmt.Println(re)
}

//主函数
func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	f, _ := os.Create("profile_file")
	pprof.StartCPUProfile(f)
	var runFunc runWorker = runWorker{}
	funcs.AddCall("test4", test.Test4)
	var startTime = time.Now().Unix()
	for i := 0; i < 10000; i++ {
		poolOne.Run(runFunc.Run, "test4", " aa ", " BB")
		poolOne.Run(runFunc.Run, "test4", " cc ", " dd")
		poolOne.Run(runFunc.Run, "test4", " ee ", " ff")
	}
	var modTime = time.Now().Unix()

	for k := 0; k < 10000; k++ {
		test.Test4(" aa ", "BB")
		test.Test4(" cc ", " dd")
		test.Test4(" ee ", " ff")
	}
	var endTime = time.Now().Unix()
	for j := 0; j < 10000; j++ {
		funcs.Call("test4", " aa ", "BB")
		funcs.Call("test4", " cc ", " dd")
		funcs.Call("test4", " ee ", " ff")
	}
	var lastTime = time.Now().Unix()
	fmt.Println(modTime - startTime)
	fmt.Println(endTime - modTime)
	fmt.Println(lastTime - endTime)

	fmt.Println(startTime, modTime, endTime)
	defer pprof.StopCPUProfile()
	//time.Sleep(time.Millisecond * 1000)
}
