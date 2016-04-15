package main

import (
	"fmt"
)

type AppendEntriesReqEv struct {
		Term int64
		LeaderId int64
		PrevLogIndex int64
		PrevLogTerm int64
		Entries []logEntry
		CommitIndex int64		
}

func (sm *StateMachine) AppendEntriesReqEventHandler ( event interface{} ) (actions []interface{}){
	cmd := event.(AppendEntriesReqEv)
	fmt.Printf("%v\n", cmd)
	switch sm.currentState {
		case "leader":
			if sm.currentTerm <= cmd.Term {
				if sm.currentTerm < cmd.Term {
					sm.votedFor = 0
				}
				sm.currentTerm = cmd.Term
				sm.currentState = "follower"
			//	actions = append(actions, StateStore{state: sm.currentState, Term: sm.currentTerm, votedFor:sm.votedFor}) // make it function
			//	actions = append(actions, Send{peerId: sm.serverId, ev: AppendEntriesReqEv{Term: cmd.Term, leaderId: cmd.leaderId, prevLogIndex: cmd.prevLogIndex, prevLogTerm: cmd.prevLogTerm, entries: cmd.entries, commitIndex: cmd.commitIndex}}) //make it function
				actions = sm.AppendEntriesReqEventHandler(event) 
			} else {
				actions = append(actions, Send{peerId: cmd.LeaderId, ev: AppendEntriesRespEv{From: sm.serverId, Term: sm.currentTerm, Success: false}})
			}
		case "follower":
			
			if sm.currentTerm <= cmd.Term {
				if sm.currentTerm < cmd.Term {
					sm.votedFor = 0
				}
				sm.currentTerm = cmd.Term
				actions = append(actions, Alarm{t: int64(ElectionTimeoutGenerator(int(sm.ElectionTimeout), int(2*sm.ElectionTimeout)))})
				actions = append(actions, StateStore{term: sm.currentTerm, votedFor:sm.votedFor})
				if ( (cmd.PrevLogTerm == 0) || ( cmd.PrevLogIndex < int64(len(sm.log)) && (sm.log[cmd.PrevLogIndex].Term == cmd.PrevLogTerm)) ) {
					k:=0
					fmt.Printf("%v In appendentriesreq: Entries to be updated->%v\n",sm.serverId, cmd.Entries)
					templog := make([]logEntry, int(int(cmd.PrevLogIndex)+1+len(cmd.Entries)))
					for i:=0;i<int(cmd.PrevLogIndex)+1;i++ {
						templog[i]=sm.log[i]
					}
					
					for j:=int(cmd.PrevLogIndex)+1;j<(int(cmd.PrevLogIndex)+1+len(cmd.Entries));j++ {
						templog[j] = cmd.Entries[k]
						actions = append(actions, LogStore{index: int64(j), command: cmd.Entries[k]})
						k++
					}
					sm.log = templog

					//fmt.Printf("%v updated log: %v\n", sm.serverId, sm.log)
					actions = append(actions, Send{peerId: cmd.LeaderId, ev: AppendEntriesRespEv{From: sm.serverId, Term: cmd.Term, Success: true, Lastindex: int64(len(sm.log)-1)}})
					if cmd.CommitIndex > sm.commitIndex {
						if int64(len(sm.log)-1) < cmd.CommitIndex {
							for i:=sm.commitIndex+int64(1);i<=int64(len(sm.log)-1);i++ {
								fmt.Printf("%v In appendentriesreq: Commit data->%v\n", sm.serverId, sm.log[i].command)
								actions = append(actions, Commit{index: i, command: sm.log[i].command, err: nil})
							}
							sm.commitIndex = int64(len(sm.log)-1)
						} else {
							for i:=sm.commitIndex+int64(1);i<=cmd.CommitIndex;i++ {
									fmt.Printf("%v In appendentriesreq: Commit data->%v\n", sm.serverId, sm.log[i].command)
								actions = append(actions, Commit{index: i, command: sm.log[i].command, err: nil})
							}
							sm.commitIndex = cmd.CommitIndex
						}						
					}
				} else {
					actions = append(actions, Send{peerId: cmd.LeaderId, ev: AppendEntriesRespEv{From: sm.serverId, Term: cmd.Term, Success: false, Lastindex: int64(len(sm.log)-1)}})
				}
				
	   		} else {
	   			actions = append(actions, Send{peerId: cmd.LeaderId, ev: AppendEntriesRespEv{From: sm.serverId, Term: sm.currentTerm, Success: false, Lastindex: int64(len(sm.log)-1)}})
	   		}
	   	case "candidate":
	   			if sm.currentTerm <= cmd.Term {
	   				if sm.currentTerm < cmd.Term {
					sm.votedFor = 0
				    }
					sm.currentTerm = cmd.Term
					sm.currentState = "follower"
					actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
					actions = sm.AppendEntriesReqEventHandler(event)
				} else {
					actions = append(actions, Send{peerId: cmd.LeaderId, ev: AppendEntriesRespEv{From: sm.serverId, Term: sm.currentTerm, Success: false, Lastindex: int64(len(sm.log)-1)}})
				}	
		default: println("Invalid state")		
}
	return actions
}
