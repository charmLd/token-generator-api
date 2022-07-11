package transformers

// ErrorTransformer is used to transform the response payload for errors.
type ErrorTransformer struct {
	CorrelationID    interface{} `json:"correlationId,omitempty"`
	Code             string      `json:"code,omitempty"`
	Message          string      `json:"message,omitempty"`
	DeveloperMessage string      `json:"developerMessage,omitempty"`
}
