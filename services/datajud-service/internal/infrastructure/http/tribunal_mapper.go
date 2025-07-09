package http

import (
	"strings"
)

// TribunalInfo informações do tribunal
type TribunalInfo struct {
	Code        string `json:"code"`         // Código do tribunal (ex: TJSP, TRF1)
	Name        string `json:"name"`         // Nome completo
	Endpoint    string `json:"endpoint"`     // Endpoint da API DataJud
	Type        string `json:"type"`         // Tipo: supremo, superior, federal, estadual, trabalho, eleitoral, militar
	Instance    int    `json:"instance"`     // Instância: 1 ou 2
	Region      string `json:"region"`       // Região/Estado
	Active      bool   `json:"active"`       // Se está ativo na API
	Description string `json:"description"`  // Descrição adicional
}

// TribunalMapper mapeador de tribunais para endpoints DataJud
type TribunalMapper struct {
	tribunals map[string]*TribunalInfo
}

// NewTribunalMapper cria novo mapeador de tribunais
func NewTribunalMapper() *TribunalMapper {
	mapper := &TribunalMapper{
		tribunals: make(map[string]*TribunalInfo),
	}
	mapper.initializeTribunals()
	return mapper
}

// initializeTribunals inicializa todos os tribunais brasileiros
func (tm *TribunalMapper) initializeTribunals() {
	// Tribunais Superiores
	tm.addSupremeTribunals()
	
	// Tribunais Regionais Federais
	tm.addFederalTribunals()
	
	// Tribunais de Justiça Estaduais
	tm.addStateTribunals()
	
	// Tribunais Regionais do Trabalho
	tm.addLaborTribunals()
	
	// Tribunais Regionais Eleitorais
	tm.addElectoralTribunals()
	
	// Tribunais Militares
	tm.addMilitaryTribunals()
}

// addSupremeTribunals adiciona tribunais superiores
func (tm *TribunalMapper) addSupremeTribunals() {
	tribunals := []*TribunalInfo{
		{
			Code:        "STF",
			Name:        "Supremo Tribunal Federal",
			Endpoint:    "api_publica_stf",
			Type:        "supremo",
			Instance:    2,
			Region:      "BR",
			Active:      true,
			Description: "Tribunal supremo do país",
		},
		{
			Code:        "STJ",
			Name:        "Superior Tribunal de Justiça",
			Endpoint:    "api_publica_stj",
			Type:        "superior",
			Instance:    2,
			Region:      "BR",
			Active:      true,
			Description: "Instância superior de justiça comum",
		},
		{
			Code:        "TST",
			Name:        "Tribunal Superior do Trabalho",
			Endpoint:    "api_publica_tst",
			Type:        "superior",
			Instance:    2,
			Region:      "BR",
			Active:      true,
			Description: "Instância superior trabalhista",
		},
		{
			Code:        "TSE",
			Name:        "Tribunal Superior Eleitoral",
			Endpoint:    "api_publica_tse",
			Type:        "superior",
			Instance:    2,
			Region:      "BR",
			Active:      true,
			Description: "Instância superior eleitoral",
		},
		{
			Code:        "STM",
			Name:        "Superior Tribunal Militar",
			Endpoint:    "api_publica_stm",
			Type:        "superior",
			Instance:    2,
			Region:      "BR",
			Active:      true,
			Description: "Instância superior militar",
		},
	}

	for _, tribunal := range tribunals {
		tm.tribunals[tribunal.Code] = tribunal
	}
}

// addFederalTribunals adiciona tribunais regionais federais
func (tm *TribunalMapper) addFederalTribunals() {
	federalTribunals := []*TribunalInfo{
		{
			Code:        "TRF1",
			Name:        "Tribunal Regional Federal da 1ª Região",
			Endpoint:    "api_publica_trf1",
			Type:        "federal",
			Instance:    2,
			Region:      "1ª Região",
			Active:      true,
			Description: "AC, AM, AP, BA, DF, GO, MA, MG, MT, PA, PI, RO, RR, TO",
		},
		{
			Code:        "TRF2",
			Name:        "Tribunal Regional Federal da 2ª Região",
			Endpoint:    "api_publica_trf2",
			Type:        "federal",
			Instance:    2,
			Region:      "2ª Região",
			Active:      true,
			Description: "ES, RJ",
		},
		{
			Code:        "TRF3",
			Name:        "Tribunal Regional Federal da 3ª Região",
			Endpoint:    "api_publica_trf3",
			Type:        "federal",
			Instance:    2,
			Region:      "3ª Região",
			Active:      true,
			Description: "MS, SP",
		},
		{
			Code:        "TRF4",
			Name:        "Tribunal Regional Federal da 4ª Região",
			Endpoint:    "api_publica_trf4",
			Type:        "federal",
			Instance:    2,
			Region:      "4ª Região",
			Active:      true,
			Description: "PR, RS, SC",
		},
		{
			Code:        "TRF5",
			Name:        "Tribunal Regional Federal da 5ª Região",
			Endpoint:    "api_publica_trf5",
			Type:        "federal",
			Instance:    2,
			Region:      "5ª Região",
			Active:      true,
			Description: "AL, CE, PB, PE, RN, SE",
		},
		{
			Code:        "TRF6",
			Name:        "Tribunal Regional Federal da 6ª Região",
			Endpoint:    "api_publica_trf6",
			Type:        "federal",
			Instance:    2,
			Region:      "6ª Região",
			Active:      true,
			Description: "MG (criado em 2022)",
		},
	}

	for _, tribunal := range federalTribunals {
		tm.tribunals[tribunal.Code] = tribunal
	}
}

// addStateTribunals adiciona tribunais de justiça estaduais
func (tm *TribunalMapper) addStateTribunals() {
	states := []struct {
		code, name, uf string
	}{
		{"TJAC", "Tribunal de Justiça do Acre", "AC"},
		{"TJAL", "Tribunal de Justiça de Alagoas", "AL"},
		{"TJAP", "Tribunal de Justiça do Amapá", "AP"},
		{"TJAM", "Tribunal de Justiça do Amazonas", "AM"},
		{"TJBA", "Tribunal de Justiça da Bahia", "BA"},
		{"TJCE", "Tribunal de Justiça do Ceará", "CE"},
		{"TJDFT", "Tribunal de Justiça do Distrito Federal e Territórios", "DF"},
		{"TJES", "Tribunal de Justiça do Espírito Santo", "ES"},
		{"TJGO", "Tribunal de Justiça de Goiás", "GO"},
		{"TJMA", "Tribunal de Justiça do Maranhão", "MA"},
		{"TJMT", "Tribunal de Justiça de Mato Grosso", "MT"},
		{"TJMS", "Tribunal de Justiça de Mato Grosso do Sul", "MS"},
		{"TJMG", "Tribunal de Justiça de Minas Gerais", "MG"},
		{"TJPA", "Tribunal de Justiça do Pará", "PA"},
		{"TJPB", "Tribunal de Justiça da Paraíba", "PB"},
		{"TJPR", "Tribunal de Justiça do Paraná", "PR"},
		{"TJPE", "Tribunal de Justiça de Pernambuco", "PE"},
		{"TJPI", "Tribunal de Justiça do Piauí", "PI"},
		{"TJRJ", "Tribunal de Justiça do Rio de Janeiro", "RJ"},
		{"TJRN", "Tribunal de Justiça do Rio Grande do Norte", "RN"},
		{"TJRS", "Tribunal de Justiça do Rio Grande do Sul", "RS"},
		{"TJRO", "Tribunal de Justiça de Rondônia", "RO"},
		{"TJRR", "Tribunal de Justiça de Roraima", "RR"},
		{"TJSC", "Tribunal de Justiça de Santa Catarina", "SC"},
		{"TJSP", "Tribunal de Justiça de São Paulo", "SP"},
		{"TJSE", "Tribunal de Justiça de Sergipe", "SE"},
		{"TJTO", "Tribunal de Justiça do Tocantins", "TO"},
	}

	for _, state := range states {
		tribunal := &TribunalInfo{
			Code:        state.code,
			Name:        state.name,
			Endpoint:    "api_publica_" + strings.ToLower(state.code),
			Type:        "estadual",
			Instance:    2,
			Region:      state.uf,
			Active:      true,
			Description: "Tribunal de Justiça estadual de " + state.uf,
		}
		tm.tribunals[tribunal.Code] = tribunal
	}
}

// addLaborTribunals adiciona tribunais regionais do trabalho
func (tm *TribunalMapper) addLaborTribunals() {
	laborTribunals := []struct {
		region int
		name   string
		states string
	}{
		{1, "Tribunal Regional do Trabalho da 1ª Região", "RJ"},
		{2, "Tribunal Regional do Trabalho da 2ª Região", "SP"},
		{3, "Tribunal Regional do Trabalho da 3ª Região", "MG"},
		{4, "Tribunal Regional do Trabalho da 4ª Região", "RS"},
		{5, "Tribunal Regional do Trabalho da 5ª Região", "BA"},
		{6, "Tribunal Regional do Trabalho da 6ª Região", "PE"},
		{7, "Tribunal Regional do Trabalho da 7ª Região", "CE"},
		{8, "Tribunal Regional do Trabalho da 8ª Região", "PA, AP"},
		{9, "Tribunal Regional do Trabalho da 9ª Região", "PR"},
		{10, "Tribunal Regional do Trabalho da 10ª Região", "DF, TO"},
		{11, "Tribunal Regional do Trabalho da 11ª Região", "AM, RR"},
		{12, "Tribunal Regional do Trabalho da 12ª Região", "SC"},
		{13, "Tribunal Regional do Trabalho da 13ª Região", "PB"},
		{14, "Tribunal Regional do Trabalho da 14ª Região", "RO, AC"},
		{15, "Tribunal Regional do Trabalho da 15ª Região", "Campinas/SP"},
		{16, "Tribunal Regional do Trabalho da 16ª Região", "MA"},
		{17, "Tribunal Regional do Trabalho da 17ª Região", "ES"},
		{18, "Tribunal Regional do Trabalho da 18ª Região", "GO"},
		{19, "Tribunal Regional do Trabalho da 19ª Região", "AL"},
		{20, "Tribunal Regional do Trabalho da 20ª Região", "SE"},
		{21, "Tribunal Regional do Trabalho da 21ª Região", "RN"},
		{22, "Tribunal Regional do Trabalho da 22ª Região", "PI"},
		{23, "Tribunal Regional do Trabalho da 23ª Região", "MT"},
		{24, "Tribunal Regional do Trabalho da 24ª Região", "MS"},
	}

	for _, trt := range laborTribunals {
		code := "TRT" + formatRegion(trt.region)
		tribunal := &TribunalInfo{
			Code:        code,
			Name:        trt.name,
			Endpoint:    "api_publica_" + strings.ToLower(code),
			Type:        "trabalho",
			Instance:    2,
			Region:      formatRegion(trt.region) + "ª Região",
			Active:      true,
			Description: trt.states,
		}
		tm.tribunals[tribunal.Code] = tribunal
	}
}

// addElectoralTribunals adiciona tribunais regionais eleitorais
func (tm *TribunalMapper) addElectoralTribunals() {
	states := []struct {
		uf, name string
	}{
		{"AC", "Tribunal Regional Eleitoral do Acre"},
		{"AL", "Tribunal Regional Eleitoral de Alagoas"},
		{"AP", "Tribunal Regional Eleitoral do Amapá"},
		{"AM", "Tribunal Regional Eleitoral do Amazonas"},
		{"BA", "Tribunal Regional Eleitoral da Bahia"},
		{"CE", "Tribunal Regional Eleitoral do Ceará"},
		{"DF", "Tribunal Regional Eleitoral do Distrito Federal"},
		{"ES", "Tribunal Regional Eleitoral do Espírito Santo"},
		{"GO", "Tribunal Regional Eleitoral de Goiás"},
		{"MA", "Tribunal Regional Eleitoral do Maranhão"},
		{"MT", "Tribunal Regional Eleitoral de Mato Grosso"},
		{"MS", "Tribunal Regional Eleitoral de Mato Grosso do Sul"},
		{"MG", "Tribunal Regional Eleitoral de Minas Gerais"},
		{"PA", "Tribunal Regional Eleitoral do Pará"},
		{"PB", "Tribunal Regional Eleitoral da Paraíba"},
		{"PR", "Tribunal Regional Eleitoral do Paraná"},
		{"PE", "Tribunal Regional Eleitoral de Pernambuco"},
		{"PI", "Tribunal Regional Eleitoral do Piauí"},
		{"RJ", "Tribunal Regional Eleitoral do Rio de Janeiro"},
		{"RN", "Tribunal Regional Eleitoral do Rio Grande do Norte"},
		{"RS", "Tribunal Regional Eleitoral do Rio Grande do Sul"},
		{"RO", "Tribunal Regional Eleitoral de Rondônia"},
		{"RR", "Tribunal Regional Eleitoral de Roraima"},
		{"SC", "Tribunal Regional Eleitoral de Santa Catarina"},
		{"SP", "Tribunal Regional Eleitoral de São Paulo"},
		{"SE", "Tribunal Regional Eleitoral de Sergipe"},
		{"TO", "Tribunal Regional Eleitoral do Tocantins"},
	}

	for _, state := range states {
		code := "TRE" + state.uf
		tribunal := &TribunalInfo{
			Code:        code,
			Name:        state.name,
			Endpoint:    "api_publica_" + strings.ToLower(code),
			Type:        "eleitoral",
			Instance:    2,
			Region:      state.uf,
			Active:      true,
			Description: "Tribunal Regional Eleitoral de " + state.uf,
		}
		tm.tribunals[tribunal.Code] = tribunal
	}
}

// addMilitaryTribunals adiciona tribunais militares
func (tm *TribunalMapper) addMilitaryTribunals() {
	militaryTribunals := []*TribunalInfo{
		{
			Code:        "TJMSP",
			Name:        "Tribunal de Justiça Militar do Estado de São Paulo",
			Endpoint:    "api_publica_tjmsp",
			Type:        "militar",
			Instance:    2,
			Region:      "SP",
			Active:      true,
			Description: "Justiça Militar estadual de São Paulo",
		},
		{
			Code:        "TJMMG",
			Name:        "Tribunal de Justiça Militar do Estado de Minas Gerais",
			Endpoint:    "api_publica_tjmmg",
			Type:        "militar",
			Instance:    2,
			Region:      "MG",
			Active:      true,
			Description: "Justiça Militar estadual de Minas Gerais",
		},
		{
			Code:        "TJMRS",
			Name:        "Tribunal de Justiça Militar do Estado do Rio Grande do Sul",
			Endpoint:    "api_publica_tjmrs",
			Type:        "militar",
			Instance:    2,
			Region:      "RS",
			Active:      true,
			Description: "Justiça Militar estadual do Rio Grande do Sul",
		},
	}

	for _, tribunal := range militaryTribunals {
		tm.tribunals[tribunal.Code] = tribunal
	}
}

// GetTribunal obtém informações de um tribunal por código
func (tm *TribunalMapper) GetTribunal(code string) *TribunalInfo {
	code = strings.ToUpper(strings.TrimSpace(code))
	
	// Busca direta
	if tribunal, exists := tm.tribunals[code]; exists {
		return tribunal
	}
	
	// Tentativas de mapeamento alternativo
	alternatives := tm.getAlternativeCodes(code)
	for _, alt := range alternatives {
		if tribunal, exists := tm.tribunals[alt]; exists {
			return tribunal
		}
	}
	
	return nil
}

// getAlternativeCodes gera códigos alternativos para mapeamento
func (tm *TribunalMapper) getAlternativeCodes(code string) []string {
	var alternatives []string
	
	// Remove caracteres especiais e espaços
	cleanCode := strings.ReplaceAll(code, " ", "")
	cleanCode = strings.ReplaceAll(cleanCode, "-", "")
	cleanCode = strings.ReplaceAll(cleanCode, "_", "")
	
	alternatives = append(alternatives, cleanCode)
	
	// Se contém números, tenta formatar
	if len(cleanCode) > 3 {
		// Ex: TRF01 -> TRF1
		if strings.Contains(cleanCode, "0") {
			alternatives = append(alternatives, strings.ReplaceAll(cleanCode, "0", ""))
		}
		
		// Ex: TRT1 -> TRT01
		if len(cleanCode) == 4 && (strings.HasPrefix(cleanCode, "TRT") || strings.HasPrefix(cleanCode, "TRE")) {
			prefix := cleanCode[:3]
			suffix := cleanCode[3:]
			if len(suffix) == 1 {
				alternatives = append(alternatives, prefix+"0"+suffix)
			}
		}
	}
	
	return alternatives
}

// GetTribunalsByType obtém tribunais por tipo
func (tm *TribunalMapper) GetTribunalsByType(tribunalType string) []*TribunalInfo {
	var result []*TribunalInfo
	
	for _, tribunal := range tm.tribunals {
		if tribunal.Type == tribunalType && tribunal.Active {
			result = append(result, tribunal)
		}
	}
	
	return result
}

// GetTribunalsByRegion obtém tribunais por região/estado
func (tm *TribunalMapper) GetTribunalsByRegion(region string) []*TribunalInfo {
	var result []*TribunalInfo
	region = strings.ToUpper(strings.TrimSpace(region))
	
	for _, tribunal := range tm.tribunals {
		if tribunal.Region == region && tribunal.Active {
			result = append(result, tribunal)
		}
	}
	
	return result
}

// GetAllTribunals obtém todos os tribunais ativos
func (tm *TribunalMapper) GetAllTribunals() []*TribunalInfo {
	var result []*TribunalInfo
	
	for _, tribunal := range tm.tribunals {
		if tribunal.Active {
			result = append(result, tribunal)
		}
	}
	
	return result
}

// IsValidTribunal verifica se um código de tribunal é válido
func (tm *TribunalMapper) IsValidTribunal(code string) bool {
	return tm.GetTribunal(code) != nil
}

// GetTribunalEndpoint obtém o endpoint da API para um tribunal
func (tm *TribunalMapper) GetTribunalEndpoint(code string) string {
	tribunal := tm.GetTribunal(code)
	if tribunal != nil {
		return tribunal.Endpoint
	}
	return ""
}

// GetTribunalStats obtém estatísticas dos tribunais
func (tm *TribunalMapper) GetTribunalStats() map[string]int {
	stats := make(map[string]int)
	
	for _, tribunal := range tm.tribunals {
		if tribunal.Active {
			stats[tribunal.Type]++
		}
	}
	
	stats["total"] = len(tm.tribunals)
	return stats
}

// formatRegion formata número da região (1 -> "1", 10 -> "10")
func formatRegion(region int) string {
	if region < 10 {
		return "0" + string(rune('0'+region))
	}
	return string(rune('0' + region/10)) + string(rune('0' + region%10))
}

// SearchTribunals busca tribunais por nome ou código
func (tm *TribunalMapper) SearchTribunals(query string) []*TribunalInfo {
	var result []*TribunalInfo
	query = strings.ToLower(strings.TrimSpace(query))
	
	if query == "" {
		return result
	}
	
	for _, tribunal := range tm.tribunals {
		if tribunal.Active {
			// Busca por código
			if strings.Contains(strings.ToLower(tribunal.Code), query) {
				result = append(result, tribunal)
				continue
			}
			
			// Busca por nome
			if strings.Contains(strings.ToLower(tribunal.Name), query) {
				result = append(result, tribunal)
				continue
			}
			
			// Busca por região
			if strings.Contains(strings.ToLower(tribunal.Region), query) {
				result = append(result, tribunal)
				continue
			}
		}
	}
	
	return result
}