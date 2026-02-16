package order_behavior

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "CREATED"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusReserved  OrderStatus = "RESERVED"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusFailed    OrderStatus = "FAILED"
)
