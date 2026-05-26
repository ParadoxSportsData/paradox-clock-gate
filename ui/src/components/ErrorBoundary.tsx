// src/components/ErrorBoundary.tsx
// PDX-28: Class component error boundary — catches render errors in GameView.

import { Component, type ReactNode } from 'react'

interface Props {
  children: ReactNode
}

interface State {
  hasError: boolean
  message: string
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false, message: '' }
  }

  static getDerivedStateFromError(error: unknown): State {
    return {
      hasError: true,
      message: error instanceof Error ? error.message : 'Unknown error',
    }
  }

  reset = () => this.setState({ hasError: false, message: '' })

  render() {
    if (this.state.hasError) {
      return (
        <div className="p-8 text-center">
          <div className="bg-red-900 border border-red-600 text-red-200 rounded-lg p-6 inline-block">
            <p className="font-semibold mb-2">Something went wrong.</p>
            <p className="text-sm text-red-300 mb-4">{this.state.message}</p>
            <button
              onClick={this.reset}
              className="bg-red-700 hover:bg-red-600 text-white px-4 py-2 rounded text-sm"
            >
              Try again
            </button>
          </div>
        </div>
      )
    }
    return this.props.children
  }
}
