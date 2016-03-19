package main

type Send struct {
	peerId int64
	ev interface{}
}

type LogStore struct {
	index int64
	command logEntry
}

type Alarm struct {
	t int64
}

type Commit struct {
	index int64
	command []byte
	err error
}

type StateStore struct {
	state string
	term int64
	votedFor int64
}

