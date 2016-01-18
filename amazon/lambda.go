package amazon

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// the format for the lambda thumbnail
type LambdaThumbnail struct {
	Bucket    string `json:"bucket"`
	Filename  string `json:"filename"`
	Thumbnail string `json:"thumbnail"`
	MaxWidth  int    `json:"max_width"`
	MaxHeight int    `json:"max_height"`
}

// the format of the lambda context response
type LambdaResponse struct {
	Success string `json:"successMessage"`
	Error   string `json:"errorMessage"`
	Width   int    `json:"thumbWidth"`
	Height  int    `json:"thumbHeight"`
}

// posts to our api endpoint
func (a *Amazon) Execute(data LambdaThumbnail) (width, height int, err error) {

	// Marshal the structs into JSON
	output, err := json.Marshal(data)
	if err != nil {
		return
	}

	svc := lambda.New(a.session)

	// params for lambda invocation
	params := &lambda.InvokeInput{
		FunctionName:   aws.String("resize_image"),
		InvocationType: aws.String(lambda.InvocationTypeRequestResponse),
		Payload:        output,
	}

	// invoke lambda function
	resp, err := svc.Invoke(params)
	if err != nil {
		return
	}

	response := LambdaResponse{}

	err = json.Unmarshal(resp.Payload, &response)
	if err != nil {
		return
	}

	if response.Error != "" {
		err = errors.New(fmt.Sprintf("Error creating thumbnail: %s", response.Error))
		return
	}

	// return our values
	width = response.Width
	height = response.Height

	return

}
