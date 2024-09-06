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
import { listDiscounts, listOrders } from '~/lib/api'

type Discount = {
  id: number
  code: string
  value: number
  type: string
  usage_limit: number
  usage_count: number
  starts_at: string
  ends_at: string
  created_at: string
  updated_at: string
  is_active: boolean
  description: string
}

function formatValue(value: number, type: string) {
  return type === 'percentage' ? `${value}%` : `$${value}`
}

export default function DiscountPage() {
  const [value, setValue] = createSignal<string | null>(null)

  const query = createQuery(() => ({
    queryKey: ['discounts'],
    queryFn: async () => {
      const { data } = await listDiscounts()
      return data as Discount[]
    },
  }))

  return (
    <div class="flex min-h-screen w-full flex-col rounded-tl-2xl bg-background">
      <div class="flex w-full flex-row items-center justify-between p-4">
        <SearchInput
          class="w-96 bg-background"
          placeholder="Search by name, SKU, price, etc."
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
            Created <span class="text-muted-foreground">24</span>
          </ToggleGroupItem>
          <ToggleGroupItem value="clothing">
            Generated <span class="text-muted-foreground">16</span>
          </ToggleGroupItem>
        </ToggleGroup>
      </div>
      <Table>
        <TableCaption>A list of your recent invoices.</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[40px]">#</TableHead>
            <TableHead>SKU</TableHead>
            <TableHead>Amount</TableHead>
            <TableHead>Usage</TableHead>
            <TableHead>Description</TableHead>
            <TableHead class="text-right">Active</TableHead>
            <TableHead class="text-right">Delete</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={query.data}>
            {(discount) => (
              <TableRow>
                <TableCell>{discount.id}</TableCell>
                <TableCell>{discount.code}</TableCell>
                <TableCell>
                  {formatValue(discount.value, discount.type)}
                </TableCell>
                <TableCell>{discount.usage_limit}</TableCell>
                <TableCell>{discount.description}</TableCell>
                <TableCell class="float-end">
                  <Switch
                    disabled={true}
                    onChange={() => {}}
                    checked={discount.is_active}
                  />
                </TableCell>
                <TableCell class="text-right">
                  <button class="text-red-600 hover:underline">Delete</button>
                </TableCell>
              </TableRow>
            )}
          </For>
        </TableBody>
      </Table>
    </div>
  )
}
