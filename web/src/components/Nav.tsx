import { useLocation } from '@solidjs/router'
import {
  IconClock,
  IconGroup,
  IconGroupTab,
  IconHome,
  IconListView,
  IconPercent,
  IconQRCode,
} from '~/components/icons'

export default function Nav() {
  const location = useLocation()
  const active = (path: string) =>
    path == location.pathname
      ? 'text-foreground size-9 flex items-center justify-center'
      : 'text-muted-foreground size-9 flex items-center justify-center'

  return (
    <nav class="bg-transparent">
      <ul class="container flex flex-col items-center px-2.5">
        <li class="flex items-center justify-center">
          <a href="/orders" class={active('/orders')}>
            <IconListView class="size-6" />
          </a>
        </li>
        <li class="flex items-center justify-center">
          <a href="/" class={active('/')}>
            <IconQRCode class="size-6" />
          </a>
        </li>
        <li class="flex items-center justify-center">
          <a href="/users" class={active('/customers')}>
            <IconGroup class="size-6" />
          </a>
        </li>
        <li class="flex items-center justify-center">
          <a href="/discount" class={active('/discount')}>
            <IconPercent class="size-3" />
          </a>
        </li>
      </ul>
    </nav>
  )
}
