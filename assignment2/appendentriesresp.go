package main

import (
	//"fmt"
)

type AppendEntriesRespEv struct {	
		from uint64
		term uint64
		success bool	
}

func (sm *StateMachine) AppendEntriesRespEventHandler ( event interface{} ) (actions []interface{}) {
	cmd := event.(AppendEntriesRespEv)
	//fmt.Printf("%v\n", cmd)
	switch sm.currentState {
		case "leader":
			if cmd.success == false {
				if sm.currentTerm < cmd.term {
					sm.currentTerm = cmd.term
					sm.votedFor = 0
					sm.currentState = "follower"
					actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
				} else {
					sm.nextIndex[cmd.from] = sm.nextIndex[cmd.from]-uint64(1)
					actions = append(actions, Send{peerId: cmd.from, ev: AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.serverId, prevLogIndex: sm.nextIndex[cmd.from]-uint64(1), prevLogTerm: sm.log[sm.nextIndex[cmd.from]-uint64(1)].term, entries: sm.log[sm.nextIndex[cmd.from]:], commitIndex: sm.commitIndex}})
				}
			} else {
				sm.nextIndex[cmd.from] = uint64(len(sm.log))
				lastCommitIndex := sm.commitIndex
				temp := uint64(1)
				for i:=0;i<len(sm.peerIds);i++ {
					if sm.matchIndex[i] > lastCommitIndex {
						for j:=0; j<len(sm.peerIds); j++ {
							if j!=i && sm.matchIndex[j]>=sm.matchIndex[i] {
								temp = temp + uint64(1)
							}		
							if temp >= sm.majority && sm.matchIndex[i] > lastCommitIndex {
								lastCommitIndex = sm.matchIndex[i]
								break
							}				
						}
						temp=1
					}
				}
					if lastCommitIndex > sm.commitIndex && sm.log[lastCommitIndex].term == sm.currentTerm {
						for i:=sm.commitIndex+uint64(1);i<=lastCommitIndex;i++ {
							actions = append(actions, Commit{index: i, command: sm.log[i].command, err: nil})
						}
						sm.commitIndex = lastCommitIndex
					}
			}
		case "follower":
			if cmd.term > sm.currentTerm {
				sm.currentTerm = cmd.term
				sm.votedFor = 0
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
			}	
		case "candidate":
		default: println("Invalid state")		
	}	
	return actions
}

