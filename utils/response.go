package utils

// ErrorResponse estrutura para a resposta de erro
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse representa a estrutura de resposta de sucesso
type SuccessResponse struct {
	Message string
}
