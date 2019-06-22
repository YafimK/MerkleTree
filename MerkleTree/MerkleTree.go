package MerkleTree

import (
	"container/list"
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

func (mt *MerkleTree) GetRoot() *MerkleNode {
	return mt.root
}

func (mt *MerkleTree) Size() int {
	return mt.size
}

func (mt *MerkleTree) InsertNodes(values [][]byte) error {
	values = padInsertedValues(values, mt.treeChildrenSize)
	hashQueue := list.New()
	for _, value := range values {
		hashValue := mt.hashFunction(value)
		hashQueue.PushBack(NewMerkleNode(hashValue, nil))
	}

	var newNode *MerkleNode
	for hashQueue.Len() > 1 {
		var childrenSum []byte
		newNode = NewMerkleNode(nil, nil)
		for count := 0; count < mt.treeChildrenSize; count++ {
			node := hashQueue.Front()
			merkleNode, isCastOk := node.Value.(*MerkleNode)
			if !isCastOk {
				log.Fatalf("merkle tree walk queue encountred unrelated node of type [%v] ", reflect.TypeOf(node.Value))
			}
			childrenSum = append(childrenSum, merkleNode.GetValue()...)
			hashQueue.Remove(node)
			newNode.AddChild(merkleNode)
		}
		childrenHashSum := mt.hashFunction(childrenSum)
		newNode.SetValue(childrenHashSum)
		hashQueue.PushBack(newNode)
		mt.size += mt.treeChildrenSize
	}
	mt.size += 1
	mt.SetRoot(newNode)

	return nil
}

func (mt *MerkleTree) SetRoot(root *MerkleNode) {
	mt.root = root
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

//NewMerkleTree receives function used to hash the inserted values and number of children per node - i.e. if set to 3, it's ternery tree
func NewMerkleTree(hashFunction MerkleTreeHashFunc, childrenPerNode int) *MerkleTree {
	return &MerkleTree{hashFunction: hashFunction, treeChildrenSize: childrenPerNode}
}

type MerkleNode struct {
	value    []byte
	children []*MerkleNode
}

func (mn *MerkleNode) String() string {
	return string(mn.value)
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
