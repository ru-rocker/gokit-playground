package lorem_grpc

import (
	gl "github.com/drhodes/golorem"
	"strings"
	"errors"
	"context"
)

var (
	ErrRequestTypeNotFound = errors.New("Request type only valid for word, sentence and paragraph")
)

// Define service interface
type Service interface {
	// generate a word with at least min letters and at most max letters.
	Lorem(ctx context.Context, requestType string, min, max int) (string, error)
}

// Implement service with empty struct
type LoremService struct {

}

// Implement service functions
func (LoremService) Lorem(_ context.Context, requestType string, min, max int) (string, error) {
	var result string
	var err error
	if strings.EqualFold(requestType, "Word") {
		result = gl.Word(min, max)
	} else if strings.EqualFold(requestType, "Sentence") {
		result = gl.Sentence(min, max)
	} else if strings.EqualFold(requestType, "Paragraph") {
		result = gl.Paragraph(min, max)
	} else {
		err = ErrRequestTypeNotFound
	}
	return result, err
}