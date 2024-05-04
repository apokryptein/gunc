package main

// import packages
import (
	"bufio"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// Flag struct
type Flags struct {
	IP     string
	Port   int
	Listen bool
}

// Flag parsing function
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

// MAIN
func main() {
	// get and parse flags
	flags := GetFlags()

	// check for listen flag and start server or client accordingly
	if flags.Listen {
		// act as server
		fmt.Printf("Starting server: %s:%d\n", flags.IP, flags.Port)
		ServerListen(flags.IP, flags.Port)
	} else {
		// act as client
		ClientConnect(flags.IP, flags.Port)
	}
}

// Server Function
func ServerListen(ipAddr string, port int) {
	// construct address string
	fullAddr := makeAddress(ipAddr, port)

	// create listener
	listener, err := net.Listen("tcp", fullAddr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	//defer listener.Close()
	defer closeListener(listener)

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

// Function to handle server requests
func handleRequest(conn net.Conn) {
	// close conn
	defer conn.Close()

	for {
		// Get data length from client
		lengthBuf := make([]byte, 4)
		_, err := conn.Read(lengthBuf)
		if err != nil && err != io.EOF {
			if errors.Is(err, syscall.ECONNRESET) {
				// if client disconnects return
				return
			} else {
				log.Fatal(err)
				break
			}
		}

		length := binary.BigEndian.Uint32(lengthBuf)
		fmt.Printf("Receiving %d bytes of data\n", length)

		// make buffer of size specified in length header for read
		data := make([]byte, length)

		// read data from client
		_, err = conn.Read(data)
		if err != nil && err != io.EOF {
			if errors.Is(err, syscall.ECONNRESET) {
				// if client disconnects return
				return
			} else {
				log.Fatal(err)
				break
			}
		}

		// send message back to client
		conn.Write(lengthBuf)
		conn.Write(data)
	}
}

// Client Function
func ClientConnect(ipAddr string, port int) {
	// construct address string
	fullAddr := makeAddress(ipAddr, port)

	// connect to server
	conn, err := net.Dial("tcp", fullAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	//defer conn.Close()
	defer closeConn(conn)

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

		// test for exit or no input
		if input == "exit" {
			fmt.Println("Exiting. Goodbye.")
			return
		} else if input == "" {
			continue
		}

		// get input length and cast to uint32
		length := uint32(len(input))

		// create length buffer
		lengthBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(lengthBuf, length)

		// send data over socket
		conn.Write(lengthBuf)
		conn.Write([]byte(input))

		conn.Read(lengthBuf)
		respLength := binary.BigEndian.Uint32(lengthBuf)

		respData := make([]byte, respLength)
		conn.Read(respData)

		fmt.Println(string(respData))

	}
}

// Creates string address for net functions
func makeAddress(ip string, port int) string {
	p := strconv.Itoa(port)
	fullAddr := ip + ":" + p
	return fullAddr
}

// Closes net.Dial connections
func closeConn(conn net.Conn) {
	conn.Close()
	fmt.Println("Connection closed.")
}

// Closes net.Listen listeners
func closeListener(listener net.Listener) {
	listener.Close()
	fmt.Println("Listener closed.")
}
