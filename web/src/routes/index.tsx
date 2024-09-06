import { createEffect, createSignal, For } from 'solid-js'
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
import { fetchProducts, getLoggedInUser } from '~/lib/api'
import { createQuery } from '@tanstack/solid-query'
import { useNavigate } from '@solidjs/router'
import { Checkbox } from '~/components/checkbox'
import { Separator } from '~/components/separator'

type Product = {
  id: number
  handle: string
  name: string
  description: string
  variants: {
    id: number
    name: string
    available: number
    prices: {
      currency_code: string
      currency_symbol: string
      price: number
      sale_price: number | null
      is_on_sale: boolean
    }[]
  }[]
  image: string
  images: string[]
  is_published: boolean
  created_at: string
  updated_at: string
  deleted_at: string | null
}

export default function IndexPage() {
  const [value, setValue] = createSignal<string | null>(null)

  const navigate = useNavigate()

  const query = createQuery(() => ({
    queryKey: ['products'],
    queryFn: async () => {
      const { data } = await fetchProducts()
      return data as Product[]
    },
  }))

  createEffect(async () => {
    console.log('checking user')
    const { data, error } = await getLoggedInUser()
    if (error) {
      navigate('/auth/login')
    }
  })

  const getSalePrice = (product: Product) => {
    const salePrice = product.variants[0].prices.find(
      (price) => price.sale_price !== null,
    )

    return salePrice ? salePrice.sale_price : null
  }

  const getAvailability = (product: Product) => {
    return product.variants[0].available
  }

  const getProductPrice = (product: Product) => {
    return product.variants[0].prices[0].price
  }

  const normalizeSrc = (src: string) => {
    return src.startsWith('/') ? src.slice(1) : src
  }

  function cdnImage({
    src,
    width,
    quality = 80,
  }: {
    src: string
    width: number
    quality?: number
  }) {
    const params = [`width=${width}`]
    if (quality) {
      params.push(`quality=${quality}`)
    }
    const paramsString = params.join(',')
    return `https://assets.clanplatform.com/cdn-cgi/image/${paramsString}/${normalizeSrc(src)}`
  }

  return (
    <div class="flex min-h-screen w-full flex-col rounded-tl-2xl bg-background">
      <div class="flex w-full flex-row items-center justify-between p-4">
        <SearchInput
          class="w-96 bg-background"
          placeholder="Start typing to search or filter products"
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
            <TableHead class="w-10">
              <Checkbox />
            </TableHead>
            <TableHead>Item</TableHead>
            <TableHead>SKU</TableHead>
            <TableHead>Category</TableHead>
            <TableHead class="text-right">Price</TableHead>
            <TableHead class="text-right">Sale</TableHead>
            <TableHead class="text-right">Availability</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={query.data as Product[]}>
            {(product) => (
              <TableRow>
                <TableCell>
                  <Checkbox />
                </TableCell>
                <TableCell>
                  <div class="flex items-center">
                    <img
                      src={cdnImage({ src: product.image, width: 150 })}
                      alt={product.name}
                      class="mr-2 size-6 rounded object-cover"
                    />
                    <span>{product.name}</span>
                  </div>
                </TableCell>
                <TableCell class="text-muted-foreground">
                  {product.handle}
                </TableCell>
                <TableCell class="text-muted-foreground">Dresses</TableCell>
                <TableCell class="text-right text-muted-foreground">
                  {product.variants[0].prices[0].currency_symbol}
                  {product.variants[0].prices[0].price}
                </TableCell>
                <TableCell class="text-right text-muted-foreground">
                  {product.variants[0].prices[0].currency_symbol}
                  {product.variants[0].prices[0].sale_price}
                </TableCell>
                <TableCell class="flex items-center justify-end gap-4 text-right text-muted-foreground">
                  {product.variants[0].available}
                  <Switch checked={product.is_published} />
                </TableCell>
              </TableRow>
            )}
          </For>
        </TableBody>
      </Table>
    </div>
  )
}
