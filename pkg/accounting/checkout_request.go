package accounting

type CheckoutRequest struct {
	UserID   string `json:"user_id"`
	DayTotal int    `json:"day_total"`
}
