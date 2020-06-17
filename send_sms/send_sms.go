package send_sms

import (
	"encoding/json"
	"fmt"
	"github.com/mogfee/common/xhttp"
	"net/http"
	"net/url"
)

type SendSMS struct {
	baseUrl  string
	userName string
	password string
}

var sms *SendSMS

func New(username, password string) {
	sms = &SendSMS{
		baseUrl:  "http://www.emailcar.net/sms/send",
		userName: username,
		password: password,
	}
}
func SendSms(mobile, template, cont string) error {
	header := http.Header{}
	header.Add("Content-type", "application/x-www-form-urlencoded")
	postData := url.Values{}
	postData.Add("api_user", sms.userName)
	postData.Add("api_pwd", sms.password)
	postData.Add("template_id", template)
	postData.Add("sms_template", cont)
	postData.Add("mobiles", mobile)
	body, err := xhttp.Request("POST", sms.baseUrl, postData.Encode(), header)
	if err != nil {
		return err
	}
	type Res struct {
		Msg    string `json:"msg"`
		Status string `json:"status"`
	}
	res := Res{}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return err
	}
	if res.Status == "success" {
		return nil
	}
	return fmt.Errorf("短信接口错误:%s", res.Msg)
}
