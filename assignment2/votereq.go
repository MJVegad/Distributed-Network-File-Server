package main

import (
	"fmt"
)

type VoteReqEv struct {	
		term uint64
		candidateId uint64
		lastLogIndex uint64
		lastLogTerm uint64	
}

func (sm *StateMachine) VoteReqEventHandler ( event interface{} ) (actions []interface{}) {
	cmd := event.(VoteReqEv)
	fmt.Printf("%v\n", cmd)
	switch sm.currentState {
		case "leader":
			if sm.currentTerm < cmd.term {
				sm.currentTerm = cmd.term
				sm.votedFor = 0
				sm.currentState = "follower"
				actions = append(actions, Alarm{t: 100})
				if ((sm.log[len(sm.log)-1].term < cmd.lastLogTerm) || (sm.log[len(sm.log)-1].term == cmd.lastLogTerm && uint64(len(sm.log)-1)<=cmd.lastLogIndex) ) {
					actions = append(actions, Send {cmd.candidateId, VoteRespEv {term: sm.currentTerm, voteGranted: true}})
					sm.votedFor = cmd.candidateId
					actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
				} else {
					actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
					actions = append(actions, Send {cmd.candidateId, VoteRespEv {term: sm.currentTerm, voteGranted: false}})
				}
			} else {
				actions = append(actions, Send {cmd.candidateId, VoteRespEv {term: sm.currentTerm, voteGranted: false}})
			}
		case "follower":
			if (sm.currentTerm < cmd.term) && (sm.votedFor == 0 || sm.votedFor == cmd.candidateId) {
				sm.currentTerm = cmd.term
				sm.votedFor = 0
				if ((sm.log[len(sm.log)-1].term < cmd.lastLogTerm) || (sm.log[len(sm.log)-1].term == cmd.lastLogTerm && uint64(len(sm.log)-1)<=cmd.lastLogIndex) ) {
					sm.votedFor = cmd.candidateId
					actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
					actions = append(actions, Send {cmd.candidateId, VoteRespEv {term: sm.currentTerm, voteGranted: true}})
					actions = append(actions, Alarm{t: 100})
				} else {
					actions = append(actions, Send {cmd.candidateId, VoteRespEv {term: sm.currentTerm, voteGranted: false}})
				}
			} else {
				actions = append(actions, Send {cmd.candidateId, VoteRespEv {term: sm.currentTerm, voteGranted: false}})
			}
		case "candidate":
			if (sm.currentTerm >= cmd.term) {
				actions = append(actions, Send {cmd.candidateId, VoteRespEv {term: sm.currentTerm, voteGranted: false}})
			} else {
				sm.currentState = "follower"
				actions = append(actions, Alarm{t: 100})
				sm.currentTerm = cmd.term
				sm.votedFor = 0
				if ((sm.log[len(sm.log)-1].term < cmd.lastLogTerm) || (sm.log[len(sm.log)-1].term == cmd.lastLogTerm && uint64(len(sm.log)-1)<=cmd.lastLogIndex) ) {
					sm.votedFor = cmd.candidateId
					actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
					actions = append(actions, Send {cmd.candidateId, VoteRespEv {term: sm.currentTerm, voteGranted: true}})
				} else {
					actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor:sm.votedFor})
					actions = append(actions, Send {cmd.candidateId, VoteRespEv {term: sm.currentTerm, voteGranted: false}})
				}	
			}
		default: println("Invalid state")				
		
	}
	return actions
}

