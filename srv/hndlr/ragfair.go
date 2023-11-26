package hndlr

import (
	"MT-GO/data"
	"MT-GO/pkg"
	"github.com/goccy/go-json"
	"log"
	"net/http"
)

type ESortType int8

const (
	ID ESortType = iota
	Priority
	Barter
	Rating
	OfferItem
	Price
	ExpirationDate
)

type ragfairOffers struct {
	Page              int            `json:"page"`
	Limit             int            `json:"limit"`
	SortType          ESortType      `json:"sortType"`
	SortDirection     int            `json:"sortDirection"`
	Currency          int            `json:"currency"`
	PriceFrom         int            `json:"priceFrom"`
	PriceTo           int            `json:"priceTo"`
	QuantityFrom      int            `json:"quantityFrom"`
	QuantityTo        int            `json:"quantityTo"`
	ConditionFrom     int            `json:"conditionFrom"`
	ConditionTo       int            `json:"conditionTo"`
	OneHourExpiration bool           `json:"oneHourExpiration"`
	RemoveBartering   bool           `json:"removeBartering"`
	OfferOwnerType    int            `json:"offerOwnerType"`
	OnlyFunctional    bool           `json:"onlyFunctional"`
	UpdateOfferCount  bool           `json:"updateOfferCount"`
	HandbookId        string         `json:"handbookId"`
	LinkedSearchId    string         `json:"linkedSearchId"`
	NeededSearchId    string         `json:"neededSearchId"`
	BuildItems        map[string]int `json:"buildItems"`
	BuildCount        int            `json:"buildCount"`
	Tm                int            `json:"tm"`
	Reload            int            `json:"reload"`
}

func RagfairFind(w http.ResponseWriter, r *http.Request) {
	ragfair := new(ragfairOffers)
	input, err := json.MarshalNoEscape(pkg.GetParsedBody(r))
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.UnmarshalNoEscape(input, &ragfair); err != nil {
		log.Fatalln(err)
	}

	flea := data.GetFlea()

	log.Println(routeNotImplemented)
	body := pkg.ApplyResponseBody(flea)
	pkg.SendZlibJSONReply(w, body)
}