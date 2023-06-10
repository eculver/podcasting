package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type Podcast struct {
	Channel struct {
		Title string `xml:"title"`
		Items []struct {
			Title string `xml:"title"`
			PubDate string `xml:"pubDate"`
			Guid string `xml:"guid"`
		} `xml:"item"`
	} `xml:"channel"`
}

type Episode struct {
	Title string
	PubDate string
	Guid string
}

type Site struct {
	Episodes map[string]Episode
}

type Episode struct {
	Guid string
	Title   string
	PubDate string
}

type ByPubDate []Episode

func (a ByPubDate) Len() int           { return len(a) }
func (a ByPubDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPubDate) Less(i, j int) bool { return strings.Compare(a[i].PubDate, a[j].PubDate) > 0 }

func main() {
	/*
		examples:
			sitegen new episode/121
			sitegen sync https://anchor.fm/s/b781db40/podcast/rss
	*/
	numEpisodes := flag.Int("n", 10, "number of episodes to show")
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: rss-parser [OPTIONS] RSS_FEED_URL")
		return
	}

	url := args[0]
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	var podcast Podcast
	xml.Unmarshal(body, &podcast)
	var episodes []Episode
	for _, item := range podcast.Channel.Items {
		episodes = append(episodes, Episode{item.Title, item.PubDate})
	}
	sort.Sort(ByPubDate(episodes))

	fmt.Println("Podcast:", podcast.Channel.Title)
	for i := 0; i < *numEpisodes && i < len(episodes); i++ {
		fmt.Println("Title:", episodes[i].Title)
		fmt.Println("Publish Date:", episodes[i].PubDate)
		fmt.Println()
	}
}

func NewEpisodeCommand(c *cli.Context) error {
	if len(c.Args().Slice()) != 1 {
		return fmt.Errorf("Usage: tdtctl new EPISODE_ID")
	}

	id := c.Args().First()
	if _, ok := site.Episodes[id]; ok {
		return fmt.Errorf("Error: Episode with ID '%s' already exists", id)
	}

	site.Episodes[id] = Episode{
		Title:    "New Episode",
		PubDate:  "",
		Guid:     "",
	}
	return saveSite(site)
}

func SyncCommand(c *cli.Context) error {
	if len(c.Args().Slice()) != 1 {
		return fmt.Errorf("Usage: tdtctl sync RSS_FEED_URL")
	}

	url := c.Args().First()
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error fetching data: %s", err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error reading data: %s", err)
	}

	var podcast Podcast
	xml.Unmarshal(body, &podcast)
	for _, item := range podcast.Channel.Items {
		id := strings.ReplaceAll(item.Guid, "/", "_")
		if _, ok := site.Episodes[id]; !ok {
			site.Episodes[id] = Episode{
				Title:    item.Title,
				PubDate:  item.PubDate,
				Guid:     item.Guid,
			}
		}
	}

	return saveSite(site)
}

func GenCommand(c *cli.Context) error {
	indexPath := c.String("index")
	outPath := c.String("out")
	if indexPath == "" || outPath == "" {
		return fmt.Errorf("Usage: tdtctl gen -index INDEX_PATH -out OUT_PATH")
