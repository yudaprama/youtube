package youtube

func GetVideoInfo(url string) (Videos, error) {
	y := NewYoutube()
	if err := y.DecodeURL(url); err != nil {
		return nil, err
	}
	return y.Videos, nil
}

func GetVideoInfoAsByte(url string) ([]byte, error) {
	info, err := GetVideoInfo(url)
	if err != nil {
		return nil, err
	}
	return info.Marshal()
}

func GetVideoInfoAsJSON(url string) (string, error) {
	asByte, err := GetVideoInfoAsByte(url)
	if err != nil {
		return "", err
	}
	return string(asByte), nil
}

func GetVideoInfoAsJSONWithoutErr(url string) string {
	str, err := GetVideoInfoAsJSON(url)
	if err != nil {
		return ""
	}
	return str
}
