package collision

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

// IsRectangleIntersectsCircle 矩形和圆是否相交
func IsRectangleIntersectsCircle(rectangle Rectangle, circle Circle) bool {
	if rectangle.center.Sub(circle.center).Len() <= circle.radius {
		// 如果矩形的中心点到圆心的向量距离小于圆的半径，那么矩形一定与圆形相交
		return true
	}

	if rectangle.leftBottom.Sub(circle.center).Len() <= circle.radius {
		// 如果左下角到圆心距离小于半径，则相交
		return true
	}

	if rectangle.rightTop.Sub(circle.center).Len() <= circle.radius {
		// 如果右上角到圆心距离小于半径，则相交
		return true
	}

	if rectangle.leftBottom.Sub(rectangle.leftTop).Mul(1.0/2.0).Sub(circle.center).Len() <= circle.radius {
		// 如果矩形左右边中点到圆心的距离小于半径，则相交
		return true
	}

	if rectangle.rightBottom.Sub(rectangle.rightTop).Mul(1.0/2.0).Sub(circle.center).Len() <= circle.radius {
		// 如果矩形左右边中点到圆心的距离小于半径，则相交
		return true
	}

	if rectangle.leftBottom.Sub(rectangle.rightBottom).Mul(1.0/2.0).Sub(circle.center).Len() <= circle.radius {
		// 如果矩形上下边中点到圆心的距离小于半径，则相交
		return true
	}

	if rectangle.leftTop.Sub(rectangle.rightTop).Mul(1.0/2.0).Sub(circle.center).Len() <= circle.radius {
		// 如果矩形上下边中点到圆心的距离小于半径，则相交
		return true
	}

	return false
}

// CalculatePointRotationValue 计算一个坐标点旋转指定角度后的新坐标
func CalculatePointRotationValue(point mgl32.Vec2, angle float32) mgl32.Vec2 {
	x, y := point.Elem()
	switch angle {
	case 0:
		return mgl32.Vec2{x, y}
	case 90:
		return mgl32.Vec2{y, -x}
	case 180:
		return mgl32.Vec2{-x, -y}
	case 270:
		return mgl32.Vec2{-y, x}
	case 360:
		return point
	}

	an := float64(math.Pi * angle / 180)

	cos := float32(math.Cos(an))
	sin := float32(math.Sin(an))

	return mgl32.Vec2{
		x*cos + y*sin,
		y*cos - x*sin,
	}
}
