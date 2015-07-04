package controllers

import (
	"fmt"
    "github.com/guotie/gogb2312"
	"github.com/hefju/WebGetMyData/model"
	"github.com/hefju/WebGetMyData/myconfig"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var DailyList []*model.Daily //每日数据数组

//开始获取数据工作
func GoWorking() {

	code := "600048"
	GetCode2015(code)

    fmt.Println(DailyList)
	model.InsertDaily(DailyList)
}

func GetCode2015(code string) {
//	url1 := myconfig.BaseAddr + code + ".phtml?year=2015&jidu=1" //600048.phtml?year=2015&jidu=1
//	url2 := myconfig.BaseAddr + code + ".phtml?year=2015&jidu=2"
	url3 := myconfig.BaseAddr + code + ".phtml?year=2015&jidu=3"
//	GetOneCode(url1)
//	GetOneCode(url2)
	GetOneCode(url3)
}

//获取一个代号的价格
func GetOneCode(url string) {

	tr_chan := make(chan string)    //存放数据行(tr)
	product_done := make(chan bool) //生产者工作完成chan
	job_done := make(chan bool)     //工作者完成chan

    //fmt.Println(url)
	filehtml := httpget(url) // FromFile() //获取html数据
    fmt.Println(filehtml)
	go FindTr_empty(filehtml, tr_chan, product_done)
	go FindTr_gray(filehtml, tr_chan, product_done)

	go handleOne(tr_chan, job_done)

	<-product_done //生产者完成标志
	<-product_done
	close(tr_chan) //关闭数据通道,通知处理程序没有数据传入了

	<-job_done //处理程序通知main, 已经处理完毕
}

//处理chan中的数据
func handleOne(tr_chan chan string, done chan bool) {
	for elem := range tr_chan {
		if HasDate(elem) {
			Extract(elem)
		}
	}
	done <- true
	// fmt.Println("Channels Closing")
}

//分析一行数据
func Extract(elem string) {

	//日期和数字需要独立提取, 因为他们的正则表达式不同
	date := GetDate(elem)     //日期
	list := make([]string, 0) //多个价格数据

	//提取价格
	var hrefRegexp = regexp.MustCompile(`"center">.*?</div>`) //<tr class="gray">
	match := hrefRegexp.FindAllString(elem, -1)
	if match != nil {
		for _, v := range match {
			tmp := strings.Replace(v, `"center">`, "", 1) //去掉前缀
			tmp2 := strings.Replace(tmp, `</div>`, "", 1) //去掉后缀
			list = append(list, tmp2)
		}
	}

	daily := &model.Daily{} //每日价格数据
	daily.DateStr = date
	daily.Open = list[0]
	daily.Highest = list[1]
	daily.Close = list[2]
	daily.Low = list[3]
	daily.Volume = list[4]
	daily.Amount = list[5]
	DailyList = append(DailyList, daily)
	//fmt.Println(daily)
}

//找到<tr class="">的标签, 这表示table的一行数据, 注意目标的table是隔行样式的
func FindTr_empty(html string, tr_chan chan string, done chan bool) {
    //var hrefRegexp = regexp.MustCompile(`<tr class="">[\s\S]*?</tr>`) //<tr >
    var hrefRegexp = regexp.MustCompile(`<tr >[\s\S]*?</tr>`)
	match := hrefRegexp.FindAllString(html, -1)
	if match != nil {
		for i, v := range match {
			//fmt.Println("[", i, "]-", v)
			_ = i
			tr_chan <- v
		}
	}

	done <- true
}

//找到<tr class="gray">的标签, 这表示table的一行数据, 注意这行数据未必是我们想要的数据, 例如表头,例如页面上其他table
func FindTr_gray(html string, tr_chan chan string, done chan bool) {
    //var hrefRegexp = regexp.MustCompile(`<tr class="gray">[\s\S]*?</tr>`) //<tr class="gray">
    var hrefRegexp = regexp.MustCompile(`<tr class="tr_2">[\s\S]*?</tr>`) //<tr class="gray">
	match := hrefRegexp.FindAllString(html, -1)
	if match != nil {
		for i, v := range match {
			//fmt.Println("[", i, "]-", v)
			_ = i
			tr_chan <- v
		}
	}
	done <- true
}

//判断是否存在日期,没有日期将不会处理
func HasDate(html string) bool {
    var hrefRegexp = regexp.MustCompile(`date=.*?'>`)
	return hrefRegexp.MatchString(html)
}

//提取日期
func GetDate(html string) string {
    var hrefRegexp = regexp.MustCompile(`date=.*?'>`)
    match := hrefRegexp.FindAllString(html, -1)
    var date string
    if match != nil {
        for _, v := range match {
            tmp := strings.Replace(v, "date=", "", 1)
            tmp2 := strings.Replace(tmp, `'>`, "", 1)
            date = tmp2
        }
    }
    return date
}

//Go正则提取html A 连接标签 (从网上拿到的例子,我的代码参考这个开始)
func ListHref(html string) {
	var hrefRegexp = regexp.MustCompile("(?m)<a.*?[^<]>.*?</a>")
	match := hrefRegexp.FindAllString(html, -1)
	fmt.Println(match)
	if match != nil {
		for i, v := range match {
			fmt.Println("[", i, "]-", v)
		}
	}
}

//从文件中读取html
func FromFile() string {
	path := "600048.html" //"600048-3.html"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	return string(bs)

}

//网络请求
func httpget(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("httpget:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("httpget-ioutil.ReadAll:", err)
	}
	//bodystr := string(body)
    output, err, _, _ := gogb2312.ConvertGB2312(body) //ConvertGB2312String接收参数为string
    if err != nil {
        fmt.Println(err)
    }

    bodystr := string(output)
	return bodystr
}
