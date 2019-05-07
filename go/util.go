package util

import (
	"net"
	"errors"
)

func GetLocalIp () (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	ip := ""
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	if 0 >= len(ip) {
		return "", errors.New("Get Ip Error!")
	} else {
		return  ip, nil
	}
}
