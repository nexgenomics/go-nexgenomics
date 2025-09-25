package nexgenomics_test

import (
	"os"
	"testing"
)

var (
	AGENTSTORE_TOKEN string
)

func init() {
	AGENTSTORE_TOKEN = os.Getenv("AGENTSTORE_TOKEN")
}

func TestNewAgentstore(t *testing.T) {
	t.Logf("agentstore token [%s]", AGENTSTORE_TOKEN)
	/*
	   h := nexgenomics.NewWebhook(TEST_WEBHOOK_TOKEN)

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
	*/
}
