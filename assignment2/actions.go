package main

type Send struct {
	peerId uint64
	ev interface{}
}

type LogStore struct {
	index uint64
	command logEntry
}

type Alarm struct {
	t uint64
}

type Commit struct {
	index uint64
	command []byte
	err error
}

type StateStore struct {
	state string
	term uint64
	votedFor uint64
}

