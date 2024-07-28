package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // Импорт драйвера SQLite
)

func main() {
	db, err := sql.Open("sqlite", "file:db_name.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		address TEXT,
		status TEXT,
		created_at TEXT
	)`)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	store := NewParcelStore(db)

	newParcel := Parcel{Client: 1, Address: "123 Main St", Status: ParcelStatusRegistered, CreatedAt: "2023-07-28T12:34:56Z"}
	id, err := store.Add(newParcel)
	if err != nil {
		fmt.Println("Error adding parcel:", err)
		return
	}
	fmt.Println("Added parcel with ID:", id)

	parcel, err := store.Get(id)
	if err != nil {
		fmt.Println("Error getting parcel:", err)
		return
	}
	fmt.Println("Got parcel:", parcel)

	parcels, err := store.GetByClient(1)
	if err != nil {
		fmt.Println("Error getting parcels by client:", err)
		return
	}
	fmt.Println("Got parcels by client:", parcels)

	err = store.SetStatus(id, ParcelStatusSent)
	if err != nil {
		fmt.Println("Error setting status:", err)
		return
	}
	fmt.Println("Set status to 'sent' for parcel with ID:", id)

	err = store.SetAddress(id, "456 Elm St")
	if err != nil {
		fmt.Println("Error setting address:", err)
		return
	}
	fmt.Println("Set address for parcel with ID:", id)

	err = store.Delete(id)
	if err != nil {
		fmt.Println("Error deleting parcel:", err)
		return
	}
	fmt.Println("Deleted parcel with ID:", id)
}
