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
    <div class="flex h-16 w-full flex-row items-center p-5">
      <p class="text-2xl font-semibold">P/</p>
      <div class="flex w-full flex-row items-center justify-between">
        <p class="text-2xl font-semibold">Goods</p>
        <SearchInput placeholder="Search workspaces" class="w-96" />
        <Select
          value={value()}
          onChange={setValue}
          options={['Apple', 'Banana', 'Blueberry', 'Grapes', 'Pineapple']}
          placeholder="Select a fruit…"
          itemComponent={(props) => (
            <SelectItem item={props.item}>{props.item.rawValue}</SelectItem>
          )}>
          <SelectTrigger aria-label="Fruit" class="w-[160px]">
            <SelectValue<string>>
              {(state) => state.selectedOption()}
            </SelectValue>
          </SelectTrigger>
          <SelectContent />
        </Select>
      </div>
    </div>
  )
}
