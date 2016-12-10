package config

// Settings holds an initialized settings with some sane defaults
var Settings = &Config{
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
}

// Config holds the main configuration data
type Config struct {
	General       General
	Prim          Prim
	CloudFlare    CloudFlare
	Akismet       Akismet
	StopForumSpam StopForumSpam
	Amazon        Amazon
	Limits        Limits
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
