package types

type ErrorResponse struct {
	Error Error `json:"info,omitempty"`
}

type Error struct {
	Message        string `json:"message,omitempty"`
	DisplayMessage string `json:"displayMessage,omitempty"`
	Code           string `json:"code,omitempty"`
	Status         string `json:"status,omitempty"`
}
