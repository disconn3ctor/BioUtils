package utils

import (
	"net"
	"net/http"
)

func GetAllFormRequestValue(r *http.Request) map[string]interface{} {
	clearMapData := make(map[string]interface{})


	// chon r.Form tamame maghadir ro to array mirikht on maghidir ro az array dar avrodam
	for i, value := range r.Form {

		clearMapData[i] = value[0]

	}

	return clearMapData
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}