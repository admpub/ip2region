package ip2region

import (
	"sync"

	"github.com/admpub/ip2region/binding/golang/xdb"
)

func (a *Ip2Region) Reload(newPath ...string) error {
	dbPath := a.dbFile
	if len(newPath) > 0 && len(newPath[0]) > 0 {
		dbPath = newPath[0]
	}
	var cBuff []byte
	var vectorIndex []byte
	var err error
	if a.memoryMode {
		cBuff, err = loadContent(dbPath)
	} else {
		vectorIndex, err = xdb.LoadVectorIndexFromFile(dbPath)
	}
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.dbBuff = cBuff
	a.vectorIndex = vectorIndex
	a.mu.Unlock()
	return nil
}

var searcherPool = sync.Pool{
	New: func() interface{} {
		return &xdb.Searcher{}
	},
}
