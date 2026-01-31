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


// UpsertPoint
func UpsertPoint(ctx context.Context, model string, pt *Point) (ID,error) {

	if model == "" {
		model = "vernon-002"
	}

	resp,e := Call(&CallCfg{
		Ctx:      ctx,
		Tenant:   "0",
		Agent:    model,
		Method: "put",
		Endpoint: "point",
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

