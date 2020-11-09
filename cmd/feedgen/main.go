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

// Channel holds data about a podcast channel
type Channel struct {
	Title         string    `yaml:"title"`
	Subtitle      string    `yaml:"subtitle"`
	Author        string    `yaml:"author"`
	Description   string    `yaml:"description"`
	Created       time.Time `yaml:"created_time"`
	Language      string    `yaml:"language"`
	Copyright     string    `yaml:"copyright"`
	Category      string    `yaml:"category"`
	Subcategories []string  `yaml:"subcategories"`
	Explicit      bool      `yaml:"explicit"`
	Contact       string    `yaml:"contact_email"`
	WebURL        string    `yaml:"web_ref"`
	AtomURL       string    `yaml:"atom_ref"`
	ItunesURL     string    `yaml:"itunes_ref"`
	ImageURL      string    `yaml:"image_url"`
}

// Meta holds the metadata for an episode
type Meta struct {
	Index     int64     `yaml:"index"`
	Title     string    `yaml:"title"`
	Published time.Time `yaml:"published_time"`
	Audio     string    `yaml:"audio"`
	Thumb     string    `yaml:"thumb"`
	Notes     string    `yaml:"notes"`
	Teaser    string    `yaml:"teaser"`
}

// Episode has all the raw data associated with an episode
type Episode struct {
	Meta       Meta
	AudioPath  string
	ThumbPath  string
	NotesPath  string
	TeaserPath string

	WebURL   string
	AudioURL string
	ThumbURL string

	notesRaw  string
	teaserRaw string
}

// Hydrate reads files for episode content to populate raw data. For
// metadata, this means reading a yaml file and populating the Meta values.
// For notes and teaser paths, this means reading the file contents as markdown
// and translating to markup. If any of these operations fail, the whole thing will
// fail and a contextual error is returned.
func (e *Episode) Hydrate() error {
	// just make sure the audio and thumbnails exist
	if _, err := os.Stat(e.AudioPath); err != nil {
		return fmt.Errorf("could not stat audio file at '%s': %s", e.AudioPath, err)
	}
	if _, err := os.Stat(e.ThumbPath); err != nil {
		return fmt.Errorf("could not stat thumb file at '%s': %s", e.ThumbPath, err)
	}

	notesRaw, err := ioutil.ReadFile(e.NotesPath)
	if err != nil {
		return fmt.Errorf("unable to read notes at '%s': %s", e.NotesPath, err)
	}
	e.notesRaw = string(notesRaw)

	teaserRaw, err := ioutil.ReadFile(e.TeaserPath)
	if err != nil {
		return fmt.Errorf("unable to read notes at '%s': %s", e.TeaserPath, err)
	}
	e.teaserRaw = string(teaserRaw)

	return nil
}

// TeaserHTML renders the teaser content as HTML.
func (e *Episode) TeaserHTML() string {
	return "teaser teaser teaser"
}

func (e *Episode) String() string {
	fmtStr := `
Index: %s
Title: %s
Published: %s
Local:
	Audio: %s
	Thumb: %s
	Notes: %s
	Teaser: %s
Remote:
	Audio: %s
	Thumb: %s`
	return fmt.Sprintf(
		fmtStr,
		e.Meta.Index,
		e.Meta.Title,
		e.Meta.Published,
		e.AudioPath,
		e.ThumbPath,
		e.NotesPath,
		e.TeaserPath,
		e.AudioURL,
		e.ThumbURL,
	)
}

func main() {
	staticDir := "/Users/evan.culver/dev/src/github.com/eculver/tdtv2/static"
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

	episodes := []Episode{}

	// collect episode data from index
	if err := filepath.Walk(conf.IndexDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("could not read %q in content directory: %v", path, err)
		}

		if strings.HasSuffix(path, ".yml") {
			// try to read and hydrate meta
			var meta Meta
			metaRaw, err := ioutil.ReadFile(path)
			if err != nil {
				// can't read the file, so skip
				return nil
			}
			if err := yaml.Unmarshal(metaRaw, &meta); err != nil {
				// can't unmarshal file as yaml, so skip
			}

			audioPath := filepath.Join(conf.AudioDir, meta.Audio)
			thumbPath := filepath.Join(conf.ThumbDir, meta.Thumb)
			notesPath := filepath.Join(conf.NotesDir, meta.Notes)
			teaserPath := filepath.Join(conf.TeaserDir, meta.Teaser)

			audioURL := conf.AudioBaseURL + meta.Audio
			thumbURL := conf.ThumbBaseURL + meta.Thumb
			webURL := conf.WebBaseURL + "episode-slug"

			episode := Episode{
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

	// once episodes are enumerated...
	//   generate thumbnails
	//   generate feed
	//   upload content (audio, thumbs, feed, etc)
	//   verify?
	//   ping feedburner?

	// populate channel info
	channelPath := filepath.Join(staticDir, "channel.yml")
	var channel Channel

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
