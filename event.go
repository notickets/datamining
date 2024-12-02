package datamining

type Event struct {
	Name     string  `json:"Name"`
	URL      string  `json:"URL"`
	Date     string  `json:"Date"`
	Venue    string  `json:"Venue"`
	Scene    *string `json:"Scene,omitempty"`
	ImageURL *string `json:"ImageURL,omitempty"`
}
