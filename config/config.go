package config

var Settings *Config

func init() {

	Settings = &Config{}

}

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

	// settings for amazon lambda
	Lambda struct {
		Thumbnail struct {
			Endpoint string
			Key      string
		}
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

		// Set default posts per page
		PostsPerPage uint
		// Set default threads per index page
		ThreadsPerPage uint
		// Add one to number because first post is included
		PostsPerThread uint

		// Max request parameter input size
		ParamMaxSize uint
	}
}
