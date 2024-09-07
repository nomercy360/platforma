import {
  createContext,
  createSignal,
  useContext,
  ParentComponent,
  createEffect,
} from 'solid-js'
import { User } from '~/routes/users'
import { useNavigate } from '@solidjs/router'
import { getLoggedInUser } from '~/lib/api'

function useProviderValue() {
  const [user, setUser] = createSignal<User | null>(null)
  return { user, setUser }
}

export type AuthState = ReturnType<typeof useProviderValue>

const AuthContext = createContext<AuthState | undefined>(undefined)

export const AuthProvider: ParentComponent = (props) => {
  const value = useProviderValue()

  const navigate = useNavigate()

  createEffect(async () => {
    const { data, error } = await getLoggedInUser()
    if (error) {
      navigate('/auth/login')
    }

    value.setUser(data)
  })

  return (
    <AuthContext.Provider value={value}>{props.children}</AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
