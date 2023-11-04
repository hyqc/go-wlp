package collision

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsCoordinatePointInCircle(t *testing.T) {
	circle := NewCircle(mgl32.Vec2{0, 0}, 2)
	point := mgl32.Vec2{0, 0}
	assert.True(t, IsCoordinatePointInCircle(point, *circle), "坐标点在圆心")

	point = mgl32.Vec2{4, 0}
	assert.False(t, IsCoordinatePointInCircle(point, *circle), "错误")

	point = mgl32.Vec2{2.4, 5}
	assert.False(t, IsCoordinatePointInCircle(point, *circle), "错误")

	point = mgl32.Vec2{1.4, 5}
	assert.False(t, IsCoordinatePointInCircle(point, *circle), "错误")

	point = mgl32.Vec2{1.4, 0.5}
	assert.True(t, IsCoordinatePointInCircle(point, *circle), "错误")
}
