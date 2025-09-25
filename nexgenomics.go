package nexgenomics

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Webhook accesses agents in the NexGenomics cloud using a webhook interface.
// The object requires an authorization token belonging to the agent you want to access.
type Webhook struct {
	Token string
}

// Ping is a trivial package test.
func Ping(s string) string {
	return fmt.Sprintf("nexgenomics ping [%s]", s)
}

// NewWebhook returns a Webhook object.
// Webhooks are always directed to a specific NexGenomics Agent, hence they require
// an authorization token which is generated for that agent.
func NewWebhook(token string) *Webhook {
	return &Webhook{
		Token: token,
	}
}

// SendSentences sends an array of sentences to the NexGenomics cloud.
func (wh *Webhook) SendSentences(sentences ...string) error {

	send_blob := func(blob []string) error {
		blob_bytes := []byte(strings.Join(blob, "\n"))
		c := resty.New()
		resp, e := c.R().
			SetHeader("Content-Type", "application/octet-stream").
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", wh.Token)).
			SetBody(blob_bytes).
			Post("https://webhook.nexgenomics.ai/wh/sentences")

		if e != nil {
			return e
		}
		if sc := resp.StatusCode(); sc == 403 {
			return fmt.Errorf("unauthorized")
		} else if sc != 200 {
			return fmt.Errorf("failed with status %d", sc)
		}

		return nil
	}

	// Chunk the incoming sentences to fit appropriate message size limits
	maxbloblen := 500_000
	bloblen := 0
	blobslice := []string{}

	for _, j := range sentences {
		bloblen += len(j)
		blobslice = append(blobslice, j)

		if bloblen >= maxbloblen {
			if e := send_blob(blobslice); e != nil {
				return e
			}

			blobslice = []string{}
			bloblen = 0
		}
	}

	if bloblen > 0 {
		if e := send_blob(blobslice); e != nil {
			return e
		}
	}

	return nil
}
