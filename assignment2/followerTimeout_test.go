package main

import (
	"testing"
)

func TestFollowerTimeout (t *testing.T) {
	//excpectedPeerIds := []uint64{2,3,4,5}
	sm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, currentTerm: 2, votedFor: 1, currentState: "follower", totalvotes: 3}
	result := sm.ProcessEvent(TimeoutEv{})
	exsm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, currentTerm: 3, votedFor: 1, currentState: "candidate", totalvotes: 1}		
	exactions := []interface{}{StateStore{state: "candidate", term: 3, votedFor:1}, Alarm{t: 100}, Send{peerId: 2, ev: VoteReqEv{term: 3, candidateId: 1, lastLogIndex: 1, lastLogTerm: 2}}, Send{peerId: 3, ev: VoteReqEv{term: 3, candidateId: 1, lastLogIndex: 1, lastLogTerm: 2}}, Send{peerId: 4, ev: VoteReqEv{term: 3, candidateId: 1, lastLogIndex: 1, lastLogTerm: 2}}, Send{peerId: 5, ev: VoteReqEv{term: 3, candidateId: 1, lastLogIndex: 1, lastLogTerm: 2}}}
	//expect1 (t, result, "candidate", 3, 1, 100)
	ExpectStateMachine(t, &sm, &exsm)
	ExpectActions (t, result, exactions)
	
} 





