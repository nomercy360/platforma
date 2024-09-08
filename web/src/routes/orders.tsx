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
import { ToggleGroup, ToggleGroupItem } from '~/components/toggle-group'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/select'
import { createQuery } from '@tanstack/solid-query'
import { listOrders } from '~/lib/api'
import { Customer } from '~/routes/customers'
import { cn } from '~/lib/utils'
import SideMenu from '~/components/side-menu'
import { createStore } from 'solid-js/store'

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

  function formatDate(date: Date) {
    const optionsDate = { month: 'long', day: 'numeric' }
    const optionsYear = { year: 'numeric' }

    const dateString = date.toLocaleDateString('en-US', optionsDate)
    const yearString = date.toLocaleDateString('en-US', optionsYear)

    return { dateString, yearString }
  }

  const [editOrder, setEditOrder] = createSignal<Order | null>(null)

  const [sideMenuIsOpen, setSideMenuIsOpen] = createSignal(false)

  async function onFormSubmit() {
    console.log('onFormSubmit')
  }

  function openSideMenu(order: Order) {
    setEditOrder(order)
    setSideMenuIsOpen(true)
  }

  function getDateString(date?: string) {
    if (!date) return ''
    const dateObj = new Date(date)

    return `${formatDate(dateObj).dateString} ${formatDate(dateObj).yearString}`
  }

  return (
    <div class="flex min-h-screen w-full flex-col rounded-tl-2xl bg-background">
      <SideMenu
        title={`${editOrder()?.id} for ${editOrder()?.customer.name}`}
        subtitle={`${getDateString(editOrder()?.created_at)} / Payed via ${editOrder()?.payment_provider}`}
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
              placeholder="John Doe"
            />
          </label>
          <label class="text-sm font-semibold" for="email">
            <input
              class="h-11 w-full rounded-lg border border-neutral-200 bg-background p-2"
              id="email"
              type="email"
              placeholder="dummy@example.com"
            />
          </label>
          <label class="text-sm font-semibold" for="password">
            <input
              class="mt-1 h-11 w-full rounded-lg border border-neutral-200 bg-background p-2"
              id="password"
              type="password"
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
          placeholder="Start typing to search or filter products"
        />
      </div>
      <div class="flex w-full flex-row items-start p-4">
        <ToggleGroup>
          <ToggleGroupItem value="new">
            New <span class="text-muted-foreground">64</span>
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
            <TableHead>ID</TableHead>
            <TableHead>Date</TableHead>
            <TableHead>Customer</TableHead>
            <TableHead>E-mail</TableHead>
            <TableHead class="text-right">Payment</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={query.data}>
            {(order) => (
              <TableRow
                class="cursor-pointer"
                onClick={() => openSideMenu(order)}>
                <TableCell>
                  <div class="flex flex-col space-y-0.5">
                    {order.id}
                    <OrderStatus text={order.status}></OrderStatus>
                  </div>
                </TableCell>
                <TableCell>
                  <div class="flex flex-col space-y-0.5">
                    {formatDate(new Date(order.created_at)).dateString}
                    <span class="text-sm text-muted-foreground">
                      {formatDate(new Date(order.created_at)).yearString}
                    </span>
                  </div>
                </TableCell>
                <TableCell>
                  <div class="flex flex-col space-y-0.5">
                    {order.customer.name}
                    <span class="text-sm text-muted-foreground">
                      {order.customer.address}
                    </span>
                  </div>
                </TableCell>
                <TableCell>
                  <div class="flex flex-col space-y-0.5">
                    {order.customer.email}
                    <span class="text-sm text-muted-foreground">
                      {order.customer.phone}
                    </span>
                  </div>
                </TableCell>
                <TableCell class="flex flex-row justify-end">
                  <div class="flex flex-col space-y-0.5">
                    <div class="flex flex-row items-center space-x-2">
                      <span>
                        {order.total.toLocaleString()} {order.currency_code}
                      </span>
                      <span
                        class={cn('size-2 rounded-full', {
                          'bg-green-500': order.payment_status === 'paid',
                          'bg-red-500': order.payment_status !== 'paid',
                        })}></span>
                    </div>
                    <span class="text-sm text-muted-foreground">
                      {order.payment_provider}
                    </span>
                  </div>
                </TableCell>
              </TableRow>
            )}
          </For>
        </TableBody>
      </Table>
    </div>
  )
}

function OrderStatus({ text }: { text: string }) {
  let color

  switch (text) {
    case 'created':
      color = 'text-green-500'
      break
    case 'Overdue':
      color = 'text-red-500'
      break
    case 'Refund':
      color = 'text-purple-500'
      break
    case 'Delivering':
      color = 'text-yellow-500'
      break
    case 'Completed':
      color = 'text-neutral-500'
      break
    default:
      break
  }

  return <span class={`text-sm capitalize ${color}`}>{text}</span>
}
