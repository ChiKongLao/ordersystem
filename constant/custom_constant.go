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
	NameUserId        = "userId"
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
	NameFood          = "food"
	NameFoodId        = "foodId"
	NameOrderId       = "orderId"
	NameHead          = "head"
	NameIsCollect     = "isCollect"
	NameSort          = "sort"
)

// 数据库的表的列名
const (
	ColumnUserName   = "user_name"
	ColumnBusinessId = "business_id"
	ColumnUserId     = "user_id"
	ColumnTableId    = "table_id"
	ColumnStatus     = "status"
	ColumnNum        = "num"
	ColumnSort       = "sort"
)

// 业务配置相关
const (
	RoleManager  = 1 // 管理员
	RoleBusiness = 2 // 商家
	RoleCustomer = 3 // 客户

	TableStatusUnknown   = 0 // 未知
	TableStatusEmpty     = 1 // 闲置
	TableStatusOrdering  = 2 // 点餐中
	TableStatusUsing     = 3 // 正在使用
	TableStatusWaitClean = 4 // 待清理
	TableStatusCleaning  = 5 // 清理中

	OrderStatusUnknown = 0 // 未知
	OrderStatusWaitPay = 1 // 待付款
	OrderStatusPaid    = 2 // 已付款
	OrderStatusSure    = 3 // 已确认
	OrderStatusFinish  = 4 // 已完成

)
