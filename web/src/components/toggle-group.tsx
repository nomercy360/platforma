import {createContext, JSX, splitProps, useContext, ValidComponent} from "solid-js"

import {PolymorphicProps} from "@kobalte/core/polymorphic"
import * as ToggleGroupPrimitive from "@kobalte/core/toggle-group"
import {VariantProps} from "class-variance-authority"

import {cn} from "~/lib/utils"
import {toggleVariants} from "~/components/toggle"

const ToggleGroupContext = createContext<VariantProps<typeof toggleVariants>>({
    size: "default",
    variant: "default"
})

type ToggleGroupRootProps = ToggleGroupPrimitive.ToggleGroupRootProps &
    VariantProps<typeof toggleVariants> & { class?: string | undefined; children?: JSX.Element }

const ToggleGroup = <T extends ValidComponent = "div">(
    props: PolymorphicProps<T, ToggleGroupRootProps>
) => {
    const [local, others] = splitProps(props as ToggleGroupRootProps, [
        "class",
        "children",
        "size",
        "variant"
    ])

    return (
        <ToggleGroupPrimitive.Root
            class={cn("flex items-center justify-center gap-1", local.class)}
            {...others}
        >
            <ToggleGroupContext.Provider
                value={{
                    get size() {
                        return local.size
                    },
                    get variant() {
                        return local.variant
                    }
                }}
            >
                {local.children}
            </ToggleGroupContext.Provider>
        </ToggleGroupPrimitive.Root>
    )
}

type ToggleGroupItemProps = ToggleGroupPrimitive.ToggleGroupItemProps &
    VariantProps<typeof toggleVariants> & { class?: string | undefined }

const ToggleGroupItem = <T extends ValidComponent = "button">(
    props: PolymorphicProps<T, ToggleGroupItemProps>
) => {
    const [local, others] = splitProps(props as ToggleGroupItemProps, ["class", "size", "variant"])
    const context = useContext(ToggleGroupContext)
    return (
        <ToggleGroupPrimitive.Item
            class={cn(
                toggleVariants({
                    size: context.size || local.size,
                    variant: context.variant || local.variant
                }),
                "hover:bg-muted hover:text-muted-foreground data-[pressed]:bg-accent data-[pressed]:text-accent-foreground rounded-xl text-sm text-foreground bg-background text-nowrap flex flex-row justify-center px-2 h-7 gap-1.5 items-center",
                local.class
            )}
            {...others}
        />
    )
}

export {ToggleGroup, ToggleGroupItem}
