package Ports

import "short-link/internal/Core/Domin"

type ShortkeyRepositoryInterface interface {
	Create(id int, uid string, isActive bool) (int, error)
	GetLast() (*Domin.ShortKey, error)
}
