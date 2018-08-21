package worker

import (
	"github.com/didi/falcon-log-agent/common/g"
	"github.com/didi/falcon-log-agent/common/dlog"
	"regexp"
	"github.com/toolkits/net"
	"strings"
	"time"
	"encoding/json"
	"bytes"
	"net/http"
	"io/ioutil"
)

type post_data struct {
	Tag	string `json:"tag"`
	Counter map[string]int `json:"counter"`
}

func CalcSumStart(){
	for {
		time.Sleep(60 * time.Second)
		go PostToUrl()
	}
}

func PostToUrl(){
	g.Sum_map_data.Lock()
	defer g.Sum_map_data.Unlock()
	bo ,err := json.Marshal(&post_data{
		Tag:g.Conf().CalcSum.Tag,
		Counter:g.Sum_map_data.Counter})
	g.InitSum()
	if err != nil {
		dlog.Error(err)
	}
	SendeMessage(bo)
	dlog.Debug(string(bo))
}

func CalcSum(line string){
	if ! g.Conf().CalcSum.Enable {
		return
	}
	line_arr := strings.Split(line,g.Conf().CalcSum.Delimiter)
	ip_str := line_arr[g.Conf().CalcSum.ArrLocation]
	if IsIP(ip_str){
		g.Sumadd(ip_str)
	}
}

func IsIP(ip string)bool{
	ip_mx ,_:= regexp.Compile(`([1-9]|[1-9]\d|2[0-5]{2}.){3}[1-9]|[1-9]\d|2[0-5]{2}`)
	if ip_mx.MatchString(ip){
		if net.IsIntranet(ip){
			return false
		}
		return true
	}
	return false
}

func SendeMessage(data []byte) {
	body := bytes.NewBuffer(data)
	res, err := http.Post(g.Conf().CalcSum.SumPushUrl, "application/json;charset=utf-8", body)
	if err != nil {
		dlog.Error("post data err ! ", err)
		return
	}
	if res.StatusCode != 200 {
		resp, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			dlog.Error(err)
			return
		}

		dlog.Debug(string(resp),res.StatusCode)
	}
}