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

const promos = [
  {
    id: 1,
    code: 'BG1SS23',
    value: 10,
    usage_limit: 10,
    type: 'percentage',
    description: '10% off on all electronics',
    is_active: true,
  },
  {
    id: 2,
    code: 'SUMMER22',
    value: 20,
    usage_limit: 100,
    type: 'percentage',
    description: '20% off on selected clothing items',
    is_active: true,
  },
  {
    id: 3,
    code: 'FLATDEAL',
    value: 15,
    usage_limit: null,
    type: 'fixed',
    description: 'Get $15 off on your next purchase',
    is_active: true,
  },
  {
    id: 4,
    code: 'WELCOMENEW',
    value: 10,
    usage_limit: 1,
    type: 'percentage',
    description: 'Welcome offer - 10% off on your first purchase',
    is_active: false,
  },
  {
    id: 5,
    code: 'FREESHIP',
    value: 0,
    usage_limit: 50,
    type: 'fixed',
    description: 'Free shipping on orders above $50',
    is_active: true,
  },
]

function formatValue(value: number, type: string) {
  return type === 'percentage' ? `${value}%` : `$${value}`
}

export default function PromoPage() {
  const [value, setValue] = createSignal<string | null>(null)

  return (
    <div class="flex min-h-screen w-full flex-col rounded-t-xl bg-secondary">
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
          <For each={promos}>
            {(promo) => (
              <TableRow>
                <TableCell>{promo.id}</TableCell>
                <TableCell>{promo.code}</TableCell>
                <TableCell>{formatValue(promo.value, promo.type)}</TableCell>
                <TableCell>{promo.usage_limit}</TableCell>
                <TableCell>{promo.description}</TableCell>
                <TableCell class="float-end">
                  <Switch checked={promo.is_active} />
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
