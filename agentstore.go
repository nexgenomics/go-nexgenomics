package nexgenomics

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Agentstore
type Agentstore struct {
	Token string
}

type Agent struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// NewAgentstore
func NewAgentstore(token string) *Agentstore {
	a := Agentstore{
		Token: token,
	}

	return &a
}

// Agents returns a list of the agents you own.
func (as *Agentstore) Agents() ([]Agent, error) {

	c := resty.New()
	resp, e := c.R().
		//SetHeader("Content-Type", "application/octet-stream").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", as.Token)).
		//SetBody(blob_bytes).
		SetResult(&[]Agent{}).
		Get("https://agentstore.nexgenomics.ai/api/agents")

	if e != nil {
		return nil, e
	}

	if sc := resp.StatusCode(); sc == 403 {
		return nil, fmt.Errorf("unauthorized")
	} else if sc != 200 {
		return nil, fmt.Errorf("failed with status %d", sc)
	}

	if agents, ok := resp.Result().(*[]Agent); ok {
		return *agents, nil
	} else {
		return nil, fmt.Errorf("unknown response")
	}

}
