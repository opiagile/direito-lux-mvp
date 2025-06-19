package report

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/domain"
	"github.com/direito-lux/report-service/internal/infrastructure/excel"
	"github.com/direito-lux/report-service/internal/infrastructure/pdf"
)

// Generator implementação do gerador de relatórios
type Generator struct {
	pdfGenerator   *pdf.PDFGenerator
	excelGenerator *excel.ExcelGenerator
	logger         *zap.Logger
}

// NewGenerator cria nova instância do gerador
func NewGenerator(logger *zap.Logger) *Generator {
	return &Generator{
		pdfGenerator:   pdf.NewPDFGenerator(logger),
		excelGenerator: excel.NewExcelGenerator(logger),
		logger:         logger,
	}
}

// GeneratePDF implementa domain.ReportGenerator
func (g *Generator) GeneratePDF(ctx context.Context, report *domain.Report, data interface{}) ([]byte, error) {
	return g.pdfGenerator.GeneratePDF(ctx, report, data)
}

// GenerateExcel implementa domain.ReportGenerator
func (g *Generator) GenerateExcel(ctx context.Context, report *domain.Report, data interface{}) ([]byte, error) {
	return g.excelGenerator.GenerateExcel(ctx, report, data)
}

// GenerateCSV implementa domain.ReportGenerator
func (g *Generator) GenerateCSV(ctx context.Context, report *domain.Report, data interface{}) ([]byte, error) {
	g.logger.Debug("Generating CSV report",
		zap.String("report_id", report.ID.String()),
		zap.String("type", string(report.Type)))

	var records [][]string

	// Adicionar cabeçalho
	records = append(records, []string{
		"Direito Lux - " + report.Title,
		"Gerado em: " + report.CreatedAt.Format("02/01/2006 15:04"),
	})
	records = append(records, []string{}) // Linha em branco

	// Processar dados baseado no tipo
	switch v := data.(type) {
	case map[string]interface{}:
		records = append(records, g.processMapToCSV(v)...)
	case []interface{}:
		records = append(records, g.processSliceToCSV(v)...)
	case []*domain.KPI:
		records = append(records, g.processKPIsToCSV(v)...)
	default:
		// Fallback genérico
		records = append(records, []string{"Dados", fmt.Sprintf("%+v", data)})
	}

	// Gerar CSV
	var output strings.Builder
	writer := csv.NewWriter(&output)

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to generate CSV: %w", err)
	}

	return []byte(output.String()), nil
}

// GenerateHTML implementa domain.ReportGenerator
func (g *Generator) GenerateHTML(ctx context.Context, report *domain.Report, data interface{}) ([]byte, error) {
	g.logger.Debug("Generating HTML report",
		zap.String("report_id", report.ID.String()),
		zap.String("type", string(report.Type)))

	var html strings.Builder

	// HTML básico
	html.WriteString(`<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + report.Title + `</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { text-align: center; border-bottom: 2px solid #ccc; padding-bottom: 20px; margin-bottom: 30px; }
        .title { font-size: 24px; font-weight: bold; color: #333; }
        .subtitle { font-size: 18px; margin: 10px 0; }
        .date { font-size: 12px; color: #666; }
        .section { margin: 20px 0; }
        .section-title { font-size: 16px; font-weight: bold; margin-bottom: 10px; background: #f5f5f5; padding: 10px; }
        table { width: 100%; border-collapse: collapse; margin: 10px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; font-weight: bold; }
        .kpi-positive { color: green; }
        .kpi-negative { color: red; }
        .footer { margin-top: 50px; text-align: center; font-size: 10px; color: #666; border-top: 1px solid #ccc; padding-top: 10px; }
    </style>
</head>
<body>`)

	// Cabeçalho
	html.WriteString(`
    <div class="header">
        <div class="title">Direito Lux - Sistema Jurídico</div>
        <div class="subtitle">` + report.Title + `</div>
        <div class="date">Gerado em: ` + report.CreatedAt.Format("02/01/2006 15:04") + `</div>
    </div>`)

	// Conteúdo baseado no tipo
	switch report.Type {
	case domain.ReportTypeExecutiveSummary:
		g.addHTMLExecutiveSummary(&html, data)
	case domain.ReportTypeProcessAnalysis:
		g.addHTMLProcessAnalysis(&html, data)
	default:
		g.addHTMLGenericContent(&html, data)
	}

	// Rodapé
	html.WriteString(`
    <div class="footer">
        <p>Relatório gerado pelo sistema Direito Lux</p>
        <p>Para mais informações, acesse: https://direito-lux.com</p>
    </div>
</body>
</html>`)

	return []byte(html.String()), nil
}

// processMapToCSV converte map para CSV
func (g *Generator) processMapToCSV(data map[string]interface{}) [][]string {
	var records [][]string

	for key, value := range data {
		switch v := value.(type) {
		case []*domain.KPI:
			records = append(records, []string{key + " (KPIs)"})
			records = append(records, g.processKPIsToCSV(v)...)
		case map[string]interface{}:
			records = append(records, []string{key})
			for k, val := range v {
				records = append(records, []string{k, fmt.Sprintf("%v", val)})
			}
		default:
			records = append(records, []string{key, fmt.Sprintf("%v", value)})
		}
		records = append(records, []string{}) // Linha em branco
	}

	return records
}

// processSliceToCSV converte slice para CSV
func (g *Generator) processSliceToCSV(data []interface{}) [][]string {
	var records [][]string

	for i, item := range data {
		records = append(records, []string{fmt.Sprintf("Item %d", i+1), fmt.Sprintf("%+v", item)})
	}

	return records
}

// processKPIsToCSV converte KPIs para CSV
func (g *Generator) processKPIsToCSV(kpis []*domain.KPI) [][]string {
	var records [][]string

	// Cabeçalho
	records = append(records, []string{
		"Indicador", "Valor Atual", "Valor Anterior", "Variação (%)", "Meta", "Status",
	})

	// Dados
	for _, kpi := range kpis {
		var prevValue, variation, target string

		if kpi.PreviousValue != nil {
			prevValue = fmt.Sprintf("%.2f", *kpi.PreviousValue)
		}

		if kpi.TrendPercentage != nil {
			variation = fmt.Sprintf("%.1f", *kpi.TrendPercentage)
		}

		if kpi.Target != nil {
			target = fmt.Sprintf("%.2f", *kpi.Target)
		}

		records = append(records, []string{
			kpi.DisplayName,
			fmt.Sprintf("%.2f %s", kpi.CurrentValue, kpi.Unit),
			prevValue,
			variation,
			target,
			kpi.Trend,
		})
	}

	return records
}

// addHTMLExecutiveSummary adiciona resumo executivo em HTML
func (g *Generator) addHTMLExecutiveSummary(html *strings.Builder, data interface{}) {
	summary, ok := data.(map[string]interface{})
	if !ok {
		g.addHTMLGenericContent(html, data)
		return
	}

	// KPIs
	if kpis, ok := summary["kpis"].([]*domain.KPI); ok {
		html.WriteString(`<div class="section">
		<div class="section-title">Indicadores Principais</div>
		<table>
			<thead>
				<tr>
					<th>Indicador</th>
					<th>Valor Atual</th>
					<th>Variação</th>
					<th>Meta</th>
					<th>Status</th>
				</tr>
			</thead>
			<tbody>`)

		for _, kpi := range kpis {
			var trendClass string
			var variation string

			if kpi.TrendPercentage != nil {
				if *kpi.TrendPercentage > 0 {
					trendClass = "kpi-positive"
					variation = fmt.Sprintf("↑ %.1f%%", *kpi.TrendPercentage)
				} else {
					trendClass = "kpi-negative"
					variation = fmt.Sprintf("↓ %.1f%%", *kpi.TrendPercentage)
				}
			}

			var target string
			if kpi.Target != nil {
				target = fmt.Sprintf("%.2f %s", *kpi.Target, kpi.Unit)
			}

			html.WriteString(fmt.Sprintf(`
				<tr>
					<td>%s</td>
					<td>%.2f %s</td>
					<td class="%s">%s</td>
					<td>%s</td>
					<td>%s</td>
				</tr>`,
				kpi.DisplayName,
				kpi.CurrentValue, kpi.Unit,
				trendClass, variation,
				target,
				kpi.Trend))
		}

		html.WriteString(`</tbody></table></div>`)
	}

	// Resumo de processos
	if processes, ok := summary["processes"].(map[string]interface{}); ok {
		html.WriteString(`<div class="section">
		<div class="section-title">Resumo de Processos</div>
		<table>`)

		for key, value := range processes {
			html.WriteString(fmt.Sprintf(`
			<tr>
				<td><strong>%s</strong></td>
				<td>%v</td>
			</tr>`, key, value))
		}

		html.WriteString(`</table></div>`)
	}
}

// addHTMLProcessAnalysis adiciona análise de processos em HTML
func (g *Generator) addHTMLProcessAnalysis(html *strings.Builder, data interface{}) {
	html.WriteString(`<div class="section">
	<div class="section-title">Análise de Processos</div>`)

	g.addHTMLGenericContent(html, data)

	html.WriteString(`</div>`)
}

// addHTMLGenericContent adiciona conteúdo genérico em HTML
func (g *Generator) addHTMLGenericContent(html *strings.Builder, data interface{}) {
	// Serializar dados como JSON formatado
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		html.WriteString(fmt.Sprintf("<pre>%+v</pre>", data))
		return
	}

	html.WriteString(fmt.Sprintf("<pre>%s</pre>", string(jsonData)))
}

// getFieldValue obtém valor de campo usando reflection
func (g *Generator) getFieldValue(obj interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}

	return field.Interface()
}