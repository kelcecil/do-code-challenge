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
	for i := range recordNames {
		record := rs.findRecord(recordNames[i])
		if record != nil {
			records = append(records, record)
		}
	}
	rs.ReadWriteLock.RUnlock()
	return records
}

func (rs *RecordSet) findRecord(recordName string) *Record {
	recordsInSet := rs.Items()
	candidateIndice := sort.Search(len(recordsInSet), func(i int) bool {
		return recordsInSet[i].PackageName >= recordName
	})
	if candidateIndice >= len(recordsInSet) {
		return nil
	}
	candidateRecord := recordsInSet[candidateIndice]
	if candidateRecord.PackageName == recordName {
		return candidateRecord
	}
	return nil
}

func (rs *RecordSet) InsertRecord(newRecords ...*Record) error {
	if rs == nil {
		return nil
	}

	rs.ReadWriteLock.Lock()
	for i := range newRecords {
		if rs.findRecord(newRecords[i].PackageName) == nil {
			rs.Records = append(rs.Records, newRecords[i])
		}
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
