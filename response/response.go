package response

type Response struct {
	Detail       string      `json:"detail,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	ResponseTime string      `json:"response_time,omitempty"`
}

type FaceLenDesc struct {
	Id          string `json:"id,omitempty" `
	Name        string `json:"name,omitempty" `
	Descriptors int    `json:"descriptors,omitempty" `
}
