package entity

// BudgetStatus represents the status of a user's budget
type BudgetStatus struct {
	CurrentSpend    float64 `json:"current_spend"`
	BudgetLimit     float64 `json:"budget_limit"`
	Remaining       float64 `json:"remaining"`
	SoftCapExceeded bool    `json:"soft_cap_exceeded"`
	HardCapExceeded bool    `json:"hard_cap_exceeded"`
}
