package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

type Args struct {
	Address  string
	Filename string
}

func loadArgs() (args Args) {
	flag.Parse()

	args.Filename = flag.Arg(0)
	args.Address = flag.Arg(1)

	return
}

func main() {
	args := loadArgs()

	bytes, err := ioutil.ReadFile(args.Filename)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	lines := strings.Split(string(bytes), "\n")
	request := new(strings.Builder)
	response := new(strings.Builder)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		} else if line[0] == '>' {
			if len(line) > 1 {
				request.WriteString(line[2:])
			} else {
				request.WriteString(line[1:])
			}
			request.WriteString("\r\n")
		} else if line[0] == '<' {
			if len(line) > 1 {
				response.WriteString(line[2:])
			} else {
				response.WriteString(line[1:])
			}
			response.WriteString("\r\n")
		}
	}

	addr, err := net.ResolveTCPAddr("tcp", args.Address)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	defer conn.Close()

	written, err := conn.Write([]byte(request.String()))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	conn.CloseWrite()

	fmt.Printf("Request:\n%s", request)
	fmt.Printf("Written: %v\n", written)

	actualResponse, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	conn.Close()

	fmt.Printf("Response:\n%s", response)
	fmt.Printf("Actual:\n%s", actualResponse)

	if string(actualResponse) == response.String() {
		fmt.Printf("Success.")
	} else {
		fmt.Printf("Does not match.")
	}
}
