package main

import (
	"fmt"
	"log"
	"encoding/gob"
	"crypto/sha256"
	"bytes"
	"encoding/json"
	"github.com/soniakeys/graph"
	"github.com/soniakeys/bits"
	//"github.com/boltdb/bolt"
)


const problemBucket = "graph"
const problemsdbFile = "problems.db"


type ProblemGraph struct {
	Hash 	[]byte
	Graph 	*graph.Undirected
}


func NewProblemGraph(nodes int, edges int) *ProblemGraph {
	// g := graph.Undirected{make(graph.AdjacencyList, 3)} // 3 nodes altogether
	// g.AddEdge(0, 1)
	// g.AddEdge(1, 2)
	// g.AddEdge(2, 0)
	g := graph.GnmUndirected(nodes, edges, nil)
	pg := ProblemGraph{[]byte{}, &g}
	pg.Hash = pg.GetHash()
	return &pg
}

func (pg *ProblemGraph) FindClique() bool {
	//we check that we have a siple (not loops nor parallels) graph
	simple, _ := pg.Graph.IsSimple()
	if !simple {
		return false
	}

	var maxCliques [][]int
	m := 0
	pg.Graph.BronKerbosch3(pg.Graph.BKPivotMaxDegree, func(c bits.Bits) bool {
		clique := c.Slice()
    	if len(clique) > m {
    		maxCliques = maxCliques[:0]
    		maxCliques = append(maxCliques, c.Slice())
    		m = len(clique)
    	} else if len(clique) == m {
    		maxCliques = append(maxCliques, c.Slice())
    	}
    	return true
	})

    fmt.Printf("Max Clique Size: %d\n", m)
	for _, c := range maxCliques {
		fmt.Println(c)
	}

	return true
}

// Hash the graph
func (pg *ProblemGraph) GetHash() []byte {
    arrBytes := []byte{}
    jsonBytes, _ := json.Marshal(pg)
    arrBytes = append(arrBytes, jsonBytes...)
    

	hash := sha256.Sum256(arrBytes)

	return hash[:]
}

// Serialize serializes the graph
func (pg *ProblemGraph) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(pg)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes a graph
func DeserializeGraph(d []byte) *ProblemGraph {
	var g ProblemGraph

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&g)
	if err != nil {
		log.Panic(err)
	}

	return &g
}

// // CreateBlockchain creates a new blockchain DB
// func CreateGraph(address) *graph.Undirected {
// 	if dbExists(problemsdbFile) {
// 		fmt.Println("Problems file already exists.")
// 		g := NewBlockchain(nodeID)
// 		return bc
// 	}
	
// 	db, err := bolt.Open(problemsdbFile, 0600, nil)
// 	if err != nil {
// 		log.Panic(err)
// 	}
	
// 	err = db.Update(func(tx *bolt.Tx) error {
// 		b, err := tx.CreateBucket([]byte(problemBucket))
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		err = b.Put(genesis.Hash, genesis.Serialize())
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		err = b.Put([]byte("l"), genesis.Hash)
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		tip = genesis.Hash
// 		return nil
// 	})
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	bc := Blockchain{tip, db}

// 	return &bc
// }


// // NewBlockchain creates a new Blockchain with genesis Block
// func OpenProblemset(nodeID string) *Blockchain {
// 	dbFile := fmt.Sprintf(dbFile, nodeID)
// 	if dbExists(dbFile) == false {
// 		fmt.Println("No existing blockchain found. Create one first.")
// 		os.Exit(1)
// 	}

// 	var tip []byte
// 	db, err := bolt.Open(dbFile, 0600, nil)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	err = db.Update(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte(blocksBucket))
// 		tip = b.Get([]byte("l"))

// 		return nil
// 	})
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	bc := Blockchain{tip, db}

// 	return &bc
// }