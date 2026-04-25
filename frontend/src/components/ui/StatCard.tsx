import * as React from "react"
import { ArrowUpRight, ArrowDownRight } from "lucide-react"
import { cn } from "@/utils"

type LucideIcon = React.ComponentType<React.ComponentProps<"svg">>

interface StatCardProps {
  title: string
  value: string
  description: string
  icon: LucideIcon
  loading?: boolean
  trend?: "up" | "down"
}

interface StatCardProps {
  title: string
  value: string
  description: string
  icon: LucideIcon
  loading?: boolean
  trend?: "up" | "down"
}

export function StatCard({ title, value, description, icon: Icon, loading, trend }: StatCardProps) {
  return (
    <div className="bg-white border-[0.5px] border-border-default rounded-[10px] p-[14px_16px]">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <p className="text-[12px] text-ink-3">{title}</p>
          {loading ? (
            <div className="h-7 w-20 bg-surface-2 rounded-[5px] animate-pulse mt-1" />
          ) : (
            <p className="text-[20px] font-medium text-ink mt-0.5">{value}</p>
          )}
          <p className="text-[10px] text-ink-3 mt-0.5">{description}</p>
        </div>
        <div className="flex h-[36px] w-[36px] items-center justify-center rounded-[8px] bg-surface-2">
          <Icon className="h-4 w-4 text-ink-2" />
        </div>
      </div>
      {trend && (
        <div className={cn("flex items-center mt-2 text-[10px]", trend === "up" ? "text-green" : "text-red")}>
          {trend === "up" ? <ArrowUpRight className="h-3 w-3 mr-0.5" /> : <ArrowDownRight className="h-3 w-3 mr-0.5" />}
          <span>{trend === "up" ? "Increasing" : "Decreasing"}</span>
        </div>
      )}
    </div>
  )
}