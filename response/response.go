package response

type Response struct{
	Detail string `json:"detail,omitempty"`
	Data interface{} `json:"data,omitempty"`
	ResponseTime string `json:"response_time,omitempty"`
	Error string `json:"error,omitempty"`

}
