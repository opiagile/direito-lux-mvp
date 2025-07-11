package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// DataJudRequestRepository implementação PostgreSQL do repositório de DataJudRequest
type DataJudRequestRepository struct {
	db *sqlx.DB
}

// NewDataJudRequestRepository cria nova instância do repositório
func NewDataJudRequestRepository(db *sqlx.DB) *DataJudRequestRepository {
	return &DataJudRequestRepository{db: db}
}

// Save salva uma requisição DataJud
func (r *DataJudRequestRepository) Save(request *domain.DataJudRequest) error {
	// Serializar parâmetros como JSON
	parametersJSON, err := json.Marshal(request.Parameters)
	if err != nil {
		return fmt.Errorf("erro ao serializar parâmetros: %w", err)
	}

	// Serializar response como JSON se existir
	var responseJSON []byte
	if request.Response != nil {
		responseJSON, err = json.Marshal(request.Response)
		if err != nil {
			return fmt.Errorf("erro ao serializar resposta: %w", err)
		}
	}

	query := `
		INSERT INTO datajud_requests (
			id, tenant_id, client_id, process_id, type, priority, status, cnpj_provider_id,
			process_number, court_id, parameters, cache_key, cache_ttl, use_cache,
			retry_count, max_retries, retry_after, circuit_breaker_key,
			requested_at, processing_at, completed_at, created_at, updated_at,
			response_data, error_message, error_code
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24, $25, $26
		)`

	_, err = r.db.Exec(query,
		request.ID, request.TenantID, request.ClientID, request.ProcessID,
		request.Type, request.Priority, request.Status, request.CNPJProviderID,
		request.ProcessNumber, request.CourtID, parametersJSON, request.CacheKey,
		request.CacheTTL, request.UseCache, request.RetryCount, request.MaxRetries,
		request.RetryAfter, request.CircuitBreakerKey, request.RequestedAt,
		request.ProcessingAt, request.CompletedAt, request.CreatedAt, request.UpdatedAt,
		responseJSON, request.ErrorMessage, request.ErrorCode,
	)

	return err
}

// FindByID encontra requisição por ID
func (r *DataJudRequestRepository) FindByID(id uuid.UUID) (*domain.DataJudRequest, error) {
	query := `
		SELECT id, tenant_id, client_id, process_id, type, priority, status, cnpj_provider_id,
			   process_number, court_id, parameters, cache_key, cache_ttl, use_cache,
			   retry_count, max_retries, retry_after, circuit_breaker_key,
			   requested_at, processing_at, completed_at, created_at, updated_at,
			   response_data, error_message, error_code
		FROM datajud_requests 
		WHERE id = $1`

	return r.scanRequest(r.db.QueryRow(query, id))
}

// FindByTenantID encontra requisições de um tenant
func (r *DataJudRequestRepository) FindByTenantID(tenantID uuid.UUID, limit, offset int) ([]*domain.DataJudRequest, error) {
	query := `
		SELECT id, tenant_id, client_id, process_id, type, priority, status, cnpj_provider_id,
			   process_number, court_id, parameters, cache_key, cache_ttl, use_cache,
			   retry_count, max_retries, retry_after, circuit_breaker_key,
			   requested_at, processing_at, completed_at, created_at, updated_at,
			   response_data, error_message, error_code
		FROM datajud_requests 
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRequests(rows)
}

// FindPendingRequests encontra requisições pendentes
func (r *DataJudRequestRepository) FindPendingRequests(limit int) ([]*domain.DataJudRequest, error) {
	query := `
		SELECT id, tenant_id, client_id, process_id, type, priority, status, cnpj_provider_id,
			   process_number, court_id, parameters, cache_key, cache_ttl, use_cache,
			   retry_count, max_retries, retry_after, circuit_breaker_key,
			   requested_at, processing_at, completed_at, created_at, updated_at,
			   response_data, error_message, error_code
		FROM datajud_requests 
		WHERE status IN ('pending', 'retrying')
		  AND (retry_after IS NULL OR retry_after <= NOW())
		ORDER BY priority DESC, requested_at ASC
		LIMIT $1`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRequests(rows)
}

// FindByStatus encontra requisições por status
func (r *DataJudRequestRepository) FindByStatus(status domain.RequestStatus, limit, offset int) ([]*domain.DataJudRequest, error) {
	query := `
		SELECT id, tenant_id, client_id, process_id, type, priority, status, cnpj_provider_id,
			   process_number, court_id, parameters, cache_key, cache_ttl, use_cache,
			   retry_count, max_retries, retry_after, circuit_breaker_key,
			   requested_at, processing_at, completed_at, created_at, updated_at,
			   response_data, error_message, error_code
		FROM datajud_requests 
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRequests(rows)
}

// FindByProcessNumber encontra requisições por número de processo
func (r *DataJudRequestRepository) FindByProcessNumber(processNumber string) ([]*domain.DataJudRequest, error) {
	query := `
		SELECT id, tenant_id, client_id, process_id, type, priority, status, cnpj_provider_id,
			   process_number, court_id, parameters, cache_key, cache_ttl, use_cache,
			   retry_count, max_retries, retry_after, circuit_breaker_key,
			   requested_at, processing_at, completed_at, created_at, updated_at,
			   response_data, error_message, error_code
		FROM datajud_requests 
		WHERE process_number = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, processNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRequests(rows)
}

// Update atualiza uma requisição
func (r *DataJudRequestRepository) Update(request *domain.DataJudRequest) error {
	// Serializar parâmetros como JSON
	parametersJSON, err := json.Marshal(request.Parameters)
	if err != nil {
		return fmt.Errorf("erro ao serializar parâmetros: %w", err)
	}

	// Serializar response como JSON se existir
	var responseJSON []byte
	if request.Response != nil {
		responseJSON, err = json.Marshal(request.Response)
		if err != nil {
			return fmt.Errorf("erro ao serializar resposta: %w", err)
		}
	}

	query := `
		UPDATE datajud_requests 
		SET tenant_id = $2,
			client_id = $3,
			process_id = $4,
			type = $5,
			priority = $6,
			status = $7,
			cnpj_provider_id = $8,
			process_number = $9,
			court_id = $10,
			parameters = $11,
			cache_key = $12,
			cache_ttl = $13,
			use_cache = $14,
			retry_count = $15,
			max_retries = $16,
			retry_after = $17,
			circuit_breaker_key = $18,
			requested_at = $19,
			processing_at = $20,
			completed_at = $21,
			updated_at = $22,
			response_data = $23,
			error_message = $24,
			error_code = $25
		WHERE id = $1`

	result, err := r.db.Exec(query,
		request.ID, request.TenantID, request.ClientID, request.ProcessID,
		request.Type, request.Priority, request.Status, request.CNPJProviderID,
		request.ProcessNumber, request.CourtID, parametersJSON, request.CacheKey,
		request.CacheTTL, request.UseCache, request.RetryCount, request.MaxRetries,
		request.RetryAfter, request.CircuitBreakerKey, request.RequestedAt,
		request.ProcessingAt, request.CompletedAt, request.UpdatedAt,
		responseJSON, request.ErrorMessage, request.ErrorCode,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("requisição %s não encontrada", request.ID)
	}

	return nil
}

// UpdateStatus atualiza apenas o status de uma requisição
func (r *DataJudRequestRepository) UpdateStatus(id uuid.UUID, status domain.RequestStatus) error {
	query := `
		UPDATE datajud_requests 
		SET status = $2, updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.Exec(query, id, status)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("requisição %s não encontrada", id)
	}

	return nil
}

// Delete remove uma requisição
func (r *DataJudRequestRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM datajud_requests WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("requisição %s não encontrada", id)
	}

	return nil
}

// CleanupOldRequests limpa requisições antigas
func (r *DataJudRequestRepository) CleanupOldRequests(olderThan time.Time) (int, error) {
	query := `
		DELETE FROM datajud_requests 
		WHERE status IN ('completed', 'failed') 
		  AND completed_at < $1`

	result, err := r.db.Exec(query, olderThan)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

// GetStatistics obtém estatísticas das requisições
func (r *DataJudRequestRepository) GetStatistics(tenantID *uuid.UUID, from, to *time.Time) (map[string]interface{}, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if tenantID != nil {
		whereClause += fmt.Sprintf(" AND tenant_id = $%d", argIndex)
		args = append(args, *tenantID)
		argIndex++
	}

	if from != nil {
		whereClause += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *from)
		argIndex++
	}

	if to != nil {
		whereClause += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *to)
		argIndex++
	}

	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_requests,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_requests,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_requests,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_requests,
			COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing_requests,
			COUNT(CASE WHEN status = 'retrying' THEN 1 END) as retrying_requests,
			AVG(EXTRACT(EPOCH FROM (completed_at - requested_at)) * 1000) as avg_response_time_ms,
			COUNT(CASE WHEN type = 'process' THEN 1 END) as process_requests,
			COUNT(CASE WHEN type = 'movement' THEN 1 END) as movement_requests,
			COUNT(CASE WHEN type = 'party' THEN 1 END) as party_requests,
			COUNT(CASE WHEN type = 'document' THEN 1 END) as document_requests,
			COUNT(CASE WHEN type = 'bulk' THEN 1 END) as bulk_requests
		FROM datajud_requests %s`, whereClause)

	stats := make(map[string]interface{})
	row := r.db.QueryRow(query, args...)

	var totalRequests, completedRequests, failedRequests, pendingRequests int
	var processingRequests, retryingRequests int
	var processRequests, movementRequests, partyRequests, documentRequests, bulkRequests int
	var avgResponseTime sql.NullFloat64

	err := row.Scan(
		&totalRequests, &completedRequests, &failedRequests, &pendingRequests,
		&processingRequests, &retryingRequests, &avgResponseTime,
		&processRequests, &movementRequests, &partyRequests, &documentRequests, &bulkRequests,
	)
	if err != nil {
		return nil, err
	}

	stats["total_requests"] = totalRequests
	stats["completed_requests"] = completedRequests
	stats["failed_requests"] = failedRequests
	stats["pending_requests"] = pendingRequests
	stats["processing_requests"] = processingRequests
	stats["retrying_requests"] = retryingRequests

	if avgResponseTime.Valid {
		stats["avg_response_time_ms"] = avgResponseTime.Float64
	} else {
		stats["avg_response_time_ms"] = 0.0
	}

	if totalRequests > 0 {
		stats["success_rate"] = float64(completedRequests) / float64(totalRequests) * 100
		stats["failure_rate"] = float64(failedRequests) / float64(totalRequests) * 100
	} else {
		stats["success_rate"] = 0.0
		stats["failure_rate"] = 0.0
	}

	stats["by_type"] = map[string]int{
		"process":  processRequests,
		"movement": movementRequests,
		"party":    partyRequests,
		"document": documentRequests,
		"bulk":     bulkRequests,
	}

	return stats, nil
}

// GetHourlyStats obtém estatísticas por hora
func (r *DataJudRequestRepository) GetHourlyStats(hours int) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			DATE_TRUNC('hour', created_at) as hour,
			COUNT(*) as total_requests,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_requests,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_requests,
			AVG(EXTRACT(EPOCH FROM (completed_at - requested_at)) * 1000) as avg_response_time_ms
		FROM datajud_requests 
		WHERE created_at >= NOW() - INTERVAL '%d hours'
		GROUP BY DATE_TRUNC('hour', created_at)
		ORDER BY hour ASC`

	rows, err := r.db.Query(fmt.Sprintf(query, hours))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := []map[string]interface{}{}
	for rows.Next() {
		var hour time.Time
		var totalRequests, completedRequests, failedRequests int
		var avgResponseTime sql.NullFloat64

		err := rows.Scan(&hour, &totalRequests, &completedRequests, &failedRequests, &avgResponseTime)
		if err != nil {
			return nil, err
		}

		stat := map[string]interface{}{
			"hour":               hour,
			"total_requests":     totalRequests,
			"completed_requests": completedRequests,
			"failed_requests":    failedRequests,
		}

		if avgResponseTime.Valid {
			stat["avg_response_time_ms"] = avgResponseTime.Float64
		} else {
			stat["avg_response_time_ms"] = 0.0
		}

		if totalRequests > 0 {
			stat["success_rate"] = float64(completedRequests) / float64(totalRequests) * 100
		} else {
			stat["success_rate"] = 0.0
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

// scanRequest converte uma linha em DataJudRequest
func (r *DataJudRequestRepository) scanRequest(row *sql.Row) (*domain.DataJudRequest, error) {
	request := &domain.DataJudRequest{}
	var parametersJSON []byte
	var responseJSON []byte

	err := row.Scan(
		&request.ID, &request.TenantID, &request.ClientID, &request.ProcessID,
		&request.Type, &request.Priority, &request.Status, &request.CNPJProviderID,
		&request.ProcessNumber, &request.CourtID, &parametersJSON, &request.CacheKey,
		&request.CacheTTL, &request.UseCache, &request.RetryCount, &request.MaxRetries,
		&request.RetryAfter, &request.CircuitBreakerKey, &request.RequestedAt,
		&request.ProcessingAt, &request.CompletedAt, &request.CreatedAt, &request.UpdatedAt,
		&responseJSON, &request.ErrorMessage, &request.ErrorCode,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Deserializar parâmetros
	if len(parametersJSON) > 0 {
		err = json.Unmarshal(parametersJSON, &request.Parameters)
		if err != nil {
			return nil, fmt.Errorf("erro ao deserializar parâmetros: %w", err)
		}
	}

	// Deserializar response se existir
	if len(responseJSON) > 0 {
		request.Response = &domain.DataJudResponse{}
		err = json.Unmarshal(responseJSON, request.Response)
		if err != nil {
			return nil, fmt.Errorf("erro ao deserializar resposta: %w", err)
		}
	}

	return request, nil
}

// scanRequests converte múltiplas linhas em DataJudRequests
func (r *DataJudRequestRepository) scanRequests(rows *sql.Rows) ([]*domain.DataJudRequest, error) {
	requests := []*domain.DataJudRequest{}

	for rows.Next() {
		request := &domain.DataJudRequest{}
		var parametersJSON []byte
		var responseJSON []byte

		err := rows.Scan(
			&request.ID, &request.TenantID, &request.ClientID, &request.ProcessID,
			&request.Type, &request.Priority, &request.Status, &request.CNPJProviderID,
			&request.ProcessNumber, &request.CourtID, &parametersJSON, &request.CacheKey,
			&request.CacheTTL, &request.UseCache, &request.RetryCount, &request.MaxRetries,
			&request.RetryAfter, &request.CircuitBreakerKey, &request.RequestedAt,
			&request.ProcessingAt, &request.CompletedAt, &request.CreatedAt, &request.UpdatedAt,
			&responseJSON, &request.ErrorMessage, &request.ErrorCode,
		)
		if err != nil {
			return nil, err
		}

		// Deserializar parâmetros
		if len(parametersJSON) > 0 {
			err = json.Unmarshal(parametersJSON, &request.Parameters)
			if err != nil {
				return nil, fmt.Errorf("erro ao deserializar parâmetros: %w", err)
			}
		}

		// Deserializar response se existir
		if len(responseJSON) > 0 {
			request.Response = &domain.DataJudResponse{}
			err = json.Unmarshal(responseJSON, request.Response)
			if err != nil {
				return nil, fmt.Errorf("erro ao deserializar resposta: %w", err)
			}
		}

		requests = append(requests, request)
	}

	return requests, nil
}