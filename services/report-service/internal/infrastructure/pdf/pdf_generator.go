package pdf

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/domain"
)

// PDFGenerator implementação do gerador de PDF
type PDFGenerator struct {
	logger *zap.Logger
}

// NewPDFGenerator cria nova instância do gerador
func NewPDFGenerator(logger *zap.Logger) *PDFGenerator {
	return &PDFGenerator{
		logger: logger,
	}
}

// GeneratePDF gera relatório em PDF
func (g *PDFGenerator) GeneratePDF(ctx context.Context, report *domain.Report, data interface{}) ([]byte, error) {
	g.logger.Debug("Generating PDF report",
		zap.String("report_id", report.ID.String()),
		zap.String("type", string(report.Type)))

	// Criar novo documento PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	
	// Adicionar página
	pdf.AddPage()
	
	// Configurar fontes
	pdf.SetFont("Arial", "", 12)
	
	// Adicionar header
	g.addHeader(pdf, report)
	
	// Adicionar conteúdo baseado no tipo
	switch report.Type {
	case domain.ReportTypeExecutiveSummary:
		g.addExecutiveSummaryContent(pdf, data)
		
	case domain.ReportTypeProcessAnalysis:
		g.addProcessAnalysisContent(pdf, data)
		
	case domain.ReportTypeProductivity:
		g.addProductivityContent(pdf, data)
		
	case domain.ReportTypeFinancial:
		g.addFinancialContent(pdf, data)
		
	case domain.ReportTypeJurisprudence:
		g.addJurisprudenceContent(pdf, data)
		
	default:
		g.addGenericContent(pdf, data)
	}
	
	// Adicionar footer
	g.addFooter(pdf, report)
	
	// Gerar bytes do PDF
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	
	return buf.Bytes(), nil
}

// addHeader adiciona cabeçalho ao PDF
func (g *PDFGenerator) addHeader(pdf *gofpdf.Fpdf, report *domain.Report) {
	// Logo e título
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(190, 10, "Direito Lux - Sistema Jurídico")
	pdf.Ln(10)
	
	// Título do relatório
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, report.Title)
	pdf.Ln(10)
	
	// Data e hora
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(190, 5, fmt.Sprintf("Gerado em: %s", time.Now().Format("02/01/2006 15:04")))
	pdf.Ln(10)
	
	// Linha separadora
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(5)
	
	// Resetar cor do texto
	pdf.SetTextColor(0, 0, 0)
}

// addFooter adiciona rodapé ao PDF
func (g *PDFGenerator) addFooter(pdf *gofpdf.Fpdf, report *domain.Report) {
	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(0, 10, fmt.Sprintf("Página %d", pdf.PageNo()), "", 0, "C", false, 0, "")
}

// addExecutiveSummaryContent adiciona conteúdo do resumo executivo
func (g *PDFGenerator) addExecutiveSummaryContent(pdf *gofpdf.Fpdf, data interface{}) {
	summary, ok := data.(map[string]interface{})
	if !ok {
		g.addGenericContent(pdf, data)
		return
	}
	
	// KPIs principais
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Indicadores Principais")
	pdf.Ln(10)
	
	if kpis, ok := summary["kpis"].([]*domain.KPI); ok {
		pdf.SetFont("Arial", "", 11)
		for _, kpi := range kpis {
			// Nome do KPI
			pdf.SetFont("Arial", "B", 11)
			pdf.Cell(60, 6, kpi.DisplayName+":")
			
			// Valor
			pdf.SetFont("Arial", "", 11)
			valueStr := fmt.Sprintf("%.2f", kpi.CurrentValue)
			if kpi.Unit != "" {
				valueStr += " " + kpi.Unit
			}
			pdf.Cell(40, 6, valueStr)
			
			// Tendência
			if kpi.TrendPercentage != nil {
				if *kpi.TrendPercentage > 0 {
					pdf.SetTextColor(0, 128, 0)
					pdf.Cell(30, 6, fmt.Sprintf("↑ %.1f%%", *kpi.TrendPercentage))
				} else {
					pdf.SetTextColor(255, 0, 0)
					pdf.Cell(30, 6, fmt.Sprintf("↓ %.1f%%", *kpi.TrendPercentage))
				}
				pdf.SetTextColor(0, 0, 0)
			}
			
			pdf.Ln(8)
		}
	}
	
	pdf.Ln(10)
	
	// Resumo de processos
	if processes, ok := summary["processes"].(map[string]interface{}); ok {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(190, 10, "Resumo de Processos")
		pdf.Ln(10)
		
		pdf.SetFont("Arial", "", 11)
		
		if total, ok := processes["total"].(int); ok {
			pdf.Cell(190, 6, fmt.Sprintf("Total de processos: %d", total))
			pdf.Ln(6)
		}
		
		if active, ok := processes["active"].(int); ok {
			pdf.Cell(190, 6, fmt.Sprintf("Processos ativos: %d", active))
			pdf.Ln(6)
		}
		
		if archived, ok := processes["archived"].(int); ok {
			pdf.Cell(190, 6, fmt.Sprintf("Processos arquivados: %d", archived))
			pdf.Ln(6)
		}
	}
}

// addProcessAnalysisContent adiciona conteúdo de análise de processos
func (g *PDFGenerator) addProcessAnalysisContent(pdf *gofpdf.Fpdf, data interface{}) {
	// Título da seção
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Análise de Processos")
	pdf.Ln(10)
	
	// Implementar análise específica
	pdf.SetFont("Arial", "", 11)
	
	// Por enquanto, conteúdo genérico
	g.addGenericContent(pdf, data)
}

// addProductivityContent adiciona conteúdo de produtividade
func (g *PDFGenerator) addProductivityContent(pdf *gofpdf.Fpdf, data interface{}) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Relatório de Produtividade")
	pdf.Ln(10)
	
	// Implementar métricas de produtividade
	pdf.SetFont("Arial", "", 11)
	
	g.addGenericContent(pdf, data)
}

// addFinancialContent adiciona conteúdo financeiro
func (g *PDFGenerator) addFinancialContent(pdf *gofpdf.Fpdf, data interface{}) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Relatório Financeiro")
	pdf.Ln(10)
	
	// Implementar análise financeira
	pdf.SetFont("Arial", "", 11)
	
	g.addGenericContent(pdf, data)
}

// addJurisprudenceContent adiciona conteúdo de jurisprudência
func (g *PDFGenerator) addJurisprudenceContent(pdf *gofpdf.Fpdf, data interface{}) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Análise de Jurisprudência")
	pdf.Ln(10)
	
	// Implementar análise de jurisprudência
	pdf.SetFont("Arial", "", 11)
	
	g.addGenericContent(pdf, data)
}

// addGenericContent adiciona conteúdo genérico
func (g *PDFGenerator) addGenericContent(pdf *gofpdf.Fpdf, data interface{}) {
	pdf.SetFont("Arial", "", 11)
	
	// Converter dados para string e adicionar
	content := fmt.Sprintf("%+v", data)
	
	// Quebrar texto em linhas
	lines := pdf.SplitLines([]byte(content), 180)
	for _, line := range lines {
		pdf.Cell(190, 6, string(line))
		pdf.Ln(6)
	}
}