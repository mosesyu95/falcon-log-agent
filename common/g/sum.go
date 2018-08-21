package g

import "sync"

var Sum_map_data *Sum_map

type Sum_map struct {
	Counter map[string]int
	sync.RWMutex
}
func InitSum(){
	Sum_map_data = &Sum_map{
		Counter:make(map[string]int,0),
	}
}

func Sumadd(ip string){
	Sum_map_data.Lock()
	defer Sum_map_data.Unlock()
	if _,ok := Sum_map_data.Counter[ip];ok {
		Sum_map_data.Counter[ip] += 1
		return
	}
	Sum_map_data.Counter[ip] = 1
	return
}