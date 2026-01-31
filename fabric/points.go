package fabric

import (
	"context"
	"fmt"
	"log"
)


type ID string

// Point
type Point struct {
	Id ID
	Embedding []float32
	Payload map[string]string
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

	t, ok := resp.Body.(map[string]any)
	if !ok {
		return "", fmt.Errorf("No reply from model")
	}

	log.Printf("VERNON %v", t)

	return "NOTHING",nil
}

