package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Global config file path
const configPath = "/etc/pram/pram.conf"

// Settings holds an initialized settings with some sane defaults
var Settings *Config

func init() {
	// Initialize default settings (these will be used if config file is not found)
	Settings = &Config{
		General: General{
			GuestPosting:     true,
			AutoRegistration: true,
		},
		Prim: Prim{
			CSS: "prim.css",
			JS:  "prim.js",
		},
		CloudFlare: CloudFlare{},
		Akismet:    Akismet{},
		StopForumSpam: StopForumSpam{
			Confidence: 40,
		},
		Amazon: Amazon{},
		Limits: Limits{
			ImageMinWidth:      100,
			ImageMinHeight:     100,
			ImageMaxWidth:      20000,
			ImageMaxHeight:     20000,
			ImageMaxSize:       20000000,
			AvatarMinWidth:     100,
			AvatarMinHeight:    100,
			AvatarMaxWidth:     1000,
			AvatarMaxHeight:    1000,
			AvatarMaxSize:      1000000,
			WebmMaxLength:      300,
			PostsMax:           800,
			CommentMaxLength:   1000,
			CommentMinLength:   3,
			TitleMaxLength:     40,
			TitleMinLength:     3,
			NameMaxLength:      20,
			NameMinLength:      3,
			TagMaxLength:       128,
			TagMinLength:       3,
			PasswordMaxLength:  128,
			PasswordMinLength:  8,
			ThumbnailMaxWidth:  200,
			ThumbnailMaxHeight: 300,
			PostsPerPage:       40,
			ThreadsPerPage:     10,
			PostsPerThread:     5,
			ParamMaxSize:       1000000,
		},
		Session: Session{
			OldSecret: "",
			NewSecret: "",
		},
	}

	// Try to load configuration from file
	LoadConfig()
}

// LoadConfig loads configuration from the config file
func LoadConfig() error {
	file, err := os.Open(configPath)
	if err != nil {
		// File not found, use default settings
		fmt.Printf("Config file not found at %s, using defaults\n", configPath)
		return err
	}
	defer file.Close()

	// Read the file content
	configData, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return err
	}

	// Create a temporary config to decode into
	tempConfig := &Config{}

	// Decode the JSON
	err = json.Unmarshal(configData, tempConfig)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		return err
	}

	// Update Settings with the loaded configuration
	Settings = tempConfig

	// Validate JWT secrets
	if Settings.Session.NewSecret == "" {
		fmt.Println("Warning: JWT NewSecret is empty in config file")
	}

	return nil
}

// Config holds the main configuration data
type Config struct {
	General       General
	Prim          Prim
	CloudFlare    CloudFlare
	Akismet       Akismet
	StopForumSpam StopForumSpam
	Scamalytics   Scamalytics
	Amazon        Amazon
	Limits        Limits
	Session       Session
}

// General options
type General struct {
	GuestPosting     bool
	AutoRegistration bool
}

// Prim holds asset names for Prim
type Prim struct {
	CSS string
	JS  string
}

// CloudFlare API settings
type CloudFlare struct {
	Configured bool
	Key        string
	Email      string
}

// Akismet settings
type Akismet struct {
	Configured bool
	Key        string
	Host       string
}

// StopForumSpam settings
type StopForumSpam struct {
	Confidence float64
}

// Scamalytics settings
type Scamalytics struct {
	Configured bool
	Key        string
	Endpoint   string
	Path       string
	Score      int
}

// Amazon holds API settings for Amazon
type Amazon struct {
	Configured bool
	Region     string
	Bucket     string
	ID         string
	Key        string
}

// Limits for various items
type Limits struct {
	// Image settings
	ImageMinWidth  int
	ImageMinHeight int
	ImageMaxWidth  int
	ImageMaxHeight int
	ImageMaxSize   int

	// avatar settings
	AvatarMinWidth  int
	AvatarMinHeight int
	AvatarMaxWidth  int
	AvatarMaxHeight int
	AvatarMaxSize   int

	// webm settings
	WebmMaxLength int

	// Max posts in a thread
	PostsMax uint

	// Lengths for posting
	CommentMaxLength int
	CommentMinLength int

	TitleMaxLength int
	TitleMinLength int

	NameMaxLength int
	NameMinLength int

	TagMaxLength int
	TagMinLength int

	PasswordMaxLength int
	PasswordMinLength int

	// Max thumbnail sizes
	ThumbnailMaxWidth  int
	ThumbnailMaxHeight int

	// Set default posts per page
	PostsPerPage uint
	// Set default threads per index page
	ThreadsPerPage uint
	// Add one to number because first post is included
	PostsPerThread uint

	// Max request parameter input size
	ParamMaxSize uint
}

// Session holds the secrets for JWT authentication
type Session struct {
	// OldSecret is used for validating existing tokens during rotation
	OldSecret string
	// NewSecret is used for signing new tokens and validating tokens
	NewSecret string
}
