package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/philmacfly/wheretoeat/pkg/venue"

	"github.com/gorilla/mux"
)

const errormessage = `<div class="alert alert-danger" role="alert">
  <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span>
  <span class="sr-only">Error:</span>
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
}

const (
	overviewActive int = 1 + iota
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

	err := sendHTTPRequest("GET", "venue/list", nil, &mp.Venues)
	if err != nil {
		mp.Default.Message = buildMessage(errormessage, "Error creating venue/list request: "+err.Error())
		showtemplate(w, tp, mp)
		return
	}
	showtemplate(w, tp, mp)
}

func venueUIViewHandler(w http.ResponseWriter, r *http.Request) {
	var vvp venueViewPage
	tp := "../../web/templates/venue/view.html"
	vvp.Default.Navbar = buildNavbar(overviewActive)
	vvp.Default.Pagename = "Venue View"

	id := r.FormValue("id")

	v := venue.Venue{VenueID: id}
	err := v.LoadFromDataLocation()
	if err != nil {
		vvp.Default.Message = buildMessage(errormessage, "Error loading venue: "+err.Error())
		showtemplate(w, tp, vvp)
		return
	}
	vvp.Venue = convertVenuetoWebVenue(v)
	showtemplate(w, tp, vvp)
}

func venueUIHandler(w http.ResponseWriter, r *http.Request) {
	a := r.FormValue("action")
	switch a {
	case "view":
		venueUIViewHandler(w, r)
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
