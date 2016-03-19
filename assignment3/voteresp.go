package main

import (
	"fmt"
)

type VoteRespEv struct {	
		Term int64
		VoteGranted bool	
}

func (sm *StateMachine) VoteRespEventHandler ( event interface{} ) (actions []interface{}) {
	cmd := event.(VoteRespEv)
	fmt.Printf("%v\n", cmd)
	switch sm.currentState {
		case "leader":
		case "follower":
			if sm.currentTerm < cmd.Term {
				sm.currentTerm = cmd.Term
				sm.votedFor = 0
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
			}	
		case "candidate":
			if cmd.VoteGranted == true {
				sm.totalvotes = sm.totalvotes + 1
				if sm.totalvotes >= sm.majority {
					sm.currentState = "leader"
					actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
					actions = append(actions, Alarm{t: sm.HeartbeatTimeout})
					for i:=0;i<len(sm.peerIds);i++ {
						sm.nextIndex[i] = int64(len(sm.log))
						sm.matchIndex[i] = 0
						if (len(sm.log)-2)<0 {
							actions = append(actions, Send{peerId: sm.peerIds[i], ev: AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.serverId, PrevLogIndex: int64(-1), PrevLogTerm: 0, Entries: nil, CommitIndex: sm.commitIndex}})				
						} else {
							actions = append(actions, Send{peerId: sm.peerIds[i], ev: AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.serverId, PrevLogIndex: int64(len(sm.log)-2), PrevLogTerm: sm.log[len(sm.log)-2].Term, Entries: nil, CommitIndex: sm.commitIndex}})
						}					
					}
				}
			} else {
				if cmd.Term > sm.currentTerm {
					sm.currentTerm = cmd.Term
					sm.votedFor = 0
					sm.currentState = "follower"
					actions = append(actions, Alarm{t: int64(ElectionTimeoutGenerator(int(sm.ElectionTimeout), int(2*sm.ElectionTimeout)))})
				} else {
					sm.novotes = sm.novotes + 1
					if sm.novotes >= sm.majority {
						sm.currentState = "follower"
						actions = append(actions, Alarm{t: int64(ElectionTimeoutGenerator(int(sm.ElectionTimeout), int(2*sm.ElectionTimeout)))})
					}
				}
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
			}
		default: println("Invalid state")	
				
	}	
	return actions
}

