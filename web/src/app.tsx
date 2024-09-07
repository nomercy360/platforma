import { Router } from '@solidjs/router'
import { FileRoutes } from '@solidjs/start/router'
import { Match, Suspense, Switch } from 'solid-js'
import Navbar from '~/components/navbar'
import './app.css'
import Controls from '~/components/controls'
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query'
import { AuthProvider } from '~/lib/auth-context'

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
              <AuthProvider>
                <div class="min-h-screen">
                  <Controls />
                  <div class="flex flex-row items-start">
                    <Navbar />
                    <div class="ml-14 flex min-h-screen w-full flex-col rounded-tl-2xl bg-background pb-10">
                      <Suspense>{props.children}</Suspense>
                    </div>
                  </div>
                </div>
              </AuthProvider>
            </Match>
          </Switch>
        )}>
        <FileRoutes />
      </Router>
    </QueryClientProvider>
  )
}
