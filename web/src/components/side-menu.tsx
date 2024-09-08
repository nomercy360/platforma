import { IconClose } from '~/components/icons'
import { Show } from 'solid-js'
import { Separator } from '~/components/separator'

type SideMenuProps = {
  isOpen: boolean
  setIsOpen: (value: boolean) => void
  children: any
  title: string
  subtitle: string
}

export default function SideMenu(props: SideMenuProps) {
  return (
    <Show when={props.isOpen}>
      <div class="fixed right-0 top-0 z-50 flex h-full w-[670px] flex-col items-start justify-start bg-secondary">
        <div class="flex w-full flex-row items-start justify-between p-4">
          <div class="space-y-1">
            <p class="text-2xl font-semibold">{props.title}</p>
            <p class="text-sm text-muted-foreground">{props.subtitle}</p>
          </div>
          <button onClick={() => props.setIsOpen(false)}>
            <IconClose class="size-5" />
          </button>
        </div>
        <Separator orientation="horizontal" />
        {props.children}
      </div>
    </Show>
  )
}
