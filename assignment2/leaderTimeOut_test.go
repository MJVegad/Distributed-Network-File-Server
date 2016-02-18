package main

import (
	"fmt"
	"testing"
)

func TestLeaderTimeout (t *testing.T) {
	excpectedPeerIds := []uint64{2,3,4,5}
	sm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, currentTerm: 2, votedFor: 1, currentState: "leader", totalvotes: 3}
	result := sm.ProcessEvent(TimeoutEv{})
	expect (t, result, excpectedPeerIds)
	
} 

func expect(t *testing.T, a []interface{}, b []uint64) {
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


