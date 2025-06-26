'use client'

import React, { Suspense } from 'react'
import { LoadingScreen } from '@/components/ui/loading-screen'

interface LoadingBoundaryProps {
  children: React.ReactNode
  fallback?: React.ReactNode
}

export function LoadingBoundary({ children, fallback }: LoadingBoundaryProps) {
  return (
    <Suspense fallback={fallback || <LoadingScreen />}>
      {children}
    </Suspense>
  )
}

// Error boundary component for chunk loading errors
interface ErrorBoundaryState {
  hasError: boolean
  error: Error | null
}

class ChunkLoadErrorBoundary extends React.Component<
  React.PropsWithChildren<{}>,
  ErrorBoundaryState
> {
  constructor(props: React.PropsWithChildren<{}>) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    // Check if it's a chunk loading error
    if (error.name === 'ChunkLoadError' || error.message.includes('Loading chunk')) {
      return { hasError: true, error }
    }
    return { hasError: false, error: null }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    // Log the error for debugging
    console.error('ChunkLoadError caught:', error, errorInfo)
    
    // Attempt to reload the page for chunk loading errors
    if (error.name === 'ChunkLoadError' || error.message.includes('Loading chunk')) {
      // Wait a bit before reloading to avoid infinite loops
      setTimeout(() => {
        window.location.reload()
      }, 1000)
    }
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="flex flex-col items-center justify-center min-h-screen p-4">
          <div className="text-center">
            <h2 className="text-xl font-semibold mb-2">Loading Error</h2>
            <p className="text-gray-600 mb-4">
              There was an error loading the application. The page will reload automatically.
            </p>
            <button
              onClick={() => window.location.reload()}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Reload Now
            </button>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}

export { ChunkLoadErrorBoundary }