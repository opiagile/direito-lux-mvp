package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"encoding/json"
	"github.com/direito-lux/process-service/internal/domain"
)

// PartyRepository implementação PostgreSQL do repositório de partes
type PartyRepository struct {
	db *sql.DB
}

// NewPartyRepository cria nova instância do repositório
func NewPartyRepository(db *sql.DB) *PartyRepository {
	return &PartyRepository{db: db}
}

// Create cria nova parte
func (r *PartyRepository) Create(party *domain.Party) error {
	query := `
		INSERT INTO parties (
			id, process_id, type, name, document, document_type, role,
			is_active, lawyer, contact, address, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)`

	// Serializar campos JSON
	var lawyerJSON []byte
	var err error
	if party.Lawyer != nil {
		lawyerJSON, err = json.Marshal(party.Lawyer)
		if err != nil {
			return fmt.Errorf("erro ao serializar lawyer: %w", err)
		}
	}

	contactJSON, err := json.Marshal(party.Contact)
	if err != nil {
		return fmt.Errorf("erro ao serializar contact: %w", err)
	}

	addressJSON, err := json.Marshal(party.Address)
	if err != nil {
		return fmt.Errorf("erro ao serializar address: %w", err)
	}

	_, err = r.db.Exec(
		query,
		party.ID,
		party.ProcessID,
		party.Type,
		party.Name,
		party.Document,
		party.DocumentType,
		party.Role,
		party.IsActive,
		lawyerJSON,
		contactJSON,
		addressJSON,
		party.CreatedAt,
		party.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao criar parte: %w", err)
	}

	return nil
}

// Update atualiza parte existente
func (r *PartyRepository) Update(party *domain.Party) error {
	query := `
		UPDATE parties SET
			name = $2, document = $3, document_type = $4, is_active = $5,
			lawyer = $6, contact = $7, address = $8, updated_at = $9
		WHERE id = $1`

	// Serializar campos JSON
	var lawyerJSON []byte
	var err error
	if party.Lawyer != nil {
		lawyerJSON, err = json.Marshal(party.Lawyer)
		if err != nil {
			return fmt.Errorf("erro ao serializar lawyer: %w", err)
		}
	}

	contactJSON, err := json.Marshal(party.Contact)
	if err != nil {
		return fmt.Errorf("erro ao serializar contact: %w", err)
	}

	addressJSON, err := json.Marshal(party.Address)
	if err != nil {
		return fmt.Errorf("erro ao serializar address: %w", err)
	}

	result, err := r.db.Exec(
		query,
		party.ID,
		party.Name,
		party.Document,
		party.DocumentType,
		party.IsActive,
		lawyerJSON,
		contactJSON,
		addressJSON,
		party.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar parte: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrPartyNotFound
	}

	return nil
}

// Delete remove parte (soft delete)
func (r *PartyRepository) Delete(id string) error {
	query := `UPDATE parties SET is_active = false, updated_at = NOW() WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar parte: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrPartyNotFound
	}

	return nil
}

// GetByID busca parte por ID
func (r *PartyRepository) GetByID(id string) (*domain.Party, error) {
	query := `
		SELECT id, process_id, type, name, document, document_type, role,
			is_active, lawyer, contact, address, created_at, updated_at
		FROM parties
		WHERE id = $1`

	return r.scanParty(r.db.QueryRow(query, id))
}

// GetByProcess busca partes por processo
func (r *PartyRepository) GetByProcess(processID string) ([]*domain.Party, error) {
	query := `
		SELECT id, process_id, type, name, document, document_type, role,
			is_active, lawyer, contact, address, created_at, updated_at
		FROM parties
		WHERE process_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query, processID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar partes: %w", err)
	}
	defer rows.Close()

	return r.scanParties(rows)
}

// GetByDocument busca partes por documento
func (r *PartyRepository) GetByDocument(document string) ([]*domain.Party, error) {
	query := `
		SELECT id, process_id, type, name, document, document_type, role,
			is_active, lawyer, contact, address, created_at, updated_at
		FROM parties
		WHERE document = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, document)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar partes por documento: %w", err)
	}
	defer rows.Close()

	return r.scanParties(rows)
}

// Search busca partes com filtros
func (r *PartyRepository) Search(filters domain.PartyFilters) ([]*domain.Party, error) {
	query, args := r.buildQuery(filters)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar partes: %w", err)
	}
	defer rows.Close()

	return r.scanParties(rows)
}

// buildQuery constrói query com filtros
func (r *PartyRepository) buildQuery(filters domain.PartyFilters) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Filtro por processo
	if filters.ProcessID != "" {
		conditions = append(conditions, fmt.Sprintf("process_id = $%d", argIndex))
		args = append(args, filters.ProcessID)
		argIndex++
	}

	// Filtros por tipo
	if len(filters.Type) > 0 {
		typePlaceholders := make([]string, len(filters.Type))
		for i, partyType := range filters.Type {
			typePlaceholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, partyType)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("type IN (%s)", strings.Join(typePlaceholders, ",")))
	}

	// Filtros por papel
	if len(filters.Role) > 0 {
		rolePlaceholders := make([]string, len(filters.Role))
		for i, role := range filters.Role {
			rolePlaceholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, role)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("role IN (%s)", strings.Join(rolePlaceholders, ",")))
	}

	// Filtro por ativo
	if filters.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filters.IsActive)
		argIndex++
	}

	// Filtro por tipo de documento
	if filters.DocumentType != "" {
		conditions = append(conditions, fmt.Sprintf("document_type = $%d", argIndex))
		args = append(args, filters.DocumentType)
		argIndex++
	}

	// Busca textual
	if filters.Search != "" {
		searchCondition := fmt.Sprintf("(name ILIKE $%d OR document ILIKE $%d)", argIndex, argIndex+1)
		searchTerm := "%" + filters.Search + "%"
		conditions = append(conditions, searchCondition)
		args = append(args, searchTerm, searchTerm)
		argIndex += 2
	}

	// Construir query final
	query := `
		SELECT id, process_id, type, name, document, document_type, role,
			is_active, lawyer, contact, address, created_at, updated_at
		FROM parties`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Ordenação
	query += " ORDER BY created_at ASC"

	// Paginação
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filters.Limit)
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", filters.Offset)
	}

	return query, args
}

// scanParty converte linha do banco para entidade Party
func (r *PartyRepository) scanParty(row *sql.Row) (*domain.Party, error) {
	var party domain.Party
	var lawyerJSON, contactJSON, addressJSON []byte
	var document, documentType sql.NullString

	err := row.Scan(
		&party.ID,
		&party.ProcessID,
		&party.Type,
		&party.Name,
		&document,
		&documentType,
		&party.Role,
		&party.IsActive,
		&lawyerJSON,
		&contactJSON,
		&addressJSON,
		&party.CreatedAt,
		&party.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPartyNotFound
		}
		return nil, fmt.Errorf("erro ao escanear parte: %w", err)
	}

	// Tratar campos nullable
	if document.Valid {
		party.Document = document.String
	}

	if documentType.Valid {
		party.DocumentType = documentType.String
	}

	// Deserializar campos JSON
	if lawyerJSON != nil {
		var lawyer domain.Lawyer
		if err := json.Unmarshal(lawyerJSON, &lawyer); err != nil {
			return nil, fmt.Errorf("erro ao deserializar lawyer: %w", err)
		}
		party.Lawyer = &lawyer
	}

	if err := json.Unmarshal(contactJSON, &party.Contact); err != nil {
		return nil, fmt.Errorf("erro ao deserializar contact: %w", err)
	}

	if err := json.Unmarshal(addressJSON, &party.Address); err != nil {
		return nil, fmt.Errorf("erro ao deserializar address: %w", err)
	}

	return &party, nil
}

// scanParties converte múltiplas linhas para entidades Party
func (r *PartyRepository) scanParties(rows *sql.Rows) ([]*domain.Party, error) {
	var parties []*domain.Party

	for rows.Next() {
		var party domain.Party
		var lawyerJSON, contactJSON, addressJSON []byte
		var document, documentType sql.NullString

		err := rows.Scan(
			&party.ID,
			&party.ProcessID,
			&party.Type,
			&party.Name,
			&document,
			&documentType,
			&party.Role,
			&party.IsActive,
			&lawyerJSON,
			&contactJSON,
			&addressJSON,
			&party.CreatedAt,
			&party.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao escanear parte: %w", err)
		}

		// Tratar campos nullable
		if document.Valid {
			party.Document = document.String
		}

		if documentType.Valid {
			party.DocumentType = documentType.String
		}

		// Deserializar campos JSON
		if lawyerJSON != nil {
			var lawyer domain.Lawyer
			if err := json.Unmarshal(lawyerJSON, &lawyer); err != nil {
				return nil, fmt.Errorf("erro ao deserializar lawyer: %w", err)
			}
			party.Lawyer = &lawyer
		}

		if err := json.Unmarshal(contactJSON, &party.Contact); err != nil {
			return nil, fmt.Errorf("erro ao deserializar contact: %w", err)
		}

		if err := json.Unmarshal(addressJSON, &party.Address); err != nil {
			return nil, fmt.Errorf("erro ao deserializar address: %w", err)
		}

		parties = append(parties, &party)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar partes: %w", err)
	}

	return parties, nil
}