package runners

import (
	"fmt"
	"io/ioutil"
	"net"

	"github.com/testdouble/http-assertion-tool/parsers"
)

type Serial struct {
	Address string
}

func (s *Serial) Run(spec *parsers.Spec) (*Result, error) {
	for _, test := range spec.Tests {
		fmt.Printf("Test request: %v", test.Request)
		fmt.Printf("Test response: %v", test.Response)

		addr, err := net.ResolveTCPAddr("tcp", s.Address)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return nil, err
		}

		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return nil, err
		}

		defer conn.Close()

		written, err := conn.Write([]byte(test.Request))
		if err != nil {
			fmt.Printf("Error: %v", err)
			return nil, err
		}

		conn.CloseWrite()

		fmt.Printf("Request:\n%s", test.Request)
		fmt.Printf("Written: %v\n", written)

		actualResponse, err := ioutil.ReadAll(conn)

		fmt.Printf("Response:\n%s", test.Response)
		fmt.Printf("Actual:\n%s", actualResponse)

		if string(actualResponse) == test.Response {
			fmt.Printf("Success.")
		} else {
			fmt.Printf("Does not match.")
		}
	}

	return nil, nil
}
