import type {Component, ComponentProps} from "solid-js"
import {splitProps} from "solid-js"

import {cn} from "~/lib/utils"
import {IconSearch} from "~/components/icons";

const Input: Component<ComponentProps<"input">> = (props) => {
    const [local, others] = splitProps(props, ["type", "class"])
    return (
        <input
            type={local.type}
            class={cn(
                "flex h-8 w-full bg-transparent text-sm file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50",
                local.class
            )}
            {...others}
        />
    )
}

const SearchInput: Component<ComponentProps<"input">> = (props) => {
    const [local, others] = splitProps(props, ["class"])

    return (
        <label class={cn("flex items-center gap-2 bg-input px-2 py-3 h-8 rounded-xl w-full", local.class)}>
            <IconSearch/>
            <Input
                type="search"
                placeholder="Search"
                class="w-full"
                {...props}
            />
        </label>
    )
}

export {Input, SearchInput}
