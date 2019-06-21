package MerkleTree

import (
	"container/list"
	"crypto/sha256"
	"log"
	"math"
	"reflect"
)

type MerkleTreeHashFunc func(data []byte) []byte

type MerkleTree struct {
	root             *MerkleNode
	hashFunction     MerkleTreeHashFunc
	treeChildrenSize int
	size             int
}

func (tt *MerkleTree) Size() int {
	return tt.size
}

func NewMerkleTree(hashFunction func(data []byte) []byte, treeChildrenSize int) *MerkleTree {
	return &MerkleTree{hashFunction: hashFunction, treeChildrenSize: treeChildrenSize}
}

type MerkleNode struct {
	value    []byte
	children []*MerkleNode
}

func NewMerkleNode(value []byte, children []*MerkleNode) *MerkleNode {
	return &MerkleNode{value: value, children: children}
}

func (mn *MerkleNode) SetValue(value []byte) {
	mn.value = value
}

func (mn *MerkleNode) AddChild(child *MerkleNode) {
	mn.children = append(mn.children, child)
}

func (mn *MerkleNode) GetValue() []byte {
	return mn.value
}

func (mn *MerkleNode) GetChildren() []*MerkleNode {
	return mn.children
}

func (tt *MerkleTree) Root() *MerkleNode {
	return tt.root
}

func (tt *MerkleTree) InsertNodes(values [][]byte) error {
	values = padInsertedValues(values, tt.treeChildrenSize)
	hashQueue := list.New()
	for _, value := range values {
		hashValue := tt.hashFunction(value)
		hashQueue.PushBack(NewMerkleNode(hashValue, nil))
	}

	var newNode *MerkleNode
	for hashQueue.Len() > 1 {
		var childrenSum []byte
		newNode = NewMerkleNode(nil, nil)

		for count := 0; count < tt.treeChildrenSize; count++ {
			node := hashQueue.Front()
			merkleNode, isCastOk := node.Value.(*MerkleNode)
			if !isCastOk {
				log.Fatalf("merkle tree walk queue encountred unrelated node of type [%v] ", reflect.TypeOf(node.Value))
			}
			childrenSum = append(childrenSum, merkleNode.GetValue()...)
			hashQueue.Remove(node)
			newNode.AddChild(merkleNode)
		}

		childrenHashSum := sha256.Sum256(childrenSum)
		newNode.SetValue(childrenHashSum[:])
		hashQueue.PushBack(newNode)
		tt.size += tt.treeChildrenSize
	}
	tt.size += 1
	tt.SetRoot(newNode)

	return nil
}

func padInsertedValues(values [][]byte, childrenSize int) [][]byte {
	requiredValueNumber := 0
	for i := 0; requiredValueNumber < len(values); i++ {
		requiredValueNumber = int(math.Pow(float64(childrenSize), float64(i)))
	}
	requiredValueNumber -= len(values)
	for i := 0; i < requiredValueNumber; i++ {
		// we don't have enough to create perfect ternary merkle tree, so we duplicate last items till we get
		// to a power of tree children size
		values = append(values, values[len(values)-1])
	}
	return values
}

func (tt *MerkleTree) SetRoot(root *MerkleNode) {
	tt.root = root
}

