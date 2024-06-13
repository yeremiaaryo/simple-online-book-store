package response

type BaseResponse struct {
	Result bool   `json:"result"`
	Error  string `json:"error,omitempty"`
}
