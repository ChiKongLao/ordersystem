package constant

// key名
const (
	Name              = "name"
	NameNum           = "num"
	NamePic           = "pic"
	NamePrice         = "price"
	NameType          = "type"
	NamePayload       = "payload"
	NameToken         = "token"
	NamePlatForm      = "platform"
	NameIsOk          = "isOk"
	NameID            = "id"
	NameUserName      = "userName"
	NamePassword      = "password"
	NameNickName      = "nickName"
	NameRole          = "role"
	NameAuthorization = "authorization"
	NameData          = "data"
	NameMsg           = "msg"
	NameDesc          = "desc"
	NameSize          = "size"
	NameCount         = "count"
	NameTotalCount    = "totalCount"
	NameTotalPerson   = "totalPerson"
	NameTotalPrice    = "totalPrice"
	NameBusinessID    = "businessId"
	NameStatus        = "status"
	NamePerson        = "person"
	NamePersonNum     = "personNum"
	NameCapacity      = "capacity"
	NameTableName     = "tableName"
	NameTableId       = "tableId"
	NameFood        = "food"
	NameFoodId      = "foodId"
	NameOrderId       = "orderId"
	NameHead          = "head"
	NameIsCollect     = "isCollect"

	NameFrom    = "From"
	NameTo      = "To"
	NameCc      = "Cc"
	NameSubject = "Subject"
)

// 数据库的表的列名
const (
	ColumnUserName   = "user_name"
	ColumnBusinessId = "business_id"
	ColumnUserId     = "user_id"
	ColumnStatus     = "status"
	ColumnNum        = "num"
)

// 业务配置相关
const (
	RoleManager  = 1 // 管理员
	RoleBusiness = 2 // 商家
	RoleCustomer = 3 // 客户

	TableStatusEmpty     = 0 // 闲置
	TableStatusUsing     = 1 // 正在使用
	TableStatusWaitClean = 2 // 待清理
	TableStatusCleaning  = 3 // 清理中

	OrderStatusUnKnown = 0 // 未知
	OrderStatusWaitPay = 1 // 待付款
	OrderStatusPaid    = 2 // 已付款
	OrderStatusSure    = 3 // 已确认
	OrderStatusFinish  = 4 // 已完成

)