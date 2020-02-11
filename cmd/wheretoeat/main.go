package main

import (
	"log"

	"github.com/kr/pretty"

	"github.com/philmacfly/wheretoeat/pkg/venue"
)

func main() {
	err := venue.SetupPlaceAPI("")
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
}
