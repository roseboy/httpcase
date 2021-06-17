package main

import (
	"fmt"
	"github.com/roseboy/httpcase/requests"
	"strings"
	"time"
)

var (
	cookie    = ""
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36"
	eid       = "PDGIU7CH7Z7SCP5MOF7QUSQIVWYOCRTNM3ONUP6TYKBWBHEWUXLSGEIO7L7LSGQMKYBNI7WQN7QAYXZXOQXNPB3G2I"
	fp        = "8cfe5ee000f9416009357f45f20f3b6f"
	skuId     = "100021367452"
)

func main() {

	http := requests.NewHttpSession().
		Header("Cookie", cookie).
		Header("User-Agent", userAgent).
		Header("Referer", "https://item.jd.com/"+skuId+".html")

SETP1:
	//获取抢购连接
	fmt.Println("获取抢购连接。。。")
	url := fmt.Sprintf("https://itemko.jd.com/itemShowBtn?callback=jQuery%d&from=pc&skuId=%s&_%d",
		time.Now().UnixNano()/1000%10000000, skuId, time.Now().UnixNano()/1000%10000000)
	txt, e := http.Get(url).Send().ReadToText()
	if e != nil {
		fmt.Println(e)
		goto SETP1
	}
	url = "https:" + txt[strings.Index(txt, "\"url\":")+7:strings.LastIndex(txt, "\"")]
	fmt.Println(url)

}
