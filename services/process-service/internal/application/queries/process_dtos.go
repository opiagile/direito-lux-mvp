package queries

import (
	"time"
	"github.com/direito-lux/process-service/internal/domain"
)

// ProcessDTO dados do processo para consultas
type ProcessDTO struct {
	ID                  string                  `json:"id"`
	TenantID            string                  `json:"tenant_id"`
	ClientID            string                  `json:"client_id"`
	Number              string                  `json:"number"`
	OriginalNumber      string                  `json:"original_number,omitempty"`
	Title               string                  `json:"title"`
	Description         string                  `json:"description,omitempty"`
	Status              domain.ProcessStatus    `json:"status"`
	Stage               domain.ProcessStage     `json:"stage"`
	Subject             ProcessSubjectDTO       `json:"subject"`
	Value               *ProcessValueDTO        `json:"value,omitempty"`
	CourtID             string                  `json:"court_id"`
	JudgeID             *string                 `json:"judge_id,omitempty"`
	Parties             []PartyDTO              `json:"parties,omitempty"`
	Monitoring          ProcessMonitoringDTO    `json:"monitoring"`
	Tags                []string                `json:"tags"`
	CustomFields        map[string]string       `json:"custom_fields,omitempty"`
	Stats               *ProcessStatsDTO        `json:"stats,omitempty"`
	LastMovementAt      *time.Time              `json:"last_movement_at,omitempty"`
	LastSyncAt          *time.Time              `json:"last_sync_at,omitempty"`
	CreatedAt           time.Time               `json:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at"`
	ArchivedAt          *time.Time              `json:"archived_at,omitempty"`
}

// ProcessSubjectDTO assunto do processo
type ProcessSubjectDTO struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	ParentCode  string `json:"parent_code,omitempty"`
}

// ProcessValueDTO valor da causa
type ProcessValueDTO struct {
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	FormattedValue  string  `json:"formatted_value"`
}

// ProcessMonitoringDTO configuração de monitoramento
type ProcessMonitoringDTO struct {
	Enabled              bool       `json:"enabled"`
	NotificationChannels []string   `json:"notification_channels"`
	Keywords             []string   `json:"keywords"`
	AutoSync             bool       `json:"auto_sync"`
	SyncIntervalHours    int        `json:"sync_interval_hours"`
	LastNotificationAt   *time.Time `json:"last_notification_at,omitempty"`
	NeedsSync            bool       `json:"needs_sync"`
}

// ProcessStatsDTO estatísticas do processo
type ProcessStatsDTO struct {
	TotalMovements     int       `json:"total_movements"`
	ImportantMovements int       `json:"important_movements"`
	LastMovementDate   *time.Time `json:"last_movement_date,omitempty"`
	AverageResponseTime string   `json:"average_response_time"`
	DaysActive         int       `json:"days_active"`
}

// ProcessListDTO lista paginada de processos
type ProcessListDTO struct {
	Processes   []ProcessDTO       `json:"processes"`
	Pagination  PaginationDTO      `json:"pagination"`
	Summary     ProcessSummaryDTO  `json:"summary"`
}

// ProcessSummaryDTO resumo da lista de processos
type ProcessSummaryDTO struct {
	TotalProcesses     int                              `json:"total_processes"`
	ByStatus          map[domain.ProcessStatus]int     `json:"by_status"`
	ByStage           map[domain.ProcessStage]int      `json:"by_stage"`
	ActiveMonitoring  int                              `json:"active_monitoring"`
	NeedingSync       int                              `json:"needing_sync"`
}

// MovementDTO dados da movimentação para consultas
type MovementDTO struct {
	ID                  string                   `json:"id"`
	ProcessID           string                   `json:"process_id"`
	ProcessNumber       string                   `json:"process_number,omitempty"`
	TenantID            string                   `json:"tenant_id"`
	Sequence            int                      `json:"sequence"`
	ExternalID          string                   `json:"external_id,omitempty"`
	Date                time.Time                `json:"date"`
	Type                domain.MovementType      `json:"type"`
	Code                string                   `json:"code"`
	Title               string                   `json:"title"`
	Description         string                   `json:"description"`
	Content             string                   `json:"content,omitempty"`
	Summary             string                   `json:"summary"`
	Judge               string                   `json:"judge,omitempty"`
	Responsible         string                   `json:"responsible,omitempty"`
	Attachments         []AttachmentDTO          `json:"attachments,omitempty"`
	RelatedParties      []string                 `json:"related_parties,omitempty"`
	IsImportant         bool                     `json:"is_important"`
	IsPublic            bool                     `json:"is_public"`
	NotificationSent    bool                     `json:"notification_sent"`
	Tags                []string                 `json:"tags"`
	Metadata            MovementMetadataDTO      `json:"metadata"`
	FormattedDate       string                   `json:"formatted_date"`
	DisplayTitle        string                   `json:"display_title"`
	ImportanceLevel     string                   `json:"importance_level"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           time.Time                `json:"updated_at"`
	SyncedAt            time.Time                `json:"synced_at"`
}

// AttachmentDTO anexo da movimentação
type AttachmentDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Type        string    `json:"type"`
	URL         string    `json:"url"`
	ExternalID  string    `json:"external_id,omitempty"`
	IsDownloaded bool     `json:"is_downloaded"`
	CreatedAt   time.Time `json:"created_at"`
}

// MovementMetadataDTO metadados da movimentação
type MovementMetadataDTO struct {
	OriginalSource string                    `json:"original_source"`
	DataJudID      string                    `json:"datajud_id,omitempty"`
	ImportBatch    string                    `json:"import_batch,omitempty"`
	Keywords       []string                  `json:"keywords"`
	Analysis       MovementAnalysisDTO       `json:"analysis"`
	CustomFields   map[string]string         `json:"custom_fields,omitempty"`
}

// MovementAnalysisDTO análise da movimentação
type MovementAnalysisDTO struct {
	Sentiment      string     `json:"sentiment"`
	Importance     int        `json:"importance"`
	Category       string     `json:"category"`
	HasDeadline    bool       `json:"has_deadline"`
	DeadlineDate   *time.Time `json:"deadline_date,omitempty"`
	RequiresAction bool       `json:"requires_action"`
	ActionType     string     `json:"action_type,omitempty"`
	Confidence     float64    `json:"confidence"`
	ProcessedBy    string     `json:"processed_by"`
	ProcessedAt    time.Time  `json:"processed_at"`
}

// MovementListDTO lista paginada de movimentações
type MovementListDTO struct {
	Movements   []MovementDTO       `json:"movements"`
	Pagination  PaginationDTO       `json:"pagination"`
	Summary     MovementSummaryDTO  `json:"summary"`
}

// MovementSummaryDTO resumo da lista de movimentações
type MovementSummaryDTO struct {
	TotalMovements      int                              `json:"total_movements"`
	ImportantMovements  int                              `json:"important_movements"`
	ByType             map[domain.MovementType]int      `json:"by_type"`
	ByMonth            map[string]int                   `json:"by_month"`
	PendingNotifications int                             `json:"pending_notifications"`
}

// PartyDTO dados da parte para consultas
type PartyDTO struct {
	ID                string            `json:"id"`
	ProcessID         string            `json:"process_id"`
	Type              domain.PartyType  `json:"type"`
	Name              string            `json:"name"`
	Document          string            `json:"document,omitempty"`
	DocumentType      string            `json:"document_type,omitempty"`
	FormattedDocument string            `json:"formatted_document,omitempty"`
	Role              domain.PartyRole  `json:"role"`
	IsActive          bool              `json:"is_active"`
	Lawyer            *LawyerDTO        `json:"lawyer,omitempty"`
	Contact           PartyContactDTO   `json:"contact"`
	Address           PartyAddressDTO   `json:"address"`
	DisplayName       string            `json:"display_name"`
	LawyerInfo        string            `json:"lawyer_info"`
	IsMainParty       bool              `json:"is_main_party"`
	IsLegalEntity     bool              `json:"is_legal_entity"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// LawyerDTO dados do advogado
type LawyerDTO struct {
	Name     string `json:"name"`
	OAB      string `json:"oab"`
	OABState string `json:"oab_state"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

// PartyContactDTO contato da parte
type PartyContactDTO struct {
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	CellPhone string `json:"cell_phone,omitempty"`
	Website   string `json:"website,omitempty"`
}

// PartyAddressDTO endereço da parte
type PartyAddressDTO struct {
	Street     string `json:"street,omitempty"`
	Number     string `json:"number,omitempty"`
	Complement string `json:"complement,omitempty"`
	District   string `json:"district,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	ZipCode    string `json:"zip_code,omitempty"`
	Country    string `json:"country,omitempty"`
}

// PartyListDTO lista paginada de partes
type PartyListDTO struct {
	Parties    []PartyDTO      `json:"parties"`
	Pagination PaginationDTO   `json:"pagination"`
	Summary    PartySummaryDTO `json:"summary"`
}

// PartySummaryDTO resumo da lista de partes
type PartySummaryDTO struct {
	TotalParties   int                            `json:"total_parties"`
	ByType        map[domain.PartyType]int       `json:"by_type"`
	ByRole        map[domain.PartyRole]int       `json:"by_role"`
	WithLawyer    int                            `json:"with_lawyer"`
	ActiveParties int                            `json:"active_parties"`
}

// PaginationDTO informações de paginação
type PaginationDTO struct {
	Page         int   `json:"page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
	TotalItems   int64 `json:"total_items"`
	HasPrevious  bool  `json:"has_previous"`
	HasNext      bool  `json:"has_next"`
}

// DashboardDTO dados do dashboard
type DashboardDTO struct {
	TenantID      string                    `json:"tenant_id"`
	ClientID      string                    `json:"client_id,omitempty"`
	Period        string                    `json:"period"`
	Overview      DashboardOverviewDTO      `json:"overview"`
	Processes     DashboardProcessesDTO     `json:"processes"`
	Movements     DashboardMovementsDTO     `json:"movements"`
	Monitoring    DashboardMonitoringDTO    `json:"monitoring"`
	Alerts        DashboardAlertsDTO        `json:"alerts"`
	Calendar      []DashboardEventDTO       `json:"calendar"`
	RecentActivity []DashboardActivityDTO   `json:"recent_activity"`
	GeneratedAt   time.Time                 `json:"generated_at"`
}

// DashboardOverviewDTO visão geral do dashboard
type DashboardOverviewDTO struct {
	TotalProcesses      int     `json:"total_processes"`
	ActiveProcesses     int     `json:"active_processes"`
	ArchivedProcesses   int     `json:"archived_processes"`
	ProcessesGrowth     float64 `json:"processes_growth"` // percentual vs período anterior
	TotalMovements      int     `json:"total_movements"`
	ImportantMovements  int     `json:"important_movements"`
	MovementsGrowth     float64 `json:"movements_growth"`
	MonitoringEnabled   int     `json:"monitoring_enabled"`
	PendingAlerts       int     `json:"pending_alerts"`
}

// DashboardProcessesDTO dados de processos do dashboard
type DashboardProcessesDTO struct {
	ByStatus           map[domain.ProcessStatus]int   `json:"by_status"`
	ByStage            map[domain.ProcessStage]int    `json:"by_stage"`
	ByCourt            map[string]int                 `json:"by_court"`
	MostActive         []ProcessActivityDTO          `json:"most_active"`
	RecentlyCreated    []ProcessDTO                  `json:"recently_created"`
	NeedingAttention   []ProcessDTO                  `json:"needing_attention"`
}

// DashboardMovementsDTO dados de movimentações do dashboard
type DashboardMovementsDTO struct {
	ByType           map[domain.MovementType]int   `json:"by_type"`
	ByImportance     map[int]int                   `json:"by_importance"`
	Timeline         []MovementTimelineDTO         `json:"timeline"`
	TopKeywords      []KeywordCountDTO             `json:"top_keywords"`
	RecentImportant  []MovementDTO                 `json:"recent_important"`
}

// DashboardMonitoringDTO dados de monitoramento do dashboard
type DashboardMonitoringDTO struct {
	EnabledProcesses    int                     `json:"enabled_processes"`
	NeedingSync         int                     `json:"needing_sync"`
	SyncErrors          int                     `json:"sync_errors"`
	NotificationsSent   int                     `json:"notifications_sent"`
	ChannelStats        map[string]int          `json:"channel_stats"`
	SyncHistory         []SyncHistoryDTO        `json:"sync_history"`
}

// DashboardAlertsDTO dados de alertas do dashboard
type DashboardAlertsDTO struct {
	TotalAlerts        int                     `json:"total_alerts"`
	UnreadAlerts       int                     `json:"unread_alerts"`
	BySeverity         map[string]int          `json:"by_severity"`
	RecentAlerts       []AlertDTO              `json:"recent_alerts"`
}

// ProcessActivityDTO atividade do processo
type ProcessActivityDTO struct {
	ProcessID      string    `json:"process_id"`
	ProcessNumber  string    `json:"process_number"`
	ProcessTitle   string    `json:"process_title"`
	MovementCount  int       `json:"movement_count"`
	LastMovement   time.Time `json:"last_movement"`
}

// MovementTimelineDTO dados da timeline de movimentações
type MovementTimelineDTO struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// KeywordCountDTO contagem de palavras-chave
type KeywordCountDTO struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
}

// SyncHistoryDTO histórico de sincronização
type SyncHistoryDTO struct {
	Date            time.Time `json:"date"`
	ProcessesCount  int       `json:"processes_count"`
	SuccessCount    int       `json:"success_count"`
	ErrorCount      int       `json:"error_count"`
	NewMovements    int       `json:"new_movements"`
	UpdatedMovements int      `json:"updated_movements"`
}

// AlertDTO alerta/notificação
type AlertDTO struct {
	ID          string    `json:"id"`
	ProcessID   string    `json:"process_id"`
	ProcessNumber string  `json:"process_number"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

// DashboardEventDTO evento do calendário
type DashboardEventDTO struct {
	ID          string    `json:"id"`
	ProcessID   string    `json:"process_id"`
	ProcessNumber string  `json:"process_number"`
	Type        string    `json:"type"` // hearing, deadline, notification
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Importance  int       `json:"importance"`
}

// DashboardActivityDTO atividade recente
type DashboardActivityDTO struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // process_created, movement_added, alert_generated
	ProcessID   string    `json:"process_id"`
	ProcessNumber string  `json:"process_number"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
}

// SearchResultDTO resultado de busca
type SearchResultDTO struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // process, movement, party
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Highlights  map[string][]string    `json:"highlights,omitempty"`
	Score       float64                `json:"score"`
	Data        map[string]interface{} `json:"data"`
}

// SearchResultsDTO lista de resultados de busca
type SearchResultsDTO struct {
	Query       string            `json:"query"`
	Results     []SearchResultDTO `json:"results"`
	Pagination  PaginationDTO     `json:"pagination"`
	Aggregations map[string]interface{} `json:"aggregations,omitempty"`
	TookMs      int64             `json:"took_ms"`
}

// TimelineEventDTO evento da timeline do processo
type TimelineEventDTO struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // movement, party_added, status_changed
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Date        time.Time              `json:"date"`
	Importance  int                    `json:"importance"`
	Data        map[string]interface{} `json:"data"`
	UserID      string                 `json:"user_id,omitempty"`
	UserName    string                 `json:"user_name,omitempty"`
}

// ProcessTimelineDTO timeline do processo
type ProcessTimelineDTO struct {
	ProcessID  string             `json:"process_id"`
	Events     []TimelineEventDTO `json:"events"`
	Pagination PaginationDTO      `json:"pagination"`
}

// CalendarEventDTO evento do calendário
type CalendarEventDTO struct {
	ID            string    `json:"id"`
	ProcessID     string    `json:"process_id"`
	ProcessNumber string    `json:"process_number"`
	ProcessTitle  string    `json:"process_title"`
	Type          string    `json:"type"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	StartDate     time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	AllDay        bool      `json:"all_day"`
	Importance    int       `json:"importance"`
	Location      string    `json:"location,omitempty"`
	Attendees     []string  `json:"attendees,omitempty"`
}

// CalendarDTO calendário com eventos
type CalendarDTO struct {
	DateFrom time.Time          `json:"date_from"`
	DateTo   time.Time          `json:"date_to"`
	Events   []CalendarEventDTO `json:"events"`
	Summary  CalendarSummaryDTO `json:"summary"`
}

// CalendarSummaryDTO resumo do calendário
type CalendarSummaryDTO struct {
	TotalEvents    int            `json:"total_events"`
	ByType         map[string]int `json:"by_type"`
	ByImportance   map[int]int    `json:"by_importance"`
	UpcomingEvents int            `json:"upcoming_events"`
}