package main

import (
	"testing"
)

func TestLeaderVoteRequest (t *testing.T) {
	//when candidate becomes leader
	sm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, 
		majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, 
		matchIndex: []uint64{uint64(1),uint64(1),uint64(1),uint64(1)}, 
		log: []logEntry{logEntry{term: 1, command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, 
		currentTerm: 2, votedFor: 1, currentState: "leader", totalvotes: 2}

	result := sm.ProcessEvent(VoteReqEv{term: 5, candidateId: 2, lastLogIndex: 3, lastLogTerm: 4})

	exsm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, 
		majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, 
		matchIndex: []uint64{uint64(1),uint64(1),uint64(1),uint64(1)}, 
		log: []logEntry{logEntry{term: 1, command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, 
		currentTerm: 5, votedFor: 2, currentState: "follower", totalvotes: 2}

	exactions := []interface{}{Alarm{t: 100}, Send {2, VoteRespEv {term: 5, voteGranted: true}},
		StateStore{state: "follower", term: 5, votedFor:2}}

	//expect (t, result, excpectedPeerIds)
	ExpectStateMachine(t, &sm, &exsm)	
	ExpectActions (t, result, exactions)
	
}

