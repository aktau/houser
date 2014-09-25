package main

import (
	"fmt"
	"io"
	"text/tabwriter"
	"text/template"

	"github.com/aktau/houser"
)

// Output the listings as text to an io.Writer, decent for console output
func PrintTabulated(w io.Writer, listings []*houser.Listing) {
	tw := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)
	defer tw.Flush()

	tw.Write([]byte("updated\ttype\tarea\trooms\tprice\tconstr.\ttitle\n"))
	for _, l := range listings {
		str := fmt.Sprintf("%.1f\t%s\t%dm²\t%.1f\t%s%d\t%d\t%s\n",
			l.DaysSinceUpdate, l.Type, l.Area, l.Rooms, currencySymbol(l.Currency), l.Price, l.ConstructionYear, l.Title)
		tw.Write([]byte(str))
	}
}

type MailTemplateData struct {
	Listings []*houser.Listing
}

var emailTpl = `
<html>
	<table>
		<thead>
			<tr>
				<th>Updated</th>
				<th>Type</th>
				<th>Area</th>
				<th>Rooms</th>
				<th>Price</th>
				<th>Constr.</th>
				<th>Title</th>
				<th>Summary</th>
				<th>Keywords</th>
			</tr>
		</thead>
		<tbody>
			{{range $index, $listing := .Listings}}
				<tr>
					<td>{{$listing.DaysSinceUpdate}}</td>
					<td>{{$listing.Type}}</td>
					<td>{{$listing.Area}}m²</td>
					<td>{{$listing.Rooms}}</td>
					<td>{{currencySymbol $listing.Currency}}{{$listing.Price}}</td>
					<td>{{$listing.ConstructionYear}}</td>
					<td><a href="{{$listing.URL}}">{{$listing.Title}}</a></td>
					<td>{{$listing.Description}}</td>
					<td>{{$listing.Keywords}}</td>
				</tr>
			{{end}}
		</tbody>
	</table>
</html>`

var funcMap = template.FuncMap{
	"currencySymbol": currencySymbol,
}

// Output the listings as html to an io.Writer, decent for websites and mail
func PrintHTML(w io.Writer, listings []*houser.Listing) error {
	t, err := template.New("email").Funcs(funcMap).Parse(emailTpl)
	if err != nil {
		return err
	}

	if t.Execute(w, MailTemplateData{Listings: listings}) != nil {
		return err
	}

	return nil
}

func currencySymbol(cur houser.Currency) string {
	switch string(cur) {
	case "EUR":
		return "€"
	case "GBP":
		return "£"
	case "USD":
		return "$"
	default:
		return "?"
	}
}
