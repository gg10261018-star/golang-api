package points

import "time"

type UserPoints struct {
	Email   string    `xorm:"pk varchar(50) 'email'"`
	Balance int       `xorm:"not null 'balance'"`
	Updated time.Time `xorm:"datetime 'updated'"`
}

type PointTransaction struct {
	ID           int64     `xorm:"pk autoincr 'id'"`
	Email        string    `xorm:"varchar(50) 'email'"`
	RequestID    string    `xorm:"varchar(64) 'request_id'"`
	Amount       int       `xorm:"'amount'"`
	Status       string    `xorm:"varchar(10) 'status'"`
	BalanceAfter int       `xorm:"'balance_after'"`
	Message      string    `xorm:"varchar(255) 'message'"`
	Created      time.Time `xorm:"datetime 'created'"`
}

type ChangeRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Amount    int    `json:"amount" binding:"required,gt=0"`
}

type ChangeResponse struct {
	Balance    int  `json:"balance"`
	Duplicated bool `json:"duplicated"`
}
