package validation

import (
	"L0/internal/models"
	"errors"
)

func ValidateOrder(o *models.Order) error {
	if o.OrderUID == "" {
		return errors.New("Order_uid required")
	}
	if o.Payment.Transaction == "" {
		return errors.New("Transaction required")
	}
	if len(o.Items) == 0 {
		return errors.New("Items empty")
	}
	return nil
}
