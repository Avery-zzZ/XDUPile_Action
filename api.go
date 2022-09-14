package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	TASK_OK = iota
	TASK_ALREADT_DONE
	TASK_NOT_LOGGED_IN
	TASK_UNKNOWN_RESPONSE
	TASK_CONNECTION_ERROR
)

var client http.Client

type CookieJar struct{
    CookieSlice []*http.Cookie
}
func (jar *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
    jar.CookieSlice = append(jar.CookieSlice, cookies...)
}
func (jar *CookieJar) Cookies(u *url.URL) []*http.Cookie {
    return jar.CookieSlice
}

func login() (bool, error) {
	const (
		login_success_msg           = "{\"e\":0,\"m\":\"操作成功\",\"d\":{}}"
		login_fail_wrong_psw_msg    = "{\"e\":1,\"m\":\"账号或密码错误\",\"d\":{}}"
		login_fail_too_much_try_msg = "{\"e\":10016,\"m\":\"错误次数已达最大上限,请稍后再试\",\"d\":{}}"
	)

	body := []byte(
		"username=" + os.Getenv("XDU_ACCOUNT") + "&" +
			"password=" + os.Getenv("XDU_PASSWORD"),
	)
	request, err := http.NewRequest(
		http.MethodPost,
		"https://xxcapp.xidian.edu.cn/uc/wap/login/check",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return false, err
	}
	request.Header.Add("content-type", "application/x-www-form-urlencoded")

	client = http.Client{Jar: &CookieJar{}}
	response, err := client.Do(request)
	if err != nil {
		return false, errors.New("connection error")
	}

	//extract response body text
	content := make([]byte, 128)
	response.Body.Read(content)
	content = bytes.TrimRight(content, string(byte(0)))

	if bytes.Equal(content, []byte(login_success_msg)) {

		return true, nil

	} else if bytes.Equal(content, []byte(login_fail_wrong_psw_msg)) ||
		bytes.Equal(content, []byte(login_fail_too_much_try_msg)) {

		return false, errors.New("wrong password or account")

	} else {

		return false, errors.New("unknow response")

	}
}

func post(task int) (bool, int, string) {
	const (
		task_success_msg         = "{\"e\":0,\"m\":\"操作成功\",\"d\":{}}"
		task_not_logged_in_msg_1 = "{\"e\": 10013,\"m\": \"用户信息已失效,请重新进入页面\",\"d\": {\"login_url\": \"https://xxcapp.xidian.edu.cn/uc/wap/login?redirect=https%3A%2F%2Fxxcapp.xidian.edu.cn%2Fncov%2Fwap%2Fdefault%2Findex\"}}"
		task_not_logged_in_msg_2 = "<!DOCTYPE html>\n<html lang=\"zh-CN\">\n    <head>\n        <meta name=\"description\" content=\"\">\n        <meta name=\"keywords\" conten"
		task_already_done_msg    = "{\"e\":1,\"m\":\"今天已经填报了\",\"d\":{}}"
	)

	var body []byte
	var url string
	switch task {
	case 0:
		body = []byte(os.Getenv("DAILY_POST_BODY"))
		url = "https://xxcapp.xidian.edu.cn/xisuncov/wap/open-report/save"
	case 1:
		body = []byte(os.Getenv("POST_BODY"))
		url = "https://xxcapp.xidian.edu.cn/ncov/wap/default/save"
	}

	request, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Println("Wrong post format. Exiting")
		os.Exit(256)
	}
	request.Header.Add("content-type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return false, TASK_CONNECTION_ERROR, ""
	}

	//extract response body text
	content := make([]byte, 128)
	response.Body.Read(content)
	content = bytes.TrimRight(content, string(byte(0)))

	if bytes.Equal(content, []byte(task_success_msg)) {

		return true, TASK_OK, ""

	} else if bytes.Equal(content, []byte(task_already_done_msg)) {

		return true, TASK_ALREADT_DONE, ""

	} else if bytes.Equal(content, []byte(task_not_logged_in_msg_1)) ||
		bytes.Equal(content, []byte(task_not_logged_in_msg_2)) {

		return false, TASK_NOT_LOGGED_IN, ""

	} else {

		return false, TASK_UNKNOWN_RESPONSE, string(content)

	}
}
