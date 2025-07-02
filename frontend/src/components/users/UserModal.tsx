'use client'

import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { toast } from 'sonner'
import { User, UserRole } from '@/types'

const userSchema = z.object({
  name: z.string().min(1, 'Nome é obrigatório').min(2, 'Nome deve ter pelo menos 2 caracteres'),
  email: z.string().min(1, 'Email é obrigatório').email('Email deve ser válido'),
  role: z.enum(['admin', 'manager', 'lawyer', 'assistant'] as const, {
    required_error: 'Role é obrigatório'
  }),
  isActive: z.boolean().default(true),
})

type UserFormData = z.infer<typeof userSchema>

interface UserModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: UserFormData) => void
  user?: User | null
  quotaInfo?: {
    used: number
    limit: number
    canAdd: boolean
  }
}

export default function UserModal({ 
  isOpen, 
  onClose, 
  onSubmit, 
  user,
  quotaInfo 
}: UserModalProps) {
  const isEditing = !!user

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    reset,
    formState: { errors, isSubmitting }
  } = useForm<UserFormData>({
    resolver: zodResolver(userSchema),
    defaultValues: {
      name: '',
      email: '',
      role: 'assistant',
      isActive: true,
    }
  })

  const selectedRole = watch('role')

  // Reset form when modal opens/closes or user changes
  useEffect(() => {
    if (isOpen) {
      if (user) {
        setValue('name', user.name)
        setValue('email', user.email)
        setValue('role', user.role)
        setValue('isActive', user.isActive)
      } else {
        reset()
      }
    }
  }, [isOpen, user, setValue, reset])

  const onFormSubmit = async (data: UserFormData) => {
    // Check quota for new users
    if (!isEditing && quotaInfo && !quotaInfo.canAdd) {
      toast.error('Limite de usuários atingido para seu plano atual')
      return
    }

    try {
      onSubmit(data)
      onClose()
      reset()
    } catch (error) {
      console.error('Erro ao salvar usuário:', error)
      toast.error('Erro ao salvar usuário')
    }
  }

  const roleOptions = [
    { value: 'admin', label: 'Administrador', description: 'Acesso total ao sistema' },
    { value: 'manager', label: 'Gerente', description: 'Gestão operacional e relatórios' },
    { value: 'lawyer', label: 'Advogado', description: 'Gestão de processos e clientes' },
    { value: 'assistant', label: 'Assistente', description: 'Acesso limitado e suporte' },
  ]

  const getSelectedRoleInfo = () => {
    return roleOptions.find(option => option.value === selectedRole)
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>
            {isEditing ? 'Editar Usuário' : 'Novo Usuário'}
          </DialogTitle>
        </DialogHeader>

        {/* Quota Warning for new users */}
        {!isEditing && quotaInfo && (
          <div className={`p-3 rounded-lg border ${!quotaInfo.canAdd ? 'bg-red-50 border-red-200' : 'bg-blue-50 border-blue-200'}`}>
            <div className={`text-sm ${!quotaInfo.canAdd ? 'text-red-800' : 'text-blue-800'}`}>
              <strong>Quota de Usuários:</strong> {quotaInfo.used} de {quotaInfo.limit === -1 ? '∞' : quotaInfo.limit} usuários utilizados
            </div>
            {!quotaInfo.canAdd && (
              <div className="mt-2 space-y-2">
                <div className="text-sm text-red-600 font-medium">
                  ⚠️ Limite de usuários atingido para seu plano atual!
                </div>
                <div className="text-xs text-red-600">
                  Para adicionar mais usuários, faça upgrade para:
                </div>
                <div className="text-xs bg-white rounded p-2 border border-red-300">
                  • <strong>Professional</strong> (até 5 usuários) - R$ 299/mês<br/>
                  • <strong>Business</strong> (até 15 usuários) - R$ 699/mês<br/>
                  • <strong>Enterprise</strong> (usuários ilimitados) - R$ 1999/mês
                </div>
                <Button 
                  type="button" 
                  variant="outline" 
                  size="sm" 
                  className="w-full mt-2 border-red-300 text-red-600 hover:bg-red-50"
                  onClick={() => window.open('/billing', '_blank')}
                >
                  🚀 Fazer Upgrade do Plano
                </Button>
              </div>
            )}
          </div>
        )}

        <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-4">
          {/* Name */}
          <div className="space-y-2">
            <Label htmlFor="name">Nome Completo</Label>
            <Input
              id="name"
              {...register('name')}
              placeholder="Digite o nome completo"
              className={errors.name ? 'border-red-500' : ''}
            />
            {errors.name && (
              <p className="text-sm text-red-600">{errors.name.message}</p>
            )}
          </div>

          {/* Email */}
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              {...register('email')}
              placeholder="usuario@empresa.com.br"
              className={errors.email ? 'border-red-500' : ''}
            />
            {errors.email && (
              <p className="text-sm text-red-600">{errors.email.message}</p>
            )}
          </div>

          {/* Role */}
          <div className="space-y-2">
            <Label htmlFor="role">Função/Role</Label>
            <Select
              value={selectedRole}
              onValueChange={(value: UserRole) => setValue('role', value)}
            >
              <SelectTrigger className={errors.role ? 'border-red-500' : ''}>
                <SelectValue placeholder="Selecione a função" />
              </SelectTrigger>
              <SelectContent>
                {roleOptions.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    <div className="flex flex-col">
                      <span className="font-medium">{option.label}</span>
                      <span className="text-xs text-muted-foreground">{option.description}</span>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {errors.role && (
              <p className="text-sm text-red-600">{errors.role.message}</p>
            )}
            
            {/* Role Description */}
            {selectedRole && (
              <div className="text-xs text-muted-foreground bg-gray-50 p-2 rounded">
                <strong>{getSelectedRoleInfo()?.label}:</strong> {getSelectedRoleInfo()?.description}
              </div>
            )}
          </div>

          {/* Status - Only for editing */}
          {isEditing && (
            <div className="space-y-2">
              <Label htmlFor="isActive">Status</Label>
              <Select
                value={watch('isActive') ? 'active' : 'inactive'}
                onValueChange={(value) => setValue('isActive', value === 'active')}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="active">
                    <div className="flex items-center space-x-2">
                      <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                      <span>Ativo</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="inactive">
                    <div className="flex items-center space-x-2">
                      <div className="w-2 h-2 bg-red-500 rounded-full"></div>
                      <span>Inativo</span>
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
          )}

          {/* Actions */}
          <div className="flex justify-end space-x-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={onClose}
            >
              Cancelar
            </Button>
            <Button
              type="submit"
              disabled={isSubmitting || (!isEditing && quotaInfo && !quotaInfo.canAdd)}
            >
              {isSubmitting ? 'Salvando...' : (isEditing ? 'Atualizar' : 'Criar Usuário')}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}