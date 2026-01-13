package fabric

import (
	"context"
	"fmt"
	_"log"
	"encoding/json"
	"strings"
	"github.com/nats-io/nats.go"
)


// CallCfg
type CallCfg struct {
	Ctx context.Context
	Tenant string
	Agent string
	Method string
	Endpoint string
	Body any
}


// Call
func Call(cfg *CallCfg) (*Response, error) {

	natsurl := get_natsurl(&ServeCfg{})
	if natsurl == "" {
		return nil,fmt.Errorf("missing identifiers")
	}

	j,e := json.Marshal(cfg.Body)
	if e != nil {
		return nil,e
	}

	nc,e := nats.Connect(natsurl)
	if e != nil {
		return nil,e
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
	msg,e := nc.RequestWithContext(cfg.Ctx, subj, j)
	if e == nil {
		var r Response
		json.Unmarshal(msg.Data, &r)
		if r.Status == 200 {
			return &r, nil
		} else {
			return &r, fmt.Errorf("Error %d", r.Status)
		}
	} else {
		return nil,fmt.Errorf("No response")
	}


}

