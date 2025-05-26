package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	appid   = "APPID"
	apiname = "APINAME"
	appkey  = "APPKEY"
	url     = "http://127.0.0.1:7081/gateway/sms/send"
)

type msgbody struct {
	Phone   string `json:"phone"`
	Content string `json:"content"`
}

type ResponseBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Exec struct {
			ExecTime int `json:"ExecTime"`
			ConnTime int `json:"ConnTime"`
		} `json:"exec"`
		Result string `json:"result"`
	} `json:"data"`
}

func SendMsg(phone, content string) error {
	signature, err := GenerateSignature(appid, apiname, appkey)
	if err != nil {
		return err
	}
	log.Printf("sigature: %v", signature)
	body := &msgbody{
		Phone:   phone,
		Content: content,
	}
	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("signature", signature)
	req.Header.Set("apiname", apiname)
	req.Header.Set("appid", appid)

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	log.Printf("发送短信请求")
	if err != nil {
		return err
	}
	//defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//fmt.Println("Response Body: ", string(respbody))
	log.Printf("短信发送接口返回: %v", string(respbody))

	var result ResponseBody
	_ = json.Unmarshal(respbody, &result)
	if result.Code == 0 {
		log.Println("发送成功")
		//fmt.Println("发送成功")
		return nil
	} else {
		log.Printf("发送失败, Err: %v", result)
		//fmt.Println("发送失败")
		return fmt.Errorf("发送失败, Err: %v", result.Msg+result.Data.Result)
	}
}

