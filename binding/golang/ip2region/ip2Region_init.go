package ip2region

import (
	"io/ioutil"
	"sync"
)

func (this *Ip2Region) InitMemorySearch() (err error) {
	this.totalBlocksOnce.Do(func() {
		this.dbBinStr, err = ioutil.ReadFile(this.dbFile)
		if err != nil {
			return
		}

		this.firstIndexPtr = getLong(this.dbBinStr, 0)
		this.lastIndexPtr = getLong(this.dbBinStr, 4)
		this.totalBlocks = (this.lastIndexPtr-this.firstIndexPtr)/INDEX_BLOCK_LENGTH + 1
	})
	return
}

func (this *Ip2Region) ResetMemoryAndBinarySearch() {
	this.totalBlocksOnce = sync.Once{}
}

func (this *Ip2Region) InitBinarySearch() (err error) {
	if len(this.dbBinStr) > 0 {
		return this.InitMemorySearch()
	}
	this.totalBlocksOnce.Do(func() {
		this.dbFileHandler.Seek(0, 0)
		superBlock := make([]byte, 8)
		this.dbFileHandler.Read(superBlock)
		this.firstIndexPtr = getLong(superBlock, 0)
		this.lastIndexPtr = getLong(superBlock, 4)
		this.totalBlocks = (this.lastIndexPtr-this.firstIndexPtr)/INDEX_BLOCK_LENGTH + 1
	})
	return
}

func (this *Ip2Region) InitBtreeSearch() (err error) {
	this.headerLenOnce.Do(func() {
		this.headerSip = []int64{}
		this.headerPtr = []int64{}

		this.dbFileHandler.Seek(8, 0)

		buffer := make([]byte, TOTAL_HEADER_LENGTH)
		this.dbFileHandler.Read(buffer)
		var idx int64
		for i := 0; i < TOTAL_HEADER_LENGTH; i += 8 {
			startIp := getLong(buffer, int64(i))
			dataPar := getLong(buffer, int64(i+4))
			if dataPar == 0 {
				break
			}

			this.headerSip = append(this.headerSip, startIp)
			this.headerPtr = append(this.headerPtr, dataPar)
			idx++
		}

		this.headerLen = idx
	})
	return
}

func (this *Ip2Region) ResetBtreeSearch() {
	this.headerLenOnce = sync.Once{}
}
