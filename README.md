Simple Markle tree implementation


```go
    package main

    func MockHashFunction(value []byte) []byte {
        return value
    }

    tree := MerkleTree.NewMerkleTree(MockHashFunction, 3)
    value1 := []byte("porky")
    value2 := []byte("dooky")
    value3 := []byte("pooky")
    values := [][]byte{value1, value2, value3}
    tree.InsertNodes(values)

```

for more usage example, visit main.go or the tests:)