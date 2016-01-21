package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func ServerTest(command string, scanner *bufio.Scanner, conn net.Conn, name string, exptime int, contents string, version int64) (output string, vers int64, content string) {

	//var output string
	//var vers int64

	switch command {
	case "write":
		// Write a file
		fmt.Fprintf(conn, "write %v %v %v\r\n%v\r\n", name, len(contents), exptime, contents)
		scanner.Scan()                  // read first line
		resp := scanner.Text()          // extract the text from the buffer
		arr := strings.Split(resp, " ") // split into OK and <version>
		//	expect(t, arr[0], "OK")
		version, err := strconv.ParseInt(arr[1], 10, 64) // parse version as number
		if err != nil {
			
			output = "Non-numeric version found"
			vers = 0
			content = ""
		} else {
			output = arr[0]
			vers = version
			content = ""
		}
	case "read":
		fmt.Fprintf(conn, "read %v\r\n", name) // try a read now
		scanner.Scan()
		output = scanner.Text()
		arr := strings.Split(scanner.Text(), " ")
		if arr[0] == "CONTENTS" {
			scanner.Scan()
			content = scanner.Text()
			vers, _ = strconv.ParseInt(arr[1], 10, 64)
		} else {
			content = ""
			vers = 0
		}
	case "delete":
		fmt.Fprintf(conn, "delete %v\r\n", name)
		scanner.Scan()
		arr := strings.Split(scanner.Text(), " ")
		output = arr[0]
		vers = 0
		content = ""
	case "cas":
		fmt.Fprintf(conn, "cas %v %v %v %v\r\n%v\r\n", name, version, len(contents), exptime, contents)
		scanner.Scan()                    // read first line
		resp1 := scanner.Text()           // extract the text from the buffer
		arr1 := strings.Split(resp1, " ") // split into OK and <version>
		
		output = arr1[0]
		version, err1 := strconv.ParseInt(arr1[1], 10, 64) // parse version as number
		if err1 != nil {
			
			output = "Non-numeric version found"
			vers = 0
			content = ""
		} else {
			output = arr1[0]
			vers = version
			content = ""
		}

		//expect(t, arr[0], "OK")
	}
	return output, vers, content
}

// Simple serial check of getting and setting
func TestTCPSimple(t *testing.T) {
	go serverMain()
	time.Sleep(1 * time.Second) // one second is enough time for the server to start
	name := "hi.txt"
	contents := "bye"
	exptime := 300000
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err.Error()) // report error through testing framework
	}

	scanner := bufio.NewScanner(conn)

	output, version, _ := ServerTest("write", scanner, conn, name, exptime, contents, 0)
	if version == 0 {
		expect(t, output, "Non-numeric version found")
	} else {
		expect(t, output, "OK")
	}

	output1, _, content1 := ServerTest("read", scanner, conn, name, exptime, contents, version)
	arr := strings.Split(output1, " ")
	expect(t, arr[0], "CONTENTS")
	expect(t, arr[1], fmt.Sprintf("%v", version)) // expect only accepts strings, convert int version to string
	expect(t, arr[2], fmt.Sprintf("%v", len(contents)))
	expect(t, content1, contents)

	output, version2, _ := ServerTest("cas", scanner, conn, name, exptime, contents, version)
	if version2 == 0 {
		expect(t, output, "Non-numeric version found")
	} else {
		expect(t, output, "OK")
	}

	output2, _, _ := ServerTest("delete", scanner, conn, name, exptime, contents, 0)
	expect(t, output2, "OK")

}

// Useful testing function
func expect(t *testing.T, a string, b string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}
