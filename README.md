# Blockchain for https://arxiv.org/abs/1708.09419

## Modules

```
go get github.com/boltdb/bolt
go get github.com/soniakeys/graph
go get github.com/soniakeys/bits
go get golang.org/x/crypto/ripemd160
go get github.com/fatih/color
```


## Launch

```
export NODE_ID=3000
go run *.go
```


## TODO

mine with solution

difficulty with solution

parallel mining with channel block to have realistic miners

block.validate should check that, if it has a solution (different from the posted problem) that solution is the best

rethink diff update?

maybe add check that graph has no better solution? (if no k+1-clique with current best k-clique is found)

add send problem graphs to server. Or maybe add ipfs implementation of problemgraphs

## BUGS
weird bug in genesis: txs hash changes after reload of blockchain, so that genesis pow check fails. the problem is the tx.serialize() that changes after closing the program and reopening it