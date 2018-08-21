package g

import "sync"

var Sum_map map[string]int
var Sum_Lock *sync.RWMutex

func InitSum(){
	Sum_Lock = new(sync.RWMutex)
	Sum_map = make(map[string]int,0)
}

func Sumadd(ip string){
	Sum_Lock.Lock()
	defer Sum_Lock.Unlock()
	if _,ok := Sum_map[ip];ok {
		Sum_map[ip] += 1
		return
	}
	Sum_map[ip] = 1
	return
}