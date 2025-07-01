'use client'

import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { X } from 'lucide-react'
import { Process, ProcessStatus, ProcessPriority } from '@/types'
import { useProcessDataStore } from '@/store'

const processSchema = z.object({
  number: z.string().min(1, 'Número do processo é obrigatório'),
  type: z.string().min(1, 'Tipo de processo é obrigatório'),
  subject: z.string().min(1, 'Assunto é obrigatório'),
  court: z.string().min(1, 'Tribunal é obrigatório'),
  status: z.enum(['active', 'suspended', 'archived', 'concluded'] as const),
  priority: z.enum(['low', 'medium', 'high', 'urgent'] as const),
  lawyer: z.string().optional(),
  estimatedValue: z.number().min(0).optional(),
  monitoring: z.boolean(),
})

type ProcessFormData = z.infer<typeof processSchema>

interface ProcessModalProps {
  isOpen: boolean
  onClose: () => void
  processId?: string
}

const statusOptions = [
  { value: 'active', label: 'Ativo' },
  { value: 'suspended', label: 'Suspenso' },
  { value: 'archived', label: 'Arquivado' },
  { value: 'concluded', label: 'Concluído' },
]

const priorityOptions = [
  { value: 'low', label: 'Baixa' },
  { value: 'medium', label: 'Média' },
  { value: 'high', label: 'Alta' },
  { value: 'urgent', label: 'Urgente' },
]

const typeOptions = [
  'Ação de Cobrança',
  'Ação Trabalhista', 
  'Divórcio Consensual',
  'Inventário',
  'Ação de Despejo',
  'Ação Declaratória',
  'Ação Ordinária',
  'Execução',
  'Mandado de Segurança',
  'Habeas Corpus',
]

const courtOptions = [
  'TJSP - 1ª Vara Cível',
  'TJSP - 2ª Vara Cível', 
  'TJSP - Vara de Família',
  'TJSP - Vara do Trabalho',
  'TRT - 2ª Região',
  'TRT - 15ª Região',
  'TRF - 3ª Região',
  'STJ - Superior Tribunal de Justiça',
  'STF - Supremo Tribunal Federal',
]

export function ProcessModal({ isOpen, onClose, processId }: ProcessModalProps) {
  const { addProcess, updateProcess, getProcessById } = useProcessDataStore()
  const [tags, setTags] = useState<string[]>([])
  const [newTag, setNewTag] = useState('')

  const isEditing = !!processId
  const editingProcess = processId ? getProcessById(processId) : null

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    watch,
    reset,
  } = useForm<ProcessFormData>({
    resolver: zodResolver(processSchema),
    defaultValues: {
      number: '',
      type: '',
      subject: '',
      court: '',
      status: 'active',
      priority: 'medium',
      lawyer: '',
      estimatedValue: 0,
      monitoring: false,
    },
  })

  // Load editing data
  useEffect(() => {
    if (isEditing && editingProcess) {
      setValue('number', editingProcess.number)
      setValue('type', editingProcess.type)
      setValue('subject', editingProcess.subject)
      setValue('court', editingProcess.court)
      setValue('status', editingProcess.status)
      setValue('priority', editingProcess.priority)
      setValue('lawyer', editingProcess.lawyer || '')
      setValue('estimatedValue', editingProcess.estimatedValue || 0)
      setValue('monitoring', editingProcess.monitoring)
      setTags(editingProcess.tags || [])
    } else {
      reset()
      setTags([])
    }
  }, [isEditing, editingProcess, setValue, reset])

  const onSubmit = (data: ProcessFormData) => {
    if (isEditing && processId) {
      // When editing, preserve existing data that's not in the form
      const existingProcess = getProcessById(processId)
      const processData = {
        ...data,
        tags,
        parties: existingProcess?.parties || [],
        movements: existingProcess?.movements || [],
        tenantId: existingProcess?.tenantId || '11111111-1111-1111-1111-111111111111',
      }
      updateProcess(processId, processData)
    } else {
      // When creating new, set defaults
      const processData = {
        ...data,
        tenantId: '11111111-1111-1111-1111-111111111111',
        parties: [],
        movements: [],
        tags,
      }
      addProcess(processData)
    }

    onClose()
  }

  const addTag = () => {
    if (newTag.trim() && !tags.includes(newTag.trim())) {
      setTags([...tags, newTag.trim()])
      setNewTag('')
    }
  }

  const removeTag = (tagToRemove: string) => {
    setTags(tags.filter(tag => tag !== tagToRemove))
  }

  const handleClose = () => {
    reset()
    setTags([])
    setNewTag('')
    onClose()
  }

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {isEditing ? 'Editar Processo' : 'Novo Processo'}
          </DialogTitle>
          <DialogDescription>
            {isEditing 
              ? 'Edite as informações do processo' 
              : 'Preencha as informações para criar um novo processo'
            }
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="number">Número do Processo *</Label>
              <Input
                id="number"
                {...register('number')}
                placeholder="5001234-20.2023.4.03.6109"
              />
              {errors.number && (
                <p className="text-sm text-red-600">{errors.number.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="type">Tipo de Processo *</Label>
              <Select value={watch('type')} onValueChange={(value) => setValue('type', value)}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione o tipo" />
                </SelectTrigger>
                <SelectContent>
                  {typeOptions.map((type) => (
                    <SelectItem key={type} value={type}>
                      {type}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.type && (
                <p className="text-sm text-red-600">{errors.type.message}</p>
              )}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="subject">Assunto *</Label>
            <Textarea
              id="subject"
              {...register('subject')}
              placeholder="Descreva o assunto do processo..."
              rows={3}
            />
            {errors.subject && (
              <p className="text-sm text-red-600">{errors.subject.message}</p>
            )}
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="court">Tribunal *</Label>
              <Select value={watch('court')} onValueChange={(value) => setValue('court', value)}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione o tribunal" />
                </SelectTrigger>
                <SelectContent>
                  {courtOptions.map((court) => (
                    <SelectItem key={court} value={court}>
                      {court}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.court && (
                <p className="text-sm text-red-600">{errors.court.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="lawyer">Advogado Responsável</Label>
              <Input
                id="lawyer"
                {...register('lawyer')}
                placeholder="Dr. João Silva"
              />
            </div>
          </div>

          <div className="grid gap-4 md:grid-cols-3">
            <div className="space-y-2">
              <Label htmlFor="status">Status</Label>
              <Select value={watch('status')} onValueChange={(value: ProcessStatus) => setValue('status', value)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {statusOptions.map((option) => (
                    <SelectItem key={option.value} value={option.value}>
                      {option.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="priority">Prioridade</Label>
              <Select value={watch('priority')} onValueChange={(value: ProcessPriority) => setValue('priority', value)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {priorityOptions.map((option) => (
                    <SelectItem key={option.value} value={option.value}>
                      {option.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="estimatedValue">Valor Estimado (R$)</Label>
              <Input
                id="estimatedValue"
                type="number"
                {...register('estimatedValue', { valueAsNumber: true })}
                placeholder="0"
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label>Tags</Label>
            <div className="flex space-x-2">
              <Input
                value={newTag}
                onChange={(e) => setNewTag(e.target.value)}
                placeholder="Adicionar tag..."
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault()
                    addTag()
                  }
                }}
              />
              <Button type="button" onClick={addTag} variant="outline">
                Adicionar
              </Button>
            </div>
            <div className="flex flex-wrap gap-2">
              {tags.map((tag) => (
                <Badge key={tag} variant="secondary" className="flex items-center gap-1">
                  {tag}
                  <button
                    type="button"
                    onClick={() => removeTag(tag)}
                    className="ml-1 hover:text-red-600"
                  >
                    <X className="w-3 h-3" />
                  </button>
                </Badge>
              ))}
            </div>
          </div>

          <div className="flex items-center space-x-2">
            <input
              type="checkbox"
              id="monitoring"
              {...register('monitoring')}
              className="rounded"
            />
            <Label htmlFor="monitoring">Ativar monitoramento automático</Label>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={handleClose}>
              Cancelar
            </Button>
            <Button type="submit">
              {isEditing ? 'Salvar Alterações' : 'Criar Processo'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}