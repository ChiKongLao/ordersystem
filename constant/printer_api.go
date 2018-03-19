package constant

// 打印机通讯api

const (
	SocketKeyBIG  = "<big>"
	SocketKeyS0XT = "<S0XT>"
	SocketKeyBMP  = "<BMP>"
	SocketKeyQRC  = "<qrc>"

	SocketKeyPing      = "*AS01#" // 心跳包, 设备id,x   X=0/1   0表示正常。1表示缺纸
	SocketKeyOrderReceive = "*AS04#" // 打印机ID为XXXX对订单号为XXXX回复
	SocketKeyOrderAccept = "*AS05#" // 打印机ID为XXXX已接受订单号为XXXX的订单
	SocketKeyOrderReject = "*AS06#" // 打印机ID为XXXX已拒绝订单号为XXXX的订单
	SocketKeyOrderTimeout = "*AS07#" // 打印机ID为XXXX没有处理订单号为XXXX的订单
	SocketKeyNetworkTime = "*AS32#" // 查询网络延时状态: 根据回复速度来检查当时网络延时状态
	SocketKeyIMEI = "*AS33#" // 查询IMEI码
	SocketKeyNetworkSetting = "*AS34#" // 设置网络参数：设置打印机ID号、IP或域名、端口号
	SocketKeyPrintSetting = "*AS35#" // 五、设置打印份数、打印速度
	SocketKeyCheckVersion = "*AS36#" // 查询打印机版本
	SocketKeyNetworkSignal = "*AS37#" // 查询网络信号值
	SocketKeyUpgradeIPAndPort = "*AS40#" // 设置远程升级IP和端口号
	SocketKeyChain = "*AS47#" // 打印联号设置: 是否打印联号
	SocketKeyClearOrder = "*AS48#" // 清空订单数据


	SocketRegexDeviceId = "A\\*\\d+\\*" // 设备id


	// A*88888888*0*AS01#
	SocketFormatPingReply = "AS02#"               // 心跳包回复
	SocketFormatCheckVersion = "AS36?#"           // 查询打印机版本

	// &!*XXXX*XXXX*<big>XXXXXX*<S0XT>XXXX*<BMP>*<qrc>XXXX#
	SocketFormatOrderSend1 = "&!*%v%s*%s#" // 下发打印数据1
	// AS38*XXXX*X#
	SocketFormatOrderReply = "AS%s*%s*%s#" // 服务器对订单号为XXXX回复

)

const(
	PrinterMaxLen = 32 	// 打印机最大长度(英文字母)
	PrinterMaxLen4_2 = 14 	// 一行四文本,第一间隔
	PrinterMaxLen4_3 = 22 	// 一行四文本,第二间隔
)