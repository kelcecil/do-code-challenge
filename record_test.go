package main

import "testing"

func TestSimpleRecordSort(t *testing.T) {
	homebrewRecord := NewRecord("homebrew")
	golangRecord := NewRecord("golang")

	recordSet := NewRecordSet()
	recordSet.InsertRecord(golangRecord, homebrewRecord)

	records := recordSet.Items()

	if records[0].PackageName != "golang" && records[1].PackageName != "homebrew" {
		t.Errorf("Packages in RecordSet are not sorted in alphabetical order.")
	}
}

func TestSimpleRecordFetch(t *testing.T) {
	homebrewRecord := NewRecord("homebrew")
	golangRecord := NewRecord("golang")

	recordSet := NewRecordSet()
	recordSet.InsertRecord(golangRecord, homebrewRecord)

	fetchedRecord := recordSet.FetchRecords("homebrew")
	if len(fetchedRecord) != 1 || fetchedRecord[0].PackageName != "homebrew" {
		t.Error("Fetching a single record failed.")
	}

	fetchedRecord = recordSet.FetchRecords("golang", "homebrew")
	if len(fetchedRecord) != 2 || fetchedRecord[0].PackageName != "golang" {
		t.Errorf("Fetching multiple records at once failed.")
	}
}
