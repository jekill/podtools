package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"net/url"

	"log"

	"github.com/jbub/podcasts"
)

func main() {

	dirName := flag.String("dir", ".", "directory with mp3 files")
	podcastName := flag.String("podcast", "Autogen", "podcast name")
	baseUrl := flag.String("base-url", "", "Base podcast url. (Example: http://example.com/mypodcast)")
	feedFileName := flag.String("feed-file", "feed.xml", "Feed file name (Example: feed.xml)")

	flag.Parse()

	if *baseUrl == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	dir, err := os.Open(*dirName)
	defer dir.Close()

	if err != nil {
		fmt.Println("Err: open error")
		return
	}

	stat, err := dir.Stat()
	if err != nil {
		fmt.Println("Err: stat error")
		return
	}

	if !stat.IsDir() {
		fmt.Println("Err: it is not directory")
		return
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Err: read directory error")
	}

	feed := createFeed(*podcastName, *baseUrl)

	for _, fi := range files {
		if !strings.HasSuffix(fi.Name(), ".mp3") {
			continue
		}

		fileToFeed(feed, fi)
	}

	writeFeedToFile(*feedFileName, feed)
}

func writeFeedToFile(fileName string, feed *podcasts.Podcast) {
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		fmt.Println("Err: feed writing error")
		return
	}

	podcastFeed, err := feed.Feed()
	if err != nil {
		fmt.Println("Err: feed writing error")
		return
	}

	podcastFeed.Write(file)

}

func createFeed(feedName string, baseURL string) *podcasts.Podcast {
	feed := &podcasts.Podcast{
		Title:       feedName,
		Description: "none",
		Link:        baseURL,
	}
	return feed
}

func fileToFeed(feed *podcasts.Podcast, mediaFile os.FileInfo) {
	fmt.Println(mediaFile.Name())
	u, err := url.Parse(fmt.Sprint(feed.Link, "/", mediaFile.Name()))
	if err != nil {
		log.Fatal("url Parse")
	}
	url := u.String()
	fmt.Println("URL:", url)
	item := &podcasts.Item{
		Title: mediaFile.Name(),
		Enclosure: &podcasts.Enclosure{
			URL:    url,
			Length: fmt.Sprint(mediaFile.Size()),
			Type:   "MP3",
		},
		GUID:    url,
		PubDate: &podcasts.PubDate{time.Now()},
	}

	feed.AddItem(item)
}
