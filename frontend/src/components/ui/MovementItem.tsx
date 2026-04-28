import { ArrowUpRight, ArrowDownRight, Minus } from "lucide-react"
import type { StockMovement } from "@/features/inventory/types"

interface MovementItemProps {
  movement: StockMovement
}

export function MovementItem({ movement }: MovementItemProps) {
  const getIcon = () => {
    switch (movement.movementType) {
      case "IN":
        return <ArrowUpRight className="h-4 w-4 text-green" />
      case "OUT":
        return <ArrowDownRight className="h-4 w-4 text-red" />
      default:
        return <Minus className="h-4 w-4 text-ink-3" />
    }
  }

  const getColor = () => {
    switch (movement.movementType) {
      case "IN":
        return "text-green"
      case "OUT":
        return "text-red"
      default:
        return "text-ink-3"
    }
  }

  return (
    <div className="flex items-center justify-between py-2 border-b-[0.5px] border-border-default last:border-b-0">
      <div className="flex items-center gap-3">
        <div className="flex h-[30px] w-[30px] items-center justify-center rounded-[7px] bg-surface-2">
          {getIcon()}
        </div>
        <div>
          <p className="text-[12px] font-medium text-ink">{movement.movementType}</p>
          <p className="text-[10px] text-ink-3">{movement.referenceType || "Manual adjustment"}</p>
        </div>
      </div>
      <div className="text-right">
        <p className={`text-[12px] font-medium ${getColor()}`}>
          {movement.movementType === "OUT" ? "-" : "+"}{movement.quantity}
        </p>
        <p className="text-[10px] text-ink-3">
          {new Date(movement.createdAt).toLocaleDateString()}
        </p>
      </div>
    </div>
  )
}