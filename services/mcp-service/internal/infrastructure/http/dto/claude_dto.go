package dto

import (
	"time"
)

// ClaudeChatRequest request para chat com Claude
type ClaudeChatRequest struct {
	Model       string                 `json:"model,omitempty"`
	Messages    []ClaudeMessage        `json:"messages" binding:"required"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	TopP        float64                `json:"top_p,omitempty"`
	TopK        int                    `json:"top_k,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	StopWords   []string               `json:"stop_sequences,omitempty"`
	System      string                 `json:"system,omitempty"`
	Tools       []ClaudeToolDefinition `json:"tools,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ClaudeMessage mensagem para Claude
type ClaudeMessage struct {
	Role    string                `json:"role"` // user, assistant
	Content []ClaudeContent       `json:"content"`
}

// ClaudeContent conteúdo da mensagem
type ClaudeContent struct {
	Type string `json:"type"` // text, image, tool_use, tool_result
	Text string `json:"text,omitempty"`
	
	// Para imagens
	Source *ClaudeImageSource `json:"source,omitempty"`
	
	// Para tool use
	ID    string                 `json:"id,omitempty"`
	Name  string                 `json:"name,omitempty"`
	Input map[string]interface{} `json:"input,omitempty"`
	
	// Para tool result
	ToolUseID string      `json:"tool_use_id,omitempty"`
	Content   interface{} `json:"content,omitempty"`
	IsError   bool        `json:"is_error,omitempty"`
}

// ClaudeImageSource fonte da imagem
type ClaudeImageSource struct {
	Type      string `json:"type"` // base64
	MediaType string `json:"media_type"` // image/jpeg, image/png, etc.
	Data      string `json:"data"`
}

// ClaudeToolDefinition definição de ferramenta para Claude
type ClaudeToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

// ClaudeChatResponse response do chat
type ClaudeChatResponse struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Role         string                 `json:"role"`
	Content      []ClaudeContent        `json:"content"`
	Model        string                 `json:"model"`
	StopReason   string                 `json:"stop_reason,omitempty"`
	StopSequence string                 `json:"stop_sequence,omitempty"`
	Usage        ClaudeUsage           `json:"usage"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

// ClaudeUsage uso de tokens
type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// ClaudeWithToolsRequest request para Claude com ferramentas
type ClaudeWithToolsRequest struct {
	SessionID    string                 `json:"session_id" binding:"required"`
	Message      string                 `json:"message" binding:"required"`
	AllowedTools []string               `json:"allowed_tools,omitempty"`
	AutoExecute  bool                   `json:"auto_execute,omitempty"`
	Options      ClaudeChatOptions      `json:"options,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ClaudeChatOptions opções do chat
type ClaudeChatOptions struct {
	Model       string  `json:"model,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	Stream      bool    `json:"stream,omitempty"`
	SaveHistory bool    `json:"save_history,omitempty"`
}

// ClaudeToolExecutionResponse response da execução com ferramentas
type ClaudeToolExecutionResponse struct {
	ID            string                    `json:"id"`
	SessionID     string                    `json:"session_id"`
	Message       ClaudeMessage             `json:"message"`
	ToolCalls     []ClaudeToolExecution     `json:"tool_calls,omitempty"`
	FinalResponse string                    `json:"final_response"`
	Usage         ClaudeUsage              `json:"usage"`
	Duration      time.Duration             `json:"duration"`
	CreatedAt     time.Time                 `json:"created_at"`
}

// ClaudeToolExecution execução de ferramenta pelo Claude
type ClaudeToolExecution struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Input      map[string]interface{} `json:"input"`
	Output     interface{}            `json:"output,omitempty"`
	Status     string                 `json:"status"`
	Error      string                 `json:"error,omitempty"`
	StartedAt  time.Time              `json:"started_at"`
	FinishedAt *time.Time             `json:"finished_at,omitempty"`
	Duration   time.Duration          `json:"duration,omitempty"`
}

// ClaudeModelsResponse lista de modelos disponíveis
type ClaudeModelsResponse struct {
	Models []ClaudeModel `json:"models"`
}

// ClaudeModel modelo do Claude
type ClaudeModel struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	MaxTokens    int       `json:"max_tokens"`
	InputPrice   float64   `json:"input_price"`  // Preço por 1M tokens de input
	OutputPrice  float64   `json:"output_price"` // Preço por 1M tokens de output
	SupportsTools bool     `json:"supports_tools"`
	SupportsImages bool    `json:"supports_images"`
	ContextWindow int      `json:"context_window"`
	CreatedAt    time.Time `json:"created_at"`
}

// ClaudeStreamResponse response de streaming
type ClaudeStreamResponse struct {
	Type    string                 `json:"type"` // message_start, content_block_start, content_block_delta, content_block_stop, message_delta, message_stop
	Message *ClaudeChatResponse    `json:"message,omitempty"`
	Index   int                    `json:"index,omitempty"`
	Delta   map[string]interface{} `json:"delta,omitempty"`
	Usage   *ClaudeUsage          `json:"usage,omitempty"`
}

// ClaudeAPIError erro da API do Claude
type ClaudeAPIError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}