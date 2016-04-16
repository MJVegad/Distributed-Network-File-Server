package main

import (
	"fmt"
)

type AppendEntriesRespEv struct {	
		From int64
		Term int64
		Success bool
		Lastindex int64	
}

func (sm *StateMachine) AppendEntriesRespEventHandler ( event interface{} ) (actions []interface{}) {
	cmd := event.(AppendEntriesRespEv)
	//fmt.Printf("%v\n", cmd)
	switch sm.currentState {
		case "leader":
			var ind int
			for i:=0;i<len(sm.peerIds);i++ {
						if(sm.peerIds[i]==cmd.From) {
							ind = i
							break
						}
					}
			if cmd.Success == false {
				if sm.currentTerm < cmd.Term {
					sm.currentTerm = cmd.Term
					sm.votedFor = 0
					sm.currentState = "follower"
					actions = append(actions, Alarm{t: int64(ElectionTimeoutGenerator(int(sm.ElectionTimeout), int(2*sm.ElectionTimeout)))})
					actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
				} else {
					sm.nextIndex[ind] = sm.nextIndex[ind]-int64(1)
					if sm.nextIndex[ind] < 0 {
						actions = append(actions, Send{peerId: cmd.From, ev: AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.serverId, PrevLogIndex: -1, PrevLogTerm: 0, Entries: nil, CommitIndex: sm.commitIndex}})
					} else if sm.nextIndex[ind] == 0 {
						actions = append(actions, Send{peerId: cmd.From, ev: AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.serverId, PrevLogIndex: -1, PrevLogTerm: 0, Entries: sm.log[0:], CommitIndex: sm.commitIndex}})
					} else {
						//fmt.Printf("%v Inside AppendEntriesResp: ind->%v, logindex->%v\n", sm.serverId, ind, sm.nextIndex[ind]-int64(1))
						actions = append(actions, Send{peerId: cmd.From, ev: AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.serverId, PrevLogIndex: sm.nextIndex[ind]-int64(1), PrevLogTerm: sm.log[sm.nextIndex[ind]-int64(1)].Term, Entries: sm.log[sm.nextIndex[ind]:], CommitIndex: sm.commitIndex}})
					}				
				}
			} else {
				sm.nextIndex[ind] = int64(cmd.Lastindex+1)
				sm.matchIndex[ind] = int64(cmd.Lastindex)
				lastCommitIndex := sm.commitIndex
				//fmt.Printf("last commit index===>%v", lastCommitIndex)
				temp := int64(1)
				for i:=0;i<len(sm.peerIds);i++ {
					if sm.matchIndex[i] > lastCommitIndex {
						//fmt.Printf("Dn't be here..!!")
						for j:=0; j<len(sm.peerIds); j++ {
							if j!=i && sm.matchIndex[j]>=sm.matchIndex[i] {
								temp = temp + int64(1)
							}		
							if temp >= sm.majority-1 && sm.matchIndex[i] > lastCommitIndex {
								lastCommitIndex = sm.matchIndex[i]
								break
							}				
						}
						temp=1
					}
				}
					if lastCommitIndex >=0 {
					if lastCommitIndex > sm.commitIndex && sm.log[lastCommitIndex].Term == sm.currentTerm {
						for i:=sm.commitIndex+int64(1);i<=lastCommitIndex;i++ {
							fmt.Printf("Leader->%v, Commit data->%v\n", sm.serverId, sm.log[i].Command)
							actions = append(actions, Commit{index: i, command: sm.log[i].Command, err: nil})
						}
						sm.commitIndex = lastCommitIndex
					}
				}
			}
		case "follower":
			if cmd.Term > sm.currentTerm {
				sm.currentTerm = cmd.Term
				sm.votedFor = 0
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
			}	
		case "candidate":
			if cmd.Term > sm.currentTerm {
				sm.currentTerm = cmd.Term
				sm.votedFor = 0
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
			}	
		default: println("Invalid state")		
	}	
	return actions
}

