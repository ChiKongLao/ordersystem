package network

// 打印机通讯api

const (

	SocketKeyPing = "A\\*%s\\*%v\\*AS01#"	// 心跳包, 设备id,x   X=0/1   0表示正常。1表示缺纸
	SocketKeyPingReply = "AS02#"		// 心跳包回复
)