package collision

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsRectangleIntersectsCircle(t *testing.T) {
	circle := NewCircle(mgl32.Vec2{0, 0}, 2)
	rectangle := NewRectangle(mgl32.Vec2{1, 2}, mgl32.Vec2{2, 4}, 0)
	assert.True(t, IsRectangleIntersectsCircle(*rectangle, *circle), "错误")
}

func TestCalculatePointRotationValue(t *testing.T) {
	p1 := mgl32.Vec2{1, 1}
	p2 := CalculatePointRotationValue(p1, 90)
	assert.Equal(t, mgl32.Vec2{1, -1}, p2, "计算错误")
}
