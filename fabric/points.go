package fabric

import (
	"context"
	"fmt"
	"log"
)

type ID string

// Point
type Point struct {
	Id        ID
	Embedding []float32
	Payload   map[string]any
	// scores, etc
	Score float32
}

// UpsertPoint.
func UpsertPoint(ctx context.Context, model string, pt *Point) (ID, error) {

	if model == "" {
		model = "vernon-002"
	}

	scfg := ServeCfg{}
	hdrs := []string{
		fmt.Sprintf("tenant:%s", get_tenant(&scfg)),
		fmt.Sprintf("agent:%s", get_agentid(&scfg)),
		fmt.Sprintf("app:%s", "X"),
	}

	resp, e := Call(&CallCfg{
		Ctx:      ctx,
		Tenant:   "0",
		Agent:    model,
		Method:   "put",
		Endpoint: "point",
		Headers:  hdrs,
		Body:     pt,
	})

	if e != nil {
		return "", e
	}

	if t, ok := resp.Body.(map[string]any); ok {
		if id, ok := t["id"].(string); ok {
			return ID(id), nil
		}
	}
	return "", fmt.Errorf("No reply from model")
}

// SearchCfg
type SearchCfg struct {
	Embedding []float32
	Limit     int
}

// SearchPoints
func SearchPoints(ctx context.Context, model string, search *SearchCfg) ([]Point, error) {

	if model == "" {
		model = "vernon-002"
	}

	scfg := ServeCfg{}
	hdrs := []string{
		fmt.Sprintf("tenant:%s", get_tenant(&scfg)),
		fmt.Sprintf("agent:%s", get_agentid(&scfg)),
		fmt.Sprintf("app:%s", "X"),
	}

	resp, e := Call(&CallCfg{
		Ctx:      ctx,
		Tenant:   "0",
		Agent:    model,
		Method:   "post",
		Endpoint: "search",
		Headers:  hdrs,
		Body:     search,
	})

	log.Printf("SEARCH %v", resp)
	log.Printf("SEARCH %v", e)

	return []Point{}, e
}
