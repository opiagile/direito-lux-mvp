'use client'

import { useState, useEffect } from 'react'
import { 
  Search as SearchIcon, 
  Filter, 
  Calendar, 
  MapPin, 
  Tag, 
  FileText, 
  Clock,
  TrendingUp,
  Eye,
  Download,
  BookOpen,
  Scale,
  Building,
  User,
  Zap,
  Star,
  Grid,
  List,
  SlidersHorizontal,
  Loader2,
  X
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { useSearchStore, SearchResult } from '@/store/search'
import { formatDate } from '@/lib/utils'

export default function SearchPage() {
  const {
    query,
    results,
    isSearching,
    recentSearches,
    filters,
    setQuery,
    setFilters,
    performSearch,
    clearResults,
    getSuggestions
  } = useSearchStore()
  
  const [searchQuery, setSearchQuery] = useState('')
  const [showFilters, setShowFilters] = useState(false)
  const [searchType, setSearchType] = useState<'all' | 'processes' | 'jurisprudence' | 'documents' | 'contacts'>('all')
  const [viewMode, setViewMode] = useState<'list' | 'grid'>('list')
  const [suggestions, setSuggestions] = useState<string[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)

  const handleSearch = async () => {
    if (!searchQuery.trim()) return
    
    // Update store query
    setQuery(searchQuery)
    
    // Apply type filter
    const typeFilters = searchType === 'all' ? [] : [
      searchType === 'processes' ? 'process' : 
      searchType === 'jurisprudence' ? 'jurisprudence' :
      searchType === 'documents' ? 'document' : 'contact'
    ]
    
    setFilters({ type: typeFilters })
    
    // Perform search
    await performSearch(searchQuery)
  }

  const handleInputChange = (value: string) => {
    setSearchQuery(value)
    
    // Get suggestions
    if (value.length > 2) {
      const newSuggestions = getSuggestions(value)
      setSuggestions(newSuggestions)
      setShowSuggestions(newSuggestions.length > 0)
    } else {
      setShowSuggestions(false)
    }
  }

  const selectSuggestion = (suggestion: string) => {
    setSearchQuery(suggestion)
    setShowSuggestions(false)
    setQuery(suggestion)
    performSearch(suggestion)
  }

  const getResultIcon = (type: SearchResult['type']) => {
    switch (type) {
      case 'process': return <Scale className="w-5 h-5 text-blue-500" />
      case 'jurisprudence': return <BookOpen className="w-5 h-5 text-purple-500" />
      case 'document': return <FileText className="w-5 h-5 text-green-500" />
      case 'contact': return <User className="w-5 h-5 text-orange-500" />
      default: return <FileText className="w-5 h-5" />
    }
  }

  const getTypeLabel = (type: SearchResult['type']): string => {
    switch (type) {
      case 'process': return 'Processo'
      case 'jurisprudence': return 'Jurisprudência'
      case 'document': return 'Documento'
      case 'contact': return 'Contato'
      default: return type
    }
  }

  const getRelevanceColor = (relevance: number): string => {
    if (relevance >= 90) return 'text-green-600'
    if (relevance >= 75) return 'text-yellow-600'
    return 'text-red-600'
  }

  // Use real search results and apply type filtering
  const filteredResults = results.filter(result => {
    if (searchType !== 'all') {
      const typeMap = {
        'processes': 'process',
        'jurisprudence': 'jurisprudence',
        'documents': 'document',
        'contacts': 'contact'
      }
      return result.type === typeMap[searchType as keyof typeof typeMap]
    }
    return true
  })

  const renderListView = () => (
    <div className="space-y-4">
      {filteredResults.map((result) => (
        <Card key={result.id} className="hover:shadow-md transition-shadow">
          <CardContent className="p-6">
            <div className="flex items-start space-x-4">
              <div className="flex-shrink-0">
                {getResultIcon(result.type)}
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h3 className="text-lg font-semibold text-primary hover:underline cursor-pointer">
                      {result.title}
                    </h3>
                    <div className="flex items-center space-x-2 mt-1 mb-2">
                      <Badge variant="outline">{getTypeLabel(result.type)}</Badge>
                      <span className="text-sm text-muted-foreground">•</span>
                      <span className="text-sm text-muted-foreground">{result.source}</span>
                      <span className="text-sm text-muted-foreground">•</span>
                      <span className="text-sm text-muted-foreground">{formatDate(result.date)}</span>
                    </div>
                    <p className="text-muted-foreground mb-3">
                      {result.description}
                    </p>
                    <div className="flex items-center space-x-2">
                      {result.tags.map((tag, index) => (
                        <Badge key={index} variant="secondary" className="text-xs">
                          {tag}
                        </Badge>
                      ))}
                    </div>
                  </div>
                  <div className="flex flex-col items-end space-y-2">
                    <div className="flex items-center space-x-1">
                      <Star className="w-4 h-4 fill-current text-yellow-500" />
                      <span className={`text-sm font-medium ${getRelevanceColor(result.relevance)}`}>
                        {result.relevance}%
                      </span>
                    </div>
                    <div className="flex space-x-1">
                      <Button variant="outline" size="sm">
                        <Eye className="w-4 h-4" />
                      </Button>
                      <Button variant="outline" size="sm">
                        <Download className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )

  const renderGridView = () => (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {filteredResults.map((result) => (
        <Card key={result.id} className="hover:shadow-md transition-shadow">
          <CardHeader className="pb-3">
            <div className="flex items-start justify-between">
              <div className="flex items-center space-x-2">
                {getResultIcon(result.type)}
                <Badge variant="outline" className="text-xs">
                  {getTypeLabel(result.type)}
                </Badge>
              </div>
              <div className="flex items-center space-x-1">
                <Star className="w-3 h-3 fill-current text-yellow-500" />
                <span className={`text-xs ${getRelevanceColor(result.relevance)}`}>
                  {result.relevance}%
                </span>
              </div>
            </div>
            <CardTitle className="text-sm line-clamp-2 cursor-pointer hover:text-primary">
              {result.title}
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <p className="text-sm text-muted-foreground mb-3 line-clamp-3">
              {result.description}
            </p>
            <div className="space-y-2">
              <div className="text-xs text-muted-foreground">
                {result.source} • {formatDate(result.date)}
              </div>
              <div className="flex flex-wrap gap-1">
                {result.tags.slice(0, 2).map((tag, index) => (
                  <Badge key={index} variant="secondary" className="text-xs">
                    {tag}
                  </Badge>
                ))}
                {result.tags.length > 2 && (
                  <Badge variant="outline" className="text-xs">
                    +{result.tags.length - 2}
                  </Badge>
                )}
              </div>
              <div className="flex justify-between pt-2">
                <Button variant="outline" size="sm">
                  <Eye className="w-3 h-3 mr-1" />
                  Ver
                </Button>
                <Button variant="outline" size="sm">
                  <Download className="w-3 h-3" />
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Busca Avançada</h1>
          <p className="text-muted-foreground">
            Encontre processos, jurisprudência, documentos e contatos
          </p>
        </div>
      </div>

      {/* Search Interface */}
      <Card>
        <CardContent className="p-6">
          {/* Main Search */}
          <div className="space-y-4">
            <div className="flex space-x-2">
              <div className="relative flex-1">
                <SearchIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-muted-foreground" />
                <Input
                  placeholder="Digite sua busca..."
                  value={searchQuery}
                  onChange={(e) => handleInputChange(e.target.value)}
                  className="pl-10 text-lg"
                  onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                />
                {/* Search Suggestions Dropdown */}
                {showSuggestions && suggestions.length > 0 && (
                  <div className="absolute top-full left-0 right-0 mt-1 bg-white border rounded-md shadow-lg z-10 max-h-48 overflow-y-auto">
                    {suggestions.map((suggestion, index) => (
                      <button
                        key={index}
                        className="w-full text-left px-4 py-2 hover:bg-muted text-sm"
                        onClick={() => selectSuggestion(suggestion)}
                      >
                        <SearchIcon className="w-3 h-3 inline mr-2 text-muted-foreground" />
                        {suggestion}
                      </button>
                    ))}
                  </div>
                )}
              </div>
              <Button onClick={handleSearch} size="lg">
                <SearchIcon className="w-4 h-4 mr-2" />
                Buscar
              </Button>
              <Button 
                variant="outline" 
                size="lg"
                onClick={() => setShowFilters(!showFilters)}
              >
                <SlidersHorizontal className="w-4 h-4 mr-2" />
                Filtros
              </Button>
            </div>

            {/* Search Type Tabs */}
            <Tabs value={searchType} onValueChange={(value) => setSearchType(value as any)}>
              <TabsList>
                <TabsTrigger value="all">Todos</TabsTrigger>
                <TabsTrigger value="processes">Processos</TabsTrigger>
                <TabsTrigger value="jurisprudence">Jurisprudência</TabsTrigger>
                <TabsTrigger value="documents">Documentos</TabsTrigger>
                <TabsTrigger value="contacts">Contatos</TabsTrigger>
              </TabsList>
            </Tabs>

            {/* Advanced Filters */}
            {showFilters && (
              <div className="border rounded-lg p-4 space-y-4 bg-muted/20">
                <h3 className="font-medium">Filtros Avançados</h3>
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                  <div className="space-y-2">
                    <Label>Data Inicial</Label>
                    <Input
                      type="date"
                      value={filters.dateRange.start}
                      onChange={(e) => setFilters({
                        ...filters,
                        dateRange: { ...filters.dateRange, start: e.target.value }
                      })}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Data Final</Label>
                    <Input
                      type="date"
                      value={filters.dateRange.end}
                      onChange={(e) => setFilters({
                        ...filters,
                        dateRange: { ...filters.dateRange, end: e.target.value }
                      })}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Tribunal</Label>
                    <select className="w-full p-2 border rounded">
                      <option value="">Todos</option>
                      <option value="TJSP">TJSP</option>
                      <option value="TRT">TRT</option>
                      <option value="STJ">STJ</option>
                      <option value="STF">STF</option>
                    </select>
                  </div>
                  <div className="space-y-2">
                    <Label>Tags</Label>
                    <Input placeholder="Ex: civil, trabalhista" />
                  </div>
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Search Suggestions */}
      {!searchQuery && (
        <div className="grid gap-6 md:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Clock className="w-5 h-5" />
                <span>Buscas Recentes</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {recentSearches.map((search, index) => (
                  <button
                    key={index}
                    className="w-full text-left p-2 hover:bg-muted rounded text-sm"
                    onClick={() => {
                      setSearchQuery(search)
                      setQuery(search)
                      performSearch(search)
                    }}
                  >
                    <SearchIcon className="w-4 h-4 inline mr-2 text-muted-foreground" />
                    {search}
                  </button>
                ))}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <TrendingUp className="w-5 h-5" />
                <span>Sugestões</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {[
                  'contratos de prestação de serviços',
                  'ações trabalhistas em SP',
                  'jurisprudência STJ responsabilidade civil',
                  'execução de honorários advocatícios',
                  'mandado de segurança contra ato administrativo'
                ].map((search, index) => (
                  <button
                    key={index}
                    className="w-full text-left p-2 hover:bg-muted rounded text-sm"
                    onClick={() => {
                      setSearchQuery(search)
                      setQuery(search)
                      performSearch(search)
                    }}
                  >
                    <Zap className="w-4 h-4 inline mr-2 text-muted-foreground" />
                    {search}
                  </button>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Search Results */}
      {(searchQuery || results.length > 0) && (
        <div className="space-y-4">
          {/* Loading State */}
          {isSearching && (
            <div className="flex items-center justify-center py-8">
              <Loader2 className="w-6 h-6 animate-spin mr-2" />
              <span>Buscando...</span>
            </div>
          )}
          
          {/* Results Header */}
          {!isSearching && (
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-lg font-semibold">
                  {searchQuery ? `Resultados da busca para "${searchQuery}"` : 'Resultados da busca'}
                </h2>
                <p className="text-sm text-muted-foreground">
                  {filteredResults.length} resultados encontrados
                </p>
              </div>
              <div className="flex items-center space-x-2">
                <Button
                  variant={viewMode === 'list' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setViewMode('list')}
                >
                  <List className="w-4 h-4" />
                </Button>
                <Button
                  variant={viewMode === 'grid' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setViewMode('grid')}
                >
                  <Grid className="w-4 h-4" />
                </Button>
              </div>
            </div>
          )}

          {/* Results */}
          {!isSearching && (
            filteredResults.length > 0 ? (
              viewMode === 'list' ? renderListView() : renderGridView()
            ) : searchQuery ? (
              <div className="text-center py-8">
                <p className="text-muted-foreground">Nenhum resultado encontrado para "{searchQuery}"</p>
                <Button 
                  variant="outline" 
                  className="mt-4"
                  onClick={() => {
                    setSearchQuery('')
                    clearResults()
                  }}
                >
                  Limpar busca
                </Button>
              </div>
            ) : null
          )}
        </div>
      )}
    </div>
  )
}