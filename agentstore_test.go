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
}
