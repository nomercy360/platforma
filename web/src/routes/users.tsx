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
import { createQuery } from '@tanstack/solid-query'
import { fetchProducts, listCustomers } from '~/lib/api'
import { Product } from '~/routes/index'

export type Customer = {
  id: number
  name: string
  email: string
  phone: string
  country: string
  address: string
  zip: string
  created_at: string
  updated_at: string
  deleted_at: string
}

export default function UsersPage() {
  const query = createQuery(() => ({
    queryKey: ['customers'],
    queryFn: async () => {
      const { data } = await listCustomers()
      return data as Customer[]
    },
  }))

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
            <TableHead>ID</TableHead>
            <TableHead>Name</TableHead>
            <TableHead>Email</TableHead>
            <TableHead>Country</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={query.data}>
            {(user) => (
              <TableRow>
                <TableCell>{user.id}</TableCell>
                <TableCell>{user.name}</TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>{user.country}</TableCell>
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
