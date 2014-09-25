package nestoria

import (
	"strconv"
	"strings"
	"time"

	"github.com/aktau/houser"
)

type NestoriaListing struct {
	Title               string  `json:"title"`
	Description         string  `json:"summary"`
	URL                 string  `json:"lister_url"`
	Comission           float64 `json:"comission,string"` // a multiplier or percentage?
	CarSpaces           *string `json:"car_spaces"`
	Source              string  `json:"datasource_name"`
	ConstructionYearRaw *string `json:"construction_year"` // actually an int, needs to be parsed separately...
	ConstructionYear    uint    `json:"-"`
	Floor               int     `json:"floor,string"`
	Guid                string  `json:"guid"`
	Rooms               float64 `json:"room_number"`
	Type                string  `json:"property_type"`
	Keywords            string  `json:"keywords"`

	Lister string `json:"string"` // the person or company who posted the listing

	Size     uint   `json:"size,string"`
	SizeUnit string `json:"size_unit"`

	Price     uint   `json:"price"`
	PriceBare uint   `json:"price_coldrent"`
	Currency  string `json:"price_currency"`

	AuctionDate     *time.Time `json:"auction_date"`
	DaysSinceUpdate float64    `json:"updated_in_days"`
}

func (n *NestoriaListing) Finish() {
	if n.ConstructionYearRaw != nil {
		constr, _ := strconv.Atoi(*n.ConstructionYearRaw)
		n.ConstructionYear = uint(constr)
	}
}

func (n *NestoriaListing) ToGeneric() *houser.Listing {
	return &houser.Listing{
		Title:            n.Title,
		Description:      n.Description,
		URL:              n.URL,
		Rooms:            n.Rooms,
		Area:             n.Size,
		ConstructionYear: n.ConstructionYear,
		Type:             houser.PropertyType(n.Type),
		DaysSinceUpdate:  n.DaysSinceUpdate,
		Price:            n.Price,
		PriceBare:        n.PriceBare,
		Currency:         houser.Currency(n.Currency),
		Deposit:          0,
		Keywords:         strings.Split(n.Keywords, ","),
	}
}
