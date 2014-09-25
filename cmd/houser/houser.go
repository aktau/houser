package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	"github.com/aktau/houser"
	"github.com/aktau/houser/nestoria"
)

var (
	mailEnabled = flag.Bool("m", false, "send the output via mail as well (needs other flags)")
	mailUser    = flag.String("u", "", "username of the mail account")
	mailPass    = flag.String("p", "", "password of the mail account")
	mailHost    = flag.String("h", "", "host smtp server (url:port)")
	mailDst     = flag.String("dst", "", "mail address of receiving user")
)

func init() {
	flag.Parse()
}

func main() {
	if err := rmain(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func rmain() error {
	fmt.Println("Welcome to Houser 0.2")

	if *mailEnabled {
		fmt.Printf("Going to send mail to %s via %s\n", *mailDst, *mailUser)
	}

	nestoria.DEBUG = true

	s, err := nestoria.New("deutschland")
	if err != nil {
		return fmt.Errorf("can't create search object: %v", err)
	}

	q := &houser.Query{
		TransactionType: "rent",
		RoomMin:         2,
		PriceMin:        500,
		PriceMax:        1600,
		AreaMin:         45,
		// UpdatedSince:    time.Now().AddDate(0, 0, -24),
	}

	locations := []string{
		"Sendling-Westpark",
		// "Maxvorstadt",
		// "lehel_muenchen",
		// "schwabing-ost",
		// "schwabing-west",
		// "au-haidhausen",
		// "bogenhausen_muenchen",
	}

	var buf bytes.Buffer
	w := io.MultiWriter(os.Stdout, &buf)

	var listings []*houser.Listing
	for _, loc := range locations {
		q.City = loc
		listings = searchAndOutput(w, s, q)
	}

	if !*mailEnabled {
		return nil
	}

	// split host:port into two variables
	split := strings.Split(*mailHost, ":")
	var port int
	host := split[0]
	if len(split) == 1 {
		port = 25
	} else {
		port, _ = strconv.Atoi(split[1])
	}

	user := *mailUser
	pass := *mailPass

	fmt.Println("sending mail from account", *mailUser, "over server", host, port)
	var body bytes.Buffer
	if err := PrintHTML(&body, listings); err != nil {
		return err
	}
	subject := "Subject: Test email from Go!\n"
	err = sendMailHTML(user, pass, host, port, subject, body.String(), []string{*mailDst})
	if err != nil {
		return fmt.Errorf("could not send mail: %v", err)
	}

	return nil
}

func searchAndOutput(w io.Writer, s houser.Repo, q *houser.Query) []*houser.Listing {
	listings, err := s.Search(q)
	if err != nil {
		fmt.Println("error while searching: ", err)
		os.Exit(1)
	}

	sbPrice := func(c1, c2 *houser.Listing) bool { return c1.Price < c2.Price }
	sbRooms := func(c1, c2 *houser.Listing) bool { return c1.Rooms < c2.Rooms }
	houser.OrderedBy(sbPrice, sbRooms).Sort(listings)

	fmt.Fprintln(w, q.City, "\n===================")
	PrintTabulated(w, listings)

	return listings
}

func sendMailHTML(user, pass, host string, port int, subject, body string, to []string) error {
	auth := smtp.PlainAuth("", user, pass, host)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + mime + body)
	return smtp.SendMail(fmt.Sprintf("%s:%d", host, port), auth, user, to, msg)
}
