package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/philmacfly/wheretoeat/pkg/venue"

	"github.com/gorilla/mux"
)

const errormessage = `<div class="alert alert-danger" role="alert">
  <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span>
  <span class="sr-only">Error:</span>
  $MESSAGE$
</div>`

const successmessage = `<div class="alert alert alert-success" role="alert">
  <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span>
  <span class="sr-only">Success:</span>
  $MESSAGE$
</div>`

var navitems [][]template.HTML

func createNavitem(name string, link string) []template.HTML {
	active := `<li class="nav-item active"><a class="nav-link" href="/` + link + `">` + name + ` <span class="sr-only">(current)</span></a></li>`
	inactive := `<li class="nav-item"><a class="nav-link" href="/` + link + `">` + name + `</a></li>`
	activehtml := template.HTML(active)
	inactivehtml := template.HTML(inactive)
	return []template.HTML{activehtml, inactivehtml}
}

func init() {
	navitems = append(navitems, createNavitem("Overview", ""))
	navitems = append(navitems, createNavitem("Add Venue", "ui/venue/?action=add"))
	navitems = append(navitems, createNavitem("Select New Venue", "ui/venue/?action=not-visited"))
}

const (
	overviewActive int = 1 + iota
	addvenueActive
	notVisitedActive
)

func buildNavbar(item int) template.HTML {
	var res template.HTML
	res = res + `<div class="collapse navbar-collapse" id="navbarCollapse">`
	res = res + `<ul class="navbar-nav mr-auto">`
	for i, n := range navitems {
		var add template.HTML
		if i+1 == item {
			add = n[0]
		} else {
			add = n[1]
		}
		res = res + add + "\n"
	}
	res = res + `</ul>`
	res = res + `</div>`
	return res
}

func buildMessage(tp string, message string) template.HTML {
	message = html.EscapeString(message)
	return template.HTML(strings.Replace(tp, "$MESSAGE$", message, -1))
}

func getNewHTTPRequest(method string, endpoint string, in io.Reader) (*http.Request, error) {
	var req *http.Request
	var err error
	url := "http://127.0.0.1:4334/api/" + endpoint
	fmt.Println(url)
	req, err = http.NewRequest(method, url, in)
	if err != nil {
		return req, err
	}
	return req, nil
}

func sendHTTPRequest(method string, endpoint string, in io.Reader, v interface{}) error {
	req, err := getNewHTTPRequest(method, endpoint, in)
	if err != nil {
		return errors.New("Error creating request: " + err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Error executing request: " + err.Error())
	}

	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		fmt.Println(resp.StatusCode)
		decoder := json.NewDecoder(resp.Body)
		var er errorResponse
		err = decoder.Decode(&er)
		if err != nil {
			data, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(data))
			return errors.New("Error while decoding Error response. Only God can help you now:" + err.Error())
		}
		return errors.New("Got negativ status code: " + er.Errormessage)
	}

	if v != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(v)
		if err != nil {
			return errors.New("Unable to decode on given Interface: " + err.Error())
		}
	}
	return nil
}

func showtemplate(w http.ResponseWriter, path string, data interface{}) {
	t, err := template.ParseFiles(path)
	if err != nil {
		fmt.Fprintln(w, "Error parsing template:", err)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Fprintln(w, "Error executing template:", err)
		return
	}
}

func mainUIHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "venue/", http.StatusSeeOther)
}

func venueUIListHandler(w http.ResponseWriter, r *http.Request) {
	var mp mainPage
	tp := "../../web/templates/main.html"
	mp.Default.Navbar = buildNavbar(overviewActive)
	mp.Default.Pagename = "Venue List"

	var vv []venue.Venue

	err := sendHTTPRequest("GET", "venue/list", nil, &vv)
	if err != nil {
		mp.Default.Message = buildMessage(errormessage, "Error creating venue/list request: "+err.Error())
		showtemplate(w, tp, mp)
		return
	}
	for _, v := range vv {
		mp.Venues = append(mp.Venues, convertVenuetoWebVenue(v))
	}
	showtemplate(w, tp, mp)
}

func venueUIViewHandler(w http.ResponseWriter, r *http.Request) {
	var vvp venueViewPage
	tp := "../../web/templates/venue/view.html"
	vvp.Default.Navbar = buildNavbar(overviewActive)
	vvp.Default.Pagename = "Venue View"

	id := r.FormValue("id")

	v := venue.Venue{}

	err := sendHTTPRequest("GET", "venue/"+id, nil, &v)
	if err != nil {
		vvp.Default.Message = buildMessage(errormessage, "Error getting venue request: "+err.Error())
		showtemplate(w, tp, vvp)
		return
	}
	vvp.Venue = convertVenuetoWebVenue(v)
	showtemplate(w, tp, vvp)
}

func venueUIAddHandler(w http.ResponseWriter, r *http.Request) {
	var vap venueAddPage
	tp := "../../web/templates/venue/add.html"
	vap.Default.Navbar = buildNavbar(addvenueActive)
	vap.Default.Pagename = "Add Venue"

	name := r.FormValue("Name")
	address := r.FormValue("Address")
	v := venue.Venue{}
	if name != "" {
		query := name
		if address != "" {
			query = query + ", " + address
		}
		err := sendHTTPRequest("GET", "venue/getfromplaces/"+query, nil, &v)
		if err != nil {
			vap.Default.Message = buildMessage(errormessage, "Error getting venue request: "+err.Error())
			showtemplate(w, tp, vap)
			return
		}
	}

	vap.Venue = convertVenuetoWebVenue(v)
	showtemplate(w, tp, vap)
}

func venueUISaveHandler(w http.ResponseWriter, r *http.Request) {
	var vap venueAddPage
	tp := "../../web/templates/venue/add.html"
	vap.Default.Navbar = buildNavbar(addvenueActive)
	vap.Default.Pagename = "Add Venue"

	var wv webVenue
	wv.Name = r.FormValue("Name")
	wv.Address = r.FormValue("Address")
	rating := r.FormValue("Rating")
	ra, _ := strconv.Atoi(rating)
	wv.Rating = ra
	wv.GooglePlaceID = r.FormValue("placesid")
	wv.Website = r.FormValue("Website")
	wv.PhoneNumber = r.FormValue("phone")
	wv.Notes = r.FormValue("Notes")
	wv.OpeningHours.Monday = r.FormValue("Monday")
	wv.OpeningHours.Tuesday = r.FormValue("Tuesday")
	wv.OpeningHours.Wednesday = r.FormValue("Wednesday")
	wv.OpeningHours.Thursday = r.FormValue("Thursday")
	wv.OpeningHours.Friday = r.FormValue("Friday")
	wv.OpeningHours.Saturday = r.FormValue("Monday")
	wv.OpeningHours.Sunday = r.FormValue("Sunday")

	v, err := convertWebVenuetoVenue(wv)
	if err != nil {
		vap.Default.Message = buildMessage(errormessage, "Error converting venue: "+err.Error())
		vap.Venue = wv
		showtemplate(w, tp, vap)
		return
	}

	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(v)

	err = sendHTTPRequest("POST", "venue", b, &v)
	if err != nil {
		vap.Default.Message = buildMessage(errormessage, "Error sending Venue request: "+err.Error())
		vap.Venue = wv
		showtemplate(w, tp, vap)
		return
	}
	vap.Default.Message = buildMessage(successmessage, "New Venue successfully added")
	vap.Venue = webVenue{}
	showtemplate(w, tp, vap)
	return
}

func venueUINotVisitedHandler(w http.ResponseWriter, r *http.Request) {
	var mp mainPage
	tp := "../../web/templates/main.html"
	mp.Default.Navbar = buildNavbar(overviewActive)
	mp.Default.Pagename = "Venue List"

	v := venue.Venue{}

	err := sendHTTPRequest("GET", "venue/notvisited", nil, &v)
	if err != nil {
		mp.Default.Message = buildMessage(errormessage, "Error getting not visited venue request: "+err.Error())
		showtemplate(w, tp, mp)
		return
	}
	http.Redirect(w, r, "?action=view&id="+v.VenueID, http.StatusTemporaryRedirect)

}

func venueUIAddVisitHandler(w http.ResponseWriter, r *http.Request) {
	var vap venueAddVisitPage
	tp := "../../web/templates/venue/add-visit.html"
	vap.Default.Navbar = buildNavbar(overviewActive)
	vap.Default.Pagename = "Venue Add Visit"

	id := r.FormValue("id")

	v := venue.Venue{}

	err := sendHTTPRequest("GET", "venue/"+id, nil, &v)
	if err != nil {
		vap.Default.Message = buildMessage(errormessage, "Error getting venue request: "+err.Error())
		showtemplate(w, tp, vap)
		return
	}
	vap.Venue = convertVenuetoWebVenue(v)

	vap.Venue.LastVisit = time.Now().Format(layoutISO)

	showtemplate(w, tp, vap)
}

func venueUIAddVisitExecuteHandler(w http.ResponseWriter, r *http.Request) {
	var vap venueAddVisitPage
	tp := "../../web/templates/venue/add-visit.html"
	vap.Default.Navbar = buildNavbar(overviewActive)
	vap.Default.Pagename = "Venue Add Visit"

	id := r.FormValue("id")
	date := r.FormValue("date")

	d, err := time.Parse(layoutISO, date)
	if err != nil {
		vap.Default.Message = buildMessage(errormessage, "Error parsing given date: "+err.Error())
		showtemplate(w, tp, vap)
		return
	}

	var req addVisitsRequest

	req.Visits = append(req.Visits, d)

	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(req)

	err = sendHTTPRequest("POST", "venue/"+id+"/addvisits", b, nil)
	if err != nil {
		vap.Default.Message = buildMessage(errormessage, "Error getting venue request: "+err.Error())
		showtemplate(w, tp, vap)
		return
	}
	http.Redirect(w, r, "?action=view&id="+id, http.StatusTemporaryRedirect)
}

func venueUIUpdateVenuesfromPlacesHandler(w http.ResponseWriter, r *http.Request) {
	var udp updateDonePage
	tp := "../../web/templates/venue/update-done.html"
	udp.Default.Navbar = buildNavbar(overviewActive)
	udp.Default.Pagename = "Update Done"

	err := sendHTTPRequest("POST", "venue/updatefromplaces", nil, nil)
	if err != nil {
		udp.Default.Message = buildMessage(errormessage, "Error getting venue request: "+err.Error())
		showtemplate(w, tp, udp)
		return
	}
	udp.Default.Message = buildMessage(successmessage, "All Venues updated")
	showtemplate(w, tp, udp)
}

func venueUIHandler(w http.ResponseWriter, r *http.Request) {
	a := r.FormValue("action")
	switch a {
	case "view":
		venueUIViewHandler(w, r)
	case "add":
		venueUIAddHandler(w, r)
	case "save":
		venueUISaveHandler(w, r)
	case "not-visited":
		venueUINotVisitedHandler(w, r)
	case "add-visit":
		venueUIAddVisitHandler(w, r)
	case "add-visit-execute":
		venueUIAddVisitExecuteHandler(w, r)
	case "update-from-places":
		venueUIUpdateVenuesfromPlacesHandler(w, r)
	default:
		venueUIListHandler(w, r)
	}
}

func getUIRouter(prefix string) *mux.Router {
	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/ui/static/", http.FileServer(http.Dir("../../web/static/"))))
	r.HandleFunc("/", mainUIHandler)
	r.HandleFunc("/venue/", venueUIHandler)
	return r
}
