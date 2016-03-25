package main

import "testing"

func TestSimpleRecordSort(t *testing.T) {
	recordSet := createNewRecordSetWithData()

	records := recordSet.Items()

	if records[0].PackageName != "golang" && records[1].PackageName != "homebrew" {
		t.Errorf("Packages in RecordSet are not sorted in alphabetical order.")
	}
}

func TestSimpleRecordFetch(t *testing.T) {
	recordSet := createNewRecordSetWithData()

	fetchedRecord := recordSet.FetchRecords("homebrew")
	if len(fetchedRecord) != 1 || fetchedRecord[0].PackageName != "homebrew" {
		t.Error("Fetching a single record failed.")
	}

	fetchedRecord = recordSet.FetchRecords("golang", "homebrew")
	if len(fetchedRecord) != 2 || fetchedRecord[0].PackageName != "golang" {
		t.Errorf("Fetching multiple records at once failed.")
	}
}

func TestInsertNewPackage(t *testing.T) {
	recordSet := createNewRecordSetWithData()
	recordSet.InsertRecord(&Record{
		PackageName: "golang",
	})
	count := 0
	records := recordSet.Items()
	for i := range records {
		if records[i].PackageName == "golang" {
			count++
		}
	}
	if count != 1 {
		t.Error("Duplicate records should not be possible.")
	}
}

func TestFindKnownInsertedDependencies(t *testing.T) {
	recordSet := createNewRecordSetWithData()

	records, ok := recordSet.FindRequiredDependencies("golo", "homebrew")
	if len(records) != 2 || !ok {
		t.Errorf("Expected to find 2 dependencies and ok = true; Received %v dependencies and ok = %v", len(records), ok)
	}
	expectedDepOne := records[0].PackageName
	expectedDepTwo := records[1].PackageName
	if expectedDepOne != "golo" || expectedDepTwo != "homebrew" {
		t.Errorf("Expected golo and homebrew deps; Received %v and %v", expectedDepOne, expectedDepTwo)
	}
}

func TestFindKnownNotInsertedDependencies(t *testing.T) {
	recordSet := createNewRecordSetWithData()

	records, ok := recordSet.FindRequiredDependencies("java", "fish")
	recordCount := len(records)
	if recordCount != 0 || ok {
		t.Errorf("Expected to find 0 dependencies and ok = false; Received %v dependencies and ok = %v", recordCount, ok)
	}
}

func TestFindMixOfDependencies(t *testing.T) {
	recordSet := createNewRecordSetWithData()

	records, ok := recordSet.FindRequiredDependencies("golo", "fish")
	recordCount := len(records)
	if recordCount != 1 || ok {
		t.Errorf("Expected to find 1 dependencies and ok = false; Received %v dependencies and ok = %v", recordCount, ok)
	}
	expectedDep := records[0].PackageName
	if expectedDep != "golo" {
		t.Errorf("Expected golo deps; Received %v", expectedDep)
	}
}

func TestInsertDuplicateRecord(t *testing.T) {
	recordSet := createNewRecordSetWithData()
	recordSet.InsertRecord(&Record{
		PackageName: "golang",
	})
	count := 0
	records := recordSet.Items()
	for i := range records {
		if records[i].PackageName == "golang" {
			count++
		}
	}
	if count != 1 {
		t.Error("Duplicate records should not be possible.")
	}
}

func TestFetchNonExistentRecord(t *testing.T) {
	recordSet := createNewRecordSetWithData()

	fetchedRecord := recordSet.FetchRecords("harveyRabbit")
	if len(fetchedRecord) != 0 {
		t.Errorf("A record was fetched that should not exist.")
	}
}

func BenchmarkInsertRecords(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createNewRecordSetWithData()
	}
}

func createNewRecordSetWithData() *RecordSet {
	homebrewRecord := NewRecord("homebrew")
	golangRecord := NewRecord("golang")
	goloRecord := NewRecord("golo")
	sdlRecord := NewRecord("sdl")

	recordSet := NewRecordSet()
	recordSet.InsertRecords(golangRecord, homebrewRecord, goloRecord, sdlRecord)
	return recordSet
}
