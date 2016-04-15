package main


import (
	"github.com/cs733-iitb/log"
	"os"
	"strconv"
	"testing"
	"time"
	"fmt"
	"runtime"
)


// Number of replicated nodes
const totRaftNodes = 5
const jsonFile = "config.json"

var peers []NetConfig

func prepareRaftNodeConfigObj() {
	peers = []NetConfig{NetConfig{100, "localhost", 8001}, NetConfig{200, "localhost", 8002}, NetConfig{300, "localhost", 8003}, NetConfig{400, "localhost", 8004}, NetConfig{500, "localhost", 8005}}
}

func TestRaftNodeBasic(t *testing.T) {
	runtime.GOMAXPROCS(1000)
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
	fmt.Println("test case executed")
}

func makeRafts() []RaftNode {
	rnArr := make([]RaftNode, totRaftNodes)
	for i := 0; i < totRaftNodes; i++ {
		initRaftStateFile("PersistentData_" + strconv.Itoa((i+1)*100))
		rnArr[i] = New(Config{peers, int64((i + 1) * 100), "PersistentData_" + strconv.Itoa((i+1)*100), 4000, 500}, jsonFile)
		//fmt.Printf("sm_messaging = %v \n", rnArr[i].sm_messaging)
		go rnArr[i].processEvents()
	}
	return rnArr
}

func getLeader(rnArr []RaftNode) int64 {
	ldrId := int64(-1)
	mapIdToVotes := make(map[int64]int)
	maj := len(rnArr)/2 + 1
	for _, rn := range rnArr {
		if rn.sm.votedFor != 0 {
			mapIdToVotes[int64(rn.LeaderId())] += 1
		}
	}
	for k, v := range mapIdToVotes {
		if v >= maj {
			fmt.Printf("Leader Elected = %v \n", k)
			ldrId = k
			break
		}
	}
	fmt.Printf("getLeader: Leader id = %v \n", ldrId)
	return ldrId
}

func getLeaderById(ldrId int64, rnArr []RaftNode) *RaftNode {
	for index, rn := range rnArr {
		if int64(rn.Id()) == ldrId {
			return &rnArr[index]
		}
	}
	return nil
}

func initRaftStateFile(logDir string) {
	cleanup(logDir)
	stateAttrsFP, err := log.Open(logDir + "/" + "statefile")
	stateAttrsFP.RegisterSampleEntry(NodePers{})
	stateAttrsFP.SetCacheSize(1)
	assert(err == nil)
	defer stateAttrsFP.Close()
	stateAttrsFP.TruncateToEnd(0)
	err1 := stateAttrsFP.Append(NodePers{0, 0, "follower"})
	//fmt.Println(err1)
	assert(err1 == nil)
	fmt.Println("file created successfully")
}

func cleanup(logDir string) {
	os.RemoveAll(logDir)
}
