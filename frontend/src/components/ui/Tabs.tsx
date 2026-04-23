import * as React from "react"
import { cn } from "@/utils"

interface TabItem {
  id: string
  label: string
  icon?: React.ComponentType<{ className?: string }>
}

interface TabsProps {
  tabs: TabItem[]
  value: string
  onChange: (value: string) => void
  className?: string
}

export const Tabs: React.FC<TabsProps> = ({ tabs, value, onChange, className }) => {
  return (
    <div
      className={cn(
        "inline-flex items-center gap-1 rounded-xl bg-muted p-1",
        className
      )}
      role="tablist"
      aria-label="Tabs"
    >
      {tabs.map((tab) => {
        const Icon = tab.icon
        const isActive = tab.id === value
        return (
          <button
            key={tab.id}
            type="button"
            role="tab"
            aria-selected={isActive}
            onClick={() => onChange(tab.id)}
            className={cn(
              "flex items-center gap-2 rounded-lg px-3 py-1.5 text-sm font-medium transition-all",
              isActive
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            {Icon && <Icon className="h-4 w-4" />}
            {tab.label}
          </button>
        )}
      })}
    </div>
  )
}
