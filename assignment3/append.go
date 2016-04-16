package main

import (
//	"fmt"
	"errors"
)

type AppendEv struct {
		Data []byte
}

func (sm *StateMachine) AppendEventHandler ( event interface{} ) (actions []interface{}) {
	cmd := event.(AppendEv)
	//fmt.Printf("Command to append on leader=> %v\n", cmd)
	switch sm.currentState {
		case "leader":
			sm.log = append(sm.log, logEntry{Term: sm.currentTerm, command: cmd.Data})
			actions = append(actions, LogStore{index: int64(len(sm.log)-1), command: sm.log[int64(len(sm.log)-1)]})
			//fmt.Printf("leader->%v, log->%v\n", sm.serverId, sm.log)
			for i:=0;i<len(sm.peerIds);i++ {
				if sm.serverId != sm.peerIds[i] {
					if sm.nextIndex[i] != 0 {
					//fmt.Printf("leader->%v, log entries sent->%v\n", sm.serverId, sm.log[sm.nextIndex[i]:])	
					actions = append(actions, Send{peerId: sm.peerIds[i], ev: AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.serverId, PrevLogIndex: sm.nextIndex[i]-1, PrevLogTerm: sm.log[sm.nextIndex[i]-1].Term, Entries: sm.log[sm.nextIndex[i]:] , CommitIndex: sm.commitIndex}})
				} else {
					//fmt.Printf("leader->%v, log entries sent->%v\n", sm.serverId, sm.log[sm.nextIndex[i]:])
					actions = append(actions, Send{peerId: sm.peerIds[i], ev: AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.serverId, PrevLogIndex: sm.nextIndex[i]-1, PrevLogTerm: 0, Entries: sm.log[sm.nextIndex[i]:] , CommitIndex: sm.commitIndex}})
				}	
				}
			}
		case "follower":
			actions = append(actions, Commit{index: int64(len(sm.log)), command: cmd.Data, err: errors.New("It's a follower, Not a leader")})
		case "candidate":
			actions = append(actions, Commit{index: int64(len(sm.log)), command: cmd.Data, err: errors.New("It's a candidate, Not a leader")})
		default: println("Invalid state")	
	}	
	return actions
}

