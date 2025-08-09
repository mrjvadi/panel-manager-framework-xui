package ext

import "context"

// XUIDeleter: حذف اینباند با ID
type XUIDeleter interface {
	DeleteInboundByID(ctx context.Context, id int) error
}
