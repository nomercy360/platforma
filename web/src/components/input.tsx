import type { Component, ComponentProps } from 'solid-js'
import { splitProps } from 'solid-js'

import { cn } from '~/lib/utils'
import { IconSearch } from '~/components/icons'

const Input: Component<ComponentProps<'input'>> = (props) => {
  const [local, others] = splitProps(props, ['type', 'class'])
  return (
    <input
      type={local.type}
      class={cn(
        'flex h-8 w-full bg-transparent font-normal file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50',
        local.class,
      )}
      {...others}
    />
  )
}

const SearchInput: Component<ComponentProps<'input'>> = (props) => {
  const [local, others] = splitProps(props, ['class'])

  return (
    <label
      class={cn(
        'flex h-8 w-full items-center gap-2 rounded-xl bg-input px-2 py-3',
        local.class,
      )}>
      <IconSearch class="size-4 text-muted-foreground" />
      <Input type="search" placeholder="Search" class="w-full" {...props} />
    </label>
  )
}

export { Input, SearchInput }
