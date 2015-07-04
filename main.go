package main

import (
	"fmt"
	"github.com/hefju/WebGetMyData/controllers"
)

func main() {

    controllers.GoWorking()
	// tr_chan := make(chan string)    //存放数据行(tr)
	// product_done := make(chan bool) //生产者工作完成chan
	// job_done := make(chan bool)     //工作者完成chan

	// // var addr string = "http://600048.phtml?year=2015&jidu=1"
	// // content := httpget(addr)

	// filehtml := FromFile() //获取html数据
	// go FindTr_empty(filehtml, tr_chan, product_done)
	// go FindTr_gray(filehtml, tr_chan, product_done)

	// go handleOne(tr_chan, job_done)

	// <-product_done //生产者完成标志
	// <-product_done
	// close(tr_chan) //关闭数据通道,通知处理程序没有数据传入了

	// <-job_done //处理程序通知main, 已经处理完毕
	// fmt.Println(DailyList)
	fmt.Println("end")
}
