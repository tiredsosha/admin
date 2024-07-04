package protocols

import (
	"net"
	"strconv"

	"github.com/tiredsosha/admin/tools/logger"
)

func SendTcp(ip string, port int, data string) {
	address := ip + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(data)); err != nil {
		logger.Error.Println(err)
	}
	logger.Debug.Printf("udp msg - %q to %q\n", data, address)

}
