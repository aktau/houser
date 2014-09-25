package houser

import (
	"sort"
	"time"
)

type Address struct {
	City     string
	Street   string
	Number   uint
	Postcode string
}

type Listing struct {
	Address

	Title            string       // short title of the listing
	Description      string       // description given by the owner
	URL              string       // URL for the listing
	Rooms            float64      // number of rooms in the property
	Area             uint         // area (in m2) of the property
	ConstructionYear uint         // year in which the property was constructed
	Type             PropertyType // the type of the listing (apartment, ...)
	Keywords         []string     // keywords identifying the property

	Posted          time.Time // timestamp of when the listing was posted
	Updated         time.Time // timestamp of when the listing was last updated
	DaysSinceUpdate float64

	PriceBare uint // bare price (per month)
	Price     uint // price with energy and other costs included (per month)
	Deposit   uint // security deposit (0 if none)
	Currency  Currency
}

type Listings []*Listing

// returns true if the Listing passes all filters
func all(listing *Listing, fs ...func(l *Listing) bool) bool {
	for _, filter := range fs {
		if !filter(listing) {
			return false
		}
	}
	return true
}

// filters a list of Listings according to the passed in filters
func (ls Listings) Filter(fs ...func(l *Listing) bool) Listings {
	var nls Listings
	for _, listing := range ls {
		if all(listing, fs...) {
			nls = append(nls, listing)
		}
	}
	return nls
}

type lessFunc func(p1, p2 *Listing) bool

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorter struct {
	listings []*Listing
	less     []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(listings []*Listing) {
	ms.listings = listings
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.listings)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.listings[i], ms.listings[j] = ms.listings[j], ms.listings[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that is either Less or
// !Less. Note that it can call the less functions twice per call. We
// could change the functions to return -1, 0, 1 and reduce the
// number of calls for greater efficiency: an exercise for the reader.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := ms.listings[i], ms.listings[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}
