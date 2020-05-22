package youtube

import "encoding/json"

type Videos []Video

func UnmarshalVideos(data []byte) (Videos, error) {
	var r Videos
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Videos) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Video struct {
	Author   string `json:"author"`
	Quality  string `json:"quality"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	Duration string `json:"duration"`
}
