package main

import (
	"testing"
)

func TestLeaderAppendEntriesResponse (t *testing.T) {
	sm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, 
		majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)},
		matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, 
		command: []byte("add")},logEntry{term: 2, command: []byte("disp")}, logEntry{term: 2, command: []byte("del")}}, currentTerm: 2, votedFor: 1, 
		currentState: "leader"}
    result := sm.ProcessEvent(AppendEntriesRespEv{from: 2, term: 2, success: true})
	exsm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, 
		majority: uint64(3), commitIndex: uint64(2), nextIndex: []uint64{uint64(3),uint64(2),uint64(2),uint64(2)}, 
		matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, 
		command: []byte("add")},logEntry{term: 2, command: []byte("disp")}, logEntry{term: 2, command: []byte("del")}}, currentTerm: 2, 
		votedFor: 1, currentState: "leader"}		
	exactions := []interface{}{Commit{index: 2, command: []byte("del"), err: nil}} 
	ExpectStateMachine(t, &sm, &exsm)
	ExpectActions (t, result, exactions)
	
} 

func TestFollowerAppendEntriesResponse (t *testing.T) {
	sm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, 
		majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)},
		matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, 
		command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, currentTerm: 2, votedFor: 1, 
		currentState: "follower"}
    result := sm.ProcessEvent(AppendEntriesRespEv{from: 2, term: 3, success: true})
	exsm := StateMachine {serverId: uint64(1), peerIds: []uint64{uint64(2),uint64(3),uint64(4),uint64(5)}, 
		majority: uint64(3), commitIndex: uint64(1), nextIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)},
		matchIndex: []uint64{uint64(2),uint64(2),uint64(2),uint64(2)}, log: []logEntry{logEntry{term: 1, 
		command: []byte("add")},logEntry{term: 2, command: []byte("disp")}}, currentTerm: 3, votedFor: 0, 
		currentState: "follower"}		
	exactions := []interface{}{StateStore{state: "follower", term: 3, votedFor:0}} 
	ExpectStateMachine(t, &sm, &exsm)
	ExpectActions (t, result, exactions)
	
} 

