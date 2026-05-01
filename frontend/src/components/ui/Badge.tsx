import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/utils"

const badgeVariants = cva(
  "inline-flex items-center rounded-md border px-2 py-0.5 text-[11px] font-medium leading-4 transition-colors focus:outline-none",
  {
    variants: {
      variant: {
        default:
          "border-transparent bg-primary/20 text-primary",
        secondary:
          "border-transparent bg-secondary text-foreground",
        destructive:
          "border-transparent bg-destructive/20 text-destructive",
        outline: "border-border bg-transparent text-muted-foreground",
        success:
          "border-transparent bg-[#22C55E]/20 text-[#22C55E]",
        warning:
          "border-transparent bg-[#F59E0B]/20 text-[#F59E0B]",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  )
}

export { Badge, badgeVariants }
