package protocols

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/tiredsosha/admin/tools/logger"
)

func sendUdp(ip string, port int, data string) {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", ip+":"+strconv.Itoa(port))
	if err != nil {
		logger.Error.Println(err)
		return
	}
	defer conn.Close()
	if _, err = conn.Write([]byte(msg)); err != nil {
		fmt.Printf("Write err %v", err)
		os.Exit(-1)
	}

	if _, err = bufio.NewReader(conn).Read(p); err != nil {
		logger.Error.Println(err)
	}

	conn.Close()
	logger.Debug.Printf("responce - %q to post req from %q\n", string(body), url)

}
