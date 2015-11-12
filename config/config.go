package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Settings *Config

type Config struct {
	General struct {
		GuestPosting     bool
		AutoRegistration bool
	}

	Akismet struct {
		// Akismet settings
		Key  string
		Host string
	}

	StopForumSpam struct {
		// Stop Forum Spam settings
		Confidence float64
	}

	// settings for amazon s3
	Amazon struct {
		Region string
		Bucket string
		Id     string
		Key    string
	}

	// settings for google storage
	Google struct {
		Auth   string
		Bucket string
		Key    string
	}

	Antispam struct {
		// Antispam Key from Prim
		AntispamKey string

		// Antispam cookie
		CookieName  string
		CookieValue string
	}

	Limits struct {
		// Image settings
		ImageMinWidth  int
		ImageMinHeight int
		ImageMaxWidth  int
		ImageMaxHeight int
		ImageMaxSize   int
		WebmMaxLength  int

		// Max posts in a thread
		PostsMax uint

		// Lengths for posting
		CommentMaxLength int
		CommentMinLength int
		TitleMaxLength   int
		TitleMinLength   int
		NameMaxLength    int
		NameMinLength    int
		TagMaxLength     int
		TagMinLength     int

		// Max thumbnail sizes
		ThumbnailMaxWidth  int
		ThumbnailMaxHeight int

		// Max request parameter input size
		ParamMaxSize uint
	}
}

func Print() {

	// Marshal the structs into JSON
	output, err := json.MarshalIndent(Settings, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", output)

}

func init() {

	Settings = &Config{}

}
