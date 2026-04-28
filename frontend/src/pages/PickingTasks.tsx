import * as React from "react"
import {
  Search,
  Play,
  CheckCircle,
  XCircle,
  Clock,
  Package,
  AlertTriangle,
  ArrowRight,
} from "lucide-react"

import {
  usePickingTasks,
  usePickingTask,
} from "@/features/pickingTasks/hooks"
import { Button } from "@/components/ui/Button"
import { Badge } from "@/components/ui/Badge"
import { Skeleton } from "@/components/ui/skeleton"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table"
import { Input } from "@/components/ui/Input"

export default function PickingTasksPage() {
  const [search, setSearch] = React.useState("")
  const [statusFilter, setStatusFilter] = React.useState<string>("")
  const [selectedTaskId, setSelectedTaskId] = React.useState<string>("")

  const { data: tasksData, isLoading } = usePickingTasks({
    status: statusFilter || undefined,
    page: 1,
    limit: 20,
  })

  const { data: selectedTask } = usePickingTask(selectedTaskId)

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "PENDING":
        return <Clock className="h-4 w-4 text-muted-foreground" />
      case "IN_PROGRESS":
        return <AlertTriangle className="h-4 w-4 text-yellow-500" />
      case "COMPLETED":
        return <CheckCircle className="h-4 w-4 text-green-500" />
      case "CANCELLED":
        return <XCircle className="h-4 w-4 text-red-500" />
      default:
        return <Package className="h-4 w-4 text-muted-foreground" />
    }
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "PENDING":
        return <Badge variant="secondary">Pending</Badge>
      case "IN_PROGRESS":
        return <Badge variant="warning">In Progress</Badge>
      case "COMPLETED":
        return <Badge variant="success">Completed</Badge>
      case "CANCELLED":
        return <Badge variant="destructive">Cancelled</Badge>
      default:
        return <Badge>{status}</Badge>
    }
  }

  const getPriorityBadge = (priority: string) => {
    switch (priority) {
      case "URGENT":
        return <Badge variant="destructive">Urgent</Badge>
      case "HIGH":
        return <Badge variant="warning">High</Badge>
      case "MEDIUM":
        return <Badge variant="secondary">Medium</Badge>
      case "LOW":
        return <Badge>Low</Badge>
      default:
        return <Badge>{priority}</Badge>
    }
  }

  const filteredTasks = tasksData?.data?.filter(
    (task) =>
      task.refCode.toLowerCase().includes(search.toLowerCase()) ||
      task.status.toLowerCase().includes(search.toLowerCase())
  ) || []

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">
            Picking Tasks
          </h1>
          <p className="text-muted-foreground mt-1">
            Manage warehouse picking operations.
          </p>
        </div>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between mb-4">
            <CardTitle className="text-lg">All Picking Tasks</CardTitle>
            <div className="flex gap-2">
              <div className="relative w-64">
                <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search by ID..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-9"
                />
              </div>
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="h-10 rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                <option value="">All Status</option>
                <option value="PENDING">Pending</option>
                <option value="IN_PROGRESS">In Progress</option>
                <option value="COMPLETED">Completed</option>
                <option value="CANCELLED">Cancelled</option>
              </select>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-2">
              <Skeleton className="h-12 w-full" />
              <Skeleton className="h-12 w-full" />
              <Skeleton className="h-12 w-full" />
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Task ID</TableHead>
                  <TableHead>Sales Order</TableHead>
                  <TableHead>Priority</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Assigned To</TableHead>
                  <TableHead>Items</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="w-12">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredTasks.map((task) => (
                  <TableRow
                    key={task.id}
                    className={selectedTaskId === task.id.toString() ? "bg-muted" : ""}
                  >
                    <TableCell className="font-medium">
                      <div className="flex items-center gap-2">
                        {getStatusIcon(task.status)}
                        {task.refCode}
                      </div>
                    </TableCell>
                    <TableCell>SO #{task.salesOrderId}</TableCell>
                    <TableCell>{getPriorityBadge(task.priority)}</TableCell>
                    <TableCell>{getStatusBadge(task.status)}</TableCell>
                    <TableCell>{task.assignedUser || "-"}</TableCell>
                    <TableCell>
                      {task.items?.length || 0} item
                      {task.items && task.items.length !== 1 ? "s" : ""}
                    </TableCell>
                    <TableCell>
                      {new Date(task.createdAt).toLocaleDateString()}
                    </TableCell>
                    <TableCell>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setSelectedTaskId(task.id.toString())}
                      >
                        View
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {selectedTask && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Task Details - {selectedTask.refCode}</CardTitle>
              {getStatusBadge(selectedTask.status)}
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="grid grid-cols-3 gap-4 text-sm">
                <div>
                  <span className="font-medium">Sales Order:</span> SO #{selectedTask.salesOrderId}
                </div>
                <div>
                  <span className="font-medium">Warehouse:</span> Warehouse #{selectedTask.warehouseId}
                </div>
                <div>
                  <span className="font-medium">Priority:</span> {getPriorityBadge(selectedTask.priority)}
                </div>
                <div>
                  <span className="font-medium">Assigned To:</span> {selectedTask.assignedUser || "-"}
                </div>
                <div>
                  <span className="font-medium">Started At:</span>{" "}
                  {selectedTask.startedAt
                    ? new Date(selectedTask.startedAt).toLocaleString()
                    : "-"}
                </div>
                <div>
                  <span className="font-medium">Completed At:</span>{" "}
                  {selectedTask.completedAt
                    ? new Date(selectedTask.completedAt).toLocaleString()
                    : "-"}
                </div>
              </div>

              {selectedTask.items && selectedTask.items.length > 0 && (
                <div className="mt-4">
                  <h4 className="font-medium mb-2">Items to Pick:</h4>
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Product</TableHead>
                        <TableHead>Location</TableHead>
                        <TableHead>Required</TableHead>
                        <TableHead>Picked</TableHead>
                        <TableHead>Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {selectedTask.items.map((item) => (
                        <TableRow key={item.id}>
                          <TableCell>
                            {item.productName || `Product #${item.productId}`}
                          </TableCell>
                          <TableCell>{item.locationCode || "-"}</TableCell>
                          <TableCell>{item.quantityRequired}</TableCell>
                          <TableCell>{item.quantityPicked}</TableCell>
                          <TableCell>
                            {item.status === "COMPLETED" ? (
                              <Badge variant="success">Picked</Badge>
                            ) : (
                              <Badge variant="secondary">Pending</Badge>
                            )}
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              )}

              {selectedTask.notes && (
                <div className="mt-4">
                  <span className="font-medium">Notes:</span> {selectedTask.notes}
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}