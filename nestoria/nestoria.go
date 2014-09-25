package nestoria

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/aktau/houser"
)

const (
	SCHEMA   = "http"
	URL_PATH = "api"
)

var DEBUG = false

var CountryToAPI = map[string]string{
	"deutschland":    "api.nestoria.de",
	"de":             "api.nestoria.de",
	"germany":        "api.nestoria.de",
	"france":         "api.nestoria.fr",
	"fr":             "api.nestoria.fr",
	"united kingdom": "api.nestoria.co.uk",
	"uk":             "api.nestoria.co.uk",
	"england":        "api.nestoria.co.uk",
	"scotland":       "api.nestoria.co.uk",
}

// Nestoria implements the houser.Repo interface
var _ houser.Repo = (*Nestoria)(nil)

func New(country string) (*Nestoria, error) {
	country = strings.ToLower(country)
	URI, ok := CountryToAPI[country]
	if !ok {
		return nil, errors.New("country " + country + " is not known, if you have an endpoint, please log a bug report. Meanwhile, you can use the NewEx() interface")
	}
	return NewEx(URI), nil
}

func NewEx(apiURI string) *Nestoria {
	return &Nestoria{uri: apiURI}
}

type Nestoria struct {
	uri string
}

func (n *Nestoria) Search(query *houser.Query) ([]*houser.Listing, error) {
	resp, err := searchListings(n.uri, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("nestoria: got unexpected status code, " + resp.Status)
	}

	if DEBUG {
		fmt.Println("status code =", resp.Status)
	}

	var r io.Reader = resp.Body
	if DEBUG {
		// print response to stdout if debugging is on
		r = io.TeeReader(resp.Body, os.Stdout)
	}

	var listings struct {
		Response struct {
			Listings []*NestoriaListing `json:"listings"`
		} `json:"response"`
	}

	if err := json.NewDecoder(r).Decode(&listings); err != nil {
		return nil, err
	}

	var clistings []*houser.Listing
	for _, l := range listings.Response.Listings {
		l.Finish()
		clistings = append(clistings, l.ToGeneric())
	}

	return clistings, nil
}

// http://api.nestoria.de/api?action=keywords&encoding=json&parameter=wut&pretty=true
// http://api.nestoria.com.br/show_example?syntax=1&name=keywords_de
func searchListings(host string, query *houser.Query) (*http.Response, error) {
	v := newQuery("search_listings")
	v.Set("place_name", query.City)
	if query.TransactionType != "" {
		v.Set("listing_type", string(query.TransactionType))
	}
	if query.RoomMin != 0 {
		v.Set("room_min", strconv.Itoa(int(query.RoomMin)))
	}
	if query.RoomMax != 0 {
		v.Set("room_max", strconv.Itoa(int(query.RoomMax)))
	}
	if query.PriceMax != 0 {
		v.Set("price_max", strconv.Itoa(int(query.PriceMax)))
	}
	if query.PriceMin != 0 {
		v.Set("price_min", strconv.Itoa(int(query.PriceMin)))
	}
	if query.AreaMin != 0 {
		v.Set("size_min", strconv.Itoa(int(query.AreaMin)))
	}
	if query.AreaMax != 0 {
		v.Set("size_max", strconv.Itoa(int(query.AreaMax)))
	}
	if !query.UpdatedSince.IsZero() {
		v.Set("updated_min", strconv.Itoa(int(query.UpdatedSince.UTC().Unix())))
	}
	return http.Get(genURL(host, v).String())
}

func keywords(host string) (*http.Response, error) {
	return http.Get(genURL(host, newQuery("keywords")).String())
}

func metadata(host string) (*http.Response, error) {
	return http.Get(genURL(host, newQuery("metadata")).String())
}

func newQuery(action string) url.Values {
	v := url.Values{}
	v.Set("action", action)
	return v
}

func genURL(host string, v url.Values) *url.URL {
	v.Set("encoding", "json")

	if DEBUG {
		v.Set("pretty", "true")
	}

	u := url.URL{
		Scheme:   SCHEMA,
		Host:     host,
		Path:     URL_PATH,
		RawQuery: v.Encode(),
	}

	return &u
}
