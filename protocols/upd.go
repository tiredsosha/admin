package protocols

import (
	"fmt"
	"net"
	"strconv"

	"github.com/tiredsosha/admin/tools/logger"
)

func SendUdp(ip string, port int, data string) {
	address := ip + ":" + strconv.Itoa(port)
	conn, err := net.Dial("udp", address)
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

func SendUDPBrights(ipPrefix string, from int, to int, port int, command string) {
	for i := from; i <= to; i++ {
		localAddr := "192.168.10.5:8016"
		targetAddr := fmt.Sprintf("%s.%d:%d", ipPrefix, i, port)

		laddr, err := net.ResolveUDPAddr("udp", localAddr)
		if err != nil {
			fmt.Println("local resolve error:", err)
			continue
		}

		raddr, err := net.ResolveUDPAddr("udp", targetAddr)
		if err != nil {
			fmt.Println("remote resolve error:", err)
			continue
		}

		conn, err := net.DialUDP("udp", laddr, raddr)
		if err != nil {
			fmt.Println("dial error:", err)
			continue
		}

		_, err = conn.Write([]byte(command))
		if err != nil {
			fmt.Println("write error:", err)
		}

		conn.Close()
	}
}
