package main

import "time"

//use https://mholt.github.io/json-to-go/ to generate the struct

type HostSearchResult struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Result Result `json:"result"`
}
type Services struct {
	Port              int    `json:"port"`
	ServiceName       string `json:"service_name"`
	TransportProtocol string `json:"transport_protocol"`
}
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type Location struct {
	Continent             string      `json:"continent"`
	Country               string      `json:"country"`
	CountryCode           string      `json:"country_code"`
	City                  string      `json:"city"`
	PostalCode            string      `json:"postal_code"`
	Timezone              string      `json:"timezone"`
	Province              string      `json:"province"`
	Coordinates           Coordinates `json:"coordinates"`
	RegisteredCountry     string      `json:"registered_country"`
	RegisteredCountryCode string      `json:"registered_country_code"`
}
type AutonomousSystem struct {
	Asn         int    `json:"asn"`
	Description string `json:"description"`
	BgpPrefix   string `json:"bgp_prefix"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
}
type Hits struct {
	IP               string           `json:"ip"`
	Services         []Services       `json:"services"`
	Location         Location         `json:"location"`
	AutonomousSystem AutonomousSystem `json:"autonomous_system"`
	LastUpdatedAt    time.Time        `json:"last_updated_at"`
}
type Links struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
}
type Result struct {
	Query    string `json:"query"`
	Total    int    `json:"total"`
	Duration int    `json:"duration"`
	Hits     []Hits `json:"hits"`
	Links    Links  `json:"links"`
}
