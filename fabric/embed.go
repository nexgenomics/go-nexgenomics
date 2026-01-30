package fabric

import (
	"context"
	"encoding/json"
)

// Embed returns an embedding for the given input string as a vector of float32.
// The size of the returned embedding is dependent on the model chosen.
func Embed(ctx context.Context, model string, data []byte) ([]float32, error) {
	if model == "" {
		model = "enfield-001"
	}

	resp, e := RawCall(&CallCfg{
		Ctx:    ctx,
		Tenant: "0",
		Agent:  model,
		Body:   data,
	})
	if e != nil {
		return nil, e
	}

	var fp []float32
	json.Unmarshal(resp.Body.([]byte), &fp)

	return fp, nil
}
