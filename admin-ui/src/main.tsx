import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import './index.css'
import App from './App.tsx'

// Clear potentially corrupted mock data from localStorage
// This fixes issues where alerts data might not be an array
try {
  const stored = localStorage.getItem('mockDataStore');
  if (stored) {
    const parsed = JSON.parse(stored);
    if (!Array.isArray(parsed.alerts) || !Array.isArray(parsed.alertRules)) {
      console.warn('Clearing corrupted mock data from localStorage');
      localStorage.removeItem('mockDataStore');
    }
  }
} catch (e) {
  console.warn('Error checking localStorage, clearing mock data');
  localStorage.removeItem('mockDataStore');
}

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      retry: 2,
      refetchOnWindowFocus: false,
    },
  },
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  </StrictMode>,
)
