package main

import (
	"github.com/cs733-iitb/log"
	"os"
	"strconv"
	//"testing"
	//"time"
	"fmt"
	//"runtime"
	"encoding/json"
	"io/ioutil"
	"strings"
)


/*
func TestRaftNodeBasic(t *testing.T) {
	//runtime.GOMAXPROCS(1000)
	prepareRaftNodeConfigObj()
	rnArr := makeRafts()

	// get leader id from a stable system
	var ldrId int64
	for {
		time.Sleep(100 * time.Millisecond)
		ldrId = getLeader(rnArr)
		if ldrId != -1 {
			break
		}
	}
	// get leader raft node object using it's id
	ldr := getLeaderById(ldrId, rnArr)

	ldr.Append([]byte("foo"))
	time.Sleep(10 * time.Second)
	for _, rn := range rnArr {
		select {
		case ci := <-rn.CommitChannel():
			if ci.Err != nil {
				t.Fatal(ci.Err)
			}
			if string(ci.Data) != "foo" {
				fmt.Printf("Expected->foo, Got->%v . ", ci.Data)
				t.Fatal("Got different data")
			}
		default:
			t.Fatal("Expected message on all nodes")
		}
	}
	fmt.Println("single append testcase passed.\n")
	SystemShutdown(rnArr, nil)
}
*/
