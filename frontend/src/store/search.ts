import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { Process } from '@/types'
import { useProcessStore } from './processes'

export interface SearchResult {
  id: string
  type: 'process' | 'jurisprudence' | 'document' | 'contact'
  title: string
  description: string
  highlight?: string
  relevance: number
  date: string
  source: string
  tags: string[]
  metadata: Record<string, any>
  originalData?: any
}

interface SearchFilters {
  type: string[]
  court: string[]
  status: string[]
  priority: string[]
  dateRange: {
    start: string
    end: string
  }
  tags: string[]
  monitoring?: boolean
}

interface SearchState {
  // Current search
  query: string
  filters: SearchFilters
  results: SearchResult[]
  isSearching: boolean
  
  // History
  recentSearches: string[]
  savedSearches: Array<{
    id: string
    name: string
    query: string
    filters: SearchFilters
  }>
  
  // Actions
  setQuery: (query: string) => void
  setFilters: (filters: Partial<SearchFilters>) => void
  performSearch: (query?: string) => Promise<SearchResult[]>
  clearResults: () => void
  addRecentSearch: (query: string) => void
  saveSearch: (name: string) => void
  removeSavedSearch: (id: string) => void
  getSuggestions: (query: string) => string[]
}

// ❌ MOCK REMOVIDO - Usar API real do Search Service (port 8086)
// Dados de jurisprudência devem vir de: GET /api/v1/search/jurisprudence

// ❌ MOCK REMOVIDO - Usar API real do Document Service
// Documentos devem vir de: GET /api/v1/documents

// ❌ MOCK REMOVIDO - Usar API real do Contact Service  
// Contatos devem vir de: GET /api/v1/contacts

export const useSearchStore = create<SearchState>()(
  persist(
    (set, get) => ({
      query: '',
      filters: {
        type: [],
        court: [],
        status: [],
        priority: [],
        dateRange: { start: '', end: '' },
        tags: [],
        monitoring: undefined
      },
      results: [],
      isSearching: false,
      recentSearches: [],
      savedSearches: [],

      setQuery: (query) => set({ query }),

      setFilters: (newFilters) => 
        set((state) => ({
          filters: { ...state.filters, ...newFilters }
        })),

      performSearch: async (searchQuery) => {
        const { query, filters } = get()
        const finalQuery = searchQuery || query
        
        if (!finalQuery.trim()) {
          set({ results: [] })
          return []
        }

        set({ isSearching: true })

        try {
          // Simulate API delay
          await new Promise(resolve => setTimeout(resolve, 500))

          const allResults: SearchResult[] = []
          const queryLower = finalQuery.toLowerCase()

          // Search processes
          const processStore = useProcessStore.getState()
          const processes = processStore.processes || []
          
          processes.forEach(process => {
            const relevance = calculateRelevance(process, queryLower, filters)
            if (relevance > 0) {
              allResults.push({
                id: `process_${process.id}`,
                type: 'process',
                title: `Processo ${process.number}`,
                description: process.subject,
                highlight: getHighlight(process, queryLower),
                relevance,
                date: process.createdAt,
                source: process.court,
                tags: process.tags,
                metadata: {
                  court: process.court,
                  status: process.status,
                  priority: process.priority,
                  monitoring: process.monitoring,
                  value: process.estimatedValue
                },
                originalData: process
              })
            }
          })

          // ❌ MOCK REMOVIDO - TODO: Implementar busca real de jurisprudência
          // await searchJurisprudenceAPI(query).then(jurisprudence => {
          //   jurisprudence.forEach(jur => allResults.push(...))
          // })
          // ❌ MOCK REMOVIDO - TODO: Implementar busca real de documentos
          // await searchDocumentsAPI(query).then(documents => {
          //   documents.forEach(doc => allResults.push(...))
          // })
          // ❌ MOCK REMOVIDO - TODO: Implementar busca real de contatos  
          // await searchContactsAPI(query).then(contacts => {
          //   contacts.forEach(contact => allResults.push(...))
          // })

          // Apply filters
          const filteredResults = applyFilters(allResults, filters)
          
          // Sort by relevance
          const sortedResults = filteredResults.sort((a, b) => b.relevance - a.relevance)

          set({ results: sortedResults, isSearching: false })
          
          // Add to recent searches
          if (finalQuery.trim()) {
            get().addRecentSearch(finalQuery.trim())
          }

          return sortedResults
        } catch (error) {
          console.error('Search error:', error)
          set({ isSearching: false, results: [] })
          return []
        }
      },

      clearResults: () => set({ results: [], query: '' }),

      addRecentSearch: (query) => {
        const current = get().recentSearches
        const updated = [
          query,
          ...current.filter(q => q !== query)
        ].slice(0, 10) // Keep only last 10
        
        set({ recentSearches: updated })
      },

      saveSearch: (name) => {
        const { query, filters } = get()
        const newSearch = {
          id: `search_${Date.now()}`,
          name,
          query,
          filters
        }
        
        set((state) => ({
          savedSearches: [...state.savedSearches, newSearch]
        }))
      },

      removeSavedSearch: (id) => {
        set((state) => ({
          savedSearches: state.savedSearches.filter(s => s.id !== id)
        }))
      },

      getSuggestions: (query) => {
        const queryLower = query.toLowerCase()
        const suggestions: string[] = []
        
        // Add suggestions from processes
        const processStore = useProcessStore.getState()
        const processes = processStore.processes || []
        
        processes.forEach(process => {
          if (process.number.toLowerCase().includes(queryLower)) {
            suggestions.push(process.number)
          }
          if (process.subject.toLowerCase().includes(queryLower)) {
            suggestions.push(process.subject.substring(0, 50) + '...')
          }
          process.tags.forEach(tag => {
            if (tag.toLowerCase().includes(queryLower)) {
              suggestions.push(tag)
            }
          })
        })

        return [...new Set(suggestions)].slice(0, 5)
      }
    }),
    {
      name: 'search-storage',
      partialize: (state) => ({
        recentSearches: state.recentSearches,
        savedSearches: state.savedSearches
      })
    }
  )
)

// Helper functions
function calculateRelevance(process: Process, query: string, filters: SearchFilters): number {
  let score = 0
  
  // Exact number match gets highest score
  if (process.number.toLowerCase().includes(query)) score += 100
  
  // Subject match
  if (process.subject.toLowerCase().includes(query)) score += 80
  
  // Court match
  if (process.court.toLowerCase().includes(query)) score += 60
  
  // Type match
  if (process.type.toLowerCase().includes(query)) score += 50
  
  // Tags match
  process.tags.forEach(tag => {
    if (tag.toLowerCase().includes(query)) score += 40
  })
  
  // Lawyer match
  if (process.lawyer?.toLowerCase().includes(query)) score += 30
  
  // Party match
  process.parties.forEach(party => {
    if (party.name.toLowerCase().includes(query)) score += 35
  })

  return score
}

function calculateJurisprudenceRelevance(jur: any, query: string): number {
  let score = 0
  
  if (jur.title.toLowerCase().includes(query)) score += 90
  if (jur.description.toLowerCase().includes(query)) score += 70
  if (jur.court.toLowerCase().includes(query)) score += 50
  
  jur.tags.forEach((tag: string) => {
    if (tag.toLowerCase().includes(query)) score += 40
  })

  return score
}

function calculateDocumentRelevance(doc: any, query: string): number {
  let score = 0
  
  if (doc.title.toLowerCase().includes(query)) score += 85
  if (doc.description.toLowerCase().includes(query)) score += 65
  if (doc.author?.toLowerCase().includes(query)) score += 45
  
  doc.tags.forEach((tag: string) => {
    if (tag.toLowerCase().includes(query)) score += 35
  })

  return score
}

function calculateContactRelevance(contact: any, query: string): number {
  let score = 0
  
  if (contact.name.toLowerCase().includes(query)) score += 95
  if (contact.description.toLowerCase().includes(query)) score += 75
  if (contact.specialty.toLowerCase().includes(query)) score += 55
  
  contact.tags.forEach((tag: string) => {
    if (tag.toLowerCase().includes(query)) score += 40
  })

  return score
}

function getHighlight(process: Process, query: string): string {
  if (process.number.toLowerCase().includes(query)) return process.number
  if (process.subject.toLowerCase().includes(query)) return process.subject.substring(0, 100) + '...'
  return process.type
}

function getJurisprudenceHighlight(jur: any, query: string): string {
  if (jur.title.toLowerCase().includes(query)) return jur.title.substring(0, 100) + '...'
  return jur.description.substring(0, 100) + '...'
}

function getDocumentHighlight(doc: any, query: string): string {
  if (doc.title.toLowerCase().includes(query)) return doc.title
  return doc.description.substring(0, 100) + '...'
}

function getContactHighlight(contact: any, query: string): string {
  if (contact.name.toLowerCase().includes(query)) return contact.name
  return contact.description.substring(0, 100) + '...'
}

function applyFilters(results: SearchResult[], filters: SearchFilters): SearchResult[] {
  return results.filter(result => {
    // Type filter
    if (filters.type.length > 0 && !filters.type.includes(result.type)) {
      return false
    }

    // Date range filter
    if (filters.dateRange.start && result.date < filters.dateRange.start) {
      return false
    }
    if (filters.dateRange.end && result.date > filters.dateRange.end) {
      return false
    }

    // Process-specific filters
    if (result.type === 'process') {
      const metadata = result.metadata
      
      if (filters.court.length > 0 && !filters.court.some(court => 
        metadata.court?.toLowerCase().includes(court.toLowerCase())
      )) {
        return false
      }
      
      if (filters.status.length > 0 && !filters.status.includes(metadata.status)) {
        return false
      }
      
      if (filters.priority.length > 0 && !filters.priority.includes(metadata.priority)) {
        return false
      }
      
      if (filters.monitoring !== undefined && metadata.monitoring !== filters.monitoring) {
        return false
      }
    }

    return true
  })
}