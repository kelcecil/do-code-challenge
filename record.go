package main

import (
	"sort"
	"sync"
)

type Record struct {
	PackageName  string
	Dependencies []*Record
}

func NewRecord(name string, records ...*Record) *Record {
	return &Record{
		PackageName:  name,
		Dependencies: records,
	}
}

type RecordSet struct {
	Records       []*Record
	ReadWriteLock *sync.RWMutex
}

func NewRecordSet() *RecordSet {
	return &RecordSet{
		Records:       []*Record{},
		ReadWriteLock: &sync.RWMutex{},
	}
}

func (rs *RecordSet) Items() []*Record {
	if rs == nil {
		return nil
	}

	if !sort.IsSorted(rs) {
		sort.Sort(rs)
	}
	return rs.Records
}

func (rs *RecordSet) FetchRecords(recordNames ...string) (records []*Record) {
	if rs == nil {
		return nil
	}

	rs.ReadWriteLock.RLock()
	recordsInSet := rs.Items()
	for i := range recordNames {
		candidateIndice := sort.Search(len(rs.Records), func(j int) bool {
			return recordsInSet[j].PackageName >= recordNames[i]
		})
		candidateRecord := recordsInSet[candidateIndice]
		if candidateRecord.PackageName == recordNames[i] {
			records = append(records, candidateRecord)
		}
	}
	rs.ReadWriteLock.RUnlock()
	return records
}

func (rs *RecordSet) InsertRecord(newRecords ...*Record) error {
	if rs == nil {
		return nil
	}

	rs.ReadWriteLock.Lock()
	for i := range newRecords {
		rs.Records = append(rs.Records, newRecords[i])
	}
	rs.ReadWriteLock.Unlock()
	return nil
}

func (rs *RecordSet) Len() int {
	return len(rs.Records)
}

func (rs *RecordSet) Swap(i, j int) {
	rs.Records[i], rs.Records[j] = rs.Records[j], rs.Records[i]
}

func (rs *RecordSet) Less(i, j int) bool {
	return rs.Records[i].PackageName < rs.Records[j].PackageName
}
