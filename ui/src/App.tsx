// src/App.tsx
// PDX-28: Top-level app shell. Routes between GameSelector list view and GameView.
// No router library — single selectedGameId state drives the two-screen flow.

import { useState } from 'react'
import { GameSelector } from './components/GameSelector'
import { GameView } from './components/GameView'

function App() {
  const [selectedGameId, setSelectedGameId] = useState<string | null>(null)

  if (selectedGameId) {
    return <GameView gameId={selectedGameId} onBack={() => setSelectedGameId(null)} />
  }

  return (
    <div className="min-h-screen bg-gray-950 text-gray-100">
      <header className="px-6 py-5 bg-gray-900 border-b border-gray-800">
        <h1 className="text-2xl font-bold tracking-tight">
          clock-gate <span className="text-blue-400 font-mono text-lg">UI</span>
        </h1>
      </header>
      <GameSelector onSelect={setSelectedGameId} />
    </div>
  )
}

export default App
