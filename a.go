package main

import (
	"regexp"
	"log"
	"github.com/toolkits/net"
)

func main(){
	ip:= "54.246.76.34"
	ip_mx ,_:= regexp.Compile(`([1-9]|[1-9]\d|2[0-5]{2}.){3}[1-9]|[1-9]\d|2[0-5]{2}`)
	if ip_mx.MatchString(ip){
		if net.IsIntranet(ip){
			log.Print(false)
		}
		log.Print(true)
	}
	log.Print(false)
}