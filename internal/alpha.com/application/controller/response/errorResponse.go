package response

type ErrorRespone struct {
	Status      int32         `json:"status"`
	ErrorDetail []ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Description string `json:"description"`
	FieldName   string `json:"fieldName"`
}

type CustomError struct {
	ErrorName  string `json:"errorName"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
