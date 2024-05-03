package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Flags struct {
	IP     string
	Port   int
	Listen bool
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

	// check for minimum number of args
	if len(os.Args) < 3 {
		flag.Usage()
		os.Exit(0)
	}

	// put flags into Flags struct
	flags := Flags{
		IP:     *ipFlag,
		Port:   *portFlag,
		Listen: *listenFlag,
	}
	return flags
}

func main() {
	// get and parse flags
	flags := GetFlags()

	// check for listen flag and start server or client accordingly
	if flags.Listen == true {
		// act as server
		fmt.Printf("Starting server: %s:%d\n", flags.IP, flags.Port)
		ServerListen(flags.IP, flags.Port)
	} else {
		// act as client
		ClientConnect(flags.IP, flags.Port)
	}
}

func ServerListen(ipAddr string, port int) {
	// create listener
	listener, err := net.Listen("tcp", makeAddress(ipAddr, port))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer listener.Close()

	// loop to listen for requests
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleRequest(conn)
	}
}

func ClientConnect(ipAddr string, port int) {

	// construct address string
	fullAddr := makeAddress(ipAddr, port)

	// connect to server
	conn, err := net.Dial("tcp", fullAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer conn.Close()

	fmt.Printf("Connected to server => %s\n", fullAddr)
	fmt.Printf("'exit' to exit\n\n")

	consoleReader := bufio.NewReader(os.Stdin)

	for {
		// print prompt
		fmt.Print(">> ")

		// get input from user
		input, _ := consoleReader.ReadString('\n')

		// remove newline
		input = strings.TrimSuffix(input, "\n")

		if input == "exit" {
			fmt.Println("Exiting. Goodbye.")
			return
		} else if input == "" {
			continue
		}

		// send to server
		_, err = conn.Write([]byte(input))
		if err != nil {
			fmt.Println(err)
			return
		}

		// buffer to get data
		buf := make([]byte, 1024)

		// read from server
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		// print response to console
		fmt.Println(string(buf[:n]))
	}
}

func handleRequest(conn net.Conn) {
	// close conn
	defer conn.Close()

	// incoming request
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
			break
		}

		data := buf[:n]

		fmt.Printf("Received %d bytes\n", n)
		_, err = conn.Write(data)

		if err != nil {
			log.Fatal(err)
			break
		}
	}
}

func makeAddress(ip string, port int) string {
	p := strconv.Itoa(port)
	fullAddr := ip + ":" + p
	return fullAddr
}
