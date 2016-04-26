package main

/*import (
	//	"errors"
	"encoding/json"
	"fmt"
	"github.com/cs733-iitb/log"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

type Peer struct {
	Id      int
	Address string
}

type ClusterConfig struct {
	Peers []Peer
}

func readJSONFile() ClusterConfig {
	content, err := ioutil.ReadFile("Config.json")
	if err != nil {
		fmt.Print("Error:", err)
	}
	var conf ClusterConfig
	err = json.Unmarshal(content, &conf)
	if err != nil {
		fmt.Print("Error:", err)
	}
	return conf
}

// This func creates a state file for storing initila state of each raft machine
func initialState(sd string) {
	stateFile := sd + "/" + "mystate"
	//	rmlog(stateFile)
	st, err := log.Open(stateFile)
	//st.SetCacheSize(50)
	assert(err == nil)
	st.RegisterSampleEntry(SMState{})
	defer st.Close()
	err = st.Append(SMState{State: "FOLLOWER", CurrentTerm: 0, VotedFor: 0})
	assert(err == nil)
}

func clearFiles(file string) {
	rmlog(file)
}

func makeNetConfig(conf ClusterConfig) []NetConfig {
	clusterconf := make([]NetConfig, len(conf.Peers))

	for j := 0; j < len(conf.Peers); j++ {
		clusterconf[j].Id = conf.Peers[j].Id
		clusterconf[j].Host = strings.Split(conf.Peers[j].Address, ":")[0]
		clusterconf[j].Port, _ = strconv.Atoi(strings.Split(conf.Peers[j].Address, ":")[1])
	}
	return clusterconf
}

func makeRafts(conf ClusterConfig) []RaftNode {
	clusterconf := makeNetConfig(conf)
	//	clusterconfig:=[]NetConfig{{Id:1, Host:"localhost", Port:8001},{Id:2, Host:"localhost", Port:8002},{Id:3, Host:"localhost", Port:8003},{Id:4, Host:"localhost", Port:8004},{Id:5, Host:"localhost", Port:8005}}
	rafts := make([]RaftNode, len(conf.Peers))
	var ld string
	var sd string
	for i := 1; i <= len(conf.Peers); i++ {

		ld = "myLogDir" + strconv.Itoa(i)
		clearFiles(ld + "/logfile")
		sd = "myStateDir" + strconv.Itoa(i)
		clearFiles(sd + "/mystate")
		initialState(sd)
		eo := 2000 + 100*i
		rc := RaftConfig{cluster: clusterconf, Id: i, LogDir: ld, StateDir: sd, ElectionTimeout: eo, HeartbeatTimeout: 500}
		rafts[i-1] = New(rc)
		go rafts[i-1].processEvents()
	}

	return rafts
}

func BringNodeUp(conf ClusterConfig, i int, rafts []RaftNode) {
	clusterconf := makeNetConfig(conf)
	ld := "myLogDir" + strconv.Itoa(i)
	sd := "myStateDir" + strconv.Itoa(i)
	eo := 2000 + 100*i
	rc := RaftConfig{cluster: clusterconf, Id: i, LogDir: ld, StateDir: sd, ElectionTimeout: eo, HeartbeatTimeout: 500}
	rafts[i-1] = New(rc)
	go rafts[i-1].processEvents()
}

func getLeaderID(rafts []RaftNode, NodeDown []int) int {
	//find the number of votes of each node. If any node has got majority of votes then it is leader. If majority has not voted for a single node then election failure. Return -1 in this case
	majority := (len(rafts) / 2) + 1
	votes := make(map[int]int)
	for i := 0; i < len(rafts); i++ {
		if IsPresentInNodeDown(NodeDown, rafts[i].Id()) {
			continue
		}
		votedFor := rafts[i].LeaderId()
		votes[votedFor] += 1
	}
	for key, value := range votes {
		if value >= majority {
			return key
		}
	}
	return 0
}

func sendMsg(t *testing.T, rafts []RaftNode, cmd string, downNode []int) {
	leaderId := 0 // leader Id = 0 indicates election failure
	for leaderId == 0 {
		leaderId = getLeaderID(rafts, downNode)
	}
	time.Sleep(100 * time.Millisecond) //sleep since it may happen that vote respose didnt reach the candidate and voted for is set
	//	fmt.Printf("Leader id : %v\n", leaderId)
	leader := rafts[leaderId-1]
	leader.Append([]byte(cmd))
}

func Expect(t *testing.T, rafts []RaftNode, ExpectedData string, NodeDown []int) {
	for _, node := range rafts {
		nodeId := node.Id()
		val := IsPresentInNodeDown(NodeDown, nodeId)
		if val == true {
			continue
		}
		select { // to avoid blocking on channel.
		case ci := <-node.CommitChannel():
			if ci.Err != nil {
				t.Fatal(ci.Err)
			}
			if string(ci.Data) != ExpectedData {
				t.Fatal("Got different data : %v, Expected: %v\n ", string(ci.Data), ExpectedData)
			}
		default:
			t.Fatal("Expected message on all nodes")
		}
		//		node.Shutdown()
	}
}

func IsPresentInNodeDown(NodeDown []int, NodeId int) bool {
	for _, i := range NodeDown {
		if i == NodeId {
			return true
		}
	}
	return false
}

func SystemShutdown(rafts []RaftNode, NotToShutDown []int) {
	for _, node := range rafts {
		if IsPresentInNodeDown(NotToShutDown, node.Id()) {
			continue
		}
		node.Shutdown()
	}
}

func TestBasic(t *testing.T) {
	//What if a client send a append to leader. It should finally be committed by all the nodes
	runtime.GOMAXPROCS(1010)
	conf := readJSONFile()
	rafts := makeRafts(conf)
	time.Sleep(100 * time.Millisecond)
	sendMsg(t, rafts, "read", nil)
	time.Sleep(10 * time.Second)
	Expect(t, rafts, "read", nil)
	SystemShutdown(rafts, nil)
	fmt.Println("Pass : Basic Test")
}

func TestMultipleAppends(t *testing.T) {
	//What if client send multiple appends to a leader. All should finally be committed by all nodes
	conf := readJSONFile()
	rafts := makeRafts(conf)
	//test 3 appends - read, write and cas
	time.Sleep(100 * time.Millisecond)
	sendMsg(t, rafts, "read", nil)
	sendMsg(t, rafts, "write", nil)
	sendMsg(t, rafts, "cas", nil)
	time.Sleep(20 * time.Second)
	Expect(t, rafts, "read", nil)
	Expect(t, rafts, "write", nil)
	Expect(t, rafts, "cas", nil)
	SystemShutdown(rafts, nil)
	fmt.Println("Pass : Multiple Appends Test")
}

func TestWithMinorityShutdown(t *testing.T) {
	// What if minority nodes goes down. System should still function correctly
	conf := readJSONFile()
	rafts := makeRafts(conf)
	//test 3 appends - read, write and cas
	time.Sleep(100 * time.Millisecond)
	sendMsg(t, rafts, "read", nil)
	time.Sleep(100 * time.Millisecond)
	// shutdown two followers. Issue write to leader
	leaderId := getLeaderID(rafts, nil)
	var foll []int
	cnt := 2
	i := 1
	for cnt != 0 && i <= len(rafts) {
		if leaderId != i {
			rafts[i-1].Shutdown()
			cnt--
			foll = append(foll, i)
		}
		i++
	}
	sendMsg(t, rafts, "write", foll)
	sendMsg(t, rafts, "cas", foll)
	time.Sleep(20 * time.Second)
	Expect(t, rafts, "read", foll)
	Expect(t, rafts, "write", foll)
	Expect(t, rafts, "cas", foll)
	SystemShutdown(rafts, foll)
	fmt.Println("Pass : Minority Shutdown Test")
}

func TestPersistentAppend(t *testing.T) {
	//What if a follower node goes down and come back later. Its log should match with leader after some time.
	conf := readJSONFile()
	rafts := makeRafts(conf)
	//test 3 appends - read, write and cas
	time.Sleep(100 * time.Millisecond)
	sendMsg(t, rafts, "read", nil)
	time.Sleep(100 * time.Millisecond)
	// shutdown one followers. Issue write to leader
	leaderId := getLeaderID(rafts, nil)
	var downNode int
	cnt := 1
	i := 1
	for cnt != 0 && i <= len(rafts) {
		if leaderId != i {
			rafts[i-1].Shutdown()
			cnt--
			downNode = i
		}
		i++
	}
	dnList := []int{downNode}
	sendMsg(t, rafts, "write", dnList)
	sendMsg(t, rafts, "cas", dnList)
	// wait for sometime and wake up the follower
	time.Sleep(1000 * time.Millisecond)
	BringNodeUp(conf, downNode, rafts)
	time.Sleep(30 * time.Second)
	Expect(t, rafts, "read", nil)
	Expect(t, rafts, "write", nil)
	Expect(t, rafts, "cas", nil)
	SystemShutdown(rafts, nil)
	fmt.Println("Pass : Test Persistent Append with follower Shutdown")
}

func TestLeaderShutdown(t *testing.T) {
	// What if a leader goes down. Re-election should be there and a new leader should be elected . If the previous leader come back again, it should be turned back to follower state and logs should be updated
	conf := readJSONFile()
	rafts := makeRafts(conf)
	//test 3 appends - read, write and cas
	time.Sleep(100 * time.Millisecond)
	sendMsg(t, rafts, "read", nil)
	time.Sleep(500 * time.Millisecond)

	// shutdown leader
	leaderId := getLeaderID(rafts, nil)
	rafts[leaderId-1].Shutdown()
	downNode := leaderId
	dnList := []int{downNode}
	// let the re-election occurs
	time.Sleep(1000 * time.Millisecond)
	leaderId = getLeaderID(rafts, dnList)
	for leaderId == 0 || leaderId == downNode {
		leaderId = getLeaderID(rafts, dnList)
	}
	sendMsg(t, rafts, "write", dnList)
	sendMsg(t, rafts, "cas", dnList)
	// wait for sometime and wake up the leader
	time.Sleep(100 * time.Millisecond)
	BringNodeUp(conf, downNode, rafts)
	time.Sleep(40 * time.Second)
	Expect(t, rafts, "read", nil)
	Expect(t, rafts, "write", nil)
	Expect(t, rafts, "cas", nil)
	SystemShutdown(rafts, nil)
	fmt.Println("Pass : Test Leader ShutDown")
}
*/
