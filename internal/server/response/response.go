package response

type Response struct {
	Status  int         `json:"status"`
	Error   string      `json:"error,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}
