package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Footer struct {
	Text string `json:"text"`
}
type Image struct {
	Url string `json:"url"`
}
type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Image       `json:"image"`
	Footer      `json:"footer"`
}
type Send struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}
type Yuu struct {
	Data `json:"data"`
}
type Data struct {
	LiveRoom `json:"live_room"`
}
type LiveRoom struct {
	LiveStatus int    `json:"liveStatus"`
	Url        string `json:"url"`
	Title      string `json:"title"`
	Cover      string `json:"cover"`
}

func main() {
	resp, err := http.Get("https://api.bilibili.com/x/space/acc/info?mid=539700")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(data))
	unjson := Yuu{}
	err = json.Unmarshal(data, &unjson)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(unjson)
	if unjson.LiveStatus == 1 {
		dd(unjson.Title, "Yuu開播啦", unjson.Url, unjson.Cover)
	} else {

	}
}

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
	data, err := json.Marshal(jsend)
	if err != nil {
		log.Println(err)
		return
	}
	_, _ = http.Post(
		"https://discord.com/api/webhooks/850744832548274186/_dcZRdFA2nuF41BtVKCu9Qb9-L4zr3rxwZFXL_VT0SelX3io-6LZCZ13Z5i_8yeZTG7i",
		"application/json",
		bytes.NewReader(data))
}
