package main

import (
	"log"
	"net/http"
	"strconv"

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
	r := web.SetupRouters("/", "data")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(c.Port), r))
}
