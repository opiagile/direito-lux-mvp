import { useUsageStore } from '@/store/usage'
import { toast } from 'sonner'

export const useUsageActions = () => {
  const { incrementUsage } = useUsageStore()

  const recordDataJudQuery = () => {
    incrementUsage('datajudQueries', 1)
    console.log('📊 DataJud query recorded')
  }

  const recordAISummary = () => {
    incrementUsage('aiSummaries', 1)
    toast.success('Resumo de IA gerado')
    console.log('🤖 AI Summary recorded')
  }

  const recordReport = () => {
    incrementUsage('reports', 1)
    toast.success('Relatório gerado')
    console.log('📋 Report recorded')
  }

  const recordMCPCommand = () => {
    incrementUsage('mcpCommands', 1)
    console.log('🤖 MCP Command recorded')
  }

  const recordProcessCreation = () => {
    incrementUsage('totalProcesses', 1)
    console.log('⚖️ Process creation recorded')
  }

  // Simulate some background usage for demonstration
  const simulateRealisticActivity = () => {
    const actions = [
      { action: recordDataJudQuery, weight: 0.4 }, // 40% chance
      { action: recordAISummary, weight: 0.2 },    // 20% chance  
      { action: recordReport, weight: 0.1 },       // 10% chance
      { action: recordMCPCommand, weight: 0.2 },   // 20% chance
    ]

    actions.forEach(({ action, weight }) => {
      if (Math.random() < weight) {
        // Delay to simulate realistic usage patterns
        setTimeout(action, Math.random() * 5000)
      }
    })
  }

  return {
    recordDataJudQuery,
    recordAISummary,
    recordReport,
    recordMCPCommand,
    recordProcessCreation,
    simulateRealisticActivity
  }
}