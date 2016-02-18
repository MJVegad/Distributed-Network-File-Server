package main

import (
	"fmt"
	"testing"
)

func TestFollowerTimeout (t *testing.T) {
	//excpectedPeerIds := []uint64{2,3,4,5}
	sm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, currentTerm: 2, votedFor: 1, currentState: "follower", totalvotes: 3}
	result := sm.ProcessEvent(TimeoutEv{})
	expect1 (t, result, "candidate", 3, 1, 100)
	
} 

func expect1(t *testing.T, a []interface{}, state string, term uint64, votedFor uint64, timeout uint64) {
	if (state != a[0].(StateStore).state || term != a[0].(StateStore).term || votedFor != a[0].(StateStore).votedFor) {
		t.Error(fmt.Sprintf("Expected state: %v, term: %v, votedFor: %v, found state: %v, term: %v, votedFor: %v ", state, term, votedFor, a[0].(StateStore).state, a[0].(StateStore).term, a[0].(StateStore).votedFor)) // t.Error is visible when running `go test -verbose`
	} 
	if (timeout != a[1].(Alarm).t) {
		t.Error(fmt.Sprintf("Expected timeout %v, found %v", timeout, a[1].(Alarm).t))
	}
	/*for i:=2;i<len(a);i++ {
			if a[i].(VoteReqEv) != act {
				t.Error(fmt.Sprintf("Expected VoteReq event, found %v event", a[i].())) // t.Error is visible when running `go test -verbose`
			}
		}*/
	
}



