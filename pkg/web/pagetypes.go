package web

import (
	"html/template"

	"github.com/philmacfly/wheretoeat/pkg/venue"
)

type defaultPage struct {
	Message  template.HTML
	Navbar   template.HTML
	Pagename string
}

type mainPage struct {
	Default defaultPage
	Venues  []venue.Venue
}
