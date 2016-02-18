package main

import (
	"fmt"
)

type VoteRespEv struct {	
		term uint64
		voteGranted bool	
}

func (sm *StateMachine) VoteRespEventHandler ( event interface{} ) (actions []interface{}) {
	cmd := event.(VoteRespEv)
	fmt.Printf("%v\n", cmd)
	switch sm.currentState {
		case "leader":
		case "follower":
			if sm.currentTerm < cmd.term {
				sm.currentTerm = cmd.term
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
			}	
		case "candidate":
			if cmd.voteGranted == true {
				sm.totalvotes = sm.totalvotes + 1
				if sm.totalvotes >= sm.majority {
					sm.currentState = "leader"
					actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
					for _, pid := range sm.peerIds {
						sm.nextIndex[pid] = uint64(len(sm.log))
						sm.matchIndex[pid] = 0
						actions = append(actions, Send{peerId: pid, ev: AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.serverId, prevLogIndex: uint64(len(sm.log)-2), prevLogTerm: sm.log[len(sm.log)-2].term, entries: nil, commitIndex: sm.commitIndex}})
					}
					actions = append(actions, Alarm{t: 100})
				}
			}
		default: println("Invalid state")	
				
	}	
	return actions
}

