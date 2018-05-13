package payment
import (
	"fmt"
	"testing"
)

func Test_makeUrl(t *testing.T) {

	oreder := WechatOrder{OrderID: "123212323222", ProductName: "I am test", PriceTotal: 1, ProductID: 10001, IP: "127.0.0.1"}
	url, err := pay.GenderPayUrl(oreder)
	fmt.Println(url, err)
}
