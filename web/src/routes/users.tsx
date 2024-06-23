import { createSignal, For } from 'solid-js'
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
import { Switch } from '~/components/switch'
import { ToggleGroup, ToggleGroupItem } from '~/components/toggle-group'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/select'

const users = [
  {
    id: 1,
    name: 'Nikita Axel',
    username: 'nikitaaxel',
    email: 'nikita•••@proton.me',
    avatar_url: 'https://source.unsplash.com/random/100x100',
    role: 'Admin',
  },
  {
    id: 2,
    name: 'Emily Rodriguez',
    username: 'emilyrodriguez',
    email: 'emily•••@gmail.com',
    avatar_url: 'https://source.unsplash.com/random/100x100',
    role: 'Manager',
  },
  {
    id: 3,
    name: 'Michael Johnson',
    username: 'mjohnson',
    email: 'michael•••@outlook.com',
    avatar_url: 'https://source.unsplash.com/random/100x100',
    role: 'Owner',
  },
  {
    id: 4,
    name: 'Sarah Thompson',
    username: 'sthompson',
    email: 'sarah•••@yahoo.com',
    avatar_url: 'https://source.unsplash.com/random/100x100',
    role: 'Admin',
  },
  {
    id: 5,
    name: 'David Wilson',
    username: 'dwilson',
    email: 'david•••@protonmail.com',
    avatar_url: 'https://source.unsplash.com/random/100x100',
    role: 'Manager',
  },
]

export default function UsersPage() {
  return (
    <div class="flex min-h-screen w-full flex-col rounded-t-xl bg-secondary">
      <div class="flex w-full flex-row items-center justify-between p-4">
        <SearchInput
          class="w-96 bg-background"
          placeholder="Search by name, email, or nickname"
        />
      </div>
      <Table>
        <TableCaption>A list of CRM users</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Username</TableHead>
            <TableHead>Email</TableHead>
            <TableHead class="text-right">Access</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={users}>
            {(user) => (
              <TableRow>
                <TableCell>
                  <div class="flex items-center">
                    <img
                      src={user.avatar_url}
                      alt={user.name}
                      class="mr-2 h-8 w-8 rounded"
                    />
                    <span>{user.name}</span>
                  </div>
                </TableCell>
                <TableCell>{user.username}</TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell class="w-24 text-right">
                  <Chip text={user.role} />
                </TableCell>
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
