import { SearchInput } from '~/components/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/select'
import { createSignal } from 'solid-js'

export default function Controls() {
  const [value, setValue] = createSignal<string | undefined>(undefined)
  return (
    <div class="flex h-16 w-full flex-row items-center bg-secondary p-5">
      <p class="text-2xl font-semibold">P/</p>
      <div class="flex w-full flex-row items-center justify-between">
        <p class="text-2xl font-semibold">Goods</p>
      </div>
    </div>
  )
}
