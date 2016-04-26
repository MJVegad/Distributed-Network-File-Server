package main

import (
	//"fmt"
)

type TimeoutEv struct {
}

func (sm *StateMachine) TimeoutEventHandler(event interface{}) (actions []interface{}) {
	//cmd := event.(TimeoutEv)
	//fmt.Printf("%v\n", cmd)
	switch sm.currentState {
	case "leader":
		for i := 0; i < len(sm.peerIds); i++ {
			if (sm.nextIndex[i] - 1) < 0 {
				fmt.Printf("%v Inside Timeout if: logindex->%v\n", sm.serverId, sm.nextIndex[i])
				actions = append(actions, Send{peerId: sm.peerIds[i], ev: AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.serverId, PrevLogIndex: sm.nextIndex[i] - 1, PrevLogTerm: 0, Entries: sm.log[sm.nextIndex[i]:], CommitIndex: sm.commitIndex}})
			} else {
				fmt.Printf("%v Inside Timeout else: logindex->%v\n", sm.serverId, sm.nextIndex[i]-1)
				actions = append(actions, Send{peerId: sm.peerIds[i], ev: AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.serverId, PrevLogIndex: sm.nextIndex[i] - 1, PrevLogTerm: sm.log[sm.nextIndex[i]-1].Term, Entries: sm.log[sm.nextIndex[i]:], CommitIndex: sm.commitIndex}})
			}
		}
		actions = append(actions, Alarm{t: sm.HeartbeatTimeout})
	case "follower":
		sm.totalvotes = 1
		sm.novotes = 0
		sm.currentTerm = sm.currentTerm + 1
		sm.currentState = "candidate"
		sm.votedFor = sm.serverId
		actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor: sm.votedFor})
		actions = append(actions, Alarm{t: int64(ElectionTimeoutGenerator(int(sm.ElectionTimeout), int(2*sm.ElectionTimeout)))})
		for _, pid := range sm.peerIds {
			if (len(sm.log) - 1) < 0 {
				fmt.Printf("VOte request sent to:%v\n",pid)
				actions = append(actions, Send{peerId: pid, ev: VoteReqEv{Term: sm.currentTerm, CandidateId: sm.serverId, LastLogIndex: 0, LastLogTerm: 0}})
			} else {
				actions = append(actions, Send{peerId: pid, ev: VoteReqEv{Term: sm.currentTerm, CandidateId: sm.serverId, LastLogIndex: int64(len(sm.log) - 1), LastLogTerm: sm.log[len(sm.log)-1].Term}})
			}
		}
	case "candidate":
		sm.totalvotes = 1
		sm.novotes = 0
		sm.currentTerm = sm.currentTerm + 1
		sm.votedFor = sm.serverId
		actions = append(actions, StateStore{state: sm.currentState, term: sm.currentTerm, votedFor: sm.votedFor})
		actions = append(actions, Alarm{t: int64(ElectionTimeoutGenerator(int(sm.ElectionTimeout), int(2*sm.ElectionTimeout)))})
		for _, pid := range sm.peerIds {
			if (len(sm.log) - 1) < 0 {
				actions = append(actions, Send{peerId: pid, ev: VoteReqEv{Term: sm.currentTerm, CandidateId: sm.serverId, LastLogIndex: 0, LastLogTerm: 0}})
			} else {
				actions = append(actions, Send{peerId: pid, ev: VoteReqEv{Term: sm.currentTerm, CandidateId: sm.serverId, LastLogIndex: int64(len(sm.log) - 1), LastLogTerm: sm.log[len(sm.log)-1].Term}})
			}
		}
	default:
		println("Invalid state")
	}
	return actions
}
