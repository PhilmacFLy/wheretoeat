package web

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/philmacfly/wheretoeat/pkg/config"

	"github.com/gorilla/mux"
	"github.com/philmacfly/wheretoeat/pkg/venue"
)

var criteriaweight config.Weight

func apierror(w http.ResponseWriter, r *http.Request, err string, httpcode int) {
	log.Println(err)
	er := errorResponse{strconv.Itoa(httpcode), err}
	j, erro := json.Marshal(&er)
	if erro != nil {
		return
	}
	http.Error(w, string(j), httpcode)
}

func mainAPIHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Main API Handler")
}

func listVenuesAPIHandler(w http.ResponseWriter, r *http.Request) {
	sb := r.FormValue("sortby")
	vv, err := venue.ListVenues()
	if err != nil {
		apierror(w, r, "Error Listing Venues: "+err.Error(), http.StatusInternalServerError)
		return
	}
	switch sb {
	case "name-desc":
		sort.Sort(venue.ByNameReverse(vv))
	case "rating":
		sort.Sort(venue.ByRating(vv))
	case "rating-desc":
		sort.Sort(venue.ByRatingReverse(vv))
	default:
		sort.Sort(venue.ByName(vv))
	}
	j, err := json.Marshal(&vv)
	if err != nil {
		apierror(w, r, "Error marshalling Venues: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getVenueAPIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	var result venue.Venue
	result.VenueID = i
	err := result.LoadFromDataLocation()
	if err != nil {
		apierror(w, r, "Error Loading Venue File: "+err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(&result)
	if err != nil {
		apierror(w, r, "Error marshalling Venue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func postVenueAPIHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var v venue.Venue
	err := decoder.Decode(&v)
	if err != nil {
		apierror(w, r, "Error decoding Venue: "+err.Error(), http.StatusBadRequest)
		return
	}
	v.VenueID = v.GenerateVenueID()
	err = v.SavetoDataLocation()
	if err != nil {
		apierror(w, r, "Error saving Venue: "+err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(&v)
	if err != nil {
		apierror(w, r, "Error marshalling Venue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func patchVenueAPIHander(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	var result venue.Venue
	result.VenueID = i
	err := result.LoadFromDataLocation()
	if err != nil {
		apierror(w, r, "Error Loading Venue File: "+err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var v venue.Venue
	err = decoder.Decode(&v)
	if err != nil {
		apierror(w, r, "Error decoding Venue: "+err.Error(), http.StatusBadRequest)
		return
	}
	v.VenueID = v.GenerateVenueID()
	err = v.SavetoDataLocation()
	if err != nil {
		apierror(w, r, "Error saving Venue: "+err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(&v)
	if err != nil {
		apierror(w, r, "Error marshalling Venue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func deleteVenueAPIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	var result venue.Venue
	result.VenueID = i
	err := result.Delete()
	if err != nil {
		apierror(w, r, "Error Deleting Venue File: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func addVisitAPIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	var result venue.Venue
	result.VenueID = i
	err := result.LoadFromDataLocation()
	if err != nil {
		apierror(w, r, "Error Loading Venue File: "+err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var a addVisitsRequest
	err = decoder.Decode(&a)
	if err != nil {
		apierror(w, r, "Error decoding Venue: "+err.Error(), http.StatusBadRequest)
		return
	}
	result.Visits = append(result.Visits, a.Visits...)
	err = result.SavetoDataLocation()
	if err != nil {
		apierror(w, r, "Error saving Venue: "+err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(&result)
	if err != nil {
		apierror(w, r, "Error marshalling Venue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getVenueFromPlacesAPIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	q := vars["query"]
	v, err := venue.GetVenubyPlaceSearch(q)
	if err != nil {
		apierror(w, r, "Error searching places api: "+err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(&v)
	if err != nil {
		apierror(w, r, "Error marshalling Venue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getNotVisitedVenue(w http.ResponseWriter, r *http.Request) {
	vv, err := venue.ListVenues()
	if err != nil {
		apierror(w, r, "Error Listing Venues: "+err.Error(), http.StatusInternalServerError)
	}
	var oo []venue.Venue
	for _, v := range vv {
		if len(v.Visits) < 1 {
			oo = append(oo, v)
		}
	}
	o := oo[rand.Intn(len(oo))]
	j, err := json.Marshal(&o)
	if err != nil {
		apierror(w, r, "Error marshalling Venue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getWeightedArray(venues []venue.Venue) []venue.Venue {
	var res []venue.Venue

	for _, v := range venues {
		lastvisit := 356
		daycount := 1
		rating := v.Rating
		if len(v.Visits) > 0 {
			lv := v.Visits[len(v.Visits)-1]
			dur := time.Now().Sub(lv)
			lastvisit = int(dur / (time.Hour * 24))
			if lastvisit < 0 {
				lastvisit = 1
			}
			daycount = len(v.Visits)
		}
		lastvisitf := float64(lastvisit) * criteriaweight.LastVisit
		daycountf := float64(daycount) * criteriaweight.DayCount
		ratingf := float64(rating) * criteriaweight.Rating
		ticketcount := int(math.Ceil((lastvisitf / daycountf) * ratingf))
		for i := 0; i < ticketcount; i++ {
			res = append(res, v)
		}
	}

	return res
}

func getNextVenuetoVisit(w http.ResponseWriter, r *http.Request) {
	new := !(strings.ToLower(r.FormValue("new")) == "")
	old := !(strings.ToLower(r.FormValue("old")) == "")
	weighted := !(strings.ToLower(r.FormValue("weighted")) == "")
	vv, err := venue.ListVenues()
	if err != nil {
		apierror(w, r, "Error Listing Venues: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var candiates []venue.Venue

	for _, v := range vv {
		if ((len(v.Visits) > 0) && old) || ((len(v.Visits) == 0) && new) {
			candiates = append(candiates, v)
		}
	}

	if len(candiates) < 1 {
		apierror(w, r, "No candidates to choose from", http.StatusInternalServerError)
		return
	}

	if weighted {
		candiates = getWeightedArray(candiates)
	}

	c := candiates[rand.Intn(len(candiates))]

	j, err := json.Marshal(&c)
	if err != nil {
		apierror(w, r, "Error marshalling Venue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func postUpdatefromPlaces(w http.ResponseWriter, r *http.Request) {
	vv, err := venue.ListVenues()
	if err != nil {
		apierror(w, r, "Error Listing Venues: "+err.Error(), http.StatusInternalServerError)
		return
	}
	for _, v := range vv {
		if v.GooglePlaceID == "" {
			continue
		}
		err := v.UpdateInfos()
		if err != nil {
			apierror(w, r, "Error Updating Venue: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = v.SavetoDataLocation()
		if err != nil {
			apierror(w, r, "Error Saving Venue: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

//SetWeights gets the criteria Weight from the config to save them for calulations later
func SetWeights(w config.Weight) {
	criteriaweight = w
}

func getAPIRouter(prefix string) *mux.Router {
	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	r.HandleFunc("/", mainAPIHandler)
	r.HandleFunc("/venue", postVenueAPIHandler).Methods("POST")
	r.HandleFunc("/venue/list", listVenuesAPIHandler).Methods("GET")
	r.HandleFunc("/venue/notvisited", getNotVisitedVenue).Methods("GET")
	r.HandleFunc("/venue/next", getNextVenuetoVisit)
	r.HandleFunc("/venue/updatefromplaces", postUpdatefromPlaces).Methods("POST")
	r.HandleFunc("/venue/getfromplaces/{query}", getVenueFromPlacesAPIHandler).Methods("GET")
	r.HandleFunc("/venue/{ID}", getVenueAPIHandler).Methods("GET")
	r.HandleFunc("/venue/{ID}", patchVenueAPIHander).Methods("PATCH")
	r.HandleFunc("/venue/{ID}", deleteVenueAPIHandler).Methods("DELETE")
	r.HandleFunc("/venue/{ID}/addvisits", addVisitAPIHandler).Methods("POST")

	return r
}
