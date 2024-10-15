package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)

	require.NoError(t, err)

	require.NotEmpty(t, id)
	parcel.Number = id

	p, err := store.Get(parcel.Number)

	require.NoError(t, err)

	require.Equal(t, parcel.Number, p.Number)
	require.Equal(t, parcel.Client, p.Client)
	require.Equal(t, parcel.Status, p.Status)
	require.Equal(t, parcel.Address, p.Address)
	require.Equal(t, parcel.CreatedAt, p.CreatedAt)

	err = store.Delete(id)

	require.NoError(t, err)

	_, err = store.Get(p.Number)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)

	require.NoError(t, err)

	require.NotEmpty(t, id)
	parcel.Number = id

	newAddress := "new test address"
	err = store.SetAddress(parcel.Number, newAddress)

	require.NoError(t, err)

	p, err := store.Get(parcel.Number)

	require.NoError(t, err)

	require.Equal(t, newAddress, p.Address)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)

	require.NoError(t, err)

	require.NotEmpty(t, id)
	parcel.Number = id

	newStatus := ParcelStatusDelivered
	err = store.SetStatus(parcel.Number, newStatus)

	require.NoError(t, err)

	p, err := store.Get(parcel.Number)
	require.NoError(t, err)

	require.Equal(t, newStatus, p.Status)
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)

	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)

	require.Equal(t, len(parcels), len(storedParcels))
	for _, parcel := range storedParcels {
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
