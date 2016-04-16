package main

import (
	"fmt"
	"testing"
)

func TestFollowerTimeout (t *testing.T) {
	//excpectedPeerIds := []int64{2,3,4,5}
	sm := StateMachine {serverId: int64(1), peerIds: []int64{int64(2),int64(3),int64(4),int64(5)}, majority: int64(3), commitIndex: int64(1), nextIndex: []int64{int64(2),int64(2),int64(2),int64(2)}, matchIndex: []int64{int64(2),int64(2),int64(2),int64(2)}, log: []logEntry{logEntry{Term: 1, Command: []byte("add")},logEntry{Term: 2, Command: []byte("disp")}}, currentTerm: 2, votedFor: 1, currentState: "follower", totalvotes: 3}
	result := sm.ProcessEvent(TimeoutEv{})
	expect1 (t, result, "candidate", 3, 1, 100)
	
} 

func expect1(t *testing.T, a []interface{}, state string, Term int64, votedFor int64, timeout int64) {
	if (state != a[0].(StateStore).state || Term != a[0].(StateStore).term || votedFor != a[0].(StateStore).votedFor) {
		t.Error(fmt.Sprintf("Expected state: %v, Term: %v, votedFor: %v, found state: %v, Term: %v, votedFor: %v ", state, Term, votedFor, a[0].(StateStore).state, a[0].(StateStore).term, a[0].(StateStore).votedFor)) // t.Error is visible when running `go test -verbose`
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


