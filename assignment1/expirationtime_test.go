package main

import (
	"bufio"
	"fmt"
	"net"
	//"sync"
	"testing"
	"time"
)

//var finalver int64

// Simple serial check of getting and setting
func TCPclient4(t *testing.T) {
	defer wg.Done()
	name := "demo"
	//contents := "bye"
	exptime := 3
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err.Error()) // report error through testing framework
	}

	scanner := bufio.NewScanner(conn)
	ServerTest("write", scanner, conn, name, exptime, "11", 0)
	time.Sleep(4 * time.Second)

	//output, _, _ = ServerTest("read", scanner, conn, name, exptime, contents, 0)
	//expectVer(t, output)

}

//var wg sync.WaitGroup

func TestExpirationTime(t *testing.T) {
	//go serverMain()
	time.Sleep(1 * time.Second)

	wg.Add(1)
	go TCPclient4(t)

	wg.Wait()

	name := "demo"
	contents := "bye"
	exptime := 300000
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err.Error()) // report error through testing framework
	}

	scanner := bufio.NewScanner(conn)

	output1, _, _ := ServerTest("read", scanner, conn, name, exptime, contents, 0)
	expectExpt(t, output1)
}

// Useful testing function
func expectExpt(t *testing.T, a string) {
	if a != "ERR_FILE_NOT_FOUND" {
		t.Error(fmt.Sprintf("Expected ERR_FILE_NOT_FOUND, found %v", a)) // t.Error is visible when running `go test -verbose`
	}
}
