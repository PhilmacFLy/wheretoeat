package venue

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math"
	"os"
	"strings"
	"time"

	"googlemaps.github.io/maps"
)

//Venue hod all information for a place to eat
type Venue struct {
	VenueID          string
	Name             string
	Address          string
	Rating           int
	GooglePlaceID    string
	OpeningHours     maps.OpeningHours
	OpeningHoursText []string
	Website          string
	PhoneNumber      string
	Notes            string
	Visted           []time.Time
}

type candidates struct {
	PlaceID string  `json:"place_id"`
	Name    string  `json:"name"`
	Rating  float64 `json:"rating"`
}

type quickQueryResponse struct {
	Candidates candidates `json:"candidates"`
	DebugLog   struct {
		Line []interface{} `json:"line"`
	} `json:"debug_log"`
	Status string `json:"status"`
}

const searchqueryfields = "formatted_address,name,place_id,rating"
const detailqueryfields = "opening_hours,website,international_phone_number"

var client *maps.Client
var searchqueryfieldsmask []maps.PlaceSearchFieldMask
var detailqueryfieldsmask []maps.PlaceDetailsFieldMask

func parseSearchFields(fields string) ([]maps.PlaceSearchFieldMask, error) {
	var res []maps.PlaceSearchFieldMask
	for _, s := range strings.Split(fields, ",") {
		f, err := maps.ParsePlaceSearchFieldMask(s)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

func parseDetailFields(fields string) ([]maps.PlaceDetailsFieldMask, error) {
	var res []maps.PlaceDetailsFieldMask
	for _, s := range strings.Split(fields, ",") {
		f, err := maps.ParsePlaceDetailsFieldMask(s)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

//SetupPlaceAPI Setups the API key for the client and the fields for the search
func SetupPlaceAPI(apikey string) error {
	var err error
	client, err = maps.NewClient(maps.WithAPIKey(apikey))
	if err != nil {
		return errors.New("Error setting up client:" + err.Error())
	}
	searchqueryfieldsmask, err = parseSearchFields(searchqueryfields)
	if err != nil {
		return errors.New("Error populating search fieldmask:" + err.Error())
	}
	detailqueryfieldsmask, err = parseDetailFields(detailqueryfields)
	if err != nil {
		return errors.New("Error populating seardetail fieldmask:" + err.Error())
	}
	return nil
}

//GetVenubyPlaceSearch takes the query and queries the googleplaces api with it
func GetVenubyPlaceSearch(query string) (Venue, error) {
	var res Venue

	searchRequest := &maps.FindPlaceFromTextRequest{
		Input:     query,
		InputType: maps.FindPlaceFromTextInputTypeTextQuery,
		Fields:    searchqueryfieldsmask,
	}

	searchResp, err := client.FindPlaceFromText(context.Background(), searchRequest)
	if err != nil {
		return res, errors.New("Error on search query:" + err.Error())
	}

	candidate := searchResp.Candidates[0]
	res.Name = candidate.Name
	res.Rating = int(math.Round(float64(candidate.Rating)))
	res.GooglePlaceID = candidate.PlaceID
	res.Address = candidate.FormattedAddress

	detailRequest := &maps.PlaceDetailsRequest{
		PlaceID: res.GooglePlaceID,
		Fields:  detailqueryfieldsmask,
	}

	detailResp, err := client.PlaceDetails(context.Background(), detailRequest)
	if err != nil {
		return res, errors.New("Error on detail query:" + err.Error())
	}

	res.OpeningHours = *detailResp.OpeningHours
	res.Website = detailResp.Website
	res.PhoneNumber = detailResp.InternationalPhoneNumber
	res.OpeningHoursText = detailResp.OpeningHours.WeekdayText

	hasher := sha1.New()
	hasher.Write([]byte(res.Name + res.Address))
	res.VenueID = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return res, nil
}

//SavetoFile save a venue to a json File
func (v *Venue) SavetoFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return errors.New("Error creating file: " + err.Error())
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(v)
	if err != nil {
		return errors.New("Error saving file: " + err.Error())
	}
	return nil
}

//LoadfromFile Loads a venue from a json File
func (v *Venue) LoadfromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return errors.New("Error opening file: " + err.Error())
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(v)
	if err != nil {
		return errors.New("Error decoding file: " + err.Error())
	}
	return nil
}
