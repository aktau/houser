package houser

import "time"

// Defines an API that can be used to retrieve listings, can then be
// concretely implemented.

type Currency string

const (
	EURO   Currency = "EUR"
	DOLLAR          = "USD"
	POUND           = "GBP"
)

type PropertyType string

const (
	APARTMENT PropertyType = "apartment"
	HOUSE                  = "house"
)

type TransactionType string

const (
	RENT  TransactionType = "rent"  // renting a flat/house
	BUY   TransactionType = "buy"   // buying a flat/house
	SHARE TransactionType = "share" // shared flat/house
)

type Query struct {
	City            string
	PriceMin        uint
	PriceMax        uint
	RoomMin         uint
	RoomMax         uint
	AreaMin         uint
	AreaMax         uint
	PropertyType    PropertyType    // type of the desired property
	TransactionType TransactionType // type of the desired transaction
	UpdatedSince    time.Time
}

type Repo interface {
	Search(query *Query) ([]*Listing, error)
}
