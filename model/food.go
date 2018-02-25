package model

// 食物
type Food struct {
	Id          int      `json:"id" xorm:"not null pk autoincr unique INT"` // 食物id
	BusinessId  int      `json:"-" xorm:"not null index INT(11)"`           // 商家id
	Name        string   `json:"name" xorm:"not null VARCHAR(255)"`
	Num         int      `json:"num" xorm:"not null INT"`
	IsRecommend bool     `json:"isRecommend" xorm:"not null default 0 TINYINT(1)"` // 是否推荐
	IsCollect   bool     `json:"isCollect"`                                        // 是否收藏
	Pic         string   `json:"pic" xorm:"VARCHAR(255)"`
	Price       float32  `json:"price" xorm:"not null FLOAT"`
	ClassifyId  int      `json:"type" xorm:"INT"`       // 种类
	//Classify    Classify `json:"-" xorm:"VARCHAR(255)"` // 种类
	Desc        string   `json:"desc" xorm:"VARCHAR(255)"`
	//Desc       string `json:"desc,omitempty" xorm:"VARCHAR(255)"`
}

type FoodResponse struct {
	Food
	//Id          int      `json:"id"` // 食物id
	//Name        string   `json:"name"`
	//Num         int      `json:"num"`
	//IsRecommend bool     `json:"isRecommend"` // 是否推荐
	//IsCollect   bool     `json:"isCollect"`   // 是否收藏
	//Pic         string   `json:"pic"`
	//Desc        string   `json:"desc"`

	Classify    Classify `json:"type"`    // 种类
}

// 是否为同一道菜. 包含同样的味道
func (item *Food) IsSameFood(targetItem *Food) bool {
	return item.Id == targetItem.Id &&
		item.ClassifyId == targetItem.ClassifyId
}

// 获取food
func (item *FoodResponse) GetFood() *Food {
	return &Food{
		Id:item.Id,
		BusinessId:item.BusinessId,
		Name:item.Name,
		Num:item.Num,
		IsRecommend:item.IsRecommend,
		IsCollect:item.IsCollect,
		Pic:item.Pic,
		Price:item.Price,
		ClassifyId:item.Classify.Id,
		Desc:item.Desc,

	}
}

// 获取food
func ConvertFoodList(list []FoodResponse) []Food {
	resList := make([]Food,0)
	for _,subItem := range list {
		resList = append(resList,*(subItem.GetFood()))
	}
	return resList
}
