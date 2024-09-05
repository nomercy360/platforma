import { Router } from '@solidjs/router'
import { FileRoutes } from '@solidjs/start/router'
import { Match, Suspense, Switch } from 'solid-js'
import Nav from '~/components/Nav'
import './app.css'
import Controls from '~/components/controls'
import {
  QueryClient,
  QueryClientProvider,
  createQuery,
} from '@tanstack/solid-query'

export const queryClient = new QueryClient()

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router
        root={(props) => (
          <Switch>
            <Match when={props.location.pathname.startsWith('/auth')}>
              <main class="flex h-screen w-full flex-col items-center justify-end bg-black/30">
                <Suspense>{props.children}</Suspense>
              </main>
            </Match>
            <Match when={props.location.pathname !== '/login'}>
              <div class="min-h-screen bg-background">
                <Controls />
                <div class="flex flex-row items-start pr-4">
                  <Nav />
                  <Suspense>{props.children}</Suspense>
                </div>
              </div>
            </Match>
          </Switch>
        )}>
        <FileRoutes />
      </Router>
    </QueryClientProvider>
  )
}
