import { createSignal, For, JSX } from 'solid-js'
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
import { listCustomers, listOrders } from '~/lib/api'
import { Customer } from '~/routes/customers'

type LineItem = {
  id: number
  cart_id: number | null
  order_id: number | null
  variant_id: number
  quantity: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  variant_name: string
  price: number
  sale_price: number | null
  product_name: string
  image_url: string
}

type Order = {
  id: number
  customer_id: number
  cart_id: number
  status: string
  payment_status: string
  shipping_status: string
  total: number
  subtotal: number
  discount_id: number | null
  currency_code: string
  metadata?: Record<string, any>
  created_at: string
  updated_at: string
  deleted_at: string | null
  payment_id: string | null
  payment_provider: string
  customer: Customer
  items: LineItem[]
}
export default function OrdersPage() {
  const [value, setValue] = createSignal<string | null>(null)

  const query = createQuery(() => ({
    queryKey: ['orders'],
    queryFn: async () => {
      const { data } = await listOrders()
      return data as Order[]
    },
  }))

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
          <For each={query.data}>
            {(order) => (
              <TableRow>
                <TableCell class="flex flex-row items-center gap-2">
                  {order.id} <Chip text={order.status} />
                </TableCell>
                <TableCell>
                  {new Date(order.created_at).toLocaleDateString()}
                </TableCell>
                <TableCell>{order.customer.address}</TableCell>
                <TableCell>{order.customer.name}</TableCell>
                <TableCell>{order.customer.email}</TableCell>
                <TableCell>{order.metadata?.comment}</TableCell>
                <TableCell class="flex flex-row justify-end gap-2">
                  <span>
                    {order.items?.length} items for{' '}
                    {order.total.toLocaleString()} {order.currency_code}
                  </span>
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
    case 'created':
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
