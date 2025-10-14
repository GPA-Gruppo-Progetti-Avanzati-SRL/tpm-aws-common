package awss3lks

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type RecalculateV4Signature struct {
	next   http.RoundTripper
	signer *v4.Signer
	cfg    Config
}

func (lt *RecalculateV4Signature) RoundTrip(req *http.Request) (*http.Response, error) {
	// store for later use
	val := req.Header.Get("Accept-Encoding")

	// delete the header so the header doesn't account for in the signature
	req.Header.Del("Accept-Encoding")

	// sign with the same date
	timeString := req.Header.Get("X-Amz-Date")
	timeDate, _ := time.Parse("20060102T150405Z", timeString)
	credProvider := credentials.NewStaticCredentialsProvider(lt.cfg.AccessKey, lt.cfg.SecretKey, "")
	creds, _ := credProvider.Retrieve(req.Context())

	//creds, _ := lt.cfg.Credentials.Retrieve(req.Context())
	err := lt.signer.SignHTTP(req.Context(), creds, req, v4.GetPayloadHash(req.Context()), "s3", lt.cfg.Region, timeDate)
	if err != nil {
		return nil, err
	}
	// Reset Accept-Encoding if desired
	req.Header.Set("Accept-Encoding", val)

	fmt.Println("AfterAdjustment")
	rrr, _ := httputil.DumpRequest(req, false)
	fmt.Println(string(rrr))

	// follows up the original round tripper
	return lt.next.RoundTrip(req)
}
