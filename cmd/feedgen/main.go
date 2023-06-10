package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/eculver/tdtv2/pkg/model"
	"github.com/eduncan911/podcast"
)

// Config configures the feed generator
type Config struct {
	IndexDir    string
	TemplateDir string
	AudioDir    string
	ThumbDir    string
	NotesDir    string
	TeaserDir   string

	WebBaseURL   string
	AudioBaseURL string
	ThumbBaseURL string
}

func main() {
	// TODO: make this configurable via CLI options or static file config
	staticDir := "/Users/evan.culver/dev/src/github.com/eculver/tdtv2/content"
	conf := &Config{
		IndexDir:    staticDir + "/index",
		TemplateDir: staticDir + "/tmpl",
		AudioDir:    staticDir + "/media/audio",
		ThumbDir:    staticDir + "/media/thumb",
		NotesDir:    staticDir + "/media/notes",
		TeaserDir:   staticDir + "/media/teaser",

		WebBaseURL:   "http://thedrunkentaoist.com/episodes/",
		AudioBaseURL: "http://www.podtrac.com/pts/redirect.mp3/p.episodes.thedrunkentaoist.com/vod/tdt.distributedio/episodes/",
		ThumbBaseURL: "http://cdn.tdt.distributed.io/media/episodes/thumbs/",
	}

	episodes := []model.Episode{}

	// collect episode data from index
	if err := filepath.Walk(conf.IndexDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("could not read %q in content directory: %v", path, err)
		}

		if strings.HasSuffix(path, ".yml") {
			// try to read and hydrate meta
			var meta model.Meta
			metaRaw, err := ioutil.ReadFile(path)
			if err != nil {
				// can't read the file, so skip
				return nil
			}
			if err := yaml.Unmarshal(metaRaw, &meta); err != nil {
				// can't unmarshal file as yaml, so skip
				return nil
			}

			audioPath := filepath.Join(conf.AudioDir, meta.Audio)
			thumbPath := filepath.Join(conf.ThumbDir, meta.Thumb)
			notesPath := filepath.Join(conf.NotesDir, meta.Notes)
			teaserPath := filepath.Join(conf.TeaserDir, meta.Teaser)

			audioURL := conf.AudioBaseURL + meta.Audio
			thumbURL := conf.ThumbBaseURL + meta.Thumb
			// TODO: populate w/ real slug
			webURL := conf.WebBaseURL + "episode-slug"

			episode := model.Episode{
				AudioPath:  audioPath,
				ThumbPath:  thumbPath,
				NotesPath:  notesPath,
				TeaserPath: teaserPath,

				AudioURL: audioURL,
				ThumbURL: thumbURL,
				WebURL:   webURL,

				Meta: meta,
			}

			if err := episode.Hydrate(); err != nil {
				fmt.Printf("WARNING: could not hydrate episode: %s\n", err)
				return nil
			}
			// TODO: discover audio, thumb
			fmt.Printf("found episode: %s\n", episode.String())

			episodes = append(episodes, episode)
		}
		return nil
	}); err != nil {
		log.Fatalf("error walking the path %q: %v\n", conf.IndexDir, err)
	}

	// TODO: integrate with rest of workflow
	// once episodes are enumerated...
	//   [ ] generate thumbnails
	//   [ ] upload content (audio, thumbs, feed, etc)
	//   [x] generate feed
	//   [ ] generate home page
	//   [ ] generate archive
	//   [ ] generate episode page

	// TODO: bonus
	//   verify feed?
	//   ping feedburner?
	//   SSL

	// populate channel info
	channelPath := filepath.Join(staticDir, "channel.yml")
	var channel model.Channel

	channelRaw, err := ioutil.ReadFile(channelPath)
	if err != nil {
		log.Fatalf("could not read channel file %s: %s", channelPath, err)
	}
	if err := yaml.Unmarshal(channelRaw, &channel); err != nil {
		// can't unmarshal file as yaml, so skip
	}

	publishedTime := channel.Created
	updatedTime := time.Now()

	// instantiate a new Podcast
	p := podcast.New(
		channel.Title,
		channel.WebURL,
		channel.Description,
		&publishedTime,
		&updatedTime,
	)

	// add some channel properties
	p.AddSubTitle(channel.Subtitle)
	p.AddAuthor(channel.Author, channel.Contact)
	p.AddAtomLink(channel.AtomURL)
	p.AddImage(channel.ImageURL)
	p.AddSummary(channel.Description)
	p.IExplicit = "yes"

	for _, e := range episodes {
		// create an Item
		item := podcast.Item{
			Title:       fmt.Sprintf("Episode %d - %s", e.Meta.Index, e.Meta.Title),
			Link:        e.WebURL,
			Description: e.TeaserHTML(),
			PubDate:     &e.Meta.Published,
		}
		item.AddImage(e.ThumbURL)
		item.AddSummary(e.TeaserHTML())
		item.AddEnclosure(e.AudioURL, podcast.MP3, 55*(e.Meta.Index+1))

		// add the Item and check for validation errors
		if _, err := p.AddItem(item); err != nil {
			log.Fatalf("error adding episode %s: %s", item.Title, err)
		}
	}
	if err := p.Encode(os.Stdout); err != nil {
		log.Fatalf("could not encode: %s", err)
	}
}
