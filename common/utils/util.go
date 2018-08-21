package utils

import (
	"os"
	"strings"

	"github.com/didi/falcon-log-agent/common/dlog"
	"io/ioutil"
	"log"
	"net"
)

func LocalHostname() (string, error) {
	ifaces, e := net.Interfaces()
	if e != nil {
		log.Print(e)
	}
	var ip_arr []string
	var ip_name []string
	var main_ip string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "vir") || strings.HasPrefix(iface.Name, "veth") || strings.HasPrefix(iface.Name, "tap") || strings.HasPrefix(iface.Name, "br") {
			continue
		}
		addrs, e := iface.Addrs()
		if e != nil {
			log.Print(e)
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ip_arr = append(ip_arr, iface.Name+":"+ip.String())
			ip_name = append(ip_name,iface.Name)
		}
	}
	if len(ip_arr) == 1 {
		main_ip = strings.Split(ip_arr[0], ":")[1]
	} else if len(ip_arr) > 1 {
		ip_device := strings.Split(ip_arr[0], ":")[0]
		for i := 1; i < len(ip_arr); i++ {
			if ip_device > strings.Split(ip_arr[i], ":")[0] {
				ip_device = strings.Split(ip_arr[i], ":")[0]
			}
		}
		allname := strings.Join(ip_name, "#")
		if strings.Count(allname, ip_device) == 1 {
			for i := 1; i < len(ip_arr); i++ {
				if strings.HasPrefix(ip_arr[i], ip_device) {
					main_ip = strings.Split(ip_arr[i], ":")[1]
				}
			}
		}else {
			net_file := "/etc/sysconfig/network-scripts/ifcfg-" + ip_device
			contents, err := ioutil.ReadFile(net_file)
			if err != nil {
				net_file = "/etc/rc.d/network.sh"
				contents, err = ioutil.ReadFile(net_file)
				for _, value := range strings.Split(string(contents), "\n") {
					if strings.HasPrefix(string(value), "AIP") {
						main_ip = strings.Trim(strings.Trim(strings.Split(string(value), "=")[1], "\""), "'")
					}
				}
			}
			for _, value := range strings.Split(string(contents), "\n") {
				if strings.Contains(string(value), "IPADDR") && len(strings.Split(string(value), "=")) == 2 {
					main_ip = strings.Trim(strings.Trim(strings.Split(string(value), "=")[1], "\""), "'")
				}
			}
		}
	}
	if main_ip == ""{
		log.Print("Got endpoint error ")
		os.Exit(1)
	}
	return main_ip, nil
}

//根据配置的时间格式，获取对应的正则匹配pattern和time包用的时间格式
func GetPatAndTimeFormat(tf string) (string, string) {
	var pat, timeFormat string
	switch tf {
	case "dd/mmm/yyyy:HH:MM:SS":
		pat = `([012][0-9]|3[01])/[JFMASOND][a-z]{2}/(2[0-9]{3}):([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "02/Jan/2006:15:04:05"
	case "dd/mmm/yyyy HH:MM:SS":
		pat = `([012][0-9]|3[01])/[JFMASOND][a-z]{2}/(2[0-9]{3})\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "02/Jan/2006 15:04:05"
	case "yyyy-mm-ddTHH:MM:SS":
		pat = `(2[0-9]{3})-(0[1-9]|1[012])-([012][0-9]|3[01])T([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "2006-01-02T15:04:05"
	case "dd-mmm-yyyy HH:MM:SS":
		pat = `([012][0-9]|3[01])-[JFMASOND][a-z]{2}-(2[0-9]{3})\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "02-Jan-2006 15:04:05"
	case "yyyy-mm-dd HH:MM:SS":
		pat = `(2[0-9]{3})-(0[1-9]|1[012])-([012][0-9]|3[01])\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "2006-01-02 15:04:05"
	case "yyyy/mm/dd HH:MM:SS":
		pat = `(2[0-9]{3})/(0[1-9]|1[012])/([012][0-9]|3[01])\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "2006/01/02 15:04:05"
	case "yyyymmdd HH:MM:SS":
		pat = `(2[0-9]{3})(0[1-9]|1[012])([012][0-9]|3[01])\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "20060102 15:04:05"
	case "mmm dd HH:MM:SS":
		pat = `[JFMASOND][a-z]{2}\s+([1-9]|[1-2][0-9]|3[01])\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "Jan 2 15:04:05"
	default:
		dlog.Errorf("match time pac failed : [timeFormat:%s]", tf)
		return "", ""
	}
	return pat, timeFormat
}
