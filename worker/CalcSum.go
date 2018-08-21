package worker

import (
	"github.com/didi/falcon-log-agent/common/g"
	"regexp"
	"github.com/toolkits/net"
	"strings"
	"time"
	"encoding/json"
	"log"
)

func CalcSumStart(){
	for {
		time.Sleep(60 * time.Second)
		go PostToUrl()
	}
}

func PostToUrl(){
	Sum_Lock.Lock()
	defer Sum_Lock.Unlock()
	body ,_ := json.Marshal(g.Sum_map)
	log.Print(body)
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
