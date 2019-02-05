# Blockchain for https://arxiv.org/abs/1708.09419


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

add go get to README


## BUGS
weird bug in genesis: txs hash changes after reload of blockchain, so that genesis pow check fails. the problem is the tx.serialize() that changes after closing the program and reopening it