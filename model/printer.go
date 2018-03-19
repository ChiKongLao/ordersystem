package model

import (
	"github.com/chikong/ordersystem/constant"
	"fmt"
	"bytes"
	"github.com/chikong/ordersystem/util"
	"strings"
	"strconv"
)

// 打印机
type Printer struct {
	Id          int     `json:"id" xorm:"not null pk autoincr unique INT"`
	Name        string  `json:"name" xorm:"not null VARCHAR(255)"`
	BusinessId  int     `json:"-" xorm:"not null index INT(11)"`

}

func MakePrinterOrderData(status int,order OrderPrint) string{
	var totalCount int
	//&!*XXXX*XXXX*<big>XXXXXX*<S0XT>XXXX*<BMP>*<qrc>XXXX#
	var content bytes.Buffer
	//content = "<big>XXXXXX*<S0XT>XXXX*<BMP>*<qrc>XXXX"
	add(&content,makeTwoContent("桌号",order.TableName))
	add(&content,makeNewLine())
	add(&content,makeTwoContent("会员",order.Customer.UserName))
	add(&content,makeNewLine())
	add(&content,makeTwoContent("时间",util.GetCurrentFormatTime()))
	add(&content,makeNewLine())
	add(&content,makeInterval())
	add(&content,makeNewLine())
	add(&content,makeFourContent("名称"," 数量","单价","金额"))
	add(&content,makeNewLine())
	add(&content,makeInterval())
	add(&content,makeNewLine())
	for _, subItem := range order.FoodList {
		add(&content,subItem.Name)
		add(&content,makeNewLine())
		add(&content,makeFourContent("",strconv.Itoa(subItem.Num),
			util.Float32ToString(subItem.Price),
			util.Float32ToString(subItem.GetTotalPrice())))

		add(&content,makeNewLine())

		totalCount += subItem.Num
	}
	add(&content,makeInterval())
	add(&content,makeNewLine())

	add(&content,makeFourContent("合计",strconv.Itoa(totalCount),
		"",util.Float32ToString(order.Price)))
	add(&content,makeNewLine())
	add(&content,makeTwoContent("应收金额",util.Float32ToString(order.Price)))
	add(&content,makeNewLine())
	add(&content,makeInterval())
	add(&content,makeNewLine())

	add(&content,makeCenterContent("欢迎您下次光临!"))
	add(&content,makeNewLine())
	add(&content,makeNewLine())
	add(&content,fmt.Sprintf("地址:%s",order.Business.Address))
	add(&content,makeNewLine())
	add(&content,fmt.Sprintf("电话:%s",order.Business.Phone))
	add(&content,makeNewLine())

	return fmt.Sprintf(constant.SocketFormatOrderSend1,status,order.OrderNo,content.String())

}


func makeTwoContent(str1, str2 string) string{
	return makeTwoContentWithMaxLen(constant.PrinterMaxLen,str1,str2)
}

func makeTwoContentWithMaxLen(maxLen int, str1, str2 string) string{
	space := maxLen - util.GetLen(str1) - util.GetLen(str2)
	return str1 + makeSpace(space) + str2
}


func makeFourContent(str1, str2, str3, str4 string) string{
	var result string
	result = makeTwoContentWithMaxLen(constant.PrinterMaxLen4_2,str1,str2)
	result = makeTwoContentWithMaxLen(constant.PrinterMaxLen4_3,result,str3)
	result = makeTwoContentWithMaxLen(constant.PrinterMaxLen,result,str4)

	return result

}

func makeCenterContent(str string) string{
	space := constant.PrinterMaxLen - util.GetLen(str)
	return makeSpace(space/2)+str

}

// 换行
func makeNewLine() string{
	return "*"
}

// 间隔符"-------------"
func makeInterval() string{
	return "--------------------------------"// 32
}

// 空格符
func makeSpace(count int) string{
	var result string
	for i:=0; i<count ;i++ {
		result = result + " "
	}
	return result
}

func TestPrinter(payload string){
	payload = strings.Replace(payload,"*","\n",-1)
	data := []rune(payload)
	index := strings.Index(payload,"桌号")
	data = data[index:]
	println(string(data))
}

func add(data *bytes.Buffer, content string) {
	data.WriteString(content)
}
