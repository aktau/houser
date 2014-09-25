package houser

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Output the listings as text to an io.Writer, ok for console output
func PrintListings(w io.Writer, listings []*Listing) {
	tw := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)
	defer tw.Flush()

	tw.Write([]byte("updated\ttype\tarea\trooms\tprice\tconstr.\ttitle\n"))
	for _, l := range listings {
		str := fmt.Sprintf("%.1f\t%s\t%dm²\t%.1f\t%d (%s)\t%d\t%s\n",
			l.DaysSinceUpdate, l.Type, l.Area, l.Rooms, l.Price, l.Currency, l.ConstructionYear, l.Title)
		tw.Write([]byte(str))
	}
}
