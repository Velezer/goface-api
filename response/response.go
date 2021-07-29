package response

type Response struct{
	StatusCode int `json:"status_code"`
	Status string `json:"status"`
	Detail string `json:"detail"`
	Data interface{} `json:"data,omitempty"`
	ResponseTime string `json:"response_time,omitempty"`
}