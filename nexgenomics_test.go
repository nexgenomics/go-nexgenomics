package nexgenomics_test

import (
	"os"
	"testing"

	"github.com/nexgenomics/go-nexgenomics"
)

var (
	TEST_WEBHOOK_TOKEN string
)

func init() {
	TEST_WEBHOOK_TOKEN = os.Getenv("TEST_WEBHOOK_TOKEN")
}

func TestPing(t *testing.T) {
	t.Logf("%s", nexgenomics.Ping("abc"))
}

func TestNewWebhook(t *testing.T) {
	h := nexgenomics.NewWebhook(TEST_WEBHOOK_TOKEN)
	t.Logf("%s", h)

	sentences := []string{
		"This is thing 1",
		"This is thing 2",
		"This is thing 3",
		"This is thing 4",
		"This is thing 5",
	}
	e := h.SendSentences(sentences...)
	if e != nil {
		t.Errorf("%s", e)
	}
}
