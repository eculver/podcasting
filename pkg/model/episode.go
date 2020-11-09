package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

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
