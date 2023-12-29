package data

type MainData struct {
	ArtistInfo []Artist
	Search     SearchInfo
}

type GetAPI struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}

type Artist struct {
	Number     int      `json:"id"`
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	Members    []string `json:"members"`
	Date       int      `json:"creationDate"`
	FirstAlbum string   `json:"firstAlbum"`
	Relation   map[string][]string
	ConcertsMap map[string]string
}

type SearchInfo struct {
	Date []int
	Locs []string
}

type Relation struct {
	Index []struct {
		Id        int                 `json:"id"`
		Relations map[string][]string `json:"datesLocations"`
		ConcertsMap map[string]string
	} `json:"index"`
}

var (
	GetWebApiData GetAPI     // APIs
	RelationInfo  Relation   // dates-locations relation
	HomePageData  MainData   // for Home page
	Searching     SearchInfo // for Search suggestions
	ArtistInfo    []Artist   // Artists information
	SearchData    []Artist   // Artists information for Search results page
	YandexMapsApiKey = "ab421167-6629-47f8-b0b0-383cb58a54c7"
	Coordsmap  map[string]string
)
