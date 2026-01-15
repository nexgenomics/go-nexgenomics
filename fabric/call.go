package fabric

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	_ "log"
	"strings"
)

// CallCfg
type CallCfg struct {
	Ctx      context.Context
	NatsUrl  string // optional
	Tenant   string
	Agent    string
	Method   string
	Endpoint string
	Headers  []string
	Body     any
}

// Call
func Call(cfg *CallCfg) (*Response, error) {

	// If the caller set a NATS endpoint in the call config, use that.
	// Otherwise check the local system. This model gives the
	// maximum flexibility.
	// This module is primarily intended for use in standard agents running
	// in the Secure Fabric, so connection parameters are already available
	// and visible to get_natsurl.
	natsurl := cfg.NatsUrl
	if natsurl == "" {
		natsurl = get_natsurl(&ServeCfg{})
		if natsurl == "" {
			return nil, fmt.Errorf("missing identifiers")
		}
	}

	hdrs := cfg.Headers
	if hdrs == nil {
		hdrs = []string{}
	}

	t2 := map[string]any{
		"headers": hdrs,
		"body":    cfg.Body,
	}
	j, e := json.Marshal(t2)
	if e != nil {
		return nil, e
	}

	nc, e := nats.Connect(natsurl)
	if e != nil {
		return nil, e
	}
	defer nc.Drain()

	// This implements the "modern" calling convention for agent-rest, where the method
	// and endpoint are baked into the subject.
	m := strings.TrimSpace(strings.ToLower(cfg.Method))
	ep := strings.TrimSpace(strings.ToLower(cfg.Endpoint))
	if ep[0] == '/' {
		ep = ep[1:]
	}

	subj := fmt.Sprintf("agent.rest.%s.%s.%s.%s", cfg.Tenant, cfg.Agent, m, ep)
	//log.Printf("Calling %v",subj)
	//log.Printf("Calling %v",j)
	msg, e := nc.RequestWithContext(cfg.Ctx, subj, j)
	if e == nil {
		var r Response
		json.Unmarshal(msg.Data, &r)
		if r.Status == 200 {
			return &r, nil
		} else {
			return &r, fmt.Errorf("Error %d", r.Status)
		}
	} else {
		return nil, fmt.Errorf("No response")
	}

}
