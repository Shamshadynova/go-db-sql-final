package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		address TEXT,
		status TEXT,
		created_at TEXT
	)`)
	require.NoError(t, err)

	return db
}

func TestAddGetDelete(t *testing.T) {
	db := setupTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// get
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, storedParcel.Client)
	require.Equal(t, parcel.Status, storedParcel.Status)
	require.Equal(t, parcel.Address, storedParcel.Address)
	require.Equal(t, parcel.CreatedAt, storedParcel.CreatedAt)

	// delete
	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err)
}

func TestSetAddress(t *testing.T) {
	db := setupTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, storedParcel.Address)
}

func TestSetStatus(t *testing.T) {
	db := setupTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// set status
	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// check
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newStatus, storedParcel.Status)
}

func TestGetByClient(t *testing.T) {
	db := setupTestDB(t)
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	// add
	for i := range parcels {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotZero(t, id)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, storedParcel := range storedParcels {
		expectedParcel, exists := parcelMap[storedParcel.Number]
		require.True(t, exists)
		require.Equal(t, expectedParcel.Client, storedParcel.Client)
		require.Equal(t, expectedParcel.Status, storedParcel.Status)
		require.Equal(t, expectedParcel.Address, storedParcel.Address)
		require.Equal(t, expectedParcel.CreatedAt, storedParcel.CreatedAt)
	}
}
