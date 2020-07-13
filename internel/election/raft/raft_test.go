package raft

import (
	"testing"
	"time"
)

func TestGetElectionRaftProvider(t *testing.T) {
	election := GetElectionRaftProvider()
	go election.Run()
	time.Sleep(time.Second * 2)
	election.Stop()
}
