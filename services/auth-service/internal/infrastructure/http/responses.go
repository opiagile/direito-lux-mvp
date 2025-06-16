package http

// ErrorResponse representa uma resposta de erro padronizada
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// MessageResponse representa uma resposta de sucesso com mensagem
type MessageResponse struct {
	Message string `json:"message"`
}