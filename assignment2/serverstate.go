package main

type logEntry struct {
		term uint64
		command []byte
}

type server_state_np struct {
		commitIndex uint64
		nextIndex []uint64
		matchIndex []uint64
}	

type server_state_p struct {
		log []logEntry
		currentTerm uint64
		votedFor uint64
		currentState string
}



