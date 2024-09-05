import { createSignal } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { signInWithPassword } from '~/lib/api'

export default function Login() {
  const [email, setEmail] = createSignal('')
  const [password, setPassword] = createSignal('')

  const navigate = useNavigate()

  const loginUser = async (e: Event) => {
    e.preventDefault()
    const { data, error } = await signInWithPassword({
      email: email(),
      password: password(),
    })

    if (error) {
      alert(error)
      return
    }

    if (data) {
      navigate('/')
    }
  }

  return (
    <div class="flex h-[90%] w-3/4 flex-col items-center justify-between rounded-t-xl bg-secondary py-16 text-center">
      <div class="flex flex-col items-center justify-center space-y-4">
        <div class="flex size-12 items-center justify-center rounded-full bg-primary text-primary-foreground">
          P/
        </div>
        <p class="text-xl">SignIn on Platform</p>
      </div>
      <form
        class="flex w-full max-w-sm flex-col space-y-4"
        onSubmit={(e) => loginUser(e)}>
        <div class="flex w-full flex-col space-y-1">
          <label for="email">
            <input
              type="email"
              id="email"
              name="email"
              placeholder="Your email"
              onChange={(e) => setEmail(e.target.value)}
              class="h-11 w-full rounded-lg bg-background px-2"
            />
          </label>
          <label for="password">
            <input
              type="password"
              id="password"
              name="password"
              placeholder="Your password"
              onChange={(e) => setPassword(e.target.value)}
              class="h-11 w-full rounded-lg bg-background px-2"
            />
          </label>
        </div>
        <button class="h-11 w-full rounded-lg bg-primary text-primary-foreground">
          Sign In
        </button>
      </form>
      <a class="hover:underline" href="/auth/register">
        Register instead
      </a>
    </div>
  )
}
