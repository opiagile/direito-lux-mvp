package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/direito-lux/process-service/internal/domain"
)

// ProcessRepository implementação PostgreSQL do repositório de processos
type ProcessRepository struct {
	db *sql.DB
}

// NewProcessRepository cria nova instância do repositório
func NewProcessRepository(db *sql.DB) *ProcessRepository {
	return &ProcessRepository{db: db}
}

// Create cria novo processo
func (r *ProcessRepository) Create(process *domain.Process) error {
	query := `
		INSERT INTO processes (
			id, tenant_id, client_id, number, original_number, title, description,
			status, stage, subject, value, court_id, judge_id, monitoring, tags,
			custom_fields, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)`

	// Serializar campos JSON
	subjectJSON, err := json.Marshal(process.Subject)
	if err != nil {
		return fmt.Errorf("erro ao serializar subject: %w", err)
	}

	var valueJSON []byte
	if process.Value != nil {
		valueJSON, err = json.Marshal(process.Value)
		if err != nil {
			return fmt.Errorf("erro ao serializar value: %w", err)
		}
	}

	monitoringJSON, err := json.Marshal(process.Monitoring)
	if err != nil {
		return fmt.Errorf("erro ao serializar monitoring: %w", err)
	}

	customFieldsJSON, err := json.Marshal(process.CustomFields)
	if err != nil {
		return fmt.Errorf("erro ao serializar custom_fields: %w", err)
	}

	_, err = r.db.Exec(
		query,
		process.ID,
		process.TenantID,
		process.ClientID,
		process.Number,
		process.OriginalNumber,
		process.Title,
		process.Description,
		process.Status,
		process.Stage,
		subjectJSON,
		valueJSON,
		process.CourtID,
		process.JudgeID,
		monitoringJSON,
		pq.Array(process.Tags),
		customFieldsJSON,
		process.CreatedAt,
		process.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao criar processo: %w", err)
	}

	return nil
}

// Update atualiza processo existente
func (r *ProcessRepository) Update(process *domain.Process) error {
	query := `
		UPDATE processes SET
			title = $2, description = $3, status = $4, stage = $5, subject = $6,
			value = $7, judge_id = $8, monitoring = $9, tags = $10, custom_fields = $11,
			last_movement_at = $12, last_sync_at = $13, updated_at = $14, archived_at = $15
		WHERE id = $1 AND tenant_id = $16`

	// Serializar campos JSON
	subjectJSON, err := json.Marshal(process.Subject)
	if err != nil {
		return fmt.Errorf("erro ao serializar subject: %w", err)
	}

	var valueJSON []byte
	if process.Value != nil {
		valueJSON, err = json.Marshal(process.Value)
		if err != nil {
			return fmt.Errorf("erro ao serializar value: %w", err)
		}
	}

	monitoringJSON, err := json.Marshal(process.Monitoring)
	if err != nil {
		return fmt.Errorf("erro ao serializar monitoring: %w", err)
	}

	customFieldsJSON, err := json.Marshal(process.CustomFields)
	if err != nil {
		return fmt.Errorf("erro ao serializar custom_fields: %w", err)
	}

	result, err := r.db.Exec(
		query,
		process.ID,
		process.Title,
		process.Description,
		process.Status,
		process.Stage,
		subjectJSON,
		valueJSON,
		process.JudgeID,
		monitoringJSON,
		pq.Array(process.Tags),
		customFieldsJSON,
		process.LastMovementAt,
		process.LastSyncAt,
		process.UpdatedAt,
		process.ArchivedAt,
		process.TenantID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar processo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrProcessNotFound
	}

	return nil
}

// Delete remove processo (soft delete)
func (r *ProcessRepository) Delete(id string) error {
	query := `
		UPDATE processes SET
			status = $2, archived_at = $3, updated_at = $4
		WHERE id = $1`

	now := time.Now()
	result, err := r.db.Exec(query, id, domain.ProcessStatusArchived, now, now)
	if err != nil {
		return fmt.Errorf("erro ao deletar processo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrProcessNotFound
	}

	return nil
}

// Archive arquiva processo
func (r *ProcessRepository) Archive(id string) error {
	return r.Delete(id) // Mesmo comportamento que delete
}

// GetByID busca processo por ID
func (r *ProcessRepository) GetByID(id string) (*domain.Process, error) {
	query := `
		SELECT id, tenant_id, client_id, number, original_number, title, description,
			status, stage, subject, value, court_id, judge_id, monitoring, tags,
			custom_fields, last_movement_at, last_sync_at, created_at, updated_at, archived_at
		FROM processes
		WHERE id = $1`

	return r.scanProcess(r.db.QueryRow(query, id))
}

// GetByNumber busca processo por número
func (r *ProcessRepository) GetByNumber(number string) (*domain.Process, error) {
	query := `
		SELECT id, tenant_id, client_id, number, original_number, title, description,
			status, stage, subject, value, court_id, judge_id, monitoring, tags,
			custom_fields, last_movement_at, last_sync_at, created_at, updated_at, archived_at
		FROM processes
		WHERE number = $1`

	return r.scanProcess(r.db.QueryRow(query, number))
}

// GetByTenant busca processos por tenant com filtros
func (r *ProcessRepository) GetByTenant(tenantID string, filters domain.ProcessFilters) ([]*domain.Process, error) {
	query, args := r.buildQuery("tenant_id = $1", tenantID, filters)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar processos: %w", err)
	}
	defer rows.Close()

	return r.scanProcesses(rows)
}

// GetByClient busca processos por cliente com filtros
func (r *ProcessRepository) GetByClient(clientID string, filters domain.ProcessFilters) ([]*domain.Process, error) {
	query, args := r.buildQuery("client_id = $1", clientID, filters)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar processos: %w", err)
	}
	defer rows.Close()

	return r.scanProcesses(rows)
}

// GetActiveForMonitoring busca processos ativos para monitoramento
func (r *ProcessRepository) GetActiveForMonitoring() ([]*domain.Process, error) {
	query := `
		SELECT id, tenant_id, client_id, number, original_number, title, description,
			status, stage, subject, value, court_id, judge_id, monitoring, tags,
			custom_fields, last_movement_at, last_sync_at, created_at, updated_at, archived_at
		FROM processes
		WHERE status = $1 AND monitoring->>'enabled' = 'true'
		ORDER BY updated_at DESC`

	rows, err := r.db.Query(query, domain.ProcessStatusActive)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar processos para monitoramento: %w", err)
	}
	defer rows.Close()

	return r.scanProcesses(rows)
}

// GetNeedingSync busca processos que precisam sincronização
func (r *ProcessRepository) GetNeedingSync() ([]*domain.Process, error) {
	query := `
		SELECT id, tenant_id, client_id, number, original_number, title, description,
			status, stage, subject, value, court_id, judge_id, monitoring, tags,
			custom_fields, last_movement_at, last_sync_at, created_at, updated_at, archived_at
		FROM processes
		WHERE status = $1 
			AND monitoring->>'enabled' = 'true'
			AND monitoring->>'auto_sync' = 'true'
			AND (
				last_sync_at IS NULL 
				OR last_sync_at < NOW() - INTERVAL '1 hour' * CAST(monitoring->>'sync_interval_hours' AS INTEGER)
			)
		ORDER BY last_sync_at ASC NULLS FIRST`

	rows, err := r.db.Query(query, domain.ProcessStatusActive)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar processos para sincronização: %w", err)
	}
	defer rows.Close()

	return r.scanProcesses(rows)
}

// CountByTenant conta processos por tenant
func (r *ProcessRepository) CountByTenant(tenantID string) (int, error) {
	query := `SELECT COUNT(*) FROM processes WHERE tenant_id = $1`
	
	var count int
	err := r.db.QueryRow(query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar processos: %w", err)
	}

	return count, nil
}

// CountByStatus conta processos por status
func (r *ProcessRepository) CountByStatus(tenantID string, status domain.ProcessStatus) (int, error) {
	query := `SELECT COUNT(*) FROM processes WHERE tenant_id = $1 AND status = $2`
	
	var count int
	err := r.db.QueryRow(query, tenantID, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar processos por status: %w", err)
	}

	return count, nil
}

// buildQuery constrói query com filtros
func (r *ProcessRepository) buildQuery(baseCondition, baseValue string, filters domain.ProcessFilters) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 2 // Começa em 2 pois $1 é usado na condição base

	conditions = append(conditions, baseCondition)
	args = append(args, baseValue)

	// Filtros por status
	if len(filters.Status) > 0 {
		statusPlaceholders := make([]string, len(filters.Status))
		for i, status := range filters.Status {
			statusPlaceholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, status)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(statusPlaceholders, ",")))
	}

	// Filtros por fase
	if len(filters.Stage) > 0 {
		stagePlaceholders := make([]string, len(filters.Stage))
		for i, stage := range filters.Stage {
			stagePlaceholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, stage)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("stage IN (%s)", strings.Join(stagePlaceholders, ",")))
	}

	// Filtro por tribunal
	if filters.CourtID != "" {
		conditions = append(conditions, fmt.Sprintf("court_id = $%d", argIndex))
		args = append(args, filters.CourtID)
		argIndex++
	}

	// Filtro por juiz
	if filters.JudgeID != "" {
		conditions = append(conditions, fmt.Sprintf("judge_id = $%d", argIndex))
		args = append(args, filters.JudgeID)
		argIndex++
	}

	// Filtro por tags
	if len(filters.Tags) > 0 {
		conditions = append(conditions, fmt.Sprintf("tags && $%d", argIndex))
		args = append(args, pq.Array(filters.Tags))
		argIndex++
	}

	// Filtros de data
	if filters.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filters.DateTo)
		argIndex++
	}

	// Busca textual
	if filters.Search != "" {
		searchCondition := fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d OR number ILIKE $%d)", argIndex, argIndex+1, argIndex+2)
		searchTerm := "%" + filters.Search + "%"
		conditions = append(conditions, searchCondition)
		args = append(args, searchTerm, searchTerm, searchTerm)
		argIndex += 3
	}

	// Construir query final
	query := `
		SELECT id, tenant_id, client_id, number, original_number, title, description,
			status, stage, subject, value, court_id, judge_id, monitoring, tags,
			custom_fields, last_movement_at, last_sync_at, created_at, updated_at, archived_at
		FROM processes`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Ordenação
	if filters.SortBy != "" {
		order := "ASC"
		if filters.SortOrder == "desc" {
			order = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", filters.SortBy, order)
	} else {
		query += " ORDER BY updated_at DESC"
	}

	// Paginação
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filters.Limit)
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", filters.Offset)
	}

	return query, args
}

// scanProcess converte linha do banco para entidade Process
func (r *ProcessRepository) scanProcess(row *sql.Row) (*domain.Process, error) {
	var process domain.Process
	var subjectJSON, valueJSON, monitoringJSON, customFieldsJSON []byte
	var tags []string
	var originalNumber sql.NullString
	var judgeID sql.NullString
	var lastMovementAt, lastSyncAt, archivedAt sql.NullTime

	err := row.Scan(
		&process.ID,
		&process.TenantID,
		&process.ClientID,
		&process.Number,
		&originalNumber,
		&process.Title,
		&process.Description,
		&process.Status,
		&process.Stage,
		&subjectJSON,
		&valueJSON,
		&process.CourtID,
		&judgeID,
		&monitoringJSON,
		pq.Array(&tags),
		&customFieldsJSON,
		&lastMovementAt,
		&lastSyncAt,
		&process.CreatedAt,
		&process.UpdatedAt,
		&archivedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProcessNotFound
		}
		return nil, fmt.Errorf("erro ao escanear processo: %w", err)
	}

	// Deserializar campos JSON
	if err := json.Unmarshal(subjectJSON, &process.Subject); err != nil {
		return nil, fmt.Errorf("erro ao deserializar subject: %w", err)
	}

	if valueJSON != nil {
		var value domain.ProcessValue
		if err := json.Unmarshal(valueJSON, &value); err != nil {
			return nil, fmt.Errorf("erro ao deserializar value: %w", err)
		}
		process.Value = &value
	}

	if err := json.Unmarshal(monitoringJSON, &process.Monitoring); err != nil {
		return nil, fmt.Errorf("erro ao deserializar monitoring: %w", err)
	}

	if err := json.Unmarshal(customFieldsJSON, &process.CustomFields); err != nil {
		return nil, fmt.Errorf("erro ao deserializar custom_fields: %w", err)
	}

	// Tratar campos nullable
	if originalNumber.Valid {
		process.OriginalNumber = originalNumber.String
	}

	if judgeID.Valid {
		process.JudgeID = &judgeID.String
	}

	if lastMovementAt.Valid {
		process.LastMovementAt = &lastMovementAt.Time
	}

	if lastSyncAt.Valid {
		process.LastSyncAt = &lastSyncAt.Time
	}

	if archivedAt.Valid {
		process.ArchivedAt = &archivedAt.Time
	}

	process.Tags = tags

	return &process, nil
}

// scanProcesses converte múltiplas linhas para entidades Process
func (r *ProcessRepository) scanProcesses(rows *sql.Rows) ([]*domain.Process, error) {
	var processes []*domain.Process

	for rows.Next() {
		var process domain.Process
		var subjectJSON, valueJSON, monitoringJSON, customFieldsJSON []byte
		var tags []string
		var originalNumber sql.NullString
		var judgeID sql.NullString
		var lastMovementAt, lastSyncAt, archivedAt sql.NullTime

		err := rows.Scan(
			&process.ID,
			&process.TenantID,
			&process.ClientID,
			&process.Number,
			&originalNumber,
			&process.Title,
			&process.Description,
			&process.Status,
			&process.Stage,
			&subjectJSON,
			&valueJSON,
			&process.CourtID,
			&judgeID,
			&monitoringJSON,
			pq.Array(&tags),
			&customFieldsJSON,
			&lastMovementAt,
			&lastSyncAt,
			&process.CreatedAt,
			&process.UpdatedAt,
			&archivedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao escanear processo: %w", err)
		}

		// Deserializar campos JSON
		if err := json.Unmarshal(subjectJSON, &process.Subject); err != nil {
			return nil, fmt.Errorf("erro ao deserializar subject: %w", err)
		}

		if valueJSON != nil {
			var value domain.ProcessValue
			if err := json.Unmarshal(valueJSON, &value); err != nil {
				return nil, fmt.Errorf("erro ao deserializar value: %w", err)
			}
			process.Value = &value
		}

		if err := json.Unmarshal(monitoringJSON, &process.Monitoring); err != nil {
			return nil, fmt.Errorf("erro ao deserializar monitoring: %w", err)
		}

		if err := json.Unmarshal(customFieldsJSON, &process.CustomFields); err != nil {
			return nil, fmt.Errorf("erro ao deserializar custom_fields: %w", err)
		}

		// Tratar campos nullable
		if originalNumber.Valid {
			process.OriginalNumber = originalNumber.String
		}

		if judgeID.Valid {
			process.JudgeID = &judgeID.String
		}

		if lastMovementAt.Valid {
			process.LastMovementAt = &lastMovementAt.Time
		}

		if lastSyncAt.Valid {
			process.LastSyncAt = &lastSyncAt.Time
		}

		if archivedAt.Valid {
			process.ArchivedAt = &archivedAt.Time
		}

		process.Tags = tags

		processes = append(processes, &process)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar processos: %w", err)
	}

	return processes, nil
}