package model

// 食物
type Food struct {
	Id          int     `json:"id" xorm:"not null pk autoincr unique INT"` // 食物id
	BusinessId  int     `json:"-" xorm:"not null index INT(11)"`           // 商家id
	Name        string  `json:"name" xorm:"not null VARCHAR(255)"`
	Num         int     `json:"num" xorm:"not null INT"`
	IsRecommend bool    `json:"isRecommend" xorm:"not null default 0 TINYINT(1)"` // 是否推荐
	IsCollect   bool    `json:"isCollect"`                                        // 是否收藏
	Pic         string  `json:"pic" xorm:"VARCHAR(255)"`
	Price       float32 `json:"price" xorm:"not null FLOAT"`
	Type        string  `json:"type" xorm:"VARCHAR(255)"`       // 种类
	ClassifyId  string  `json:"classifyId" xorm:"VARCHAR(255)"` // 分类
	Desc        string  `json:"desc" xorm:"VARCHAR(255)"`
	//Desc       string `json:"desc,omitempty" xorm:"VARCHAR(255)"`
	SaleCount int `json:"saleCount" xorm:"INT"` // 月销量
}

type FoodResponse struct {
	Food
	//Classify Classify `json:"-"` // 种类
	SelectedCount int `json:"selectedCount"` // 购物车中已选择的个数
}

type FoodResponseSlice []FoodResponse
type FoodResponseMap map[string][]FoodResponse

// 是否为同一道菜. 包含同样的味道
func (item *Food) IsSameFood(targetItem Food) bool {
	return item.Id == targetItem.Id &&
		item.ClassifyId == targetItem.ClassifyId
}

// 获取food
func (item *FoodResponse) GetFood() *Food {
	return &item.Food
}

// 获取food
func ConvertFoodList(list []FoodResponse) []Food {
	resList := make([]Food, 0)
	for _, subItem := range list {
		resList = append(resList, *(subItem.GetFood()))
	}
	return resList
}

// 获取food
func ConvertFoodResponseList(list []Food) []FoodResponse {
	resList := make([]FoodResponse, 0)
	for _, subItem := range list {
		resList = append(resList, FoodResponse{Food: subItem, SelectedCount: 0,})
	}
	return resList
}

//func (list FoodResponseSlice) Len() int {
//	return len(list)
//}
//func (list FoodResponseSlice) Less(i, j int) bool {
//	return list[i].Classify.Sort < list[j].Classify.Sort
//}
//func (list FoodResponseSlice) Swap(i, j int) {
//	list[i], list[j] = list[j], list[i]
//}
//
//func (list FoodResponseMap) Len() int {
//	return len(list)
//}
//func (list FoodResponseMap) Less(i, j string) bool {
//	return list[i][0].Classify.Sort <list[j][0].Classify.Sort
//}
//func (list FoodResponseMap) Swap(i, j string) {
//	list[i], list[j] = list[j], list[i]
//}
