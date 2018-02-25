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
	FoodList   []Food  `json:"list" xorm:"not null"`                 // 菜单
}

type OrderResponse struct {
	Id         int     `json:"id"`
	TableId    int     `json:"-"`
	PersonNum  int     `json:"personNum"`
	Price      float32 `json:"price"`
	Status     int     `json:"status"`
	Time       string  `json:"time"`
	BusinessId int     `json:"-"`
	UserId     int     `json:"userId"` // 下单的用户id
	FoodList   []Food  `json:"list"`   // 菜单

	TableName string `json:"tableName"`
}

type OrderListResponse struct {
	List        []OrderResponse `json:"list"`
	TotalPerson int             `json:"totalPerson"`
	TotalPrice  float32         `json:"totalPrice"`
	TotalCount  int         `json:"totalCount"`
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
		TableId:    response.TableId,
		PersonNum:  response.PersonNum,
		Price:      response.Price,
		Status:     response.Status,
		Time:       response.Time,
		BusinessId: response.BusinessId,
		UserId:     response.UserId,
		FoodList:   response.FoodList,
	}
}
