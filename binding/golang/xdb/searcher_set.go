package xdb

import "os"

func (s *Searcher) SetVectorIndex(vectorIndex []byte) *Searcher {
	s.vectorIndex = vectorIndex
	return s
}

func (s *Searcher) SetHandle(handle *os.File) *Searcher {
	s.handle = handle
	return s
}

func (s *Searcher) SetContentBuff(buff []byte) *Searcher {
	s.contentBuff = buff
	return s
}

func (s *Searcher) SetVersion(version *Version) *Searcher {
	s.version = version
	return s
}

func (s *Searcher) Reset() *Searcher {
	s.version = nil
	s.contentBuff = nil
	s.vectorIndex = nil
	s.ioCount = 0
	s.handle = nil
	return s
}
