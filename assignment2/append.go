package main

import (
	"errors"
)

type AppendEv struct {
		data []byte
}

func (sm *StateMachine) AppendEventHandler ( event interface{} ) (actions []interface{}) {
	cmd := event.(AppendEv)
	///fmt.Printf("%v\n", cmd)
	switch sm.currentState {
		case "leader":
			sm.log = append(sm.log, logEntry{term: sm.currentTerm, command: cmd.data})
			actions = append(actions, LogStore{index: uint64(len(sm.log)-1), command: sm.log[uint64(len(sm.log)-1)]})
			/*for _, pid := range sm.peerIds {
				actions = append(actions, Send{peerId: pid, ev: AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.serverId, prevLogIndex: uint64(len(sm.log)-2), prevLogTerm: sm.log[uint64(len(sm.log)-2)].term, entries: sm.log[uint64(len(sm.log)-1):] , commitIndex: sm.commitIndex}})
			}*/
			for i:=0;i<len(sm.peerIds);i++ {
				if sm.serverId != sm.peerIds[i] {
					actions = append(actions, Send{peerId: sm.peerIds[i], ev: AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.serverId, prevLogIndex: sm.nextIndex[i]-1, prevLogTerm: sm.log[sm.nextIndex[i]-1].term, entries: sm.log[sm.nextIndex[i]:] , commitIndex: sm.commitIndex}})
				}
			}
		case "follower":
			actions = append(actions, Commit{index: uint64(len(sm.log)), command: cmd.data, err: errors.New("It's a follower, Not a leader")})
		case "candidate":
			actions = append(actions, Commit{index: uint64(len(sm.log)), command: cmd.data, err: errors.New("It's a candidate, Not a leader")})
		default: println("Invalid state")	
	}	
	return actions
}

