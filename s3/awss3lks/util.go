package awss3lks

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	smithy_http "github.com/aws/smithy-go/transport/http"
	"github.com/rs/zerolog/log"
)

const (
	alphabet            = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-._"
	replacementAlphabet = "-_0"
)

func wellFormBlobName(n string) (string, error) {
	if n == "" {
		return n, nil
	}

	var out strings.Builder
	var err error
	lenReplacementAlphabet := int32(len(replacementAlphabet))
	for _, c := range n {
		if strings.Index(alphabet, string(c)) >= 0 {
			_, err = out.WriteString(string(c))
		} else {
			if c > lenReplacementAlphabet {
				c = c % lenReplacementAlphabet
			} else {
				c = lenReplacementAlphabet % c
			}
			err = out.WriteByte(replacementAlphabet[c])
		}

		if err != nil {
			return "", err
		}
	}

	return out.String(), nil
}

func MapError(err error) (int, string) {
	const semLogContext = "aws-util::map-error"

	if err == nil {
		return http.StatusOK, http.StatusText(http.StatusOK)
	}

	statusCode := http.StatusInternalServerError
	statusText := err.Error()
	var noKey *types.NoSuchKey
	var noSuchBucket *types.NoSuchBucket
	var opError *smithy.OperationError
	var httpErr *smithy_http.ResponseError
	switch {
	case errors.As(err, &noKey):
		statusCode = http.StatusNotFound
		statusText = "blob key non trovata"
	case errors.As(err, &noSuchBucket):
		statusCode = http.StatusNotFound
		statusText = "bucket non trovato"
	case errors.As(err, &httpErr):
		log.Error().Err(httpErr.Err).Msg(semLogContext)
		statusCode = httpErr.HTTPStatusCode()
		statusText = httpErr.Err.Error()
	case errors.As(err, &opError):
		log.Error().Err(opError.Err).Msg(semLogContext + " - operation-error")
	default:
		log.Error().Err(err).Str("error-type", fmt.Sprintf("%T", err)).Msg(semLogContext)
	}

	return statusCode, statusText
}
