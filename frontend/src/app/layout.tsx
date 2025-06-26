import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { Providers } from '@/components/providers'
import { ChunkLoadErrorBoundary } from '@/components/loading-boundary'
import { Toaster } from 'sonner'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Direito Lux - Gestão Jurídica Inteligente',
  description: 'Plataforma completa para gestão de processos jurídicos com IA',
  keywords: 'gestão jurídica, processos, advogados, IA, automação legal',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="pt-BR" suppressHydrationWarning>
      <body className={inter.className}>
        <ChunkLoadErrorBoundary>
          <Providers>
            {children}
            <Toaster 
              position="top-right" 
              richColors 
              closeButton
              expand
            />
          </Providers>
        </ChunkLoadErrorBoundary>
      </body>
    </html>
  )
}