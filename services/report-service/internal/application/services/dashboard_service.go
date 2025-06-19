package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/domain"
)

// DashboardService serviço de aplicação para dashboards
type DashboardService struct {
	dashboardRepo domain.DashboardRepository
	kpiRepo       domain.KPIRepository
	dataCollector domain.DataCollector
	eventBus      domain.EventBus
	logger        *zap.Logger
}

// NewDashboardService cria nova instância do serviço
func NewDashboardService(
	dashboardRepo domain.DashboardRepository,
	kpiRepo domain.KPIRepository,
	dataCollector domain.DataCollector,
	eventBus domain.EventBus,
	logger *zap.Logger,
) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
		kpiRepo:       kpiRepo,
		dataCollector: dataCollector,
		eventBus:      eventBus,
		logger:        logger,
	}
}

// CreateDashboardRequest requisição para criar dashboard
type CreateDashboardRequest struct {
	Title           string                 `json:"title" validate:"required"`
	Description     string                 `json:"description,omitempty"`
	IsPublic        bool                   `json:"is_public"`
	IsDefault       bool                   `json:"is_default"`
	Layout          map[string]interface{} `json:"layout,omitempty"`
	RefreshInterval *int                   `json:"refresh_interval,omitempty"`
}

// CreateDashboard cria um novo dashboard
func (s *DashboardService) CreateDashboard(ctx context.Context, req *CreateDashboardRequest) (*domain.Dashboard, error) {
	tenantID := domain.MustGetTenantID(ctx)
	userID := domain.MustGetUserID(ctx)

	s.logger.Info("Creating dashboard",
		zap.String("tenant_id", tenantID.String()),
		zap.String("user_id", userID.String()),
		zap.String("title", req.Title))

	// Se marcando como default, remover default de outros
	if req.IsDefault {
		if err := s.unsetDefaultDashboards(ctx, tenantID); err != nil {
			return nil, err
		}
	}

	dashboard := &domain.Dashboard{
		ID:              uuid.New(),
		TenantID:        tenantID,
		Title:           req.Title,
		Description:     req.Description,
		IsPublic:        req.IsPublic,
		IsDefault:       req.IsDefault,
		Layout:          req.Layout,
		RefreshInterval: req.RefreshInterval,
		Widgets:         []domain.DashboardWidget{},
		CreatedBy:       userID,
		UpdatedBy:       userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.dashboardRepo.Create(ctx, dashboard); err != nil {
		s.logger.Error("Failed to create dashboard", zap.Error(err))
		return nil, fmt.Errorf("failed to create dashboard: %w", err)
	}

	// Publicar evento
	event := domain.DashboardCreatedEvent{
		DashboardEvent: domain.DashboardEvent{
			EventID:     uuid.New(),
			EventType:   "dashboard.created",
			DashboardID: dashboard.ID,
			TenantID:    tenantID,
			UserID:      userID,
			OccurredAt:  time.Now(),
		},
		Title:       dashboard.Title,
		IsPublic:    dashboard.IsPublic,
		WidgetCount: 0,
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish event", zap.Error(err))
	}

	return dashboard, nil
}

// AddWidgetRequest requisição para adicionar widget
type AddWidgetRequest struct {
	Type            string                 `json:"type" validate:"required"`
	Title           string                 `json:"title" validate:"required"`
	DataSource      string                 `json:"data_source" validate:"required"`
	ChartType       string                 `json:"chart_type,omitempty"`
	Parameters      map[string]interface{} `json:"parameters,omitempty"`
	Filters         map[string]interface{} `json:"filters,omitempty"`
	Position        domain.WidgetPosition  `json:"position"`
	Size            domain.WidgetSize      `json:"size"`
	RefreshInterval *int                   `json:"refresh_interval,omitempty"`
}

// AddWidget adiciona um widget ao dashboard
func (s *DashboardService) AddWidget(ctx context.Context, dashboardID uuid.UUID, req *AddWidgetRequest) (*domain.DashboardWidget, error) {
	tenantID := domain.MustGetTenantID(ctx)

	// Verificar se dashboard existe e pertence ao tenant
	dashboard, err := s.dashboardRepo.GetByID(ctx, dashboardID)
	if err != nil {
		return nil, err
	}
	
	if dashboard.TenantID != tenantID {
		return nil, domain.ErrUnauthorized
	}

	// Obter widgets existentes para calcular ordem
	widgets, err := s.dashboardRepo.GetWidgets(ctx, dashboardID)
	if err != nil {
		return nil, err
	}

	// Verificar limite de widgets por plano
	plan, _ := domain.GetPlan(ctx)
	if err := s.checkWidgetLimit(plan, len(widgets)); err != nil {
		return nil, err
	}

	widget := &domain.DashboardWidget{
		ID:              uuid.New(),
		DashboardID:     dashboardID,
		Type:            req.Type,
		Title:           req.Title,
		DataSource:      req.DataSource,
		ChartType:       req.ChartType,
		Parameters:      req.Parameters,
		Filters:         req.Filters,
		Position:        req.Position,
		Size:            req.Size,
		RefreshInterval: req.RefreshInterval,
		IsVisible:       true,
		Order:           len(widgets) + 1,
	}

	if err := s.dashboardRepo.AddWidget(ctx, widget); err != nil {
		s.logger.Error("Failed to add widget", zap.Error(err))
		return nil, fmt.Errorf("failed to add widget: %w", err)
	}

	// Publicar evento
	event := domain.WidgetAddedEvent{
		DashboardEvent: domain.DashboardEvent{
			EventID:     uuid.New(),
			EventType:   "widget.added",
			DashboardID: dashboardID,
			TenantID:    tenantID,
			UserID:      domain.MustGetUserID(ctx),
			OccurredAt:  time.Now(),
		},
		WidgetID:   widget.ID,
		WidgetType: widget.Type,
		DataSource: widget.DataSource,
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish event", zap.Error(err))
	}

	return widget, nil
}

// GetDashboard busca dashboard por ID
func (s *DashboardService) GetDashboard(ctx context.Context, id uuid.UUID) (*domain.Dashboard, error) {
	dashboard, err := s.dashboardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verificar permissão
	tenantID := domain.MustGetTenantID(ctx)
	if dashboard.TenantID != tenantID {
		return nil, domain.ErrUnauthorized
	}

	// Carregar widgets
	widgets, err := s.dashboardRepo.GetWidgets(ctx, id)
	if err != nil {
		s.logger.Warn("Failed to load widgets", zap.Error(err))
		widgets = []domain.DashboardWidget{}
	}
	
	dashboard.Widgets = widgets
	return dashboard, nil
}

// GetDashboardData obtém dados para renderizar o dashboard
func (s *DashboardService) GetDashboardData(ctx context.Context, dashboardID uuid.UUID) (map[string]interface{}, error) {
	dashboard, err := s.GetDashboard(ctx, dashboardID)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["dashboard"] = dashboard
	
	// Coletar dados para cada widget
	widgetData := make(map[string]interface{})
	for _, widget := range dashboard.Widgets {
		if widget.IsVisible {
			wData, err := s.dataCollector.CollectDashboardData(ctx, &widget)
			if err != nil {
				s.logger.Warn("Failed to collect widget data",
					zap.String("widget_id", widget.ID.String()),
					zap.Error(err))
				continue
			}
			widgetData[widget.ID.String()] = wData
		}
	}
	
	data["widgets"] = widgetData

	// Adicionar KPIs se é dashboard executivo
	if dashboard.IsDefault {
		kpis, err := s.kpiRepo.GetByTenantID(ctx, dashboard.TenantID)
		if err == nil {
			data["kpis"] = kpis
		}
	}

	return data, nil
}

// ListDashboards lista dashboards do tenant
func (s *DashboardService) ListDashboards(ctx context.Context) ([]*domain.Dashboard, error) {
	tenantID := domain.MustGetTenantID(ctx)
	return s.dashboardRepo.GetByTenantID(ctx, tenantID)
}

// UpdateDashboard atualiza um dashboard
func (s *DashboardService) UpdateDashboard(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*domain.Dashboard, error) {
	dashboard, err := s.GetDashboard(ctx, id)
	if err != nil {
		return nil, err
	}

	// Aplicar atualizações
	if title, ok := updates["title"].(string); ok {
		dashboard.Title = title
	}
	if description, ok := updates["description"].(string); ok {
		dashboard.Description = description
	}
	if isPublic, ok := updates["is_public"].(bool); ok {
		dashboard.IsPublic = isPublic
	}
	if layout, ok := updates["layout"].(map[string]interface{}); ok {
		dashboard.Layout = layout
	}

	dashboard.UpdatedBy = domain.MustGetUserID(ctx)
	dashboard.UpdatedAt = time.Now()

	if err := s.dashboardRepo.Update(ctx, dashboard); err != nil {
		return nil, err
	}

	// Publicar evento
	event := domain.DashboardUpdatedEvent{
		DashboardEvent: domain.DashboardEvent{
			EventID:     uuid.New(),
			EventType:   "dashboard.updated",
			DashboardID: dashboard.ID,
			TenantID:    dashboard.TenantID,
			UserID:      dashboard.UpdatedBy,
			OccurredAt:  time.Now(),
		},
		Changes: updates,
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish event", zap.Error(err))
	}

	return dashboard, nil
}

// DeleteDashboard exclui um dashboard
func (s *DashboardService) DeleteDashboard(ctx context.Context, id uuid.UUID) error {
	dashboard, err := s.GetDashboard(ctx, id)
	if err != nil {
		return err
	}

	// Não permitir excluir dashboard default
	if dashboard.IsDefault {
		return fmt.Errorf("cannot delete default dashboard")
	}

	return s.dashboardRepo.Delete(ctx, id)
}

// CalculateKPIs calcula KPIs do tenant
func (s *DashboardService) CalculateKPIs(ctx context.Context) error {
	tenantID := domain.MustGetTenantID(ctx)
	
	s.logger.Info("Calculating KPIs", zap.String("tenant_id", tenantID.String()))

	kpis, err := s.dataCollector.CalculateKPIs(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to calculate KPIs: %w", err)
	}

	// Atualizar ou criar KPIs
	for _, kpi := range kpis {
		existing, err := s.kpiRepo.GetByID(ctx, kpi.ID)
		if err == nil {
			// Atualizar KPI existente
			kpi.PreviousValue = &existing.CurrentValue
			if err := s.kpiRepo.Update(ctx, kpi); err != nil {
				s.logger.Error("Failed to update KPI", 
					zap.String("kpi_id", kpi.ID.String()),
					zap.Error(err))
			}
		} else {
			// Criar novo KPI
			if err := s.kpiRepo.Create(ctx, kpi); err != nil {
				s.logger.Error("Failed to create KPI",
					zap.String("kpi_name", kpi.Name),
					zap.Error(err))
			}
		}

		// Publicar evento
		event := domain.KPICalculatedEvent{
			KPIEvent: domain.KPIEvent{
				EventID:    uuid.New(),
				EventType:  "kpi.calculated",
				KPIID:      kpi.ID,
				TenantID:   tenantID,
				OccurredAt: time.Now(),
			},
			Name:          kpi.Name,
			CurrentValue:  kpi.CurrentValue,
			PreviousValue: kpi.PreviousValue,
			Trend:         kpi.Trend,
		}

		if err := s.eventBus.Publish(ctx, event); err != nil {
			s.logger.Error("Failed to publish KPI event", zap.Error(err))
		}

		// Verificar alertas
		s.checkKPIAlerts(ctx, kpi)
	}

	return nil
}

// checkKPIAlerts verifica se KPI deve disparar alertas
func (s *DashboardService) checkKPIAlerts(ctx context.Context, kpi *domain.KPI) {
	// Verificar se KPI tem target e está fora
	if kpi.Target != nil {
		deviation := (kpi.CurrentValue - *kpi.Target) / *kpi.Target * 100
		
		// Alerta se desvio maior que 20%
		if deviation > 20 || deviation < -20 {
			event := domain.KPIAlertEvent{
				KPIEvent: domain.KPIEvent{
					EventID:    uuid.New(),
					EventType:  "kpi.alert",
					KPIID:      kpi.ID,
					TenantID:   kpi.TenantID,
					OccurredAt: time.Now(),
				},
				AlertType:    "target_deviation",
				CurrentValue: kpi.CurrentValue,
				Threshold:    *kpi.Target,
				Message:      fmt.Sprintf("KPI %s is %.1f%% away from target", kpi.DisplayName, deviation),
			}

			if err := s.eventBus.Publish(ctx, event); err != nil {
				s.logger.Error("Failed to publish KPI alert", zap.Error(err))
			}
		}
	}
}

// unsetDefaultDashboards remove flag default de todos dashboards
func (s *DashboardService) unsetDefaultDashboards(ctx context.Context, tenantID uuid.UUID) error {
	dashboards, err := s.dashboardRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}

	for _, dashboard := range dashboards {
		if dashboard.IsDefault {
			dashboard.IsDefault = false
			dashboard.UpdatedAt = time.Now()
			if err := s.dashboardRepo.Update(ctx, dashboard); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkWidgetLimit verifica limite de widgets por plano
func (s *DashboardService) checkWidgetLimit(plan string, currentCount int) error {
	var limit int
	switch plan {
	case "starter":
		limit = 5
	case "professional":
		limit = 10
	case "business":
		limit = 20
	case "enterprise":
		return nil // unlimited
	default:
		limit = 5
	}

	if currentCount >= limit {
		return domain.ErrMaxWidgetsReached.WithDetail("limit", limit)
	}

	return nil
}