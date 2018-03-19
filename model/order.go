package model

// 订单
type Order struct {
	Id         int     `json:"id" xorm:"not null pk autoincr unique INT(11)"`
	OrderNo    string  `json:"orderNo" xorm:"unique VARCHAR(35)"` // 订单号
	TableId    int     `json:"tableId" xorm:"not null  INT(11)"`
	PersonNum  int     `json:"personNum" xorm:"not null INT(11)"`
	Price      float32 `json:"price" xorm:"not null FLOAT"`
	Status     int     `json:"status" xorm:"INT(11)"`
	CreateTime int64   `json:"createTime" xorm:"not null BIGINT(20)"`
	UpdateTime int64   `json:"-" xorm:"not null BIGINT(20)"`
	BusinessId int     `json:"-" xorm:"not null index INT(11)"`
	UserId     int     `json:"userId" xorm:"not null index INT(11)"` // 下单的用户id
	FoodList   []Food  `json:"list" xorm:"not null"`                 // 菜单
}

type OrderResponse struct {
	Order            `xorm:"extends"`
	TableName string `json:"tableName"`
	FoodCount int    `json:"count"`
}

// 打印用
type OrderPrint struct {
	OrderResponse `xorm:"extends"`
	Customer User `json:"customer"`
	Business User `json:"business"`
}

type OrderListResponse struct {
	List        []OrderResponse `json:"list"`
	TotalPerson int             `json:"totalPerson"`
	TotalPrice  float32         `json:"totalPrice"`
	TotalCount  int             `json:"totalCount"`
}

// 转化成回调的数据格式
func ConvertOrderResponseData(list []OrderResponse) *OrderListResponse {

	var personCount int
	var priceCount float32
	for _, subItem := range list {
		priceCount += subItem.Price
		personCount += subItem.PersonNum
	}

	return &OrderListResponse{
		List:        list,
		TotalPerson: personCount,
		TotalPrice:  priceCount,
		TotalCount:  len(list),
	}
}

func ConvertOrderResponseToOrder(response OrderResponse) *Order {
	return &Order{
		Id:         response.Id,
		OrderNo:    response.OrderNo,
		TableId:    response.TableId,
		PersonNum:  response.PersonNum,
		Price:      response.Price,
		Status:     response.Status,
		CreateTime: response.CreateTime,
		UpdateTime: response.UpdateTime,
		BusinessId: response.BusinessId,
		UserId:     response.UserId,
		FoodList:   response.FoodList,
	}
}
