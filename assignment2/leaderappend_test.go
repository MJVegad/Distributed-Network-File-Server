package main

import (
	//"fmt"
	"testing"
)

func TestLeaderAppend (t *testing.T) {
	//excpectedPeerIds := []uint64{2,3,4,5}
	sm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, currentTerm: 2, votedFor: 1, currentState: "leader", totalvotes: 3}
	result := sm.ProcessEvent(AppendEv{data: []byte("del")})
	//exsm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, currentTerm: 3, votedFor: 1, currentState: "candidate", totalvotes: 1}		
	exactions := []interface{}{LogStore{index: 2, command: logEntry{term: 2, command: []byte("del")}}, Send{peerId: 2, ev: AppendEntriesReqEv{term: 2, leaderId: 1, prevLogIndex: 1, prevLogTerm: 2, entries: []logEntry{logEntry{term: 2, command: []byte("del")}} , commitIndex: 1}}, Send{peerId: 3, ev: AppendEntriesReqEv{term: 2, leaderId: 1, prevLogIndex: 1, prevLogTerm: 2, entries: []logEntry{logEntry{term: 2, command: []byte("del")}} , commitIndex: 1}}, Send{peerId: 4, ev: AppendEntriesReqEv{term: 2, leaderId: 1, prevLogIndex: 1, prevLogTerm: 2, entries: []logEntry{logEntry{term: 2, command: []byte("del")}} , commitIndex: 1}}, Send{peerId: 5, ev: AppendEntriesReqEv{term: 2, leaderId: 1, prevLogIndex: 1, prevLogTerm: 2, entries: []logEntry{logEntry{term: 2, command: []byte("del")}} , commitIndex: 1}}}
	//expect1 (t, result, "candidate", 3, 1, 100)
	//ExpectStateMachine(t, &sm, &exsm)
	ExpectActions (t, result, exactions)
	
} 





