package order

import (
	"app/myproj/pkg/workflow"
	"time"
)

func NewOrderWorkflow(orderService *OrderService) *workflow.WorkflowDefinition {
	return &workflow.WorkflowDefinition{
		ID:           "order_workflow",
		InitialState: "order_input",
		Timeout:      30 * time.Minute,
		States: map[workflow.State][]workflow.StateTransition{
			"order_input": {
				{
					Event:   "submit_order",
					ToState: "payment_pending",
					Action: func(data map[string]interface{}) error {
						userID := data["user_id"].(uint)
						productID := data["product_id"].(uint)
						quantity := data["quantity"].(int)

						order, err := orderService.CreateOrder(userID, productID, quantity)
						if err != nil {
							return err
						}

						data["order_id"] = order.ID
						return nil
					},
				},
			},
			"payment_pending": {
				{
					Event:   "payment_confirmed",
					ToState: "processing",
					Action: func(data map[string]interface{}) error {
						orderID := data["order_id"].(uint)
						return orderService.db.Model(&Order{}).
							Where("id = ?", orderID).
							Update("status", "processing").Error
					},
				},
			},
			"processing": {
				{
					Event:   "order_fulfilled",
					ToState: "completed",
					Action: func(data map[string]interface{}) error {
						orderID := data["order_id"].(uint)
						return orderService.db.Model(&Order{}).
							Where("id = ?", orderID).
							Update("status", "completed").Error
					},
				},
			},
		},
	}
}
