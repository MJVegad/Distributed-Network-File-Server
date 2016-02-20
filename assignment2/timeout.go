package main

import (
	//"fmt"
)

type TimeoutEv struct {
	
}

func (sm *StateMachine) TimeoutEventHandler ( event interface{} ) (actions []interface{}) {
	//cmd := event.(TimeoutEv)
	//fmt.Printf("%v\n", cmd)
	switch sm.currentState {
		case "leader":
			for i:=0;i<len(sm.peerIds);i++ {
				actions = append(actions, Send{peerId: sm.peerIds[i], ev: AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.serverId, prevLogIndex: sm.nextIndex[i]-1, prevLogTerm: sm.log[sm.nextIndex[i]-1].term, entries: sm.log[sm.nextIndex[i]:], commitIndex: sm.commitIndex}})
			}
		case "follower":
			sm.totalvotes = 1
			sm.novotes = 0
			sm.currentTerm = sm.currentTerm + 1
			sm.currentState = "candidate"
			sm.votedFor = sm.serverId
			actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
			actions = append(actions, Alarm{t: 100})
			for _, pid := range sm.peerIds {
				actions = append(actions, Send{peerId: pid, ev: VoteReqEv{term: sm.currentTerm, candidateId: sm.serverId, lastLogIndex: uint64(len(sm.log)-1), lastLogTerm: sm.log[len(sm.log)-1].term}})
			}
		case "candidate":
			sm.totalvotes = 1
			sm.novotes = 0
			sm.currentTerm = sm.currentTerm + 1
			sm.votedFor = sm.serverId
			actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
			actions = append(actions, Alarm{t: 100})
			for _, pid := range sm.peerIds {
				actions = append(actions, Send{peerId: pid, ev: VoteReqEv{term: sm.currentTerm, candidateId: sm.serverId, lastLogIndex: uint64(len(sm.log)-1), lastLogTerm: sm.log[len(sm.log)-1].term}})
			}
		default: println("Invalid state")	
	}	
	return actions
}

