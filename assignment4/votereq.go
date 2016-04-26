package main

import (
//"fmt"
)

type VoteReqEv struct {
	Term         int64
	CandidateId  int64
	LastLogIndex int64
	LastLogTerm  int64
}

func (sm *StateMachine) VoteReqEventHandler(event interface{}) (actions []interface{}) {
	cmd := event.(VoteReqEv)
	//fmt.Printf("%v\n", cmd)
	switch sm.currentState {
	case "leader":
		if sm.currentTerm < cmd.Term {
			sm.currentTerm = cmd.Term
			sm.votedFor = 0
			sm.currentState = "follower"
			actions = append(actions, Alarm{t: int64(ElectionTimeoutGenerator(int(sm.ElectionTimeout), int(2*sm.ElectionTimeout)))})
			if (len(sm.log)-1 < 0) || ((sm.log[len(sm.log)-1].Term < cmd.LastLogTerm) || (sm.log[len(sm.log)-1].Term == cmd.LastLogTerm && int64(len(sm.log)-1) <= cmd.LastLogIndex)) {
				actions = append(actions, Send{cmd.CandidateId, VoteRespEv{Term: sm.currentTerm, VoteGranted: true}})
				sm.votedFor = cmd.CandidateId
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor: sm.votedFor})
			} else {
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor: sm.votedFor})
				actions = append(actions, Send{cmd.CandidateId, VoteRespEv{Term: sm.currentTerm, VoteGranted: false}})
			}
		} else {
			actions = append(actions, Send{cmd.CandidateId, VoteRespEv{Term: sm.currentTerm, VoteGranted: false}})
		}
	case "follower":
		if (sm.currentTerm < cmd.Term) || ((sm.currentTerm == cmd.Term) && (sm.votedFor == 0 || sm.votedFor == cmd.CandidateId)) {
			var temp bool
			if sm.currentTerm < cmd.Term {
				sm.votedFor = 0
				sm.currentTerm = cmd.Term
				temp = true
			}
			if (len(sm.log)-1 < 0) || ((sm.log[len(sm.log)-1].Term < cmd.LastLogTerm) || (sm.log[len(sm.log)-1].Term == cmd.LastLogTerm && int64(len(sm.log)-1) <= cmd.LastLogIndex)) {
				sm.votedFor = cmd.CandidateId
				temp = true
				actions = append(actions, Send{cmd.CandidateId, VoteRespEv{Term: sm.currentTerm, VoteGranted: true}})
				actions = append(actions, Alarm{t: int64(ElectionTimeoutGenerator(int(sm.ElectionTimeout), int(2*sm.ElectionTimeout)))})
			} else {
				actions = append(actions, Send{cmd.CandidateId, VoteRespEv{Term: sm.currentTerm, VoteGranted: false}})
			}
			if temp == true {
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor: sm.votedFor})
			}
		} else {
			actions = append(actions, Send{cmd.CandidateId, VoteRespEv{Term: sm.currentTerm, VoteGranted: false}})
		}
	case "candidate":
		if sm.currentTerm >= cmd.Term {
			actions = append(actions, Send{cmd.CandidateId, VoteRespEv{Term: sm.currentTerm, VoteGranted: false}})
		} else {
			sm.currentState = "follower"
			actions = append(actions, Alarm{t: int64(ElectionTimeoutGenerator(int(sm.ElectionTimeout), int(2*sm.ElectionTimeout)))})
			sm.currentTerm = cmd.Term
			sm.votedFor = 0
			if (len(sm.log)-1 < 0) || ((sm.log[len(sm.log)-1].Term < cmd.LastLogTerm) || (sm.log[len(sm.log)-1].Term == cmd.LastLogTerm && int64(len(sm.log)-1) <= cmd.LastLogIndex)) {
				sm.votedFor = cmd.CandidateId
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor: sm.votedFor})
				actions = append(actions, Send{cmd.CandidateId, VoteRespEv{Term: sm.currentTerm, VoteGranted: true}})
			} else {
				actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor: sm.votedFor})
				actions = append(actions, Send{cmd.CandidateId, VoteRespEv{Term: sm.currentTerm, VoteGranted: false}})
			}
		}
	default:
		println("Invalid state")

	}
	return actions
}
