# youtube : get youtube's video information

## Overview [![GoDoc](https://godoc.org/github.com/yudaprama/youtube?status.svg)](https://godoc.org/github.com/yudaprama/youtube) [![Go Report Card](https://goreportcard.com/badge/github.com/yudaprama/youtube)](https://goreportcard.com/report/github.com/yudaprama/youtube)

This package give you youtube video information such as title, author, quality, video type, and download/streaming url from given youtube video link. 

This won't download video. The video information can be used by frontend technology (mobile or web) to play/stream/download the video

## Install

```shell script
go get github.com/yudaprama/youtube
```

## Example

```go
import "yudaprama/youtube"

yv := "https://www.youtube.com/watch?v=f6kdp27TYZs"
infos, err := youtube.GetVideoInfo(yv)
if err != nil {
	fmt.Println(err)
}
// print as json
log.Printf("Videos: %s\n", infos)
// print each field
for _, info := range infos {
	fmt.Printf("Title [%s], Author [%s], Quality [%s] Type [%s], Streaming URL [%s]\n", info.Title, info.Author, info.Quality, info.Type, info.URL)
}

```

## Author

ToDo.

## License

MIT.
