package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

var finalver int64

// Simple serial check of getting and setting
func TCPclient1(t *testing.T) {
	defer wg.Done()
	name := "demo"
	//contents := "bye"
	exptime := 300000
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err.Error()) // report error through testing framework
	}

	scanner := bufio.NewScanner(conn)
	ServerTest("write", scanner, conn, name, exptime, "11", 0)
	ServerTest("write", scanner, conn, name, exptime, "12", 0)
	ServerTest("write", scanner, conn, name, exptime, "13", 0)
}

func TCPclient2(t *testing.T) {
	defer wg.Done()
	name := "demo"
	//contents := "bye"
	exptime := 300000
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err.Error()) // report error through testing framework
	}

	scanner := bufio.NewScanner(conn)
	ServerTest("write", scanner, conn, name, exptime, "21", 0)
	ServerTest("write", scanner, conn, name, exptime, "22", 0)
	ServerTest("write", scanner, conn, name, exptime, "23", 0)
}

func TCPclient3(t *testing.T) {
	defer wg.Done()
	name := "demo"
	//contents := "bye"
	exptime := 300000
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err.Error()) // report error through testing framework
	}

	scanner := bufio.NewScanner(conn)
	ServerTest("write", scanner, conn, name, exptime, "31", 0)
	ServerTest("write", scanner, conn, name, exptime, "32", 0)
	ServerTest("write", scanner, conn, name, exptime, "33", 0)
}

var wg sync.WaitGroup

func TestBlindWrites(t *testing.T) {
	//go serverMain()
	time.Sleep(1 * time.Second)

	wg.Add(1)
	go TCPclient1(t)
	wg.Add(1)
	go TCPclient2(t)
	wg.Add(1)
	go TCPclient3(t)

	wg.Wait()

	name := "demo"
	contents := "bye"
	exptime := 300000
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err.Error()) // report error through testing framework
	}

	scanner := bufio.NewScanner(conn)

	_, _, content1 := ServerTest("read", scanner, conn, name, exptime, contents, 0)
	expectVer(t, content1)
}

// Useful testing function
func expectVer(t *testing.T, a string) {
	if a != "13" && a != "23" && a != "33" {
		t.Error(fmt.Sprintf("Expected 13/23/33, found %v", a)) // t.Error is visible when running `go test -verbose`
	}
}
