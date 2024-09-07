import { useLocation } from '@solidjs/router'
import {
  IconGroup,
  IconListView,
  IconPercent,
  IconPerson,
  IconQRCode,
} from '~/components/icons'
import { useAuth } from '~/lib/auth-context'
import { createEffect } from 'solid-js'

export default function Navbar() {
  const location = useLocation()
  const active = (path: string) =>
    path == location.pathname
      ? 'text-foreground size-9 flex items-center justify-center'
      : 'text-muted-foreground size-9 flex items-center justify-center'

  const { user } = useAuth()

  createEffect(() => {
    if (user) {
      console.log('User:', user)
    }
  })

  return (
    <nav class="fixed left-0 top-0 z-50 flex h-screen w-14 flex-col items-center justify-start bg-secondary">
      <p class="text-2xl pt-4 font-semibold">P/</p>
      <div class="flex h-full flex-col items-center justify-between pb-6 pt-4">
        <ul class="container flex flex-col items-center space-y-2 px-2.5">
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
            <a href="/customers" class={active('/customers')}>
              <IconGroup class="size-6" />
            </a>
          </li>
          <li class="flex items-center justify-center">
            <a href="/discount" class={active('/discount')}>
              <IconPercent class="size-3" />
            </a>
          </li>
          <li class="flex items-center justify-center">
            <a href="/users" class={active('/users')}>
              <IconPerson class="size-6" />
            </a>
          </li>
        </ul>
        <a href="/profile" class="flex items-center justify-center">
          <img
            src={user()?.avatar_url}
            alt="avatar"
            class="size-8 rounded-full"
          />
        </a>
      </div>
    </nav>
  )
}
