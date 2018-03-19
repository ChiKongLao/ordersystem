package services

import (
	"github.com/kataras/iris/websocket"
	"github.com/sirupsen/logrus"
	"regexp"
	"fmt"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/manager"
	"github.com/chikong/ordersystem/model"
	"errors"
	"github.com/chikong/ordersystem/constant"
	"strconv"
	"strings"
	"github.com/chikong/ordersystem/network"
	"time"
)


type PrinterService interface {
	GetPrinterList() (int, []model.Printer, error)
	InsertPrinter(businessId int, name string) (int, error)
	UpdatePrinter(id int, businessId int, name string) (int, error)
	DeletePrinter(businessId int) (int, error)

	SendOrder(order model.OrderPrint)

	HandleConnection(c websocket.Connection)
	HandlePayload(payload string) (string,string)


}

func NewPrinterService() PrinterService {
	mDeviceIdReg, _ = regexp.Compile(constant.SocketRegexDeviceId)
	service := &printerService{

	}
	go func() {
		time.Sleep(3 * time.Second)
		network.GetWebServer().OnConnection(func(c websocket.Connection) {
			service.HandleConnection(c)
		})
	}()

	return service
}

type printerService struct {

	}

var mDeviceIdReg *regexp.Regexp

func(s *printerService) HandleConnection(c websocket.Connection) {

	c.OnMessage(func(bytes []byte) {
		payload := string(bytes)
		println(payload)
		reply,event := s.handlePayload(payload)
		if reply == "" {
			return
		}
		c.EmitMessage([]byte(reply))
		logrus.Debugf("回复:(%s) %s => %s",event,reply,payload)
	})

}

// 服务器下发打印数据
func(s *printerService) SendOrder(order model.OrderPrint) {
	//model.TestPrinter(model.MakePrinterOrderData(1,order))
	//_, printer, err := s.GetPrinterByUserId(order.BusinessId)
	//if err != nil {
	//	return
	//}
	//network.GetWebServer().

	//list := network.GetWebServer().GetConnections()
	//logrus.Infoln(list)
}

func(s *printerService) HandlePayload(payload string) (string,string){
	reply,event := s.handlePayload(payload)
	if reply == "" {
		return "",""
	}
	logrus.Debugf("回复:(%s) %s => %s",event,reply,payload)
	return reply,event
}

func(s *printerService) handlePayload(payload string) (string,string) {

	event := ""
	defer func() {
		if event == ""{
			logrus.Warnln("未找到对应的事件:",payload)
			return
		}
		logrus.Debugf("接收:(%s) %s",event,payload)
	}()

	ok := false
	if ok = strings.Contains(payload,constant.SocketKeyPing); ok {
		event = "心跳"
		return s.handlePing(payload),event
	}
	if ok = strings.Contains(payload,constant.SocketKeyNetworkTime); ok {
		event = "查询网络延时状态"
		return s.handlePing(payload),event
	}
	if ok = strings.Contains(payload,constant.SocketKeyCheckVersion); ok {
		event = "查询打印机版本"
		return s.handlePing(payload),event
	}
	if ok = strings.Contains(payload,constant.SocketKeyIMEI); ok {
		event = "查询IMEI码"
		return s.handlePing(payload),event
	}
	if ok = strings.Contains(payload,constant.SocketKeyPrintSetting); ok {
		event = "设置打印"
		return s.handlePing(payload),event
	}
	if ok = strings.Contains(payload,constant.SocketKeyNetwordSetting); ok {
		event = "设置网络参数"
		return s.handlePing(payload),event
	}
	if ok = strings.Contains(payload,constant.SocketKeyNetwordSignal); ok {
		event = "查询网络信号值"
		return s.handlePing(payload),event
	}
	if ok = strings.Contains(payload,constant.SocketKeyNetwordSignal); ok {
		event = "查询网络信号值"
		return s.handlePing(payload),event
	}
	if ok = strings.Contains(payload,constant.SocketKeyClearOrder); ok {
		event = "清空订单数据"
		return s.handlePing(payload),event
	}
	if ok = strings.Contains(payload,constant.SocketKeyChain); ok {
		event = "打印联号设置"
		return s.handlePing(payload),event
	}
	if ok = strings.Contains(payload,constant.SocketKeyUpgradeIPAndPort); ok {
		event = "设置远程升级IP和端口号"
		return s.handlePing(payload),event
	}
	if ok = strings.Contains(payload,constant.SocketKeyOrderReceive) ||
			strings.Contains(payload,constant.SocketKeyOrderAccept) ||
			strings.Contains(payload,constant.SocketKeyOrderReject) ||
			strings.Contains(payload,constant.SocketKeyOrderTimeout); ok {
		event = "打印机:"
		if strings.Contains(payload,constant.SocketKeyOrderReceive) {
			event = event + "已收到"
		}else if strings.Contains(payload,constant.SocketKeyOrderAccept){
			event = event + "接受"
		}else 	if strings.Contains(payload,constant.SocketKeyOrderReject){
			event = event + "拒绝"
		}else 	if strings.Contains(payload,constant.SocketKeyOrderTimeout){
			event = event + "超时"
		}
		return s.handlePrinterReceive(payload),event
	}

	return "",""
}

// 处理心跳,A*88888888*0*AS01#
func(s *printerService) handlePing(payload string) string {
	data,size := getRuneAndSize(payload)
	status,_ := strconv.Atoi(string(data[size-7:size-6]))
	if status == 1{ // 缺纸
		//_, userId, err := s.GetUserIdByPrinterId(getDeviceName(payload))
		//if err == nil {
			//network.SendPrinterMessage(userId,status)
		//}
	}
	return constant.SocketFormatPingReply
}

// 处理打印机回复, A*13302920661*2345*AS04#
func(s *printerService) handlePrinterReceive(payload string) string {
	deviceId := getDeviceName(payload)
	newPayload := strings.Replace(payload,deviceId,"",1)
	data,size := getRuneAndSize(newPayload)
	orderNo := string(data[3:size-6])

	replyStatus := ""
	printStatus := ""
	if strings.Contains(payload,constant.SocketKeyOrderReceive) {
		replyStatus = "38"
		printStatus = "1"
	}else if strings.Contains(payload,constant.SocketKeyOrderAccept){
		replyStatus = "39"
		printStatus = ""
	}else 	if strings.Contains(payload,constant.SocketKeyOrderReject){
		return ""
	}else 	if strings.Contains(payload,constant.SocketKeyOrderTimeout){
		return ""
	}
	return fmt.Sprintf(constant.SocketFormatOrderReply,replyStatus,orderNo,printStatus)

}

// 获取设备id
func getDeviceName(payload string) string{
	result := mDeviceIdReg.FindString(payload)
	data, size := getRuneAndSize(result)
	return string(data[2:size-1])
}

func getRuneAndSize(payload string) ([]rune,int)  {
	return []rune(payload),len(payload)
}











// 获取打印机列表
func (s *printerService) GetPrinterList() (int, []model.Printer, error) {

	list := make([]model.Printer, 0)

	err := manager.DBEngine.Find(&list)
	if err != nil {
		logrus.Errorf("获取打印机失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取打印机失败")
	}

	return iris.StatusOK, list, nil
}


// 获取打印机对应的用户
func (s *printerService) GetPrinterByUserId(id int) (int, int, error) {

	item := new(model.Printer)

	ok, err := manager.DBEngine.Where(fmt.Sprintf("%s=?",constant.ColumnBusinessId), id).
		Get(item)
	if err != nil{
		logrus.Errorf("获取用户对应的打印机失败: %s", err)
		return iris.StatusInternalServerError, 0, errors.New("用户对应的打印机失败")
	}
	if !ok {
		return iris.StatusNotFound, 0, errors.New("没有该打印机")
	}

	return iris.StatusOK, item.Id, nil
}
// 获取打印机对应的用户
func (s *printerService) GetUserIdByPrinterId(name string) (int, int, error) {

	item := new(model.User)

	ok, err := manager.DBEngine.Table("user").Select("user.id").
		Join("INNER", "`printer`", "printer.business_id = user.id").
		Where("printer.name=?", name).
		Get(item)
	if err != nil{
		logrus.Errorf("获取打印机对应的用户失败: %s", err)
		return iris.StatusInternalServerError, 0, errors.New("获取打印机对应的用户失败")
	}
	if !ok {
		return iris.StatusNotFound, 0, errors.New("没有该用户")
	}

	return iris.StatusOK, item.Id, nil
}

// 添加打印机
func (s *printerService) InsertPrinter(businessId int, name string) (int, error) {
	if name == "" || businessId == 0 {
		return iris.StatusBadRequest, errors.New("信息不能为空")
	}
	_, err := manager.DBEngine.InsertOne(&model.Printer{
		BusinessId: businessId,
		Name:       name,
	})
	if err != nil {
		logrus.Errorf("添加打印机失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加打印机失败")
	}
	return iris.StatusOK, nil
}

// 修改打印机
func (s *printerService) UpdatePrinter(id int, businessId int, name string) (int, error) {
	if name == "" || businessId == 0 || id == 0 {
		return iris.StatusBadRequest, errors.New("信息不能为空")
	}

	_, err := manager.DBEngine.Cols(constant.ColumnBusinessId,constant.Name).Where(
		fmt.Sprintf("%s=?", constant.NameID),
		businessId).Update(model.Printer{
			BusinessId:businessId,
			Name:name,
	})
	if err != nil {
		logrus.Errorf("修改打印机失败: %s", err)
		return iris.StatusInternalServerError, errors.New("修改打印机失败")
	}

	return iris.StatusOK, nil
}

// 删除打印机
func (s *printerService) DeletePrinter(id int) (int, error) {
	_, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.NameID),
		id).Delete(new(model.Printer))
	if err != nil {
		logrus.Errorf("删除打印机失败: %s", err)
		return iris.StatusInternalServerError, errors.New("删除打印机失败")
	}
	return iris.StatusOK, nil
}
