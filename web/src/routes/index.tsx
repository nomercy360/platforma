import { createResource, createSignal, For } from 'solid-js'
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
import { fetchProducts } from '~/lib/api'

export type Product = {
  id: number
  handle: string
  name: string
  description: string
  variants: {
    id: number
    name: string
    available: number
  }[]
  image: string
  images: string[]
  currency_code: string
  currency_symbol: string
  price: number
  deleted_at: string | null
}

export default function IndexPage() {
  const [value, setValue] = createSignal<string | null>(null)

  const [products] = createResource<Product[]>(() => fetchProducts(), {
    initialValue: [],
  })

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
            <TableHead class="w-[40px]">#</TableHead>
            <TableHead>Item</TableHead>
            <TableHead>SKU</TableHead>
            <TableHead>Category</TableHead>
            <TableHead class="text-right">Price</TableHead>
            <TableHead class="text-right">On sale</TableHead>
            <TableHead class="text-right">Stock</TableHead>
            <TableHead class="text-right">Availability</TableHead>
            <TableHead class="text-right">View</TableHead>
            <TableHead class="text-right">Delete</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={products()}>
            {(product) => (
              <TableRow>
                <TableCell>{product.id}</TableCell>
                <TableCell>
                  <div class="flex items-center">
                    <img
                      src={product.image}
                      alt={product.name}
                      class="mr-2 h-8 w-8 rounded"
                    />
                    <span>{product.name}</span>
                  </div>
                </TableCell>
                <TableCell>{product.handle}</TableCell>
                <TableCell>Dresses</TableCell>
                <TableCell class="text-right">{product.price}</TableCell>
                <TableCell class="text-right">{product.sale_price}</TableCell>
                <TableCell class="text-right">{product.stock}</TableCell>
                <TableCell class="text-right">{product.availability}</TableCell>
                <TableCell class="float-end">
                  <Switch checked={product.deleted_at === null} />
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
