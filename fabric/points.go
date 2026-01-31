package fabric

import (
	"context"
	"fmt"
)


type ID string

// Point
type Point struct {
	Id ID
	Embedding []float32
	Payload map[string]any
	// scores, etc
}


// UpsertPoint.
func UpsertPoint(ctx context.Context, model string, pt *Point) (ID,error) {

	if model == "" {
		model = "vernon-002"
	}

	scfg := ServeCfg{}
	hdrs := []string {
		fmt.Sprintf("tenant:%s", get_tenant(&scfg)),
		fmt.Sprintf("agent:%s", get_agentid(&scfg)),
		fmt.Sprintf("app:%s", "X"),
	}

	resp,e := Call(&CallCfg{
		Ctx:      ctx,
		Tenant:   "0",
		Agent:    model,
		Method: "put",
		Endpoint: "point",
		Headers: hdrs,
		Body: pt,
	})

	if e != nil {
		return "", e
	}

	if t, ok := resp.Body.(map[string]any); ok {
		if id,ok := t["id"].(string); ok {
			return ID(id), nil
		}
	}
	return "", fmt.Errorf("No reply from model")
}

