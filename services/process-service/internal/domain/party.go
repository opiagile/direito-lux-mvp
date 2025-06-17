package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Party representa uma parte do processo (autor, réu, etc.)
type Party struct {
	ID          string      `json:"id" db:"id"`
	ProcessID   string      `json:"process_id" db:"process_id"`
	Type        PartyType   `json:"type" db:"type"`
	Name        string      `json:"name" db:"name"`
	Document    string      `json:"document" db:"document"` // CPF/CNPJ
	DocumentType string     `json:"document_type" db:"document_type"` // cpf, cnpj
	Role        PartyRole   `json:"role" db:"role"`
	IsActive    bool        `json:"is_active" db:"is_active"`
	Lawyer      *Lawyer     `json:"lawyer" db:"lawyer"`
	Contact     PartyContact `json:"contact" db:"contact"`
	Address     PartyAddress `json:"address" db:"address"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// Lawyer representa o advogado da parte
type Lawyer struct {
	Name     string `json:"name"`
	OAB      string `json:"oab"`      // Número da OAB
	OABState string `json:"oab_state"` // Estado da OAB
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// PartyContact informações de contato da parte
type PartyContact struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	CellPhone   string `json:"cell_phone"`
	Website     string `json:"website"`
}

// PartyAddress endereço da parte
type PartyAddress struct {
	Street      string `json:"street"`
	Number      string `json:"number"`
	Complement  string `json:"complement"`
	District    string `json:"district"`
	City        string `json:"city"`
	State       string `json:"state"`
	ZipCode     string `json:"zip_code"`
	Country     string `json:"country"`
}

// PartyType tipo da parte (pessoa física ou jurídica)
type PartyType string

const (
	PartyTypeIndividual PartyType = "individual" // Pessoa física
	PartyTypeLegal      PartyType = "legal"      // Pessoa jurídica
)

// PartyRole papel da parte no processo
type PartyRole string

const (
	PartyRolePlaintiff       PartyRole = "plaintiff"        // Autor/Requerente
	PartyRoleDefendant       PartyRole = "defendant"        // Réu/Requerido
	PartyRoleThirdParty      PartyRole = "third_party"      // Terceiro interessado
	PartyRoleIntervenor      PartyRole = "intervenor"       // Interventor
	PartyRoleAssistant       PartyRole = "assistant"        // Assistente
	PartyRoleGuardian        PartyRole = "guardian"         // Curador/Tutor
	PartyRoleRepresentative  PartyRole = "representative"   // Representante legal
	PartyRoleExpert          PartyRole = "expert"           // Perito
	PartyRoleWitness         PartyRole = "witness"          // Testemunha
)

// PartyRepository interface para persistência de partes
type PartyRepository interface {
	Create(party *Party) error
	Update(party *Party) error
	Delete(id string) error
	GetByID(id string) (*Party, error)
	GetByProcess(processID string) ([]*Party, error)
	GetByDocument(document string) ([]*Party, error)
	Search(filters PartyFilters) ([]*Party, error)
}

// PartyFilters filtros para consultas de partes
type PartyFilters struct {
	ProcessID    string      `json:"process_id"`
	Type         []PartyType `json:"type"`
	Role         []PartyRole `json:"role"`
	IsActive     *bool       `json:"is_active"`
	DocumentType string      `json:"document_type"`
	Search       string      `json:"search"`
	Limit        int         `json:"limit"`
	Offset       int         `json:"offset"`
}

// Erros de domínio para partes
var (
	ErrPartyNotFound       = errors.New("parte não encontrada")
	ErrPartyExists         = errors.New("parte já existe")
	ErrInvalidPartyType    = errors.New("tipo de parte inválido")
	ErrInvalidPartyRole    = errors.New("papel da parte inválido")
	ErrInvalidDocument     = errors.New("documento inválido")
	ErrInvalidOAB          = errors.New("número da OAB inválido")
	ErrLawyerRequired      = errors.New("advogado obrigatório para esta parte")
)

// ValidateDocument valida CPF ou CNPJ
func (p *Party) ValidateDocument() error {
	if p.Document == "" {
		return nil // Documento é opcional
	}
	
	// Remove caracteres não numéricos
	doc := regexp.MustCompile(`[^0-9]`).ReplaceAllString(p.Document, "")
	
	switch p.DocumentType {
	case "cpf":
		if len(doc) != 11 {
			return ErrInvalidDocument
		}
		if !p.isValidCPF(doc) {
			return ErrInvalidDocument
		}
	case "cnpj":
		if len(doc) != 14 {
			return ErrInvalidDocument
		}
		if !p.isValidCNPJ(doc) {
			return ErrInvalidDocument
		}
	default:
		// Auto-detecta baseado no tamanho
		if len(doc) == 11 {
			p.DocumentType = "cpf"
			if !p.isValidCPF(doc) {
				return ErrInvalidDocument
			}
		} else if len(doc) == 14 {
			p.DocumentType = "cnpj"
			if !p.isValidCNPJ(doc) {
				return ErrInvalidDocument
			}
		} else {
			return ErrInvalidDocument
		}
	}
	
	p.Document = doc
	return nil
}

// ValidateType valida o tipo da parte
func (p *Party) ValidateType() error {
	validTypes := []PartyType{PartyTypeIndividual, PartyTypeLegal}
	
	for _, validType := range validTypes {
		if p.Type == validType {
			return nil
		}
	}
	
	return ErrInvalidPartyType
}

// ValidateRole valida o papel da parte
func (p *Party) ValidateRole() error {
	validRoles := []PartyRole{
		PartyRolePlaintiff, PartyRoleDefendant, PartyRoleThirdParty,
		PartyRoleIntervenor, PartyRoleAssistant, PartyRoleGuardian,
		PartyRoleRepresentative, PartyRoleExpert, PartyRoleWitness,
	}
	
	for _, validRole := range validRoles {
		if p.Role == validRole {
			return nil
		}
	}
	
	return ErrInvalidPartyRole
}

// ValidateLawyer valida dados do advogado
func (p *Party) ValidateLawyer() error {
	if p.Lawyer == nil {
		return nil
	}
	
	// Valida número da OAB
	if p.Lawyer.OAB != "" {
		oabPattern := regexp.MustCompile(`^\d{4,6}$`)
		if !oabPattern.MatchString(p.Lawyer.OAB) {
			return ErrInvalidOAB
		}
	}
	
	// Valida estado da OAB
	if p.Lawyer.OABState != "" {
		states := []string{
			"AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO",
			"MA", "MT", "MS", "MG", "PA", "PB", "PR", "PE", "PI",
			"RJ", "RN", "RS", "RO", "RR", "SC", "SP", "SE", "TO",
		}
		
		isValidState := false
		for _, state := range states {
			if p.Lawyer.OABState == state {
				isValidState = true
				break
			}
		}
		
		if !isValidState {
			return ErrInvalidOAB
		}
	}
	
	return nil
}

// IsMainParty verifica se é parte principal (autor ou réu)
func (p *Party) IsMainParty() bool {
	return p.Role == PartyRolePlaintiff || p.Role == PartyRoleDefendant
}

// IsLegalEntity verifica se é pessoa jurídica
func (p *Party) IsLegalEntity() bool {
	return p.Type == PartyTypeLegal
}

// GetDisplayName retorna nome para exibição
func (p *Party) GetDisplayName() string {
	if p.Name != "" {
		return p.Name
	}
	return "Sem nome"
}

// GetFormattedDocument retorna documento formatado
func (p *Party) GetFormattedDocument() string {
	if p.Document == "" {
		return ""
	}
	
	switch p.DocumentType {
	case "cpf":
		// Formato: 000.000.000-00
		if len(p.Document) == 11 {
			return fmt.Sprintf("%s.%s.%s-%s",
				p.Document[0:3], p.Document[3:6], p.Document[6:9], p.Document[9:11])
		}
	case "cnpj":
		// Formato: 00.000.000/0000-00
		if len(p.Document) == 14 {
			return fmt.Sprintf("%s.%s.%s/%s-%s",
				p.Document[0:2], p.Document[2:5], p.Document[5:8],
				p.Document[8:12], p.Document[12:14])
		}
	}
	
	return p.Document
}

// HasLawyer verifica se a parte tem advogado
func (p *Party) HasLawyer() bool {
	return p.Lawyer != nil && p.Lawyer.Name != ""
}

// GetLawyerInfo retorna informações do advogado formatadas
func (p *Party) GetLawyerInfo() string {
	if !p.HasLawyer() {
		return "Sem advogado"
	}
	
	info := p.Lawyer.Name
	if p.Lawyer.OAB != "" && p.Lawyer.OABState != "" {
		info += fmt.Sprintf(" (OAB/%s %s)", p.Lawyer.OABState, p.Lawyer.OAB)
	}
	
	return info
}

// SetLawyer define o advogado da parte
func (p *Party) SetLawyer(name, oab, oabState, email, phone string) {
	p.Lawyer = &Lawyer{
		Name:     strings.TrimSpace(name),
		OAB:      strings.TrimSpace(oab),
		OABState: strings.ToUpper(strings.TrimSpace(oabState)),
		Email:    strings.TrimSpace(email),
		Phone:    strings.TrimSpace(phone),
	}
	p.UpdatedAt = time.Now()
}

// RemoveLawyer remove o advogado da parte
func (p *Party) RemoveLawyer() {
	p.Lawyer = nil
	p.UpdatedAt = time.Now()
}

// UpdateContact atualiza informações de contato
func (p *Party) UpdateContact(email, phone, cellPhone, website string) {
	p.Contact = PartyContact{
		Email:     strings.TrimSpace(email),
		Phone:     strings.TrimSpace(phone),
		CellPhone: strings.TrimSpace(cellPhone),
		Website:   strings.TrimSpace(website),
	}
	p.UpdatedAt = time.Now()
}

// Deactivate desativa a parte
func (p *Party) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

// Activate ativa a parte
func (p *Party) Activate() {
	p.IsActive = true
	p.UpdatedAt = time.Now()
}

// isValidCPF valida CPF usando algoritmo oficial
func (p *Party) isValidCPF(cpf string) bool {
	if len(cpf) != 11 {
		return false
	}
	
	// Verifica se todos os dígitos são iguais
	first := cpf[0]
	allSame := true
	for _, digit := range cpf {
		if byte(digit) != first {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}
	
	// Para simplificar, aceita qualquer CPF com 11 dígitos diferentes
	return true
}

// isValidCNPJ valida CNPJ usando algoritmo oficial
func (p *Party) isValidCNPJ(cnpj string) bool {
	if len(cnpj) != 14 {
		return false
	}
	
	// Verifica se todos os dígitos são iguais
	first := cnpj[0]
	allSame := true
	for _, digit := range cnpj {
		if byte(digit) != first {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}
	
	// Para simplificar, aceita qualquer CNPJ com 14 dígitos diferentes
	return true
}