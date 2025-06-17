package application

import (
	"context"
	"github.com/direito-lux/process-service/internal/application/commands"
	"github.com/direito-lux/process-service/internal/application/queries"
	"github.com/direito-lux/process-service/internal/domain"
)

// ProcessService serviço principal da aplicação implementando CQRS
type ProcessService struct {
	// Command Handlers
	processCommandHandler  *commands.ProcessCommandHandler
	movementCommandHandler *commands.MovementCommandHandler
	syncCommandHandler     *commands.SyncCommandHandler
	partyCommandHandler    *commands.PartyCommandHandler

	// Query Handlers
	processQueryHandler   *queries.ProcessQueryHandler
	movementQueryHandler  *queries.MovementQueryHandler
	dashboardQueryHandler *queries.DashboardQueryHandler
}

// NewProcessService cria nova instância do serviço
func NewProcessService(
	processRepo domain.ProcessRepository,
	movementRepo domain.MovementRepository,
	partyRepo domain.PartyRepository,
	eventPublisher domain.EventPublisher,
) *ProcessService {
	return &ProcessService{
		// Command Handlers
		processCommandHandler: commands.NewProcessCommandHandler(
			processRepo, movementRepo, partyRepo, eventPublisher,
		),
		movementCommandHandler: commands.NewMovementCommandHandler(
			movementRepo, processRepo, eventPublisher,
		),
		syncCommandHandler: commands.NewSyncCommandHandler(
			processRepo, movementRepo, eventPublisher,
		),
		partyCommandHandler: commands.NewPartyCommandHandler(
			partyRepo, processRepo, eventPublisher,
		),

		// Query Handlers
		processQueryHandler: queries.NewProcessQueryHandler(
			processRepo, movementRepo, partyRepo,
		),
		movementQueryHandler: queries.NewMovementQueryHandler(
			movementRepo, processRepo,
		),
		dashboardQueryHandler: queries.NewDashboardQueryHandler(
			processRepo, movementRepo, partyRepo,
		),
	}
}

// === PROCESS COMMANDS ===

// CreateProcess executa comando de criação de processo
func (s *ProcessService) CreateProcess(ctx context.Context, cmd *commands.CreateProcessCommand) error {
	return s.processCommandHandler.HandleCreateProcess(ctx, cmd)
}

// UpdateProcess executa comando de atualização de processo
func (s *ProcessService) UpdateProcess(ctx context.Context, cmd *commands.UpdateProcessCommand) error {
	return s.processCommandHandler.HandleUpdateProcess(ctx, cmd)
}

// ArchiveProcess executa comando de arquivamento de processo
func (s *ProcessService) ArchiveProcess(ctx context.Context, cmd *commands.ArchiveProcessCommand) error {
	return s.processCommandHandler.HandleArchiveProcess(ctx, cmd)
}

// ReactivateProcess executa comando de reativação de processo
func (s *ProcessService) ReactivateProcess(ctx context.Context, cmd *commands.ReactivateProcessCommand) error {
	return s.processCommandHandler.HandleReactivateProcess(ctx, cmd)
}

// EnableMonitoring executa comando de habilitação do monitoramento
func (s *ProcessService) EnableMonitoring(ctx context.Context, cmd *commands.EnableMonitoringCommand) error {
	return s.processCommandHandler.HandleEnableMonitoring(ctx, cmd)
}

// DisableMonitoring executa comando de desabilitação do monitoramento
func (s *ProcessService) DisableMonitoring(ctx context.Context, cmd *commands.DisableMonitoringCommand) error {
	return s.processCommandHandler.HandleDisableMonitoring(ctx, cmd)
}

// === MOVEMENT COMMANDS ===

// CreateMovement executa comando de criação de movimentação
func (s *ProcessService) CreateMovement(ctx context.Context, cmd *commands.CreateMovementCommand) error {
	return s.movementCommandHandler.HandleCreateMovement(ctx, cmd)
}

// UpdateMovement executa comando de atualização de movimentação
func (s *ProcessService) UpdateMovement(ctx context.Context, cmd *commands.UpdateMovementCommand) error {
	return s.movementCommandHandler.HandleUpdateMovement(ctx, cmd)
}

// AnalyzeMovement executa comando de análise de movimentação
func (s *ProcessService) AnalyzeMovement(ctx context.Context, cmd *commands.AnalyzeMovementCommand) error {
	return s.movementCommandHandler.HandleAnalyzeMovement(ctx, cmd)
}

// === SYNC COMMANDS ===

// SyncProcess executa comando de sincronização de processo
func (s *ProcessService) SyncProcess(ctx context.Context, cmd *commands.SyncProcessCommand) error {
	return s.syncCommandHandler.HandleSyncProcess(ctx, cmd)
}

// BatchSyncProcesses executa comando de sincronização em lote
func (s *ProcessService) BatchSyncProcesses(ctx context.Context, cmd *commands.BatchSyncProcessesCommand) error {
	return s.syncCommandHandler.HandleBatchSyncProcesses(ctx, cmd)
}

// === PARTY COMMANDS ===

// AddParty executa comando de adição de parte
func (s *ProcessService) AddParty(ctx context.Context, cmd *commands.AddPartyCommand) error {
	return s.partyCommandHandler.HandleAddParty(ctx, cmd)
}

// UpdateParty executa comando de atualização de parte
func (s *ProcessService) UpdateParty(ctx context.Context, cmd *commands.UpdatePartyCommand) error {
	return s.partyCommandHandler.HandleUpdateParty(ctx, cmd)
}

// RemoveParty executa comando de remoção de parte
func (s *ProcessService) RemoveParty(ctx context.Context, cmd *commands.RemovePartyCommand) error {
	return s.partyCommandHandler.HandleRemoveParty(ctx, cmd)
}

// === PROCESS QUERIES ===

// GetProcess executa query de busca de processo
func (s *ProcessService) GetProcess(ctx context.Context, query *queries.ProcessQuery) (*queries.ProcessDTO, error) {
	return s.processQueryHandler.HandleGetProcess(ctx, query)
}

// ListProcesses executa query de listagem de processos
func (s *ProcessService) ListProcesses(ctx context.Context, query *queries.ProcessListQuery) (*queries.ProcessListDTO, error) {
	return s.processQueryHandler.HandleListProcesses(ctx, query)
}

// GetProcessStats executa query de estatísticas de processos
func (s *ProcessService) GetProcessStats(ctx context.Context, query *queries.ProcessStatsQuery) (*queries.ProcessSummaryDTO, error) {
	return s.processQueryHandler.HandleGetProcessStats(ctx, query)
}

// GetMonitoringProcesses executa query de processos para monitoramento
func (s *ProcessService) GetMonitoringProcesses(ctx context.Context, query *queries.ProcessMonitoringQuery) (*queries.ProcessListDTO, error) {
	return s.processQueryHandler.HandleGetMonitoringProcesses(ctx, query)
}

// === MOVEMENT QUERIES ===

// GetMovement executa query de busca de movimentação
func (s *ProcessService) GetMovement(ctx context.Context, query *queries.MovementQuery) (*queries.MovementDTO, error) {
	return s.movementQueryHandler.HandleGetMovement(ctx, query)
}

// ListMovements executa query de listagem de movimentações
func (s *ProcessService) ListMovements(ctx context.Context, query *queries.MovementListQuery) (*queries.MovementListDTO, error) {
	return s.movementQueryHandler.HandleListMovements(ctx, query)
}

// SearchMovements executa query de busca textual em movimentações
func (s *ProcessService) SearchMovements(ctx context.Context, query *queries.MovementSearchQuery) (*queries.SearchResultsDTO, error) {
	return s.movementQueryHandler.HandleSearchMovements(ctx, query)
}

// GetMovementStats executa query de estatísticas de movimentações
func (s *ProcessService) GetMovementStats(ctx context.Context, query *queries.MovementStatsQuery) (*queries.MovementSummaryDTO, error) {
	return s.movementQueryHandler.HandleGetMovementStats(ctx, query)
}

// === DASHBOARD QUERIES ===

// GetDashboard executa query do dashboard
func (s *ProcessService) GetDashboard(ctx context.Context, query *queries.DashboardQuery) (*queries.DashboardDTO, error) {
	return s.dashboardQueryHandler.HandleGetDashboard(ctx, query)
}

// === HEALTH CHECKS ===

// HealthCheck verifica saúde do serviço
func (s *ProcessService) HealthCheck(ctx context.Context) error {
	// Implementação básica de health check
	// Em produção, verificaria conexões com banco, message broker, etc.
	return nil
}

// === UTILITY METHODS ===

// GetServiceInfo retorna informações do serviço
func (s *ProcessService) GetServiceInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":    "process-service",
		"version": "1.0.0",
		"pattern": "CQRS",
		"features": []string{
			"process_management",
			"movement_tracking",
			"monitoring",
			"search",
			"dashboard",
			"event_sourcing",
		},
	}
}

// === COMMAND FACTORY METHODS ===

// CreateProcessCommandBuilder helper para criar comando de processo
func (s *ProcessService) CreateProcessCommandBuilder() *ProcessCommandBuilder {
	return &ProcessCommandBuilder{}
}

// ProcessCommandBuilder builder para comando de criação de processo
type ProcessCommandBuilder struct {
	cmd commands.CreateProcessCommand
}

// WithTenant define tenant
func (b *ProcessCommandBuilder) WithTenant(tenantID string) *ProcessCommandBuilder {
	b.cmd.TenantID = tenantID
	return b
}

// WithClient define cliente
func (b *ProcessCommandBuilder) WithClient(clientID string) *ProcessCommandBuilder {
	b.cmd.ClientID = clientID
	return b
}

// WithNumber define número do processo
func (b *ProcessCommandBuilder) WithNumber(number string) *ProcessCommandBuilder {
	b.cmd.Number = number
	return b
}

// WithTitle define título
func (b *ProcessCommandBuilder) WithTitle(title string) *ProcessCommandBuilder {
	b.cmd.Title = title
	return b
}

// WithDescription define descrição
func (b *ProcessCommandBuilder) WithDescription(description string) *ProcessCommandBuilder {
	b.cmd.Description = description
	return b
}

// WithStatus define status
func (b *ProcessCommandBuilder) WithStatus(status domain.ProcessStatus) *ProcessCommandBuilder {
	b.cmd.Status = status
	return b
}

// WithStage define fase
func (b *ProcessCommandBuilder) WithStage(stage domain.ProcessStage) *ProcessCommandBuilder {
	b.cmd.Stage = stage
	return b
}

// WithSubject define assunto
func (b *ProcessCommandBuilder) WithSubject(code, description, parentCode string) *ProcessCommandBuilder {
	b.cmd.Subject = commands.CreateProcessSubjectCommand{
		Code:        code,
		Description: description,
		ParentCode:  parentCode,
	}
	return b
}

// WithCourt define tribunal
func (b *ProcessCommandBuilder) WithCourt(courtID string) *ProcessCommandBuilder {
	b.cmd.CourtID = courtID
	return b
}

// WithJudge define juiz
func (b *ProcessCommandBuilder) WithJudge(judgeID string) *ProcessCommandBuilder {
	b.cmd.JudgeID = &judgeID
	return b
}

// WithValue define valor da causa
func (b *ProcessCommandBuilder) WithValue(amount float64, currency string) *ProcessCommandBuilder {
	b.cmd.Value = &commands.CreateProcessValueCommand{
		Amount:   amount,
		Currency: currency,
	}
	return b
}

// WithMonitoring define configuração de monitoramento
func (b *ProcessCommandBuilder) WithMonitoring(enabled bool, channels []string, syncInterval int) *ProcessCommandBuilder {
	b.cmd.Monitoring = commands.CreateMonitoringCommand{
		Enabled:              enabled,
		NotificationChannels: channels,
		AutoSync:             enabled,
		SyncIntervalHours:    syncInterval,
	}
	return b
}

// WithCreatedBy define criador
func (b *ProcessCommandBuilder) WithCreatedBy(userID string) *ProcessCommandBuilder {
	b.cmd.CreatedBy = userID
	return b
}

// AddParty adiciona parte ao comando
func (b *ProcessCommandBuilder) AddParty(party commands.CreatePartyCommand) *ProcessCommandBuilder {
	b.cmd.Parties = append(b.cmd.Parties, party)
	return b
}

// AddTag adiciona tag
func (b *ProcessCommandBuilder) AddTag(tag string) *ProcessCommandBuilder {
	b.cmd.Tags = append(b.cmd.Tags, tag)
	return b
}

// SetCustomField define campo customizado
func (b *ProcessCommandBuilder) SetCustomField(key, value string) *ProcessCommandBuilder {
	if b.cmd.CustomFields == nil {
		b.cmd.CustomFields = make(map[string]string)
	}
	b.cmd.CustomFields[key] = value
	return b
}

// Build constrói o comando
func (b *ProcessCommandBuilder) Build() *commands.CreateProcessCommand {
	return &b.cmd
}

// === QUERY FACTORY METHODS ===

// CreateProcessListQueryBuilder helper para criar query de listagem
func (s *ProcessService) CreateProcessListQueryBuilder() *ProcessListQueryBuilder {
	return &ProcessListQueryBuilder{
		query: queries.ProcessListQuery{
			Page:     1,
			PageSize: 20,
		},
	}
}

// ProcessListQueryBuilder builder para query de listagem de processos
type ProcessListQueryBuilder struct {
	query queries.ProcessListQuery
}

// WithTenant define tenant
func (b *ProcessListQueryBuilder) WithTenant(tenantID string) *ProcessListQueryBuilder {
	b.query.TenantID = tenantID
	return b
}

// WithClient define cliente
func (b *ProcessListQueryBuilder) WithClient(clientID string) *ProcessListQueryBuilder {
	b.query.ClientID = clientID
	return b
}

// WithStatus filtra por status
func (b *ProcessListQueryBuilder) WithStatus(status ...domain.ProcessStatus) *ProcessListQueryBuilder {
	b.query.Status = status
	return b
}

// WithStage filtra por fase
func (b *ProcessListQueryBuilder) WithStage(stage ...domain.ProcessStage) *ProcessListQueryBuilder {
	b.query.Stage = stage
	return b
}

// WithCourt filtra por tribunal
func (b *ProcessListQueryBuilder) WithCourt(courtID string) *ProcessListQueryBuilder {
	b.query.CourtID = courtID
	return b
}

// WithSearch define busca textual
func (b *ProcessListQueryBuilder) WithSearch(search string) *ProcessListQueryBuilder {
	b.query.Search = search
	return b
}

// WithPagination define paginação
func (b *ProcessListQueryBuilder) WithPagination(page, pageSize int) *ProcessListQueryBuilder {
	b.query.Page = page
	b.query.PageSize = pageSize
	return b
}

// WithSorting define ordenação
func (b *ProcessListQueryBuilder) WithSorting(sortBy, sortOrder string) *ProcessListQueryBuilder {
	b.query.SortBy = sortBy
	b.query.SortOrder = sortOrder
	return b
}

// Build constrói a query
func (b *ProcessListQueryBuilder) Build() *queries.ProcessListQuery {
	return &b.query
}