package tfstate

import (
	"fmt"
)

func (s *TFState) ServiceVCLBackends(serviceID string) ([]map[string]interface{}, error) {
	q := fmt.Sprintf(`.resources[] | select(.type=="fastly_service_vcl") | .instances[] | select(.attributes.id=="%s") | .attributes.backend`, serviceID)

	r, err := s.Query(q)
	if err != nil {
		return nil, err
	}

	raw, ok := r.Value.([]interface{})
	if !ok {
		// If backend is null/missing in state, treat as empty.
		return nil, nil
	}

	out := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, m)
	}

	return out, nil
}
