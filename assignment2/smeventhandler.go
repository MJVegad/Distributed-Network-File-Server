package main

import (
	
)

func main () {
	var sm StateMachine
	sm.ProcessEvent(AppendEntriesRespEv{from : 2, term: 4, success: true})
}

func (sm *StateMachine) ProcessEvent (ev interface{}) {
		switch ev.(type) {
			case AppendEv:
				sm.AppendEventHandler (ev)
			case AppendEntriesReqEv:
				sm.AppendEntriesReqEventHandler (ev)
			case AppendEntriesRespEv:
				sm.AppendEntriesRespEventHandler (ev)
			case TimeoutEv:
				sm.TimeoutEventHandler (ev)
			case VoteReqEv:
				sm.VoteReqEventHandler (ev)
			case VoteRespEv:
				sm.VoteRespEventHandler (ev)
			default: println("unrecognized event")										
		}
	
}

