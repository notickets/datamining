package datamining

type Event struct {
	Name     string  `json:"eventName"`
	URL      string  `json:"eventURL"`
	Date     string  `json:"eventDate"`
	Venue    string  `json:"eventVenue"`
	Scene    *string `json:"eventScene,omitempty"`
	ImageURL *string `json:"eventImageURL,omitempty"`
}
