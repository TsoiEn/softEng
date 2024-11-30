package consensus

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

type RaftNode struct {
	NodeID        string
	State         State
	CurrentTerm   int
	VotedFor      string
	Log           []Block
	CommitIndex   int
	LastApplied   int
	Peers         []string
	LeaderID      string
	ElectionTimer *time.Timer
	Mutex         sync.Mutex
}

type Block struct {
	Index     int
	Data      string
	Timestamp time.Time
	Hash      []byte
	PrevHash  []byte
}

func NewRaftNode(nodeID string, peers []string) *RaftNode {
	node := &RaftNode{
		NodeID:      nodeID,
		State:       Follower,
		CurrentTerm: 0,
		VotedFor:    "",
		Log:         []Block{},
		CommitIndex: 0,
		LastApplied: 0,
		Peers:       peers,
		LeaderID:    "",
	}
	node.ResetElectionTimer()
	return node
}

func (node *RaftNode) ResetElectionTimer() {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()
	if node.ElectionTimer != nil {
		node.ElectionTimer.Stop()
	}
	node.ElectionTimer = time.AfterFunc(time.Duration(150+rand.Intn(150))*time.Millisecond, func() {
		node.startElection()
	})
}

func (node *RaftNode) startElection() {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()
	node.State = Candidate
	node.CurrentTerm++
	node.VotedFor = node.NodeID
	votes := 1

	fmt.Printf("Node %s starts election for term %d\n", node.NodeID, node.CurrentTerm)

	// Request votes from peers
	for _, peer := range node.Peers {
		go func(peer string) {
			// Simulated RPC call for RequestVote
			success := node.requestVote(peer)
			if success {
				node.Mutex.Lock()
				votes++
				if votes > len(node.Peers)/2 && node.State == Candidate {
					node.becomeLeader()
				}
				node.Mutex.Unlock()
			}
		}(peer)
	}
}

func (node *RaftNode) becomeLeader() {
	node.State = Leader
	node.LeaderID = node.NodeID
	fmt.Printf("Node %s became the leader for term %d\n", node.NodeID, node.CurrentTerm)
	// Send heartbeat to maintain leadership
	go node.sendHeartbeats()

	// Example call to appendEntries to avoid unused function error
	node.appendEntries([]Block{
		{Index: 1, Data: "Initial entry", Timestamp: time.Now()},
	})
}

func (node *RaftNode) requestVote(peer string) bool {
	// Simulated logic for requesting a vote
	// In a real system, this would be a network call
	fmt.Printf("Node %s requests vote from %s\n", node.NodeID, peer)
	return true
}

func (node *RaftNode) appendEntries(entries []Block) bool {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	if node.State != Leader {
		return false
	}

	// Append entries to log and replicate to peers
	node.Log = append(node.Log, entries...)
	fmt.Printf("Leader %s appended entries and replicated to peers\n", node.NodeID)

	for _, peer := range node.Peers {
		go func(peer string) {
			node.replicateLog(peer)
		}(peer)
	}

	return true
}

func (node *RaftNode) replicateLog(peer string) {
	// Simulated logic for log replication
	fmt.Printf("Replicating log to peer %s\n", peer)
	// Assume successful replication
}

func (node *RaftNode) sendHeartbeats() {
	for node.State == Leader {
		fmt.Printf("Leader %s sending heartbeats\n", node.NodeID)
		for _, peer := range node.Peers {
			go func(peer string) {
				// Simulated RPC call for heartbeat
				node.sendHeartbeat(peer)
			}(peer)
		}
		time.Sleep(50 * time.Millisecond) // Send heartbeats periodically
	}
}

func (node *RaftNode) sendHeartbeat(peer string) {
	// Simulated logic for heartbeat
	fmt.Printf("Heartbeat sent to peer %s\n", peer)
}

// count of recussions: 1
