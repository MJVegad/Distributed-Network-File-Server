package main

import ()

func (sm *StateMachine) ProcessEvent(ev interface{}) (actions []interface{}) {
	switch ev.(type) {
	case AppendEv:
		actions = sm.AppendEventHandler(ev)
	case AppendEntriesReqEv:
		actions = sm.AppendEntriesReqEventHandler(ev)
	case AppendEntriesRespEv:
		actions = sm.AppendEntriesRespEventHandler(ev)
	case TimeoutEv:
		actions = sm.TimeoutEventHandler(ev)
	case VoteReqEv:
		actions = sm.VoteReqEventHandler(ev)
	case VoteRespEv:
		actions = sm.VoteRespEventHandler(ev)
	default:
		println("unrecognized event")
	}
	return actions
}
