package web

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/philmacfly/wheretoeat/pkg/venue"
)

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
	fmt.Println(sb)
	vv, err := venue.ListVenues()
	if err != nil {
		apierror(w, r, "Error Listing Venues: "+err.Error(), http.StatusInternalServerError)
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

func getAPIRouter(prefix string) *mux.Router {
	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	r.HandleFunc("/", mainAPIHandler)
	r.HandleFunc("/venue", postVenueAPIHandler).Methods("POST")
	r.HandleFunc("/venue/list", listVenuesAPIHandler).Methods("GET")
	r.HandleFunc("/venue/notvisited", getNotVisitedVenue).Methods("GET")
	r.HandleFunc("/venue/getfromplaces/{query}", getVenueFromPlacesAPIHandler).Methods("GET")
	r.HandleFunc("/venue/{ID}", getVenueAPIHandler).Methods("GET")
	r.HandleFunc("/venue/{ID}", patchVenueAPIHander).Methods("PATCH")
	r.HandleFunc("/venue/{ID}", deleteVenueAPIHandler).Methods("DELETE")
	r.HandleFunc("/venue/{ID}/addvisits", addVisitAPIHandler).Methods("POST")

	return r
}
