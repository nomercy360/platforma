import { createSignal, For, Show } from 'solid-js'
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '~/components/table'
import { SearchInput } from '~/components/input'
import { createQuery } from '@tanstack/solid-query'
import { createUser, listCustomers, listUsers } from '~/lib/api'
import { IconClose, IconPlus } from '~/components/icons'

export type User = {
  id: number
  name: string
  email: string
  role: string
}

export default function UsersPage() {
  const [sideMenuIsOpen, setSideMenuIsOpen] = createSignal(false)

  const [user, setUser] = createSignal<{
    email: string
    password: string
    name: string | null
  }>({
    email: '',
    password: '',
    name: null,
  })

  const query = createQuery(() => ({
    queryKey: ['customers'],
    queryFn: async () => {
      const { data } = await listUsers()
      return data as User[]
    },
  }))

  async function onFormSubmit() {
    const { email, password, name } = user()
    await createUser({ email, password, name })
    await query.refetch()
    setSideMenuIsOpen(false)
  }

  return (
    <div class="flex min-h-screen w-full flex-col rounded-tl-2xl bg-background">
      <Show when={sideMenuIsOpen()}>
        <div class="absolute right-0 top-0 z-50 flex h-full w-[670px] flex-col items-start justify-start space-y-4 bg-secondary p-4">
          <div class="flex w-full flex-row items-center justify-between">
            <p class="text-2xl font-semibold">Add User</p>
            <button onClick={() => setSideMenuIsOpen(false)}>
              <IconClose class="size-5" />
            </button>
          </div>
          <form
            class="flex w-full flex-col space-y-3"
            onSubmit={(e) => {
              e.preventDefault()
              onFormSubmit()
            }}>
            <label class="text-sm font-semibold" for="name">
              <input
                class="mt-1 h-11 w-full rounded-lg border border-neutral-200 bg-background p-2"
                id="name"
                type="text"
                value={user().name || ''}
                onInput={(e) =>
                  setUser({
                    ...user(),
                    name: e.currentTarget.value,
                  })
                }
                placeholder="John Doe"
              />
            </label>
            <label class="text-sm font-semibold" for="email">
              <input
                class="h-11 w-full rounded-lg border border-neutral-200 bg-background p-2"
                id="email"
                type="email"
                value={user().email}
                onInput={(e) =>
                  setUser({
                    ...user(),
                    email: e.currentTarget.value,
                  })
                }
                placeholder="dummy@example.com"
              />
            </label>
            <label class="text-sm font-semibold" for="password">
              <input
                class="mt-1 h-11 w-full rounded-lg border border-neutral-200 bg-background p-2"
                id="password"
                type="password"
                value={user().password}
                onInput={(e) =>
                  setUser({
                    ...user(),
                    password: e.currentTarget.value,
                  })
                }
                placeholder="********"
              />
            </label>
            <button
              class="mt-4 w-full rounded bg-primary p-2 text-white"
              type="submit">
              Add User
            </button>
          </form>
        </div>
      </Show>
      <div class="flex w-full flex-row items-center justify-between p-4">
        <SearchInput
          class="w-96 bg-background"
          placeholder="Search by name, email, or nickname"
        />
        <button
          class="flex size-4 items-center justify-center"
          onClick={() => setSideMenuIsOpen(true)}>
          <IconPlus class="size-4" />
        </button>
      </div>
      <Table>
        <TableCaption>A list of CRM users</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Name</TableHead>
            <TableHead>Email</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={query.data}>
            {(user) => (
              <TableRow>
                <TableCell>{user.id}</TableCell>
                <TableCell>{user.name}</TableCell>
                <TableCell>{user.email}</TableCell>
              </TableRow>
            )}
          </For>
        </TableBody>
      </Table>
    </div>
  )
}

function Chip({ text }: { text: string }) {
  let color

  switch (text) {
    case 'Admin':
      color = 'bg-red-200 text-red-800'
      break
    case 'Owner':
      color = 'bg-green-200 text-green-800'
      break
    case 'Manager':
      color = 'bg-blue-200 text-blue-800'
      break
    default:
      break
  }

  return (
    <div
      class={`flex h-6 items-center justify-center rounded-full px-2 py-1 text-sm ${color}`}>
      {text}
    </div>
  )
}
