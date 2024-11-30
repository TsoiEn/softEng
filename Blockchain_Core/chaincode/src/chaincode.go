package src

import (
	"fmt"
	"log"

	"github.com/TsoiEn/softEng/Blockchain_Core/chaincode/consensus"
	model "github.com/TsoiEn/softEng/Blockchain_Core/chaincode/src/model"
)

type Blockchain struct {
	Blocks   []*model.Block      // Blockchain blocks
	RaftNode *consensus.RaftNode // Raft node for consensus
}

// Initialize the blockchain with a genesis block.
func (bc *Blockchain) InitLedger() error {
	// Step 1: Create the genesis block.
	genesisBlock := model.CreateBlock(1, []byte("Genesis Block"), []byte(""))

	// Step 2: Propose the genesis block via Raft.
	success := bc.RaftNode.ProposeBlock(genesisBlock)
	if !success {
		return fmt.Errorf("failed to propose genesis block via Raft")
	}

	// Step 3: Append the genesis block after consensus.
	bc.Blocks = append(bc.Blocks, genesisBlock)
	log.Printf("Genesis block initialized via Raft: %+v", genesisBlock)
	return nil
}

// Create a new block and add it to the blockchain.
func (bc *Blockchain) CreateBlock(data string) error {
	// Ensure the blockchain is initialized.
	if len(bc.Blocks) == 0 {
		return fmt.Errorf("blockchain is not initialized")
	}

	// Step 1: Get the last block.
	lastBlock := bc.Blocks[len(bc.Blocks)-1]

	// Step 2: Create a new block using the previous block's hash.
	newBlock := model.CreateBlock(lastBlock.Index+1, []byte(data), lastBlock.Hash)

	// Step 3: Propose the new block via Raft.
	success := bc.RaftNode.ProposeBlock(newBlock)
	if !success {
		return fmt.Errorf("failed to propose block via Raft")
	}

	// Step 4: Append the new block to the blockchain after consensus.
	bc.Blocks = append(bc.Blocks, newBlock)
	log.Printf("New block added via Raft: %+v", newBlock)
	return nil
}
