package web

import (
	"errors"
	"html/template"
	"strings"
	"time"

	"github.com/philmacfly/wheretoeat/pkg/venue"
	"googlemaps.github.io/maps"
)

const layoutISO = "2006-01-02"

type defaultPage struct {
	Message  template.HTML
	Navbar   template.HTML
	Pagename string
}

type mainPage struct {
	Default defaultPage
	Venues  []webVenue
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
	LastVisit     string
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
	result.LastVisit = ""
	if len(v.Visits) > 0 {
		result.LastVisit = v.Visits[len(v.Visits)-1].Format(layoutISO)
	}

	return result
}

func convertOpeningHours(day time.Weekday, hours string) (maps.OpeningHoursPeriod, error) {
	var result maps.OpeningHoursPeriod

	input := strings.Replace(hours, " ", "", -1)
	input = strings.Replace(input, ":", "", -1)
	if hours == "" {
		return result, nil
	}
	periods := strings.Split(input, ";")
	if len(periods) < 1 {
		return result, errors.New("No valid periods found. Watch the formatting")
	}
	for _, p := range periods {
		if p == "" {
			break
		}
		ocs := strings.Split(p, "-")
		if len(ocs) != 2 {
			return result, errors.New("No valid opening/closing found. Watch the formatting")
		}
		result.Open.Day = day
		result.Open.Time = ocs[0]
		result.Close.Day = day
		result.Close.Time = ocs[1]
	}
	return result, nil
}

func convertWebVenuetoVenue(wv webVenue) (venue.Venue, error) {
	result := venue.Venue{VenueID: wv.VenueID, Name: wv.Name, Address: wv.Address,
		Rating: wv.Rating, GooglePlaceID: wv.GooglePlaceID, Website: wv.Website,
		PhoneNumber: wv.PhoneNumber, Notes: wv.Notes, Visits: wv.Visits}
	var ocs []maps.OpeningHoursPeriod
	oc, err := convertOpeningHours(time.Monday, wv.OpeningHours.Monday)
	if err != nil {
		return result, errors.New("Error converting Opening Hours:" + err.Error())
	}
	ocs = append(ocs, oc)
	oc, err = convertOpeningHours(time.Tuesday, wv.OpeningHours.Tuesday)
	if err != nil {
		return result, errors.New("Error converting Opening Hours:" + err.Error())
	}
	ocs = append(ocs, oc)
	oc, err = convertOpeningHours(time.Wednesday, wv.OpeningHours.Wednesday)
	if err != nil {
		return result, errors.New("Error converting Opening Hours:" + err.Error())
	}
	ocs = append(ocs, oc)
	oc, err = convertOpeningHours(time.Thursday, wv.OpeningHours.Thursday)
	if err != nil {
		return result, errors.New("Error converting Opening Hours:" + err.Error())
	}
	ocs = append(ocs, oc)
	oc, err = convertOpeningHours(time.Friday, wv.OpeningHours.Friday)
	if err != nil {
		return result, errors.New("Error converting Opening Hours:" + err.Error())
	}
	ocs = append(ocs, oc)
	oc, err = convertOpeningHours(time.Saturday, wv.OpeningHours.Saturday)
	if err != nil {
		return result, errors.New("Error converting Opening Hours:" + err.Error())
	}
	ocs = append(ocs, oc)
	oc, err = convertOpeningHours(time.Sunday, wv.OpeningHours.Sunday)
	if err != nil {
		return result, errors.New("Error converting Opening Hours:" + err.Error())
	}
	ocs = append(ocs, oc)
	for _, o := range ocs {
		if o.Open.Time != "" {
			result.OpeningHours.Periods = append(result.OpeningHours.Periods, o)
		}
	}
	return result, nil
}

type venueViewPage struct {
	Default defaultPage
	Venue   webVenue
}

type venueAddPage struct {
	Default defaultPage
	Venue   webVenue
}

type venueAddVisitPage struct {
	Default defaultPage
	Venue   webVenue
}

type updateDonePage struct {
	Default defaultPage
}

type nextOptionsPage struct {
	Default defaultPage
}
