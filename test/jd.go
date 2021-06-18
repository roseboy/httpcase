package main

import (
	"fmt"
	"github.com/roseboy/httpcase/requests"
	"strings"
	"time"
)

var (
	cookie    = "__jdu=1594199934904129553323; shshshfpb=w60i8yNhbiRzfnD9e4t2I1Q%3D%3D; shshshfpa=2e514eec-1b75-2399-71bb-f5ad62686fad-1595305577; __jdv=76161171|direct|-|none|-|1623981088689; areaId=19; ipLoc-djd=19-1607-3155-0; o2State={%22webp%22:true%2C%22lastvisit%22:1623981089637}; shshshfp=2fa2108ceb071a9a0bb3730bf4ff3dad; TrackID=1TPBhJqrRCx8DDcXDdzi9v3CWSqlA-uYz6O_Q8SIq2vhzHHqIr7I4-qlKfTR8l1mb-6C13ppdPlGDS8w92JOZZDbwQ-WwYS0IBmwRq946QZEuFCYLxBLB-zvlBnLOs7Ld; thor=1F36F995536EB4D0126FC888D79637A24C7BB2AA57CC1024FA8940BE0A25CCAEE55210A2FAFCCA41ADC06F83245AA8BA0362146508962243B20EF105672E968D4E69FFA4EF151628C36090CA422E114F662B6C462D4A9588B90F0D65D2FD44336A4C8447B3321955BE8EE61E7E3513FBBC6551EB9E27F1A2A6EB2A9F4E8ECF250E2C43297FF9159BA9976586CDDDF84C; pinId=E_jjK1kgjIyFa5WGA0VXOg; pin=18315911332_p; unick=%E5%BC%BA%E4%B8%9C-; ceshi3.com=201; _tp=MVaPqdMeZa8kKA385TCuOA%3D%3D; _pst=18315911332_p; __jda=76161171.1594199934904129553323.1594199935.1621221966.1623981089.12; __jdb=76161171.4.1594199934904129553323|12.1623981089; __jdc=76161171; shshshsID=f540fe0e674351355a4ac81b348fbe0d_2_1623981108583; 3AB9D23F7A4B3C9B=GXFRQGEMMM325SAXB7RUTUJ3GWHRFHMCNR7IYQTA6CCNCN5WKLJPCFK6WNF5VB7BDKZJRTMAX2TYPEMM7T5QVNWUGM"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36"
	eid       = "PDGIU7CH7Z7SCP5MOF7QUSQIVWYOCRTNM3ONUP6TYKBWBHEWUXLSGEIO7L7LSGQMKYBNI7WQN7QAYXZXOQXNPB3G2I"
	fp        = "8cfe5ee000f9416009357f45f20f3b6f"
	skuId     = "100021367452"
)

func main() {

	http := requests.NewHttpSession().
		//Header("Cookie", cookie).
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
