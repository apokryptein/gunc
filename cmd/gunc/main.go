package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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

	// put flags into Flags struct
	flags := Flags{
		IP:     *ipFlag,
		Port:   *portFlag,
		Listen: *listenFlag,
	}
	return flags
}

func main() {
	flags := GetFlags()
	if flags.Listen == true {
		fmt.Printf("Starting server: %s:%d\n", flags.IP, flags.Port)
		ServerListen(flags.IP, flags.Port)
	} else {
		fmt.Printf("Connecting: %s:%d\n", flags.IP, flags.Port)
		ClientConnect(flags.IP, flags.Port)
	}

}

func ServerListen(ipAddr string, port int) {

	listener, err := net.Listen("tcp", makeAddress(ipAddr, port))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func ClientConnect(ipAddr string, port int) {

	tcpServer, err := net.ResolveTCPAddr("tcp", makeAddress(ipAddr, port))

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpServer)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer conn.Close()

	for {
		fmt.Print("# ")
		reader := bufio.NewReader(os.Stdin)
		// ReadString will block until the delimiter is entered
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// remove the delimeter from the string
		input = strings.TrimSuffix(input, "\n")

		_, err = conn.Write([]byte(input))
		if err != nil {
			fmt.Println("Write data failed:", err.Error())
			break
		}

		// buffer to get data
		received := make([]byte, 1024)
		fmt.Println("receiving now")
		n, err := conn.Read(received)
		fmt.Println("just after receiving")
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %s\n", err)
			}
			fmt.Println("some other kind of error")
			log.Fatal(err)
			break
		}
		fmt.Println("should have received from server")
		fmt.Println(string(received[:n]))
	}
}

func handleRequest(conn net.Conn) {
	// incoming request
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// write data to response
	time := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("Your message is: %v. Received time: %v", string(buffer[0:n]), time)
	conn.Write([]byte(responseStr))

	// close conn
	conn.Close()
}

func makeAddress(ip string, port int) string {
	p := strconv.Itoa(port)
	fullAddr := ip + ":" + p
	return fullAddr
}
