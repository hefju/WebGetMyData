package tools
import (
    "os"
    "fmt"
    "encoding/csv"
    "io"
    "github.com/hefju/WebGetMyData/model"
)
//传入参数sh,或者sz, 来导入csv数据
func Input(stype string ){
    filename:="sh.csv"
    if stype=="sz"{
        filename="sz.csv"
    }

    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()
    reader := csv.NewReader(file)

    list:=make([]*model.Stock,0)
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Println("Error:", err)
            return
        }
        stock := new(model.Stock)
        stock.Scode2 = record[0]
        stock.Sname = record[1]
        list=append(list,stock)
        //fmt.Println(stock)
    }
    model.InsertStock(list)
}
