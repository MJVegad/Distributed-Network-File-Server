package main

import (
	
)

type VoteRespEv struct {	
		term uint64
		voteGranted bool	
}

func (sm *StateMachine) VoteRespEventHandler ( event interface{} ) (actions []interface{}) {
	cmd := event.(VoteRespEv)
	//fmt.Printf("%v\n", cmd)
	switch sm.currentState {
		case "leader":
		case "follower":
			if sm.currentTerm < cmd.term {
				sm.currentTerm = cmd.term
				sm.votedFor = 0
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
			}	
		case "candidate":
			if cmd.voteGranted == true {
				sm.totalvotes = sm.totalvotes + 1
				if sm.totalvotes >= sm.majority {
					sm.currentState = "leader"
					actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
					actions = append(actions, Alarm{t: 100})
					for i:=0;i<len(sm.peerIds);i++ {
						sm.nextIndex[i] = uint64(len(sm.log))
						sm.matchIndex[i] = 0
						actions = append(actions, Send{peerId: sm.peerIds[i], ev: AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.serverId, prevLogIndex: uint64(len(sm.log)-2), prevLogTerm: sm.log[len(sm.log)-2].term, entries: nil, commitIndex: sm.commitIndex}})
					}
				}
			} else {
				if cmd.term > sm.currentTerm {
					sm.currentTerm = cmd.term
					sm.votedFor = 0
					sm.currentState = "follower"
					actions = append(actions, Alarm{t: 100})
				} else {
					sm.novotes = sm.novotes + 1
					if sm.novotes >= sm.majority {
						sm.currentState = "follower"
						actions = append(actions, Alarm{t: 100})
					}
				}
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
			
			}
		default: println("Invalid state")	
				
	}	
	return actions
}

