package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Youtube is object for youtube
type Youtube struct {
	Videos    Videos
	VideoID   string
	videoInfo string
}

// NewYoutube will initialize youtube object
func NewYoutube() *Youtube {
	return &Youtube{}
}

// DecodeURL will decode youtube URL to retrieval video information.
func (y *Youtube) DecodeURL(url string) error {
	err := y.findVideoID(url)
	if err != nil {
		return fmt.Errorf("findVideoID error=%s", err)
	}
	err = y.getVideoInfo()
	if err != nil {
		return fmt.Errorf("getVideoInfo error=%s", err)
	}
	err = y.parseVideoInfo()
	if err != nil {
		return fmt.Errorf("parse video info failed, err=%s", err)
	}
	return nil
}

func (y *Youtube) parseVideoInfo() error {
	answer, err := url.ParseQuery(y.videoInfo)
	if err != nil {
		return err
	}
	status, ok := answer["status"]
	if !ok {
		err = fmt.Errorf("got no response status")
		return err
	}
	if status[0] == "fail" {
		reason, ok := answer["reason"]
		if ok {
			err = fmt.Errorf("'fail', reason: '%s'", reason[0])
		} else {
			err = errors.New(fmt.Sprint("'fail', no reason given"))
		}
		return err
	}
	if status[0] != "ok" {
		err = fmt.Errorf("non-success response (status: '%s')", status)
		return err
	}
	// read the streams map
	streamMap, ok := answer["player_response"]
	if !ok {
		err = errors.New(fmt.Sprint("no stream map found."))
		return err
	}
	// Get video title and author.
	title, author := getVideoTitleAuthor(answer)
	// get video info
	var vi VideoInfo
	if err := json.Unmarshal([]byte(streamMap[0]), &vi); err != nil {
		return errors.New("Player response json data has changed.")
	}
	// Get video download link
	if vi.PlayabilityStatus.Status == "UNPLAYABLE" {
		// Cannot playback on embedded video screen, could not download.
		return errors.New(fmt.Sprint("Cannot playback and download, reason:", vi.PlayabilityStatus.Reason))
	}
	var videos Videos
	for _, streamRaw := range vi.StreamingData.Formats {
		if streamRaw.MimeType == "" {
			// An error occurred while decoding one of the video's stream
			continue
		}
		streamUrl := streamRaw.URL
		if streamUrl == "" {
			cipher := streamRaw.Cipher
			decipheredUrl, err := y.decipher(cipher)
			if err != nil {
				return err
			}
			streamUrl = decipheredUrl
		}
		videos = append(videos, Video{
			Author:   author,
			Quality:  streamRaw.Quality,
			Title:    title,
			Type:     streamRaw.MimeType,
			URL:      streamUrl,
			Duration: streamRaw.ApproxDurationMs,
		})
	}
	y.Videos = videos
	if len(y.Videos) == 0 {
		return errors.New(fmt.Sprint("no stream list found"))
	}
	return nil
}

// setup http client
func GetHTTPClient() *http.Client {
	httpTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpClient := &http.Client{Transport: httpTransport}
	return httpClient
}

func (y *Youtube) getVideoInfo() error {
	id := "https://youtube.googleapis.com/v/" + y.VideoID
	u := "https://youtube.com/get_video_info?video_id=" + y.VideoID + "&eurl=" + id
	httpClient := GetHTTPClient()
	resp, err := httpClient.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	y.videoInfo = string(body)
	return nil
}

func (y *Youtube) findVideoID(url string) error {
	videoID := url
	if strings.Contains(videoID, "youtu") || strings.ContainsAny(videoID, "\"?&/<%=") {
		reList := []*regexp.Regexp{
			regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`([^"&?/=%]{11})`),
		}
		for _, re := range reList {
			if isMatch := re.MatchString(videoID); isMatch {
				subs := re.FindStringSubmatch(videoID)
				videoID = subs[1]
			}
		}
	}
	y.VideoID = videoID
	if strings.ContainsAny(videoID, "?&/<%=") {
		return errors.New("invalid characters in video id")
	}
	if len(videoID) < 10 {
		return errors.New("the video id must be at least 10 characters long")
	}
	return nil
}

func getVideoTitleAuthor(in url.Values) (string, string) {
	playResponse, ok := in["player_response"]
	if !ok {
		return "", ""
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(playResponse[0]), &m); err != nil {
		panic(err)
	}
	s := m["videoDetails"]
	maps := s.(map[string]interface{})
	// fmt.Println("-->", myMap["title"], "oooo:", myMap["author"])
	if title, ok := maps["title"]; ok {
		if author, ok := maps["author"]; ok {
			return title.(string), author.(string)
		}
	}
	return "", ""
}
