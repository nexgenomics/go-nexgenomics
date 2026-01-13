package fabric

import (
	"context"
	"fmt"
)

// Genomicize passes a string to a genomic language model (specified by the caller),
// and returns a string.
func Genomicize(ctx context.Context, model string, prompt string) (string, error) {

	resp, e := Call(&CallCfg{
		Ctx:      ctx,
		Tenant:   "0",
		Agent:    model,
		Method:   "put",
		Endpoint: "chat",
		Headers:  []string{},
		Body:     map[string]any{"prompt": prompt},
	})
	if e != nil {
		return "", e
	}

	t, ok := resp.Body.(map[string]any)
	if !ok {
		return "", fmt.Errorf("No reply from model")
	}
	return t["reply"].(string), nil
}
