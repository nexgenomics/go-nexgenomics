package nexgenomics_test

import (
	"os"
	"testing"

	"github.com/nexgenomics/go-nexgenomics"
)

var (
	AGENTSTORE_TOKEN string
)

func init() {
	AGENTSTORE_TOKEN = os.Getenv("AGENTSTORE_TOKEN")
}

func TestNewAgentstore(t *testing.T) {
	t.Logf("agentstore token [%s]", AGENTSTORE_TOKEN)
	as := nexgenomics.NewAgentstore(AGENTSTORE_TOKEN)

	agents, e := as.Agents()
	if e != nil {
		t.Errorf("%s", e)
	}

	for i, a := range agents {
		t.Logf("%d) %s", i, a)
	}
}
