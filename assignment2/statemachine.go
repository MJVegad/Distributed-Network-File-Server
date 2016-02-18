package main

import (

)

type logEntry struct {
		term uint64
		command []byte
}

/*
type server_config struct {
		serverId uint64
		peerIds []uint64
		majority uint64
}

type server_state_np struct {
//		gotVotes uint64
		commitIndex uint64
		nextIndex []uint64
		matchIndex []uint64
}	

type server_state_p struct {
		log []logEntry
		currentTerm uint64
		votedFor uint64
//		logIndex uint64
		currentState string
}
*/

type StateMachine struct {
		serverId uint64
		peerIds []uint64
		majority uint64
		commitIndex uint64
		nextIndex []uint64
		matchIndex []uint64
		log []logEntry
		currentTerm uint64
		votedFor uint64
		currentState string
		totalvotes uint64
} 

