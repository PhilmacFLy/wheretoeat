package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kr/pretty"

	"github.com/philmacfly/wheretoeat/pkg/config"
	"github.com/philmacfly/wheretoeat/pkg/venue"
	"github.com/philmacfly/wheretoeat/pkg/web"
)

func main() {
	c, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	err = venue.SetupPlaceAPI(c.GoogleAPIKey)
	if err != nil {
		log.Fatal("Error setting up Places API:", err)
	}
	v, err := venue.GetVenubyPlaceSearch("Pizza Latina")
	if err != nil {
		log.Fatal("Error getting place", err)
	}
	pretty.Println(v)
	v.SavetoFile("test.json")
	var v2 venue.Venue
	v2.LoadfromFile("test.json")
	pretty.Println(v2)
	r := web.GetRouters("/")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(c.Port), r))
}
