import { useLocation } from '@solidjs/router'
import {
  IconClock,
  IconGroup,
  IconGroupTab,
  IconHome,
  IconListView,
} from '~/components/icons'

export default function Nav() {
  const location = useLocation()
  const active = (path: string) =>
    path == location.pathname ? 'text-foreground' : 'text-muted-foreground'
  return (
    <nav class="bg-transparent">
      <ul class="container flex flex-col items-center px-2.5">
        <li class="flex size-9 items-center justify-center">
          <a href="/" class={active('/')}>
            <IconHome class="size-6" />
          </a>
        </li>
        <li class="flex size-9 items-center justify-center">
          <a href="/users" class={active('/users')}>
            <IconGroup class="size-6" />
          </a>
        </li>
        <li class="flex size-9 items-center justify-center">
          <a href="/orders" class={active('/orders')}>
            <IconListView class="size-6" />
          </a>
        </li>
        <li class="flex size-9 items-center justify-center">
          <a href="/promo" class={active('/promo')}>
            <IconGroupTab class="size-6" />
          </a>
        </li>
      </ul>
    </nav>
  )
}
