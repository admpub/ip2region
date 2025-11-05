package ip2region

import (
	"fmt"
	"os"
	"sync"

	"github.com/admpub/ip2region/v3/binding/golang/xdb"
)

type Ip2Region struct {
	dbFile      string
	dbBuff      []byte
	vectorIndex []byte
	memoryMode  bool
	dbVer       *xdb.Version
	mu          sync.RWMutex
}

func loadContent(dbPath string) ([]byte, error) {
	return xdb.LoadContentFromFile(dbPath)
}

func loadHeader(dbPath string) (*xdb.Version, error) {
	// auto-detect the ip version from the xdb header
	header, err := xdb.LoadHeaderFromFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load header from `%s`: %s", dbPath, err)
	}

	var version *xdb.Version
	version, err = xdb.VersionFromHeader(header)
	if err != nil {
		return nil, fmt.Errorf("failed to detect IP version from `%s`: %s", dbPath, err)
	}
	return version, err
}

func New(dbPath string, memoryMode bool) (*Ip2Region, error) {
	var cBuff []byte
	var vectorIndex []byte
	var err error
	if memoryMode {
		cBuff, err = loadContent(dbPath)
	} else {
		vectorIndex, err = xdb.LoadVectorIndexFromFile(dbPath)
	}
	if err != nil {
		return nil, err
	}
	var version *xdb.Version
	version, err = loadHeader(dbPath)
	if err != nil {
		return nil, err
	}
	//searcher, err := xdb.NewWithFileOnly(path)
	return &Ip2Region{
		dbFile:      dbPath,
		dbBuff:      cBuff,
		vectorIndex: vectorIndex,
		dbVer:       version,
		memoryMode:  memoryMode,
	}, nil
}

func (a *Ip2Region) DBBuff() []byte {
	a.mu.RLock()
	buff := a.dbBuff
	a.mu.RUnlock()
	return buff
}

func (a *Ip2Region) VectorIndex() []byte {
	a.mu.RLock()
	vectorIndex := a.vectorIndex
	a.mu.RUnlock()
	return vectorIndex
}

func (a *Ip2Region) Close() error {
	return nil
}

func (a *Ip2Region) MemorySearchString(ipStr string) (result string, err error) {
	searcher := searcherPool.Get().(*xdb.Searcher)
	if a.memoryMode {
		searcher.SetContentBuff(a.DBBuff())
	} else {
		var handle *os.File
		handle, err = os.OpenFile(a.dbFile, os.O_RDONLY, 0600)
		if err != nil {
			return
		}
		searcher.SetHandle(handle)
		searcher.SetVectorIndex(a.VectorIndex())
	}
	searcher.SetVersion(a.dbVer)
	result, err = searcher.SearchByStr(ipStr)
	searcher.Close()
	searcher.Reset()
	searcherPool.Put(searcher)
	return
}

func (a *Ip2Region) MemorySearch(ipStr string) (ipInfo IpInfo, err error) {
	var result string
	result, err = a.MemorySearchString(ipStr)
	ipInfo.Parse(result)
	return
}
