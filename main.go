package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type (
	// Footer 建立Embed內的Footer結構
	Footer struct {
		Text string `json:"text"`
	}

	// Image 建立Embed內的Image結構
	Image struct {
		Url string `json:"url"`
	}

	// Embed 建造Send內Embed的結構
	Embed struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		Image       `json:"image"`
		Footer      `json:"footer"`
	}
	// Send 建立一個叫Send的結構
	Send struct {
		Content string  `json:"content"`
		Embeds  []Embed `json:"embeds"`
	}
	// Yuu 建立一個Yuu結構以獲取所需要的資訊
	Yuu struct {
		Data `json:"data"`
	}
	Data struct {
		LiveRoom `json:"live_room"`
	}
	LiveRoom struct {
		LiveStatus int    `json:"liveStatus"`
		Url        string `json:"url"`
		Title      string `json:"title"`
		Cover      string `json:"cover"`
	}
)

//獲取直播間資訊的方法
func main() {
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
	//使用josn包內的Unmarshal解析unjson內的資料為string
	err = json.Unmarshal(data, &unjson)
	if err != nil {
		log.Println(err)
		return
	}
	//判斷Yuu直播間是否已開播
	if unjson.LiveStatus == 1 {
		dd(unjson.Title, "Yuu開播啦", unjson.Url, unjson.Cover)
	} else {

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
				Image:       Image{Url: image},
				Footer:      Footer{Text: "快來看"},
			},
		},
	}
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
