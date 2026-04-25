package constants

const (
	RoleSuperAdmin = "super-admin"
	RoleAdmin      = "admin"

	// Movement Types
	MovementReceipt     = "RECEIPT"
	MovementPutaway     = "PUTAWAY"
	MovementPick        = "PICK"
	MovementPack        = "PACK"
	MovementShip        = "SHIP"
	MovementTransferIn  = "TRANSFER_IN"
	MovementTransferOut = "TRANSFER_OUT"
	MovementAdjustment  = "ADJUSTMENT"
	MovementReservation = "RESERVATION"
	MovementReturn      = "RETURN"
	MovementDamage      = "DAMAGE"
	MovementExpiry      = "EXPIRY"

	// Reference Types
	ReferenceManualAdjustment = "MANUAL_ADJUSTMENT"
	ReferencePurchaseOrder    = "PURCHASE_ORDER"
	ReferenceSalesOrder       = "SALES_ORDER"
	ReferenceTransfer         = "TRANSFER"

	// Status Values
	StatusPending            = "PENDING"
	StatusProcessing         = "PROCESSING"
	StatusShipped            = "SHIPPED"
	StatusCancelled          = "CANCELLED"
	StatusUnallocated        = "UNALLOCATED"
	StatusPartiallyAllocated = "PARTIALLY_ALLOCATED"
	StatusFullyAllocated     = "FULLY_ALLOCATED"

	// Prefixes
	PrefixSalesOrder = "SO"
)
