package venue

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
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
	Visits           []time.Time
}

//ByName is for sorting Venues by Name
type ByName []Venue

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return strings.Compare(a[i].Name, a[j].Name) == -1 }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//ByNameReverse is for sorting Venues by Name
type ByNameReverse []Venue

func (a ByNameReverse) Len() int           { return len(a) }
func (a ByNameReverse) Less(i, j int) bool { return strings.Compare(a[i].Name, a[j].Name) == 1 }
func (a ByNameReverse) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//ByRating sorts Venues by Rating
type ByRating []Venue

func (a ByRating) Len() int           { return len(a) }
func (a ByRating) Less(i, j int) bool { return a[i].Rating < a[j].Rating }
func (a ByRating) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//ByRatingReverse sorts Venues by Rating
type ByRatingReverse []Venue

func (a ByRatingReverse) Len() int           { return len(a) }
func (a ByRatingReverse) Less(i, j int) bool { return a[i].Rating > a[j].Rating }
func (a ByRatingReverse) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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

var datafolder string

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

//SetDataFolder sets the folder where the json Files reside
func SetDataFolder(venuedatafolder string) {
	datafolder = venuedatafolder + string(os.PathSeparator)
}

//GenerateVenueID takes the Name and the Adress and builds the id from it
func (v *Venue) GenerateVenueID() string {
	hasher := sha1.New()
	hasher.Write([]byte(v.Name + v.Address))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
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

	if len(searchResp.Candidates) == 0 {
		return res, nil
	}

	candidate := searchResp.Candidates[0]
	res.Name = candidate.Name
	res.Rating = int(math.Round(float64(candidate.Rating)))
	res.GooglePlaceID = candidate.PlaceID
	res.Address = candidate.FormattedAddress

	err = res.UpdateInfos()
	if err != nil {
		return res, errors.New("Error updating details:" + err.Error())
	}

	return res, nil
}

func (v *Venue) getJSONFile() string {
	return filepath.Join(datafolder, v.VenueID) + ".json"
}

//LoadFromDataLocation loads the JSON File of the Venue from the Data Location
func (v *Venue) LoadFromDataLocation() error {
	return v.loadfromFile(v.getJSONFile())
}

//SavetoDataLocation saves the JSON File of the Venue to the Data Location
func (v *Venue) SavetoDataLocation() error {
	return v.savetoFile(v.getJSONFile())
}

//SavetoFile save a venue to a json File
func (v *Venue) savetoFile(filename string) error {
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
func (v *Venue) loadfromFile(filename string) error {
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

//Delete Removes the Venue file from the Drive
func (v *Venue) Delete() error {
	err := os.Remove(v.getJSONFile())
	if err != nil {
		return errors.New("Error deleting file: " + err.Error())
	}
	return nil
}

//UpdateInfos updates volatile Infos of a Venue (Opening Hours, Website, Phone Number)
func (v *Venue) UpdateInfos() error {
	detailRequest := &maps.PlaceDetailsRequest{
		PlaceID: v.GooglePlaceID,
		Fields:  detailqueryfieldsmask,
	}

	detailResp, err := client.PlaceDetails(context.Background(), detailRequest)
	if err != nil {
		return errors.New("Error on detail query:" + err.Error())
	}

	v.OpeningHours = *detailResp.OpeningHours
	v.Website = detailResp.Website
	v.PhoneNumber = detailResp.InternationalPhoneNumber
	v.OpeningHoursText = detailResp.OpeningHours.WeekdayText

	return nil
}

//ListVenues gives back a slice with all venues in a folder
func ListVenues() ([]Venue, error) {
	var result []Venue
	files, err := ioutil.ReadDir(datafolder)
	if err != nil {
		return result, errors.New("Error reading folder: " + err.Error())
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		extension := filepath.Ext(f.Name())
		if strings.Compare(extension, ".json") != 0 {
			continue
		}
		var v Venue
		err := v.loadfromFile(filepath.Join(datafolder, f.Name()))
		if err != nil {
			return result, errors.New("Error loading one venue: " + err.Error())
		}
		result = append(result, v)
	}
	return result, nil
}
