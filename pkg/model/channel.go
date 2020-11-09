package model

import "time"

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
