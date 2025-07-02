import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { User } from '@/types'
import { useAuthStore } from '@/store'

// API client para usuários
const usersAPI = {
  list: async (): Promise<User[]> => {
    const { tenant } = useAuthStore.getState()
    if (!tenant) throw new Error('Tenant não encontrado')

    const token = localStorage.getItem('auth_token')
    const response = await fetch('http://localhost:8081/api/v1/users', {
      headers: {
        'Content-Type': 'application/json',
        'X-Tenant-ID': tenant.id,
        ...(token && { 'Authorization': `Bearer ${token}` })
      }
    })

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    const users = await response.json()
    
    // Converter formato backend → frontend
    return users.map((user: any): User => ({
      id: user.id,
      email: user.email,
      name: `${user.first_name} ${user.last_name}`.trim(),
      role: user.role === 'operator' ? 'lawyer' : user.role, // Mapear operator → lawyer
      tenantId: tenant.id, // Usar tenant ID do auth store
      isActive: user.status === 'active',
      createdAt: user.created_at,
      updatedAt: user.created_at // Backend não tem updated_at separado
    }))
  }
}

// Hook para listar usuários reais da API
export const useUsersAPI = () => {
  const { tenant } = useAuthStore()

  return useQuery({
    queryKey: ['users-api', tenant?.id],
    queryFn: usersAPI.list,
    enabled: !!tenant,
    staleTime: 30 * 1000, // 30 segundos
    retry: 2,
    refetchOnWindowFocus: false
  })
}

export default useUsersAPI
