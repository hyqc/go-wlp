package collision

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRectangle(t *testing.T) {

}

func TestCalculateRectangleCenter(t *testing.T) {
	leftBottom := mgl32.Vec2{0, 0}
	rightTop := mgl32.Vec2{3, 4}
	center := CalculateRectangleCenter(leftBottom, rightTop)
	assert.Equal(t, mgl32.Vec2{1.5, 2}, center, "计算矩形中心点坐标错误")

	leftBottom = mgl32.Vec2{0, 0}
	rightTop = mgl32.Vec2{-3, -4}
	center = CalculateRectangleCenter(leftBottom, rightTop)
	assert.Equal(t, mgl32.Vec2{-1.5, -2}, center, "计算矩形中心点坐标错误")
}
