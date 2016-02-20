package main

import (
	"testing"
)

func TestLeaderTimeout (t *testing.T) {
	//excpectedPeerIds := []uint64{2,3,4,5}
	sm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, 
		majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, 
		matchIndex: []uint64{uint64(1),uint64(1),uint64(1),uint64(1)}, 
		log: []logEntry{logEntry{term: 1, command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, 
		currentTerm: 2, votedFor: 1, currentState: "leader", totalvotes: 3}
	exactions := []interface{}{Send{peerId: 2, ev: AppendEntriesReqEv{term: 2, leaderId: 1, prevLogIndex: 1, 
		prevLogTerm: 2, entries: []logEntry{}, commitIndex: 1}}, 
		Send{peerId: 3, ev: AppendEntriesReqEv{term: 2, leaderId: 1, prevLogIndex: 1, prevLogTerm: 2, 
		entries: []logEntry{}, commitIndex: 1}}, 
		Send{peerId: 4, ev: AppendEntriesReqEv{term: 2, leaderId: 1, prevLogIndex: 1, prevLogTerm: 2, 
		entries: []logEntry{}, commitIndex: 1}}, 
		Send{peerId: 5, ev: AppendEntriesReqEv{term: 2, leaderId: 1, prevLogIndex: 1, prevLogTerm: 2, 
		entries: []logEntry{}, commitIndex: 1}}}
	result := sm.ProcessEvent(TimeoutEv{})
	//expect (t, result, excpectedPeerIds)
	ExpectActions (t, result, exactions)
	
} 



