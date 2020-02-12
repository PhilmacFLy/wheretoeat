package main

import (
	"log"
	"net/http"

	"github.com/kr/pretty"

	"github.com/philmacfly/wheretoeat/pkg/venue"
	"github.com/philmacfly/wheretoeat/pkg/web"
)

func main() {
	err := venue.SetupPlaceAPI("AIzaSyC9El18L_ADXJ9Z_r_gzQK9Ez8Cn4j7_o4")
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
	log.Fatal(http.ListenAndServe(":4334", r))
}
