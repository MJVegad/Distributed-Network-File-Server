package main

import (
	"math/rand"
)

type logEntry struct {
		Term int64
		command []byte
}



func ElectionTimeoutGenerator (min int, max int) int {
	return min + rand.Intn(max-min)
}

type StateMachine struct {
		serverId int64
		peerIds []int64
		majority int64
		commitIndex int64
		nextIndex []int64
		matchIndex []int64
		log []logEntry
		currentTerm int64
		votedFor int64
		currentState string
		totalvotes int64
		novotes int64
		ElectionTimeout  int64
	    HeartbeatTimeout int64
} 

