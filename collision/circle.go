package collision

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Circle struct {
	center mgl32.Vec2
	radius float32
}

func NewCircle(center mgl32.Vec2, radius float32) *Circle {
	return &Circle{
		center: center,
		radius: radius,
	}
}

// IsCoordinatePointInCircle 坐标点是否在坐标系的圆中
// 坐标点a和圆心b的向量长度是否小于等于圆的半径值
func IsCoordinatePointInCircle(point mgl32.Vec2, circle Circle) bool {
	return point.Sub(circle.center).Len() <= circle.radius
}
