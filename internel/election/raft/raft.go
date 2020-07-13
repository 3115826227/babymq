package raft

import (
	"encoding/json"
	"fmt"
	"sync"
)

type RaftState uint32

const (
	Follower RaftState = iota
	Candidate
	Leader
	Shutdown
)

type ServerMeta struct {
}

type ElectionRaftProvider struct {
	stateLock sync.Mutex
	state     RaftState

	leaderLock   sync.Mutex
	leader       ServerMeta
	shutdownChan chan bool

	ReceiveResponse chan ReceiveResponse
	SendRequest     chan SendRequest
}

type MessageType uint32

const (
	Vote MessageType = iota
)

type SendRequest struct {
	MessageType MessageType
	Content     []byte
}

type ReceiveResponse struct {
	MessageType MessageType
	Content     []byte
}

type VoteMeta struct {
	Term    uint32
	Version uint32
}

var election *ElectionRaftProvider

func init() {
	election = &ElectionRaftProvider{
		stateLock:    sync.Mutex{},
		state:        Candidate,
		leaderLock:   sync.Mutex{},
		leader:       ServerMeta{},
		shutdownChan: make(chan bool, 1),
	}
}

func GetElectionRaftProvider() *ElectionRaftProvider {
	return election
}

func (election *ElectionRaftProvider) Run() {
	for {
		select {
		case <-election.shutdownChan:
			election.setLeader(ServerMeta{})
			return
		default:
		}

		switch election.state {
		case Leader:
			election.runLeader()
		case Follower:
			election.runFollower()
		case Candidate:
			election.runCandidate()
		}
	}
}

func (election *ElectionRaftProvider) Stop() {
	election.shutdownChan <- true
}

func (election *ElectionRaftProvider) runFollower() {
	fmt.Println("follower")

	election.setState(Candidate)
}

func (election *ElectionRaftProvider) runCandidate() {
	fmt.Println("candidate")
	for election.getState() == Candidate {
		select {
		case resp := <-election.ReceiveResponse:
			switch resp.MessageType {
			case Vote:
				election.vote(resp)
			}
		}
	}
}

func (election *ElectionRaftProvider) runLeader() {
	fmt.Println("leader")
	election.setState(Follower)
}

func (election *ElectionRaftProvider) processRPC() {

}

func (election *ElectionRaftProvider) vote(resp ReceiveResponse) {
	fmt.Println(resp)
	var vote = VoteMeta{
		Term:    1,
		Version: 1,
	}
	content, err := json.Marshal(vote)
	if err != nil {
		return
	}
	election.SendRequest <- SendRequest{
		MessageType: Vote,
		Content:     content,
	}
	election.setState(Leader)
}

func (election *ElectionRaftProvider) setLeader(server ServerMeta) {
	election.leaderLock.Lock()
	defer election.leaderLock.Unlock()
	election.leader = server
}

func (election *ElectionRaftProvider) getState() RaftState {
	return election.state
}

func (election *ElectionRaftProvider) setState(state RaftState) {
	election.stateLock.Lock()
	defer election.stateLock.Unlock()
	election.state = state
}
