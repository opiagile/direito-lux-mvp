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

// MovementRepository implementação PostgreSQL do repositório de movimentações
type MovementRepository struct {
	db *sql.DB
}

// NewMovementRepository cria nova instância do repositório
func NewMovementRepository(db *sql.DB) *MovementRepository {
	return &MovementRepository{db: db}
}

// Create cria nova movimentação
func (r *MovementRepository) Create(movement *domain.Movement) error {
	query := `
		INSERT INTO movements (
			id, process_id, tenant_id, sequence, external_id, date, type, code,
			title, description, content, judge, responsible, attachments,
			related_parties, is_important, is_public, notification_sent, tags,
			metadata, created_at, updated_at, synced_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23
		)`

	// Serializar campos JSON
	attachmentsJSON, err := json.Marshal(movement.Attachments)
	if err != nil {
		return fmt.Errorf("erro ao serializar attachments: %w", err)
	}

	metadataJSON, err := json.Marshal(movement.Metadata)
	if err != nil {
		return fmt.Errorf("erro ao serializar metadata: %w", err)
	}

	_, err = r.db.Exec(
		query,
		movement.ID,
		movement.ProcessID,
		movement.TenantID,
		movement.Sequence,
		movement.ExternalID,
		movement.Date,
		movement.Type,
		movement.Code,
		movement.Title,
		movement.Description,
		movement.Content,
		movement.Judge,
		movement.Responsible,
		attachmentsJSON,
		pq.Array(movement.RelatedParties),
		movement.IsImportant,
		movement.IsPublic,
		movement.NotificationSent,
		pq.Array(movement.Tags),
		metadataJSON,
		movement.CreatedAt,
		movement.UpdatedAt,
		movement.SyncedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao criar movimentação: %w", err)
	}

	return nil
}

// CreateBatch cria múltiplas movimentações em uma transação
func (r *MovementRepository) CreateBatch(movements []*domain.Movement) error {
	if len(movements) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO movements (
			id, process_id, tenant_id, sequence, external_id, date, type, code,
			title, description, content, judge, responsible, attachments,
			related_parties, is_important, is_public, notification_sent, tags,
			metadata, created_at, updated_at, synced_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23
		)`)
	if err != nil {
		return fmt.Errorf("erro ao preparar statement: %w", err)
	}
	defer stmt.Close()

	for _, movement := range movements {
		// Serializar campos JSON
		attachmentsJSON, err := json.Marshal(movement.Attachments)
		if err != nil {
			return fmt.Errorf("erro ao serializar attachments: %w", err)
		}

		metadataJSON, err := json.Marshal(movement.Metadata)
		if err != nil {
			return fmt.Errorf("erro ao serializar metadata: %w", err)
		}

		_, err = stmt.Exec(
			movement.ID,
			movement.ProcessID,
			movement.TenantID,
			movement.Sequence,
			movement.ExternalID,
			movement.Date,
			movement.Type,
			movement.Code,
			movement.Title,
			movement.Description,
			movement.Content,
			movement.Judge,
			movement.Responsible,
			attachmentsJSON,
			pq.Array(movement.RelatedParties),
			movement.IsImportant,
			movement.IsPublic,
			movement.NotificationSent,
			pq.Array(movement.Tags),
			metadataJSON,
			movement.CreatedAt,
			movement.UpdatedAt,
			movement.SyncedAt,
		)

		if err != nil {
			return fmt.Errorf("erro ao criar movimentação em lote: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil
}

// Update atualiza movimentação existente
func (r *MovementRepository) Update(movement *domain.Movement) error {
	query := `
		UPDATE movements SET
			title = $2, description = $3, content = $4, judge = $5, responsible = $6,
			attachments = $7, is_important = $8, notification_sent = $9, tags = $10,
			metadata = $11, updated_at = $12
		WHERE id = $1 AND tenant_id = $13`

	// Serializar campos JSON
	attachmentsJSON, err := json.Marshal(movement.Attachments)
	if err != nil {
		return fmt.Errorf("erro ao serializar attachments: %w", err)
	}

	metadataJSON, err := json.Marshal(movement.Metadata)
	if err != nil {
		return fmt.Errorf("erro ao serializar metadata: %w", err)
	}

	result, err := r.db.Exec(
		query,
		movement.ID,
		movement.Title,
		movement.Description,
		movement.Content,
		movement.Judge,
		movement.Responsible,
		attachmentsJSON,
		movement.IsImportant,
		movement.NotificationSent,
		pq.Array(movement.Tags),
		metadataJSON,
		movement.UpdatedAt,
		movement.TenantID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar movimentação: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrMovementNotFound
	}

	return nil
}

// Delete remove movimentação
func (r *MovementRepository) Delete(id string) error {
	query := `DELETE FROM movements WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar movimentação: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrMovementNotFound
	}

	return nil
}

// GetByID busca movimentação por ID
func (r *MovementRepository) GetByID(id string) (*domain.Movement, error) {
	query := `
		SELECT id, process_id, tenant_id, sequence, external_id, date, type, code,
			title, description, content, judge, responsible, attachments,
			related_parties, is_important, is_public, notification_sent, tags,
			metadata, created_at, updated_at, synced_at
		FROM movements
		WHERE id = $1`

	return r.scanMovement(r.db.QueryRow(query, id))
}

// GetByProcess busca movimentações por processo com filtros
func (r *MovementRepository) GetByProcess(processID string, filters domain.MovementFilters) ([]*domain.Movement, error) {
	query, args := r.buildQuery("process_id = $1", processID, filters)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar movimentações: %w", err)
	}
	defer rows.Close()

	return r.scanMovements(rows)
}

// GetByTenant busca movimentações por tenant com filtros
func (r *MovementRepository) GetByTenant(tenantID string, filters domain.MovementFilters) ([]*domain.Movement, error) {
	query, args := r.buildQuery("tenant_id = $1", tenantID, filters)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar movimentações: %w", err)
	}
	defer rows.Close()

	return r.scanMovements(rows)
}

// GetByExternalID busca movimentação por ID externo
func (r *MovementRepository) GetByExternalID(externalID string) (*domain.Movement, error) {
	query := `
		SELECT id, process_id, tenant_id, sequence, external_id, date, type, code,
			title, description, content, judge, responsible, attachments,
			related_parties, is_important, is_public, notification_sent, tags,
			metadata, created_at, updated_at, synced_at
		FROM movements
		WHERE external_id = $1`

	return r.scanMovement(r.db.QueryRow(query, externalID))
}

// GetImportantMovements busca movimentações importantes por processo
func (r *MovementRepository) GetImportantMovements(processID string) ([]*domain.Movement, error) {
	query := `
		SELECT id, process_id, tenant_id, sequence, external_id, date, type, code,
			title, description, content, judge, responsible, attachments,
			related_parties, is_important, is_public, notification_sent, tags,
			metadata, created_at, updated_at, synced_at
		FROM movements
		WHERE process_id = $1 AND is_important = true
		ORDER BY date DESC`

	rows, err := r.db.Query(query, processID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar movimentações importantes: %w", err)
	}
	defer rows.Close()

	return r.scanMovements(rows)
}

// GetRecentMovements busca movimentações recentes por tenant
func (r *MovementRepository) GetRecentMovements(tenantID string, limit int) ([]*domain.Movement, error) {
	query := `
		SELECT id, process_id, tenant_id, sequence, external_id, date, type, code,
			title, description, content, judge, responsible, attachments,
			related_parties, is_important, is_public, notification_sent, tags,
			metadata, created_at, updated_at, synced_at
		FROM movements
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.db.Query(query, tenantID, limit)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar movimentações recentes: %w", err)
	}
	defer rows.Close()

	return r.scanMovements(rows)
}

// GetPendingNotifications busca movimentações com notificações pendentes
func (r *MovementRepository) GetPendingNotifications(tenantID string) ([]*domain.Movement, error) {
	query := `
		SELECT id, process_id, tenant_id, sequence, external_id, date, type, code,
			title, description, content, judge, responsible, attachments,
			related_parties, is_important, is_public, notification_sent, tags,
			metadata, created_at, updated_at, synced_at
		FROM movements
		WHERE tenant_id = $1 AND is_important = true AND notification_sent = false
		ORDER BY date DESC`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar notificações pendentes: %w", err)
	}
	defer rows.Close()

	return r.scanMovements(rows)
}

// SearchByContent busca movimentações por conteúdo textual
func (r *MovementRepository) SearchByContent(tenantID string, query string, filters domain.MovementFilters) ([]*domain.Movement, error) {
	// Adicionar busca textual aos filtros
	filters.Search = query
	querySQL, args := r.buildQuery("tenant_id = $1", tenantID, filters)
	
	rows, err := r.db.Query(querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("erro na busca por conteúdo: %w", err)
	}
	defer rows.Close()

	return r.scanMovements(rows)
}

// GetMovementsByKeywords busca movimentações por palavras-chave
func (r *MovementRepository) GetMovementsByKeywords(tenantID string, keywords []string) ([]*domain.Movement, error) {
	query := `
		SELECT id, process_id, tenant_id, sequence, external_id, date, type, code,
			title, description, content, judge, responsible, attachments,
			related_parties, is_important, is_public, notification_sent, tags,
			metadata, created_at, updated_at, synced_at
		FROM movements
		WHERE tenant_id = $1 AND (
			title ILIKE ANY($2) OR 
			description ILIKE ANY($2) OR 
			content ILIKE ANY($2)
		)
		ORDER BY date DESC`

	// Preparar keywords com wildcards
	keywordPatterns := make([]string, len(keywords))
	for i, keyword := range keywords {
		keywordPatterns[i] = "%" + keyword + "%"
	}

	rows, err := r.db.Query(query, tenantID, pq.Array(keywordPatterns))
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar por keywords: %w", err)
	}
	defer rows.Close()

	return r.scanMovements(rows)
}

// GetStatistics busca estatísticas de movimentações
func (r *MovementRepository) GetStatistics(tenantID string, dateFrom, dateTo time.Time) (*domain.MovementStatistics, error) {
	// Query para estatísticas básicas
	statsQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN is_important THEN 1 END) as important,
			COUNT(CASE WHEN notification_sent THEN 1 END) as notifications_sent
		FROM movements
		WHERE tenant_id = $1 AND date BETWEEN $2 AND $3`

	var stats domain.MovementStatistics
	err := r.db.QueryRow(statsQuery, tenantID, dateFrom, dateTo).Scan(
		&stats.TotalMovements,
		&stats.ImportantMovements,
		&stats.NotificationsSent,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estatísticas básicas: %w", err)
	}

	// Query para movimentações por tipo
	typeQuery := `
		SELECT type, COUNT(*)
		FROM movements
		WHERE tenant_id = $1 AND date BETWEEN $2 AND $3
		GROUP BY type`

	rows, err := r.db.Query(typeQuery, tenantID, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estatísticas por tipo: %w", err)
	}
	defer rows.Close()

	stats.MovementsByType = make(map[domain.MovementType]int)
	for rows.Next() {
		var movType domain.MovementType
		var count int
		if err := rows.Scan(&movType, &count); err != nil {
			return nil, fmt.Errorf("erro ao escanear tipo: %w", err)
		}
		stats.MovementsByType[movType] = count
	}

	// Query para movimentações por mês
	monthQuery := `
		SELECT 
			TO_CHAR(date, 'YYYY-MM') as month,
			COUNT(*)
		FROM movements
		WHERE tenant_id = $1 AND date BETWEEN $2 AND $3
		GROUP BY TO_CHAR(date, 'YYYY-MM')
		ORDER BY month`

	rows, err = r.db.Query(monthQuery, tenantID, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estatísticas por mês: %w", err)
	}
	defer rows.Close()

	stats.MovementsByMonth = make(map[string]int)
	for rows.Next() {
		var month string
		var count int
		if err := rows.Scan(&month, &count); err != nil {
			return nil, fmt.Errorf("erro ao escanear mês: %w", err)
		}
		stats.MovementsByMonth[month] = count
	}

	// Calcular média por processo (simplificado)
	if stats.TotalMovements > 0 {
		processCountQuery := `
			SELECT COUNT(DISTINCT process_id)
			FROM movements
			WHERE tenant_id = $1 AND date BETWEEN $2 AND $3`

		var processCount int
		err := r.db.QueryRow(processCountQuery, tenantID, dateFrom, dateTo).Scan(&processCount)
		if err == nil && processCount > 0 {
			stats.AveragePerProcess = float64(stats.TotalMovements) / float64(processCount)
		}
	}

	return &stats, nil
}

// buildQuery constrói query com filtros
func (r *MovementRepository) buildQuery(baseCondition, baseValue string, filters domain.MovementFilters) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 2 // Começa em 2 pois $1 é usado na condição base

	conditions = append(conditions, baseCondition)
	args = append(args, baseValue)

	// Filtros por tipo
	if len(filters.Type) > 0 {
		typePlaceholders := make([]string, len(filters.Type))
		for i, moveType := range filters.Type {
			typePlaceholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, moveType)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("type IN (%s)", strings.Join(typePlaceholders, ",")))
	}

	// Filtro por importância
	if filters.IsImportant != nil {
		conditions = append(conditions, fmt.Sprintf("is_important = $%d", argIndex))
		args = append(args, *filters.IsImportant)
		argIndex++
	}

	// Filtro por público
	if filters.IsPublic != nil {
		conditions = append(conditions, fmt.Sprintf("is_public = $%d", argIndex))
		args = append(args, *filters.IsPublic)
		argIndex++
	}

	// Filtro por notificação
	if filters.HasNotification != nil {
		conditions = append(conditions, fmt.Sprintf("notification_sent = $%d", argIndex))
		args = append(args, *filters.HasNotification)
		argIndex++
	}

	// Filtros de data
	if filters.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("date >= $%d", argIndex))
		args = append(args, *filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("date <= $%d", argIndex))
		args = append(args, *filters.DateTo)
		argIndex++
	}

	// Filtro por juiz
	if filters.Judge != "" {
		conditions = append(conditions, fmt.Sprintf("judge ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Judge+"%")
		argIndex++
	}

	// Filtro por tags
	if len(filters.Tags) > 0 {
		conditions = append(conditions, fmt.Sprintf("tags && $%d", argIndex))
		args = append(args, pq.Array(filters.Tags))
		argIndex++
	}

	// Busca textual
	if filters.Search != "" {
		searchCondition := fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d OR content ILIKE $%d)", argIndex, argIndex+1, argIndex+2)
		searchTerm := "%" + filters.Search + "%"
		conditions = append(conditions, searchCondition)
		args = append(args, searchTerm, searchTerm, searchTerm)
		argIndex += 3
	}

	// Construir query final
	query := `
		SELECT id, process_id, tenant_id, sequence, external_id, date, type, code,
			title, description, content, judge, responsible, attachments,
			related_parties, is_important, is_public, notification_sent, tags,
			metadata, created_at, updated_at, synced_at
		FROM movements`

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
		query += " ORDER BY date DESC"
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

// scanMovement converte linha do banco para entidade Movement
func (r *MovementRepository) scanMovement(row *sql.Row) (*domain.Movement, error) {
	var movement domain.Movement
	var attachmentsJSON, metadataJSON []byte
	var relatedParties, tags []string
	var externalID, judge, responsible, content sql.NullString

	err := row.Scan(
		&movement.ID,
		&movement.ProcessID,
		&movement.TenantID,
		&movement.Sequence,
		&externalID,
		&movement.Date,
		&movement.Type,
		&movement.Code,
		&movement.Title,
		&movement.Description,
		&content,
		&judge,
		&responsible,
		&attachmentsJSON,
		pq.Array(&relatedParties),
		&movement.IsImportant,
		&movement.IsPublic,
		&movement.NotificationSent,
		pq.Array(&tags),
		&metadataJSON,
		&movement.CreatedAt,
		&movement.UpdatedAt,
		&movement.SyncedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrMovementNotFound
		}
		return nil, fmt.Errorf("erro ao escanear movimentação: %w", err)
	}

	// Deserializar campos JSON
	if err := json.Unmarshal(attachmentsJSON, &movement.Attachments); err != nil {
		return nil, fmt.Errorf("erro ao deserializar attachments: %w", err)
	}

	if err := json.Unmarshal(metadataJSON, &movement.Metadata); err != nil {
		return nil, fmt.Errorf("erro ao deserializar metadata: %w", err)
	}

	// Tratar campos nullable
	if externalID.Valid {
		movement.ExternalID = externalID.String
	}

	if judge.Valid {
		movement.Judge = judge.String
	}

	if responsible.Valid {
		movement.Responsible = responsible.String
	}

	if content.Valid {
		movement.Content = content.String
	}

	movement.RelatedParties = relatedParties
	movement.Tags = tags

	return &movement, nil
}

// scanMovements converte múltiplas linhas para entidades Movement
func (r *MovementRepository) scanMovements(rows *sql.Rows) ([]*domain.Movement, error) {
	var movements []*domain.Movement

	for rows.Next() {
		var movement domain.Movement
		var attachmentsJSON, metadataJSON []byte
		var relatedParties, tags []string
		var externalID, judge, responsible, content sql.NullString

		err := rows.Scan(
			&movement.ID,
			&movement.ProcessID,
			&movement.TenantID,
			&movement.Sequence,
			&externalID,
			&movement.Date,
			&movement.Type,
			&movement.Code,
			&movement.Title,
			&movement.Description,
			&content,
			&judge,
			&responsible,
			&attachmentsJSON,
			pq.Array(&relatedParties),
			&movement.IsImportant,
			&movement.IsPublic,
			&movement.NotificationSent,
			pq.Array(&tags),
			&metadataJSON,
			&movement.CreatedAt,
			&movement.UpdatedAt,
			&movement.SyncedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao escanear movimentação: %w", err)
		}

		// Deserializar campos JSON
		if err := json.Unmarshal(attachmentsJSON, &movement.Attachments); err != nil {
			return nil, fmt.Errorf("erro ao deserializar attachments: %w", err)
		}

		if err := json.Unmarshal(metadataJSON, &movement.Metadata); err != nil {
			return nil, fmt.Errorf("erro ao deserializar metadata: %w", err)
		}

		// Tratar campos nullable
		if externalID.Valid {
			movement.ExternalID = externalID.String
		}

		if judge.Valid {
			movement.Judge = judge.String
		}

		if responsible.Valid {
			movement.Responsible = responsible.String
		}

		if content.Valid {
			movement.Content = content.String
		}

		movement.RelatedParties = relatedParties
		movement.Tags = tags

		movements = append(movements, &movement)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar movimentações: %w", err)
	}

	return movements, nil
}