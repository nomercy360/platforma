import { createSignal } from 'solid-js'

export function useTableSelection<T extends { id: string | number }>() {
  const [selected, setSelected] = createSignal<Array<T['id']>>([])

  const toggleSelection = (id: T['id']) => {
    setSelected((prev) =>
      prev.includes(id)
        ? prev.filter((itemId) => itemId !== id)
        : [...prev, id],
    )
  }

  const toggleSelectAll = (items: T[]) => {
    setSelected((prev) =>
      prev.length === items.length ? [] : items.map((item) => item.id),
    )
  }

  return {
    selected,
    toggleSelection,
    toggleSelectAll,
  }
}
