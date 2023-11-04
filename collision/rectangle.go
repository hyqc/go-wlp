package collision

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Rectangle struct {
	center        mgl32.Vec2
	leftBottom    mgl32.Vec2 // 左下角坐标(x1,y1)
	leftTop       mgl32.Vec2 // 左上角坐标(x1,y2)
	rightTop      mgl32.Vec2 // 右上角坐标(x2,y2)
	rightBottom   mgl32.Vec2 // 右下角坐标(x2,y1)
	rotationAngle float32    // 矩形在坐标系中的旋转角度，顺时针方向，如45°
}

func NewRectangle(leftBottom, rightTop mgl32.Vec2, rotationAngle float32) *Rectangle {
	x1, y1 := leftBottom.Elem()
	x2, y2 := rightTop.Elem()
	rect := &Rectangle{
		center:        mgl32.Vec2{},
		leftBottom:    leftBottom,
		leftTop:       mgl32.Vec2{x1, y2},
		rightTop:      rightTop,
		rightBottom:   mgl32.Vec2{x2, y1},
		rotationAngle: rotationAngle,
	}
	rect.leftBottom = CalculatePointRotationValue(rect.leftBottom, rotationAngle)
	rect.leftTop = CalculatePointRotationValue(rect.leftTop, rotationAngle)
	rect.rightTop = CalculatePointRotationValue(rect.rightTop, rotationAngle)
	rect.rightBottom = CalculatePointRotationValue(rect.rightBottom, rotationAngle)
	rect.center = CalculateRectangleCenter(rect.leftBottom, rect.rightTop)
	return rect
}

// CalculateRectangleCenter 计算矩形的中心坐标点
func CalculateRectangleCenter(leftBottom, rightTop mgl32.Vec2) mgl32.Vec2 {
	// a(x1,y1)
	// b(x2,y2)
	// ab = (x2-x1,y2-y1)
	// ab/2 = (x2-x1,y2-y1)/2 = ((x2-x1)/2,(y2-y1)/2)
	// center = a + ab/2
	return leftBottom.Add(rightTop.Sub(leftBottom).Mul(1.0 / 2.0))
}
