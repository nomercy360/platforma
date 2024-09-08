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
import { createUser, listUsers } from '~/lib/api'
import { IconClose, IconPlus } from '~/components/icons'
import { Checkbox } from '~/components/checkbox'
import { useTableSelection } from '~/lib/table'
import SideMenu from '~/components/side-menu'

export type User = {
  id: number
  name: string
  email: string
  avatar_url: string
  created_at: string
  updated_at: string
  deleted_at: string | null
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
    queryKey: ['users'],
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

  const { selected, toggleSelection, toggleSelectAll } = useTableSelection()

  return (
    <div class="flex min-h-screen w-full flex-col rounded-tl-2xl bg-background">
      <SideMenu
        title="Create new Platform User"
        subtitle="Invite will be sent to the user's email"
        isOpen={sideMenuIsOpen()}
        setIsOpen={setSideMenuIsOpen}>
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
      </SideMenu>
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
            <TableHead class="w-8">
              <Checkbox
                checked={selected().length === query.data?.length}
                onChange={() => toggleSelectAll(query.data!)}
              />
            </TableHead>
            <TableHead>Nickname</TableHead>
            <TableHead>Email</TableHead>
            <TableHead class="text-right">Role</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={query.data}>
            {(user) => (
              <TableRow>
                <TableCell>
                  <Checkbox
                    onChange={() => toggleSelection(user.id)}
                    checked={selected().includes(user.id)}
                  />
                </TableCell>
                <TableCell>
                  <div class="flex items-center">
                    <img
                      src={user.avatar_url}
                      alt={user.name}
                      class="mr-2 size-6 rounded object-cover"
                    />
                    <span>{user.name}</span>
                  </div>
                </TableCell>
                <TableCell class="text-muted-foreground">
                  {user.email}
                </TableCell>
                <TableCell class="text-right">{user.role}</TableCell>
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
