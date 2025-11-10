package cache

import "auth-management/pkg/enum"

type RefreshData struct {
	UserId string    `json:"user_id"`
	Role   enum.ROLE `json:"role"`
}
