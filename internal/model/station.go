package model

type Station struct {
	Name     string `json:"name"`
	URL      string `json:"url_resolved"`
	Country  string `json:"country"`
	Tags     string `json:"tags"`
	Bitrate  int    `json:"bitrate"`
}
