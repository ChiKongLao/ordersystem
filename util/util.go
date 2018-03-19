package util

import (
	"net"
	"github.com/chikong/ordersystem/constant"
	"github.com/satori/go.uuid"
	"time"
	"unicode"
	"strconv"
)

// 获取本机IP
func GetLocalIP() string{
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String() +":"+ constant.SystemHost
				}

			}
		}
	}
	return "localhost:"+constant.SystemHost
}
// 获取本机IP,带http
func GetLocalIPWithHttp() string{
	return "http://"+GetLocalIP()
}

// 获取UUID
func GetUUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}

// 系统当前时间戳,毫秒
func GetCurrentTime() int64 {
	return time.Now().UnixNano()/1000/1000
	//return strconv.FormatInt(time.Now().UnixNano()/1000/1000, 10)
}

func GetCurrentFormatTime() string {
	return time.Now().Format("20060102150405")
}
// 获取今天零时的时间戳
func GetTodayZeroTime() int64{
	now := time.Now()
	t, _ := time.ParseInLocation("2006-01-02", now.Format("2006-01-02"), time.Local)
	return t.Unix()*1000
}

// 计算长度,中文长度为2
func GetLen(str string) int{
	size := 0
	for _, r := range str {
		if IsChineseChar(r) {
			size += 2
		}else {
			size++
		}
	}
	return size
}

// 判断字符是否为中文字符
func IsChineseChar(r rune) bool {
	return unicode.Is(unicode.Scripts["Han"], r)
}

func Float32ToString(value float32) string {
	return strconv.FormatFloat(float64(value),'f',-1,32)

}
