package main

import (
	"testing"
)

func TestFollowerAppendEntriesRequest (t *testing.T) {
	//excpectedPeerIds := []uint64{2,3,4,5}
	sm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, 
		majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)},
		matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, 
		command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, currentTerm: 2, votedFor: 1, 
		currentState: "follower"}
    result := sm.ProcessEvent(AppendEntriesReqEv{term: 3, leaderId: 2, prevLogIndex: 1, prevLogTerm: 2, 
		entries: []logEntry{logEntry{term: 3, command: []byte("del")}}, commitIndex: 1})
	exsm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, 
		majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, 
		matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, 
		command: []byte("add")},logEntry{term: 2, command: []byte("disp")}, logEntry{term: 3, command: []byte("del")}}, currentTerm: 3, 
		votedFor: 0, currentState: "follower"}		
	exactions := []interface{}{Alarm{t: 100}, StateStore{state: "follower", term: 3, votedFor:0}, 
		LogStore{index: uint64(2), command: logEntry{term: 3, command: []byte("del")}}, Send{peerId: 2, ev: AppendEntriesRespEv{from: 1, term: 3, success: true}}} 
	//	Send{peerId: 2, ev: VoteReqEv{term: 3, candidateId: 1, lastLogIndex: 1, lastLogTerm: 2}}, Send{peerId: 3, ev: VoteReqEv{term: 3, candidateId: 1, lastLogIndex: 1, lastLogTerm: 2}}, Send{peerId: 4, ev: VoteReqEv{term: 3, candidateId: 1, lastLogIndex: 1, lastLogTerm: 2}}, Send{peerId: 5, ev: VoteReqEv{term: 3, candidateId: 1, lastLogIndex: 1, lastLogTerm: 2}}}
	//expect1 (t, result, "candidate", 3, 1, 100)
	ExpectStateMachine(t, &sm, &exsm)
	ExpectActions (t, result, exactions)
	
} 

