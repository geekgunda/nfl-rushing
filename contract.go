package nflrushing

// Stat holds a single rushing stat with all related properties
type Stat struct {
	Player             string      `json:"Player"`
	Team               string      `json:"Team"`
	Position           string      `json:"Pos"`
	Attempts           int         `json:"Att"`
	AttemptsPerGameAvg float64     `json:"Att/G"`
	Yards              interface{} `json:"Yds"`
	YardsParsed        int         `json:"-"`
	AvgYardsPerAttempt float64     `json:"Yds/G"`
	Touchdowns         int         `json:"TD"`
	Longest            interface{} `json:"Lng"`
	LongestParsed      int         `json:"-"`
	IsLngTouchdown     bool        `json:"-"`
	FirstDown          int         `json:"1st"`
	PerFirstDown       float64     `json:"1st%"`
	Yards20            int         `json:"20+"`
	Yards40            int         `json:"40+"`
	Fumbles            int         `json:"FUM"`
}

// Stats holds multiple stat in a slice
type Stats []Stat

// Request holds filters and options used to trim the resultset
// as per user's request
type Request struct {
	PlayerFilter             string
	SortFilter               string
	OrderFilter              string
	ResponseFilter           string
	IsDefaultSortSelected    bool
	IsYardsSortSelected      bool
	IsLongestSortSelected    bool
	IsTouchdownsSortSelected bool
	IsOrderAsc               bool
	IsOrderDesc              bool
}

// TemplateData holds data used by template to populate webpage
type TemplateData struct {
	Request
	Stats
}
