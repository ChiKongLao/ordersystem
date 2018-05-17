package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"gopkg.in/chanxuehong/wechat.v2/oauth2"
	"github.com/chikong/ordersystem/constant"
)
var (
	sessionStorage                 = session.New(20*60, 60*60)
	oauth2Endpoint oauth2.Endpoint = mpoauth2.NewEndpoint(constant.WechatAppId,constant.WechatAppSecret )
)

func init() {
	//http.HandleFunc("/page1", Page1Handler)
	//http.HandleFunc("/page2", Page2Handler)

	Page1Handler()
}

// 建立必要的 session, 然后跳转到授权页面
func Page1Handler() {
	//sid := sid.New()
	state := string(rand.NewHex())

	//if err := sessionStorage.Add(sid, state); err != nil {
	//	log.Println(err)
	//	return
	//}

	AuthCodeURL := mpoauth2.AuthCodeURL(constant.WechatAppId, constant.WechatOauth2RedirectURI, constant.WechatOauth2Scope, state)
	log.Println("AuthCodeURL:", AuthCodeURL)

}

// 授权后回调页面
func Page2Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)

	cookie, err := r.Cookie("sid")
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	session, err := sessionStorage.Get(cookie.Value)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	savedState := session.(string) // 一般是要序列化的, 这里保存在内存所以可以这么做

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	code := queryValues.Get("code")
	if code == "" {
		log.Println("用户禁止授权")
		return
	}

	queryState := queryValues.Get("state")
	if queryState == "" {
		log.Println("state 参数为空")
		return
	}
	if savedState != queryState {
		str := fmt.Sprintf("state 不匹配, session 中的为 %q, url 传递过来的是 %q", savedState, queryState)
		io.WriteString(w, str)
		log.Println(str)
		return
	}

	oauth2Client := oauth2.Client{
		Endpoint: oauth2Endpoint,
	}
	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}
	log.Printf("token: %+v\r\n", token)

	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(userinfo)
	log.Printf("userinfo: %+v\r\n", userinfo)
	return
}

func main() {
	fmt.Println(http.ListenAndServe(":8080", nil))
}