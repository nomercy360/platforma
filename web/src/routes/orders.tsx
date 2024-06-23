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

const orders = [
  {
    id: 3234,
    status: 'New',
    date: 'March 24',
    delivery_to: 'Moscow, Russia',
    customer: 'Nikita Axel',
    email: 'nikita---@proton.me',
    notes: 'Call before delivery',
    items: [
      {
        name: 'Kitik Sweater',
        price: 840,
        payment_status: 'yellow',
      },
    ],
  },
  {
    id: 3234,
    status: 'Overdue',
    date: 'March 20',
    delivery_to: 'Ekaterinburg, RU',
    customer: 'Nikita Axel',
    email: 'nikita---@proton.me',
    notes: 'Call before delivery',
    items: [
      {
        name: 'Kitik Sweater',
        price: 1234.46,
        payment_status: 'red',
      },
    ],
  },
  {
    id: 3234,
    status: 'Refund',
    date: 'March 19',
    delivery_to: 'Tbilisi, GE',
    customer: 'Nikita Axel',
    email: 'nikita---@proton.me',
    notes: 'Call before delivery',
    items: [
      {
        name: 'Kitik Sweater',
        price: 1234.46,
        payment_status: 'red',
      },
    ],
  },
  {
    id: 3234,
    status: 'Delivering',
    date: 'March 19',
    delivery_to: 'Moscow, RU',
    customer: 'Nikita Axel',
    email: 'nikita---@proton.me',
    notes: 'Call before delivery',
    items: [
      {
        name: 'Kitik Sweater',
        price: 1234.46,
        payment_status: 'yellow',
      },
    ],
  },
  {
    id: 3234,
    status: 'Delivering',
    date: 'March 19',
    delivery_to: 'Bucharest, RO',
    customer: 'Nikita Axel',
    email: 'nikita---@proton.me',
    notes: 'Call before delivery',
    items: [
      {
        name: 'Kitik Sweater',
        price: 1234.46,
        payment_status: 'yellow',
      },
      {
        name: 'Kitik Sweater',
        price: 1234.46,
        payment_status: 'green',
      },
    ],
  },
  {
    id: 3234,
    status: 'Completed',
    date: 'March 19',
    delivery_to: 'Bucharest, RO',
    customer: 'Nikita Axel',
    email: 'nikita---@proton.me',
    notes: 'Call before delivery',
    items: [
      {
        name: 'Kitik Sweater',
        price: 1234.46,
        payment_status: 'green',
      },
      {
        name: 'Kitik Sweater',
        price: 1234.46,
        payment_status: 'green',
      },
    ],
  },
]

export default function OrdersPage() {
  const [value, setValue] = createSignal<string | null>(null)

  return (
    <div class="flex min-h-screen w-full flex-col rounded-t-xl bg-secondary">
      <div class="flex w-full flex-row items-center justify-between p-4">
        <SearchInput
          class="w-96 bg-background"
          placeholder="Search by number, name, location, email, etc."
        />
        <Select
          value={value()}
          onChange={setValue}
          options={[
            'Id',
            'Item',
            'SKU',
            'Category',
            'Price',
            'Stock',
            'Availability',
          ]}
          placeholder={'Sort by...'}
          itemComponent={(props) => (
            <SelectItem item={props.item}>{props.item.rawValue}</SelectItem>
          )}>
          <SelectTrigger aria-label="Fruit" class="w-[160px] bg-background">
            <SelectValue<string>>
              {(state) => 'Sort by ' + state.selectedOption()}
            </SelectValue>
          </SelectTrigger>
          <SelectContent />
        </Select>
      </div>
      <div class="flex w-full flex-row items-start p-4">
        <ToggleGroup>
          <ToggleGroupItem value="all">
            All <span class="text-muted-foreground">64</span>
          </ToggleGroupItem>
          <ToggleGroupItem value="gadgets">
            Gadgets <span class="text-muted-foreground">24</span>
          </ToggleGroupItem>
          <ToggleGroupItem value="clothing">
            Clothing <span class="text-muted-foreground">16</span>
          </ToggleGroupItem>
          <ToggleGroupItem value="electronics">
            Electronics <span class="text-muted-foreground">12</span>
          </ToggleGroupItem>
          <ToggleGroupItem value="appliances">
            Appliances <span class="text-muted-foreground">8</span>
          </ToggleGroupItem>
        </ToggleGroup>
      </div>
      <Table>
        <TableCaption>A list of your recent products.</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead class="w-40">#</TableHead>
            <TableHead>Date</TableHead>
            <TableHead>Delivery to</TableHead>
            <TableHead>Customer</TableHead>
            <TableHead>E-mail</TableHead>
            <TableHead>Notes</TableHead>
            <TableHead class="text-right">Items and payment</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={orders}>
            {(order) => (
              <TableRow>
                <TableCell class="flex flex-row items-center gap-2">
                  {order.id} <Chip text={order.status} />
                </TableCell>
                <TableCell>{order.date}</TableCell>
                <TableCell>{order.delivery_to}</TableCell>
                <TableCell>{order.customer}</TableCell>
                <TableCell>{order.email}</TableCell>
                <TableCell>{order.notes}</TableCell>
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
    case 'New':
      color = 'bg-green-200 text-green-800'
      break
    case 'Overdue':
      color = 'bg-red-200 text-red-800'
      break
    case 'Refund':
      color = 'bg-purple-200 text-purple-800'
      break
    case 'Delivering':
      color = 'bg-yellow-200 text-yellow-800'
      break
    case 'Completed':
      color = 'bg-neutral-200 text-neutral-800'
      break
    default:
      break
  }

  return (
    <span
      class={`flex h-6 items-center justify-center rounded-full px-2 py-1 text-sm ${color}`}>
      {text}
    </span>
  )
}
