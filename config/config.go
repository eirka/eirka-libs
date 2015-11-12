package config

import (
	"fmt"
	"strings"
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
	fmt.Println(strings.Repeat("*", 60))
	fmt.Printf("%-20v\n\n", "Global Settings")
	fmt.Println(strings.Repeat("*", 60))
	fmt.Printf("%-20v\n\n", "General")
	fmt.Printf("%-20v%40v\n", "Guest Posting", Settings.General.GuestPosting)
	fmt.Printf("%-20v%40v\n", "Auto Registration", Settings.General.AutoRegistration)
	fmt.Println(strings.Repeat("*", 60))
	fmt.Printf("%-20v\n\n", "Antispam")
	fmt.Printf("%-20v%40v\n", "Key", Settings.Antispam.AntispamKey)
	fmt.Printf("%-20v%40v\n", "Cookie Name", Settings.Antispam.CookieName)
	fmt.Printf("%-20v%40v\n", "Cookie Value", Settings.Antispam.CookieValue)
	fmt.Println(strings.Repeat("*", 60))
	fmt.Printf("%-20v\n\n", "Limits")
	fmt.Printf("%-20v\n\n", "Images")
	fmt.Printf("%-20v%40v\n", "Min Width", Settings.Limits.ImageMinWidth)
	fmt.Printf("%-20v%40v\n", "Min Height", Settings.Limits.ImageMinHeight)
	fmt.Printf("%-20v%40v\n", "Max Width", Settings.Limits.ImageMaxWidth)
	fmt.Printf("%-20v%40v\n", "Max Height", Settings.Limits.ImageMaxHeight)
	fmt.Printf("%-20v%40v\n", "Max Size", Settings.Limits.ImageMaxSize)
	fmt.Printf("%-20v%40v\n", "WebM Length", Settings.Limits.WebmMaxLength)
	fmt.Printf("%-20v%40v\n", "Thumb Max Width", Settings.Limits.ThumbnailMaxWidth)
	fmt.Printf("%-20v%40v\n", "Thumb Max Height", Settings.Limits.ThumbnailMaxHeight)
	fmt.Printf("\n%-20v\n\n", "Posting")
	fmt.Printf("%-20v%40v\n", "Thread Max Posts", Settings.Limits.PostsMax)
	fmt.Printf("%-20v%40v\n", "Comment Max", Settings.Limits.CommentMaxLength)
	fmt.Printf("%-20v%40v\n", "Comment Min", Settings.Limits.CommentMinLength)
	fmt.Printf("%-20v%40v\n", "Title Max", Settings.Limits.TitleMaxLength)
	fmt.Printf("%-20v%40v\n", "Title Min", Settings.Limits.TitleMinLength)
	fmt.Printf("%-20v%40v\n", "Name Max", Settings.Limits.NameMaxLength)
	fmt.Printf("%-20v%40v\n", "Name Min", Settings.Limits.NameMinLength)
	fmt.Printf("%-20v%40v\n", "Tag Max", Settings.Limits.TagMaxLength)
	fmt.Printf("%-20v%40v\n", "Tag Min", Settings.Limits.TagMinLength)
	fmt.Printf("%-20v%40v\n", "Max Param Size", Settings.Limits.ParamMaxSize)
	fmt.Println(strings.Repeat("*", 60))
	fmt.Printf("%-20v\n\n", "Akismet")
	fmt.Printf("%-20v%40v\n", "Key", Settings.Akismet.Key)
	fmt.Printf("%-20v%40v\n", "Host", Settings.Akismet.Host)
	fmt.Println(strings.Repeat("*", 60))
	fmt.Printf("%-20v\n\n", "Stop Forum Spam")
	fmt.Printf("%-20v%40v\n", "Confidence", Settings.StopForumSpam.Confidence)
	fmt.Println(strings.Repeat("*", 60))
	fmt.Printf("%-20v\n\n", "Amazon")
	fmt.Printf("%-20v%40v\n", "Region", Settings.Amazon.Region)
	fmt.Printf("%-20v%40v\n", "Bucket", Settings.Amazon.Bucket)
	fmt.Printf("%-20v%40v\n", "Id", Settings.Amazon.Id)
	fmt.Printf("%-20v%40v\n", "Key", Settings.Amazon.Key)
	fmt.Println(strings.Repeat("*", 60))
	fmt.Printf("%-20v\n\n", "Google")
	fmt.Printf("%-20v%40v\n", "Auth", Settings.Google.Auth)
	fmt.Printf("%-20v%40v\n", "Bucket", Settings.Google.Bucket)
	fmt.Printf("%-20v%40v\n", "Key", Settings.Google.Key)
	fmt.Println(strings.Repeat("*", 60))
}

func init() {

	Settings = &Config{}

}
