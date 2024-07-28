package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func getTestParcelStore(t *testing.T) ParcelStore {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		address TEXT,
		status TEXT,
		created_at TEXT
	)`)
	if err != nil {
		t.Fatal(err)
	}

	return NewParcelStore(db)
}

func clearTable(db *sql.DB) {
	db.Exec("DELETE FROM parcel")
}

func TestAddAndGetParcel(t *testing.T) {
	store := getTestParcelStore(t)
	clearTable(store.db)

	newParcel := Parcel{Client: 1, Address: "123 Main St", Status: ParcelStatusRegistered, CreatedAt: "2023-07-28T12:34:56Z"}
	id, err := store.Add(newParcel)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parcel, err := store.Get(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if parcel.Client != newParcel.Client || parcel.Address != newParcel.Address || parcel.Status != newParcel.Status || parcel.CreatedAt != newParcel.CreatedAt {
		t.Fatalf("expected %v, got %v", newParcel, parcel)
	}
}

func TestGetByClient(t *testing.T) {
	store := getTestParcelStore(t)
	clearTable(store.db)

	parcels := []Parcel{
		{Client: 1, Address: "Address 1", Status: ParcelStatusRegistered, CreatedAt: "2023-07-28T12:34:56Z"},
		{Client: 1, Address: "Address 2", Status: ParcelStatusSent, CreatedAt: "2023-07-28T12:35:56Z"},
		{Client: 2, Address: "Address 3", Status: ParcelStatusDelivered, CreatedAt: "2023-07-28T12:36:56Z"},
	}

	for _, p := range parcels {
		_, err := store.Add(p)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}

	result, err := store.GetByClient(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 parcels, got %v", len(result))
	}
}

func TestSetStatus(t *testing.T) {
	store := getTestParcelStore(t)
	clearTable(store.db)

	newParcel := Parcel{Client: 1, Address: "123 Main St", Status: ParcelStatusRegistered, CreatedAt: "2023-07-28T12:34:56Z"}
	id, err := store.Add(newParcel)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = store.SetStatus(id, ParcelStatusSent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parcel, err := store.Get(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if parcel.Status != ParcelStatusSent {
		t.Fatalf("expected status %v, got %v", ParcelStatusSent, parcel.Status)
	}
}
