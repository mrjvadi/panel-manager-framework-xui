package xui

// اگر Remark/Port ندی، درایور مقادیر امن تولید می‌کند.
// اگر Client ندی، کلاینت ساخته نمی‌شود.
type CloneInboundOptions struct {
	Remark *string
	Port   *int
	Client *ClientCreate
}
