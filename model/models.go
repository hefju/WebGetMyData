package model

import (
	"fmt"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

//每日价格数据
type Daily struct {
	Id      int64
	Scode   string //代号
	DateStr string //日期
	Open    string //开盘价
	Highest string //最高价
	Close   string //收盘价
	Low     string //最低价
	Volume  string //交易量(股)
	Amount  string //交易金额(元)
	//    DateStr time.Time //日期
	//    Open    float64 //开盘价
	//    Highest float64 //最高价
	//    Close   float64 //收盘价
	//    Low     float64 //最低价
	//    Volume  float64 //交易量(股)
	//    Amount  float64 //交易金额(元)
}
type Daily2 struct {
	Id      int64
	Scode   string    //代号
	DateStr time.Time //日期
	Open    float64   //开盘价
	Highest float64   //最高价
	Close   float64   //收盘价
	Low     float64   //最低价
	Volume  float64   //交易量(股)
	Amount  float64   //交易金额(元)
}

func GetDaily(code string) []Daily {
	list := make([]Daily, 0)
	engine.Where("scode=?", code).Find(&list)
	return list
}

func InsertDaily(list []*Daily) {
	session := engine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		fmt.Println(err)
	}
	count := 0
	for _, item := range list {
		_, err := session.Insert(item)
		if err != nil {
			session.Rollback()
			fmt.Println("Insert err:", err)
			return
		}
		count++
	}
	err = session.Commit()
	if err != nil {
		return
	}
	fmt.Println("Insert result:", count)
}

type Stock struct {
	Id     int64
	Scode  string //纯数字
	Scode2 string //字母+数字
	Sname  string
	Stype  string
	Status int //-2表示没有用的取不到信息的,-1表示未设置,0表示正常,
}

func GetStockName(stype string) []*Stock {
	list := make([]*Stock, 0)
	engine.Where("stype=? and id>=1050 and id<1100", stype).Find(&list)
	return list
}
func InsertStock(list []*Stock) {
	session := engine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		fmt.Println(err)
	}
	count := 0
	for _, item := range list {
		_, err := session.Insert(item)
		if err != nil {
			session.Rollback()
			fmt.Println("Insert err:", err)
			return
		}
		count++
	}
	err = session.Commit()
	if err != nil {
		return
	}
	fmt.Println("Insert result:", count)
}

var engine *xorm.Engine

func init() {
	var err error
	//	engine, err = xorm.NewEngine("odbc", "driver={SQL Server};Server=192.168.1.200; Database=charge; uid=sa; pwd=123;")
	engine, err = xorm.NewEngine("sqlite3", "./WebGetMyData.db")

	if err != nil {
		log.Fatalln("xorm create error", err)
	}
	//engine.ShowSQL = true
	engine.SetMapper(core.SameMapper{})
	// engine.CreateTables(new(tp_charge_billing))
	//err = engine.Sync2(new(Daily)) //, new(Group))
	err = engine.Sync2(new(Daily), new(Stock)) //, new(Group))
	if err != nil {
		log.Fatalln("xorm sync error", err)
	}
}
