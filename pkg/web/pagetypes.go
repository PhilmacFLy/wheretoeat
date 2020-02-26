package web

import (
	"html/template"
	"time"

	"github.com/philmacfly/wheretoeat/pkg/venue"
)

type defaultPage struct {
	Message  template.HTML
	Navbar   template.HTML
	Pagename string
}

type mainPage struct {
	Default defaultPage
	Venues  []venue.Venue
}

type webOpeningHours struct {
	Monday    string
	Tuesday   string
	Wednesday string
	Thursday  string
	Friday    string
	Saturday  string
	Sunday    string
}

type webVenue struct {
	VenueID       string
	Name          string
	Address       string
	Rating        int
	GooglePlaceID string
	OpeningHours  webOpeningHours
	Website       string
	PhoneNumber   string
	Notes         string
	Visits        []time.Time
}

func convertVenuetoWebVenue(v venue.Venue) webVenue {
	result := webVenue{VenueID: v.VenueID, Name: v.Name, Address: v.Address,
		Rating: v.Rating, GooglePlaceID: v.GooglePlaceID, Website: v.Website,
		PhoneNumber: v.PhoneNumber, Notes: v.Notes, Visits: v.Visits}

	for _, value := range v.OpeningHours.Periods {
		time := value.Open.Time + "-" + value.Close.Time + ";"
		switch value.Open.Day {
		case 0:
			result.OpeningHours.Sunday = result.OpeningHours.Sunday + time
		case 1:
			result.OpeningHours.Monday = result.OpeningHours.Monday + time
		case 2:
			result.OpeningHours.Tuesday = result.OpeningHours.Tuesday + time
		case 3:
			result.OpeningHours.Wednesday = result.OpeningHours.Wednesday + time
		case 4:
			result.OpeningHours.Thursday = result.OpeningHours.Thursday + time
		case 5:
			result.OpeningHours.Friday = result.OpeningHours.Friday + time
		case 6:
			result.OpeningHours.Saturday = result.OpeningHours.Saturday + time
		}
	}
	return result
}

type venueViewPage struct {
	Default defaultPage
	Venue   webVenue
}
