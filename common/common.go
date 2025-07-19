package common

type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"details"`
}

type ResponsePagination struct {
	Success bool        `json:"success"`
	Items   interface{} `json:"data"`
	Message string      `json:"details"`
	Total   uint        `json:"total"`
	Page    uint        `json:"page"`
	Size    uint        `json:"size"`
	Pages   uint        `json:"pages"`
}

// Generic function to filter the map based on a list of allowed keys.
func FilterMapByKeys(input map[string]any, allowedKeys []string) map[string]any {
	filtered := make(map[string]any)

	for _, key := range allowedKeys {
		if value, ok := input[key]; ok {
			filtered[key] = value
		}
	}

	return filtered
}

// Generic function to filter the map based on a list of allowed keys.
func FilterSearchTerms(input map[string]any, allowedKeys []string) []string {
	filtered := make([]string, 0)

	for _, key := range allowedKeys {
		if value, ok := input[key]; ok {
			filtered = append(filtered, value.(string))
		}
	}

	return filtered
}
