package main

/*
import (
	//"github.com/cs733-iitb/log"
	//"os"
	"strconv"
	"testing"
	"time"
	"fmt"
	//"runtime"
)

func sendMsg(t *testing.T, rafts []RaftNode, cmd string, downNode []int) {
	leaderId := 0 // leader Id = 0 indicates election failure
	for leaderId == 0 {
		leaderId = getLeaderID(rafts, downNode)
	}
	time.Sleep(100 * time.Millisecond) //sleep since it may happen that vote respose didnt reach the candidate and voted for is set
	//	fmt.Printf("Leader id : %v\n", leaderId)
	for _, node2 := range rafts {
		if node2.sm.serverId == int64(leaderId) {
				leader := node2
				leader.Append([]byte(cmd))
				break
		}
	}

}

func Expect(t *testing.T, rafts []RaftNode, ExpectedData string, NodeDown []int) {
	for _, node := range rafts {
		nodeId := node.Id()
		val := IsPresentInNodeDown(NodeDown, nodeId)
		if val == true {
			fmt.Printf("%v Node Down.\n", nodeId)
			continue
		}
		select { // to avoid blocking on channel.
		case ci := <-node.CommitChannel():
		//case ci := <-node.commitch:
			fmt.Printf("Node->%v, data on commitch->%v\n", nodeId, string(ci.Data))
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

func SystemShutdown(rafts []RaftNode, NotToShutDown []int) {
	for _, node := range rafts {
		if IsPresentInNodeDown(NotToShutDown, node.Id()) {
			continue
		}
		node.Shutdown()
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

func BringNodeUp(i int, rafts []RaftNode) {
	//clusterconf := makeNetConfig(conf)
	prepareRaftNodeConfigObj()
	ld := "PersistentData_" + strconv.Itoa(i)
	eo := 2000 + 100*i
	rc := Config{cluster: peers, Id: int64(i), LogDir: ld, ElectionTimeout: int64(eo), HeartbeatTimeout: int64(500)}
	for j := 0; j < len(rafts); j++ {
		if rafts[j].sm.serverId == int64(i) {
			rafts[j] = New(rc, jsonFile)
			go rafts[j].processEvents()
			break
		}
	}

}

func TestMultipleAppends(t *testing.T) {
	//What if client send multiple appends to a leader. All should finally be committed by all nodes
	//conf := readJSONFile()
	//rafts := makeRafts(conf)
	prepareRaftNodeConfigObj()
	rafts := makeRafts()
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
	prepareRaftNodeConfigObj()
	rafts := makeRafts()
	//test 3 appends - read, write and cas
	time.Sleep(100 * time.Millisecond)
	sendMsg(t, rafts, "read", nil)
	time.Sleep(100 * time.Millisecond)
	// shutdown two followers. Issue write to leader
	leaderId := getLeaderID(rafts, nil)
	var foll []int
	cnt := 2
	i := 0
	for cnt != 0 && i < len(rafts) {
		if int64(leaderId) != rafts[i].sm.serverId {
			rafts[i].Shutdown()
			//fmt.Printf("%v Node shutdown.\n", rafts[i].sm.serverId)
			cnt--
			foll = append(foll, int(rafts[i].sm.serverId))
		}
		i++
	}
	sendMsg(t, rafts, "write", foll)
	sendMsg(t, rafts, "cas", foll)
	time.Sleep(30 * time.Second)
	Expect(t, rafts, "read", foll)
	Expect(t, rafts, "write", foll)
	Expect(t, rafts, "cas", foll)
	SystemShutdown(rafts, foll)
	fmt.Println("Pass : Minority Shutdown Test")
}

func TestPersistentAppend(t *testing.T) {
	//What if a follower node goes down and come back later. Its log should match with leader after some time.
	prepareRaftNodeConfigObj()
	rafts := makeRafts()
	//test 3 appends - read, write and cas
	time.Sleep(100 * time.Millisecond)
	sendMsg(t, rafts, "read", nil)
	time.Sleep(100 * time.Millisecond)
	// shutdown one followers. Issue write to leader
	leaderId := getLeaderID(rafts, nil)
	var downNode int
	cnt := 1
	i := 0
	for cnt != 0 && i < len(rafts) {
		if int64(leaderId) != rafts[i].sm.serverId {
			rafts[i].Shutdown()
			cnt--
			downNode = int(rafts[i].sm.serverId)
		}
		i++
	}
	dnList := []int{downNode}
	sendMsg(t, rafts, "write", dnList)
	sendMsg(t, rafts, "cas", dnList)
	// wait for sometime and wake up the follower
	time.Sleep(1000 * time.Millisecond)
	BringNodeUp(downNode, rafts)
	time.Sleep(30 * time.Second)
	Expect(t, rafts, "read", nil)
	Expect(t, rafts, "write", nil)
	Expect(t, rafts, "cas", nil)
	SystemShutdown(rafts, nil)
	fmt.Println("Pass : Test Persistent Append with follower Shutdown")
}


func TestLeaderShutdown(t *testing.T) {
	// What if a leader goes down. Re-election should be there and a new leader should be elected . If the previous leader come back again, it should be turned back to follower state and logs should be updated
	prepareRaftNodeConfigObj()
	rafts := makeRafts()
	//test 3 appends - read, write and cas
	time.Sleep(100 * time.Millisecond)
	sendMsg(t, rafts, "read", nil)
	time.Sleep(500 * time.Millisecond)

	// shutdown leader
	leaderId := getLeaderID(rafts, nil)
	for j := 0; j < len(rafts); j++ {
		if rafts[j].sm.serverId == int64(leaderId) {
				rafts[j].Shutdown()
				break;
		}
	}

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
	BringNodeUp(downNode, rafts)
	time.Sleep(40 * time.Second)
	Expect(t, rafts, "read", nil)
	Expect(t, rafts, "write", nil)
	Expect(t, rafts, "cas", nil)
	SystemShutdown(rafts, nil)
	fmt.Println("Pass : Test Leader ShutDown")
}

*/
