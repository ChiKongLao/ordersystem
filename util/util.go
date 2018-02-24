package util

import (
	"net"
	"github.com/chikong/ordersystem/constant"
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

