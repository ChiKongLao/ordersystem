package model

// 订单
type Order struct {
	Id         int     `json:"id" xorm:"not null pk autoincr unique INT(11)"`
	TableId    int     `json:"-" xorm:"not null  INT(11)"`
	PersonNum  int     `json:"personNum" xorm:"not null INT(11)"`
	Price      float32 `json:"price" xorm:"not null FLOAT"`
	Status     int     `json:"status" xorm:"INT(11)"`
	Time       string  `json:"time" xorm:"not null VARCHAR(25)"`
	BusinessId int     `json:"-" xorm:"not null index INT(11)"`
	UserId     int     `json:"userId" xorm:"not null index INT(11)"` // 下单的用户id
	FoodList   []Food  `json:"foodList" xorm:"not null"`             // 菜单
}

type OrderResponse struct {
	Order
	TableName string `json:"tableName"`
}

type OrderListResponse struct {
	List      []OrderResponse `json:"list"`
	TotalPerson int `json:"totalPerson"`
	TotalPrice	float32 `json:"totalPrice"`
}

// 转化成回调的数据格式
func ConversionOrderResponseData(list []OrderResponse) *OrderListResponse {

	var personCount int
	var priceCount float32
	for _, subItem := range list {
		priceCount += subItem.Price
		personCount += subItem.PersonNum
	}

	return &OrderListResponse{
		List:list,
		TotalPerson:personCount,
		TotalPrice:priceCount,
	}
}
