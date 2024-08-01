package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		address TEXT,
		status TEXT,
		created_at TEXT
	)`)
	assert.NoError(t, err)

	return db
}

func TestAddGetDelete(t *testing.T) {
	db := setupTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.Greater(t, id, 0, "ID should be greater than zero")

	// get
	storedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, parcel.Client, storedParcel.Client)
	assert.Equal(t, parcel.Status, storedParcel.Status)
	assert.Equal(t, parcel.Address, storedParcel.Address)
	assert.Equal(t, parcel.CreatedAt, storedParcel.CreatedAt)

	// delete
	err = store.Delete(id)
	assert.NoError(t, err)

	_, err = store.Get(id)
	assert.Error(t, err)
}

func TestSetAddress(t *testing.T) {
	db := setupTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.Greater(t, id, 0, "ID should be greater than zero")

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	assert.NoError(t, err)

	// check
	storedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, newAddress, storedParcel.Address)
}

func TestSetStatus(t *testing.T) {
	db := setupTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.Greater(t, id, 0, "ID should be greater than zero")

	// set status
	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	assert.NoError(t, err)

	// check
	storedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, newStatus, storedParcel.Status)
}

func TestGetByClient(t *testing.T) {
	db := setupTestDB(t)
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	// add
	for i := range parcels {
		id, err := store.Add(parcels[i])
		assert.NoError(t, err)
		assert.Greater(t, id, 0, "ID should be greater than zero")
		parcels[i].Number = id
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	assert.NoError(t, err)

	// check
	assert.ElementsMatch(t, parcels, storedParcels)
}
