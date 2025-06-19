package excel

import (
	"context"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/domain"
)

// ExcelGenerator implementação do gerador de Excel
type ExcelGenerator struct {
	logger *zap.Logger
}

// NewExcelGenerator cria nova instância do gerador
func NewExcelGenerator(logger *zap.Logger) *ExcelGenerator {
	return &ExcelGenerator{
		logger: logger,
	}
}

// GenerateExcel gera relatório em Excel
func (g *ExcelGenerator) GenerateExcel(ctx context.Context, report *domain.Report, data interface{}) ([]byte, error) {
	g.logger.Debug("Generating Excel report",
		zap.String("report_id", report.ID.String()),
		zap.String("type", string(report.Type)))

	// Criar novo arquivo Excel
	f := excelize.NewFile()
	
	// Criar planilha principal
	sheet := "Relatório"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %w", err)
	}
	
	// Definir como planilha ativa
	f.SetActiveSheet(index)
	
	// Adicionar header
	g.addHeader(f, sheet, report)
	
	// Adicionar conteúdo baseado no tipo
	switch report.Type {
	case domain.ReportTypeExecutiveSummary:
		g.addExecutiveSummaryContent(f, sheet, data)
		
	case domain.ReportTypeProcessAnalysis:
		g.addProcessAnalysisContent(f, sheet, data)
		
	case domain.ReportTypeProductivity:
		g.addProductivityContent(f, sheet, data)
		
	case domain.ReportTypeFinancial:
		g.addFinancialContent(f, sheet, data)
		
	case domain.ReportTypeJurisprudence:
		g.addJurisprudenceContent(f, sheet, data)
		
	default:
		g.addGenericContent(f, sheet, data)
	}
	
	// Ajustar largura das colunas
	g.autoFitColumns(f, sheet)
	
	// Gerar bytes do Excel
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Excel: %w", err)
	}
	
	return buffer.Bytes(), nil
}

// addHeader adiciona cabeçalho ao Excel
func (g *ExcelGenerator) addHeader(f *excelize.File, sheet string, report *domain.Report) {
	// Estilo do título principal
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 20,
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	
	// Estilo do subtítulo
	subtitleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 14,
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	
	// Estilo da data
	dateStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  10,
			Color: "666666",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	
	// Título principal
	f.MergeCell(sheet, "A1", "F1")
	f.SetCellValue(sheet, "A1", "Direito Lux - Sistema Jurídico")
	f.SetCellStyle(sheet, "A1", "F1", titleStyle)
	
	// Título do relatório
	f.MergeCell(sheet, "A2", "F2")
	f.SetCellValue(sheet, "A2", report.Title)
	f.SetCellStyle(sheet, "A2", "F2", subtitleStyle)
	
	// Data e hora
	f.MergeCell(sheet, "A3", "F3")
	f.SetCellValue(sheet, "A3", fmt.Sprintf("Gerado em: %s", time.Now().Format("02/01/2006 15:04")))
	f.SetCellStyle(sheet, "A3", "F3", dateStyle)
	
	// Linha em branco
	f.SetRowHeight(sheet, 4, 10)
}

// addExecutiveSummaryContent adiciona conteúdo do resumo executivo
func (g *ExcelGenerator) addExecutiveSummaryContent(f *excelize.File, sheet string, data interface{}) {
	summary, ok := data.(map[string]interface{})
	if !ok {
		g.addGenericContent(f, sheet, data)
		return
	}
	
	row := 5
	
	// Estilo do cabeçalho de seção
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 12,
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"E0E0E0"},
			Pattern: 1,
		},
	})
	
	// KPIs principais
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Indicadores Principais")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("F%d", row), headerStyle)
	row += 2
	
	// Cabeçalhos da tabela de KPIs
	headers := []string{"Indicador", "Valor Atual", "Valor Anterior", "Variação", "Meta", "Status"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c%d", 'A'+i, row)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	row++
	
	// Dados dos KPIs
	if kpis, ok := summary["kpis"].([]*domain.KPI); ok {
		for _, kpi := range kpis {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), kpi.DisplayName)
			
			// Valor atual
			value := fmt.Sprintf("%.2f", kpi.CurrentValue)
			if kpi.Unit != "" {
				value += " " + kpi.Unit
			}
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), value)
			
			// Valor anterior
			if kpi.PreviousValue != nil {
				prevValue := fmt.Sprintf("%.2f", *kpi.PreviousValue)
				if kpi.Unit != "" {
					prevValue += " " + kpi.Unit
				}
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), prevValue)
			}
			
			// Variação
			if kpi.TrendPercentage != nil {
				variation := fmt.Sprintf("%.1f%%", *kpi.TrendPercentage)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), variation)
				
				// Aplicar cor baseada na tendência
				var cellStyle int
				if *kpi.TrendPercentage > 0 {
					cellStyle, _ = f.NewStyle(&excelize.Style{
						Font: &excelize.Font{Color: "008000"},
					})
				} else {
					cellStyle, _ = f.NewStyle(&excelize.Style{
						Font: &excelize.Font{Color: "FF0000"},
					})
				}
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), cellStyle)
			}
			
			// Meta
			if kpi.Target != nil {
				target := fmt.Sprintf("%.2f", *kpi.Target)
				if kpi.Unit != "" {
					target += " " + kpi.Unit
				}
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), target)
			}
			
			// Status
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), kpi.Trend)
			
			row++
		}
	}
	
	row += 2
	
	// Resumo de processos
	if processes, ok := summary["processes"].(map[string]interface{}); ok {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Resumo de Processos")
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("F%d", row), headerStyle)
		row += 2
		
		// Dados de processos
		processData := [][]interface{}{
			{"Métrica", "Valor"},
			{"Total de processos", processes["total"]},
			{"Processos ativos", processes["active"]},
			{"Processos arquivados", processes["archived"]},
			{"Processos suspensos", processes["suspended"]},
			{"Novos este mês", processes["new_this_month"]},
		}
		
		for _, rowData := range processData {
			for col, value := range rowData {
				cell := fmt.Sprintf("%c%d", 'A'+col, row)
				f.SetCellValue(sheet, cell, value)
			}
			row++
		}
	}
}

// addProcessAnalysisContent adiciona conteúdo de análise de processos
func (g *ExcelGenerator) addProcessAnalysisContent(f *excelize.File, sheet string, data interface{}) {
	row := 5
	
	// Título da seção
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Análise de Processos")
	row += 2
	
	// Por enquanto, adicionar conteúdo genérico
	g.addGenericContent(f, sheet, data)
}

// addProductivityContent adiciona conteúdo de produtividade
func (g *ExcelGenerator) addProductivityContent(f *excelize.File, sheet string, data interface{}) {
	row := 5
	
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Relatório de Produtividade")
	row += 2
	
	g.addGenericContent(f, sheet, data)
}

// addFinancialContent adiciona conteúdo financeiro
func (g *ExcelGenerator) addFinancialContent(f *excelize.File, sheet string, data interface{}) {
	row := 5
	
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Relatório Financeiro")
	row += 2
	
	g.addGenericContent(f, sheet, data)
}

// addJurisprudenceContent adiciona conteúdo de jurisprudência
func (g *ExcelGenerator) addJurisprudenceContent(f *excelize.File, sheet string, data interface{}) {
	row := 5
	
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Análise de Jurisprudência")
	row += 2
	
	g.addGenericContent(f, sheet, data)
}

// addGenericContent adiciona conteúdo genérico
func (g *ExcelGenerator) addGenericContent(f *excelize.File, sheet string, data interface{}) {
	row := 5
	
	// Converter dados para string e adicionar
	content := fmt.Sprintf("%+v", data)
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), content)
}

// autoFitColumns ajusta automaticamente a largura das colunas
func (g *ExcelGenerator) autoFitColumns(f *excelize.File, sheet string) {
	// Definir larguras padrão
	columns := []struct {
		col   string
		width float64
	}{
		{"A", 30},
		{"B", 20},
		{"C", 20},
		{"D", 15},
		{"E", 15},
		{"F", 15},
	}
	
	for _, col := range columns {
		f.SetColWidth(sheet, col.col, col.col, col.width)
	}
}