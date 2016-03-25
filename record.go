package main

import (
	"errors"
	"sort"
	"sync"
)

type Record struct {
	PackageName  string
	Dependencies []*Record
	DependedOn   []*Record
}

var (
	DEPENDENCY_NOT_AVAILABLE = errors.New("Dependency is not available")
)

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

func (rs *RecordSet) InsertRecords(newRecords ...*Record) (errs []error) {
	for i := range newRecords {
		err := rs.InsertRecord(newRecords[i])
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (rs *RecordSet) InsertRecord(newRecord *Record) error {
	if rs == nil {
		return nil
	}

	rs.ReadWriteLock.Lock()
	if rs.findRecord(newRecord.PackageName) == nil {
		rs.Records = append(rs.Records, newRecord)
	}
	rs.ReadWriteLock.Unlock()
	return nil
}

func (rs *RecordSet) InsertPackage(pkgName string, dependencies ...string) error {
	if rs == nil {
		return nil
	}

	rs.ReadWriteLock.Lock()
	depPackages, ok := rs.FindRequiredDependencies(dependencies...)
	if !ok {
		return DEPENDENCY_NOT_AVAILABLE
	}
	if rs.findRecord(pkgName) == nil {
		newRecord := NewRecord(pkgName, depPackages...)
		rs.Records = append(rs.Records, newRecord)
	}
	rs.ReadWriteLock.Unlock()
	return nil
}

func (rs *RecordSet) FindRequiredDependencies(dependencies ...string) (foundDependencies []*Record, noMissingDeps bool) {
	noMissingDeps = true
	for i := range dependencies {
		record := rs.findRecord(dependencies[i])
		if record != nil {
			foundDependencies = append(foundDependencies, record)
		} else {
			noMissingDeps = false
		}
	}
	return foundDependencies, noMissingDeps
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
