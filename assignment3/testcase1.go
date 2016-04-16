package main

import (
	"fmt"
	"testing"
	"reflect"
)

func TestLeaderTimeout (t *testing.T) {
	//excpectedPeerIds := []int64{2,3,4,5}
	sm := StateMachine {serverId: int64(1), peerIds: []int64{int64(2),int64(3),int64(4),int64(5)}, majority: int64(3), commitIndex: int64(1), nextIndex: []int64{int64(2),int64(2),int64(2),int64(2)}, matchIndex: []int64{int64(2),int64(2),int64(2),int64(2)}, log: []logEntry{logEntry{Term: 1, Command: []byte("add")},logEntry{Term: 2, Command: []byte("disp")}}, currentTerm: 2, votedFor: 1, currentState: "leader", totalvotes: 3}
	exactions := []interface{}{Send{peerId: 2, ev: AppendEntriesReqEv{Term: 2, LeaderId: 1, PrevLogIndex: 0, PrevLogTerm: 1, 
		Entries: nil, CommitIndex: 1}}, Send{peerId: 3, ev: AppendEntriesReqEv{Term: 2, LeaderId: 1, PrevLogIndex: 0, 
		PrevLogTerm: 1, Entries: nil, CommitIndex: 1}}, Send{peerId: 4, ev: AppendEntriesReqEv{Term: 2, LeaderId: 1, 
		PrevLogIndex: 0, PrevLogTerm: 1, Entries: nil, CommitIndex: 1}}, Send{peerId: 5, ev: AppendEntriesReqEv{Term: 2,
		LeaderId: 1, PrevLogIndex: 0, PrevLogTerm: 1, Entries: nil, CommitIndex: 1}}}
	result := sm.ProcessEvent(TimeoutEv{})
	//expect (t, result, excpectedPeerIds)
	ExpectActions (t, result, exactions)
	
} 

func expect(t *testing.T, a []interface{}, b []int64) {
	if (len(a) != len(b)) {
		t.Error(fmt.Sprintf("Expected %v elements, found %v elements", len(b), len(a))) // t.Error is visible when running `go test -verbose`
	} else {
		for i:=0;i<len(a);i++ {
			if a[i].(Send).peerId != b[i] {
				t.Error(fmt.Sprintf("Expected %v id, found %v id", b[i], a[i].(Send).peerId)) // t.Error is visible when running `go test -verbose`
			}
		}
	}
}

func ExpectActions(t *testing.T, a []interface{}, b []interface{}) {
	if (len(a) != len(b)) {
		t.Error(fmt.Sprintf("Expected %v elements, found %v elements", len(b), len(a))) // t.Error is visible when running `go test -verbose`
	} else {
		
		if ( len(CompareActions(a, b)) == 0 && len(CompareActions(b, a)) == 0 ) {
		e1 := 1
		e2 := 1
		i := 0
		for ;i<len(a);i++ {
			e1 = 0
			for j:=0;j<len(b);j++ {
				if reflect.DeepEqual(a[i], b[j]) {
					e1 = 1
					break;
				}
			}
			e2 = e2 & e1
			if e2 == 0 {
				break;
			}
		}
		if e2 == 0 {
			t.Fatal(fmt.Sprintf("%v not found in the expected list of actions", a[i]))
		}
		} else {
			t.Fatal(fmt.Sprintf("%v not found in the expected list of actions", 9))
		}
		
	}
}

func CompareActions(X, Y []interface{}) []interface{} {
    m := make(map[reflect.Type]int)

    for _, y := range Y {
        m[reflect.TypeOf(y)]++
    }

    var ret []interface{}
    for _, x := range X {
        if m[reflect.TypeOf(x)] > 0 {
            m[reflect.TypeOf(x)]--
            continue
        }
        ret = append(ret, x)
    }

    return ret
}
 //&& a.currentState==b.currentState && a.currentTerm==b.currentTerm && a.majority==b.majority && a.novotes==b.novotes && a.serverId==b.serverId && a.totalvotes==b.totalvotes && a.votedFor==b.votedFor && reflect.DeepEqual(a.log, b.log))
func ExpectStateMachine(t *testing.T, a *StateMachine, b *StateMachine) {
	if (a.commitIndex != b.commitIndex){
		t.Fatal(fmt.Sprintf("Expected %v commitIndex, found %v commitIndex", b.commitIndex, a.commitIndex)) // t.Error is visible when running `go test -verbose`
	} else if (a.currentTerm!=b.currentTerm) {
			t.Fatal(fmt.Sprintf("Expected %v currentTerm, found %v currentTerm", b.currentTerm, a.currentTerm)) // t.Error is visible when running `go test -verbose`
	} else if (a.majority!=b.majority) {
			t.Fatal(fmt.Sprintf("Expected %v majority, found %v majority", b.majority, a.majority)) // t.Error is visible when running `go test -verbose`
	} else if (a.serverId!=b.serverId) {
			t.Fatal(fmt.Sprintf("Expected %v serverId, found %v serverId", b.serverId, a.serverId)) // t.Error is visible when running `go test -verbose`
	} else if (a.novotes!=b.novotes) {
			t.Fatal(fmt.Sprintf("Expected %v novotes, found %v novotes", b.novotes, a.novotes)) // t.Error is visible when running `go test -verbose`
	} else if (a.totalvotes!=b.totalvotes) {
			t.Fatal(fmt.Sprintf("Expected %v totalvotes, found %v totalvotes", b.totalvotes, a.totalvotes)) // t.Error is visible when running `go test -verbose`
	} else if (a.votedFor!=b.votedFor) {
			t.Fatal(fmt.Sprintf("Expected %v votedFor, found %v votedFor", b.votedFor, a.votedFor)) // t.Error is visible when running `go test -verbose`
	} else if (reflect.DeepEqual(a.log, b.log)==false) {
			t.Fatal("Log not as expected.") // t.Error is visible when running `go test -verbose`
	} else if (reflect.DeepEqual(a.nextIndex, b.nextIndex)) {
			t.Fatal("nextIndex not as expected") // t.Error is visible when running `go test -verbose`
	} else if (reflect.DeepEqual(a.matchIndex, b.matchIndex)) {
			t.Fatal("matchIndex not as expected") // t.Error is visible when running `go test -verbose`
	} 
}


/*func TestLeaderAppend (t *testing.T) {
	//excpectedPeerIds := []int64{2,3,4,5}
	sm := StateMachine {serverId: int64(1), peerIds: []int64{int64(2),int64(3),int64(4),int64(5)}, majority: int64(3), commitIndex: int64(1), nextIndex: []int64{int64(2),int64(2),int64(2),int64(2)}, matchIndex: []int64{int64(2),int64(2),int64(2),int64(2)}, log: []logEntry{logEntry{Term: 1, command: []byte("add")},logEntry{Term: 2, command: []byte("disp")}}, currentTerm: 2, votedFor: 1, currentState: "leader", totalvotes: 3}
	result := sm.ProcessEvent(AppendEv{data: []byte("del")})
	//exsm := StateMachine {serverId: int64(1), peerIds: []int64{int64(2),int64(3),int64(4),int64(5)}, majority: int64(3), commitIndex: int64(1), nextIndex: []int64{int64(2),int64(2),int64(2),int64(2)}, matchIndex: []int64{int64(2),int64(2),int64(2),int64(2)}, log: []logEntry{logEntry{Term: 1, command: []byte("add")},logEntry{Term: 2, command: []byte("disp")}}, currentTerm: 3, votedFor: 1, currentState: "candidate", totalvotes: 1}		
	exactions := []interface{}{LogStore{index: 2, command: []byte("del")}, Send{peerId: 2, ev: AppendEntriesReqEv{Term: 2, leaderId: 2, prevLogIndex: 1, prevLogTerm: 2, entries: []logEntry{logEntry{Term: 2, command: []byte("del")}} , commitIndex: 1}}, Send{peerId: 3, ev: AppendEntriesReqEv{Term: 2, leaderId: 2, prevLogIndex: 1, prevLogTerm: 2, entries: []logEntry{logEntry{Term: 2, command: []byte("del")}} , commitIndex: 1}}, Send{peerId: 4, ev: AppendEntriesReqEv{Term: 2, leaderId: 2, prevLogIndex: 1, prevLogTerm: 2, entries: []logEntry{logEntry{Term: 2, command: []byte("del")}} , commitIndex: 1}}, Send{peerId: 5, ev: AppendEntriesReqEv{Term: 2, leaderId: 2, prevLogIndex: 1, prevLogTerm: 2, entries: []logEntry{logEntry{Term: 2, command: []byte("del")}} , commitIndex: 1}}}
	//expect1 (t, result, "candidate", 3, 1, 100)
	//ExpectStateMachine(t, &sm, &exsm)
	ExpectActions (t, result, exactions)
	
} */
