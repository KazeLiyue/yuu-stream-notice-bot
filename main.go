package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type (
	// Embed 建造Send內Embed的結構
	Embed struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		Image       struct {
			Url string `json:"url"`
		} `json:"image"`
		Footer struct {
			Text string `json:"text"`
		} `json:"footer"`
	}
	// Send 建立一個叫Send的結構
	Send struct {
		Content string  `json:"content"`
		Embeds  []Embed `json:"embeds"`
	}
	// Yuu 建立一個Yuu結構以獲取所需要的資訊
	Yuu struct {
		Data struct {
			LiveRoom struct {
				LiveStatus int    `json:"liveStatus"`
				Url        string `json:"url"`
				Title      string `json:"title"`
				Cover      string `json:"cover"`
			} `json:"live_room"`
		} `json:"data"`
	}
)

var streaming bool

//獲取直播間資訊的方法
func main() {
	for {
		job()
		time.Sleep(5 * time.Minute)
	}
}

func job() {
	//使用http包內的Get方法獲取資料 (mid=後的數字 可更改喜歡的直播主主頁網址後的數字)
	resp, err := http.Get("https://api.bilibili.com/x/space/acc/info?mid=539700")
	if err != nil {
		log.Println(err)
		return
	}
	//結束main方法後清理resp
	defer resp.Body.Close()
	//讀取resp.Body內的數據
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	//定義unjson為Yuu
	unjson := Yuu{}
	//使用json包內的Unmarshal解析unjson內的資料為struct
	err = json.Unmarshal(data, &unjson)
	if err != nil {
		log.Println(err)
		return
	}
	//判斷Yuu直播間是否已開播
	if unjson.Data.LiveRoom.LiveStatus == 1 {
		//streaming行為紀錄上一次的狀態
		if streaming == false {
			dd(unjson.Data.LiveRoom.Title, "Yuu開播啦", unjson.Data.LiveRoom.Url, unjson.Data.LiveRoom.Cover)
		}
		streaming = true
	} else {
		streaming = false
	}

}

//製作一個發送訊息的方法
func dd(title string, des string, url string, image string) {
	jsend := Send{
		Content: "Yuu直播通知!",
		Embeds: []Embed{
			{
				Title:       title,
				Description: des,
				Url:         url,
			},
		},
	}
	jsend.Embeds[0].Image.Url = image
	jsend.Embeds[0].Footer.Text = "快來看"
	//將string類型轉換成json格式
	data, err := json.Marshal(jsend)
	if err != nil {
		log.Println(err)
		return
	}
	//使用http包內的Post方法使Discord機器人發送消息 (網址可以成自己的Discord伺服器的webhook)
	_, _ = http.Post(
		"https://discord.com/api/webhooks/850744832548274186/_dcZRdFA2nuF41BtVKCu9Qb9-L4zr3rxwZFXL_VT0SelX3io-6LZCZ13Z5i_8yeZTG7i",
		"application/json",
		//使用bytes包內的NewReader方法處理Post方法返回的第三個值
		bytes.NewReader(data))
}
