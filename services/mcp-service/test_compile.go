package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	
	"github.com/direito-lux/mcp-service/internal/infrastructure/config"
	"github.com/direito-lux/mcp-service/internal/infrastructure/http/dto"
	"github.com/direito-lux/mcp-service/internal/infrastructure/http/handlers"
	"github.com/direito-lux/mcp-service/internal/domain"
)

func main() {
	// Teste básico de compilação
	fmt.Println("Testing MCP Service compilation...")

	// Testar config
	cfg := &config.Config{}
	fmt.Printf("Config loaded: %v\n", cfg != nil)

	// Testar DTOs
	sessionReq := &dto.CreateSessionRequest{
		Channel:  "web",
		UserID:   "test",
		TenantID: "test",
	}
	fmt.Printf("DTO created: %v\n", sessionReq != nil)

	// Testar handlers
	handler := handlers.CreateMCPSessionV2()
	fmt.Printf("Handler created: %v\n", handler != nil)

	// Testar domain
	session := domain.NewMCPSession(
		uuid.New(), // tenantID
		uuid.New(), // userID
		"web",      // channel
		"test-123", // externalID
	)
	fmt.Printf("Domain entity created: %v\n", session != nil)

	log.Println("MCP Service compilation test successful!")
}