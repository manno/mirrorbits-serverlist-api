package mirror

type Mirror struct {
	ID             string
	HttpURL        string
	Latitude       float32
	Longitude      float32
	ContinentCode  string
	CountryCodes   string
	LastSync       int64
	Asnum          int
	SponsorURL     string
	SponsorLogoURL string
	SponsorName    string
	FileCount      int
}
