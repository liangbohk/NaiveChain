package core

import (
	"crypto/sha256"
	"log"
)

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNode *MerkleNode
	CheckData []byte
}

func NewMerkleTree(data [][]byte) *MerkleTree {

	if len(data) == 0 {
		log.Panic("empty data")
	}

	var merkleNodes []*MerkleNode

	for _, dataItem := range data {
		merkleNode := NewMerkleNode(nil, nil, dataItem)
		merkleNodes = append(merkleNodes, merkleNode)
	}

	nodes := GenMerkleTree(merkleNodes)
	return &MerkleTree{nodes[0]}
}

func GenMerkleTree(nodes []*MerkleNode) []*MerkleNode {
	if len(nodes) == 0 {
		log.Panic("no nodes provided")
	} else if len(nodes) == 1 {
		return nodes
	} else {
		var newNodes []*MerkleNode
		if len(nodes)%2 != 0 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}
		for i := 0; i < len(nodes); i = i + 2 {
			newNodes = append(newNodes, NewMerkleNode(nodes[i], nodes[i+1], nil))
		}
		nodes = newNodes

	}
	return GenMerkleTree(nodes)
}

func NewMerkleNode(leftNode, rightNode *MerkleNode, data []byte) *MerkleNode {

	node := &MerkleNode{}
	if leftNode == nil && rightNode == nil {
		hash := sha256.Sum256(data)
		node.CheckData = hash[:]
	} else {
		hash := append(leftNode.CheckData, rightNode.CheckData...)
		node.CheckData = hash[:]
	}
	node.LeftNode = leftNode
	node.RightNode = rightNode

	return node
}
