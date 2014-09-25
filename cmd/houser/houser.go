package main

import (
	"fmt"
	"os"

	"github.com/aktau/houser"
	"github.com/aktau/houser/nestoria"
)

func searchAndPrint(s houser.Repo, q *houser.Query) []*houser.Listing {
	listings, err := s.Search(q)
	if err != nil {
		fmt.Println("error while searching: ", err)
		os.Exit(1)
	}

	sbPrice := func(c1, c2 *houser.Listing) bool { return c1.Price < c2.Price }
	sbRooms := func(c1, c2 *houser.Listing) bool { return c1.Rooms < c2.Rooms }
	houser.OrderedBy(sbPrice, sbRooms).Sort(listings)

	fmt.Println(q.City, "\n===================")
	houser.PrintListings(os.Stdout, listings)

	return listings
}

func main() {
	fmt.Println("Welcome to Houser 0.1")

	// nestoria.DEBUG = true

	s, err := nestoria.New("deutschland")
	if err != nil {
		fmt.Println("can't create search object: ", err)
		os.Exit(1)
	}

	q := &houser.Query{
		TransactionType: "rent",
		RoomMin:         2,
		PriceMin:        500,
		PriceMax:        1600,
		AreaMin:         45,
		// UpdatedSince:    time.Now().AddDate(0, 0, -24),
	}

	q.City = "Sendling-Westpark "
	searchAndPrint(s, q)
	q.City = "Maxvorstadt"
	searchAndPrint(s, q)
	q.City = "lehel_muenchen"
	searchAndPrint(s, q)
	q.City = "schwabing-ost"
	searchAndPrint(s, q)
	q.City = "schwabing-west"
	searchAndPrint(s, q)
	q.City = "au-haidhausen"
	searchAndPrint(s, q)
	q.City = "bogenhausen_muenchen"
	searchAndPrint(s, q)
}
