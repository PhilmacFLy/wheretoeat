package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/philmacfly/wheretoeat/pkg/config"
	"github.com/philmacfly/wheretoeat/pkg/venue"
	"github.com/philmacfly/wheretoeat/pkg/web"
)

func main() {
	rand.Seed(time.Now().Unix())
	c, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	err = venue.SetupPlaceAPI(c.GoogleAPIKey)
	if err != nil {
		log.Fatal("Error setting up Places API:", err)
	}
	venue.SetDataFolder("data")
	web.SetWeights(c.Weight)
	r := web.SetupRouters("/")
	log.Fatal(http.ListenAndServe(c.Host+":"+strconv.Itoa(c.Port), r))
}
