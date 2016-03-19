package main

import (
	"github.com/cs733-iitb/log"
)

type Node interface {
	Append ([]byte)
	CommitChannel() <- chan CommitInfo
	CommittedIndex() int
	Get(index int) (error, []byte)
	Id()
	LeaderId() int
	Shutdown()
}

func (rn *RaftNode) Append (data []byte) {
	rn.eventch <- AppendEv{Data: data}
}

func (rn *RaftNode) CommitChannel() <- chan CommitInfo {
	return rn.commitch
}

func (rn *RaftNode) CommittedIndex() int {
	return int(rn.sm.commitIndex)
}

func (rn *RaftNode) Get(index int) (err1 error, data []byte) {
	lg, err := log.Open(rn.logfile)
	lg.RegisterSampleEntry(logEntry{})
	assert(err == nil)
	defer lg.Close()
	val,err1 := lg.Get(int64(index))
	if err1==nil {
		data = val.(logEntry).command
		return err1,data
	} else {
		return err1, data
	}
}

func (rn *RaftNode) Id() int {
	return int(rn.sm.serverId)
}

func (rn *RaftNode) LeaderId() int {
	return int(rn.sm.votedFor)
}

func (rn *RaftNode) Shutdown() {
	rn.sm_messaging.Close()
}