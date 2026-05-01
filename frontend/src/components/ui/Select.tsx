import * as React from "react"
import { cn } from "@/utils"

const Select = React.forwardRef<
  HTMLSelectElement,
  React.SelectHTMLAttributes<HTMLSelectElement>
>(({ className, children, ...props }, ref) => {
  return (
    <select
      ref={ref}
      className={cn(
        "flex h-9 w-full rounded-[6px] border border-input bg-[#111827] px-3 py-1 text-sm text-foreground transition-colors focus-visible:border-primary focus-visible:outline-none focus-visible:ring-0 disabled:cursor-not-allowed disabled:opacity-50",
        className
      )}
      {...props}
    >
      {children}
    </select>
  )
})
Select.displayName = "Select"

export { Select }
