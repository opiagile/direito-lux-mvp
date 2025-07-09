package domain

import (
	"time"
)

// ProcessResponseData_Legacy dados de resposta para consulta de processo (renamed to avoid conflict)
type ProcessResponseData_Legacy struct {
	Found   bool         `json:"found"`
	Process *ProcessInfo `json:"process,omitempty"`
}


// MovementResponseData_Legacy dados de resposta para movimentações (renamed to avoid conflict)
type MovementResponseData_Legacy struct {
	Total     int             `json:"total"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	Movements []*MovementInfo `json:"movements"`
}

// MovementInfo informações de movimentação
type MovementInfo struct {
	ID          string     `json:"id"`
	Date        *time.Time `json:"date"`
	Code        string     `json:"code"`
	Description string     `json:"description"`
	Complement  string     `json:"complement"`
	Type        string     `json:"type"`
	Responsible string     `json:"responsible"`
	HasDocument bool       `json:"has_document"`
	DocumentURL string     `json:"document_url"`
}

// PartyResponseData_Legacy dados de resposta para partes (renamed to avoid conflict)
type PartyResponseData_Legacy struct {
	Total   int          `json:"total"`
	Parties []*PartyInfo `json:"parties"`
}

// PartyInfo informações da parte
type PartyInfo struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Document       string `json:"document"`
	DocumentType   string `json:"document_type"`
	Type           string `json:"type"`
	Role           string `json:"role"`
	LawyerName     string `json:"lawyer_name"`
	LawyerOAB      string `json:"lawyer_oab"`
	LawyerUF       string `json:"lawyer_uf"`
	IsLegalEntity  bool   `json:"is_legal_entity"`
}

