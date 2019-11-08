package amazon

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/eirka/eirka-libs/config"
)

// Amazon holds the session
type Amazon struct {
	session *session.Session
}

// New will create aws session with credentials
func New() (amazon *Amazon) {

	// new credentials from settings
	creds := credentials.NewStaticCredentials(config.Settings.Amazon.ID, config.Settings.Amazon.Key, "")

	// create our session
	// TODO: handle error
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(config.Settings.Amazon.Region),
		Credentials: creds,
		MaxRetries:  aws.Int(10),
	})

	return &Amazon{
		session: sess,
	}

}
