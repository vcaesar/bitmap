package bitmap

import (
	"fmt"
	"testing"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/tt"
)

func TestBitmap(t *testing.T) {
	bit := robotgo.CaptureScreen()
	defer robotgo.FreeBitmap(bit)
	tt.NotNil(t, bit)
	e := Save(bit, "robot_test.png")
	tt.Nil(t, e)

	bit0 := robotgo.CaptureScreen(10, 10, 20, 20)
	defer robotgo.FreeBitmap(bit0)
	x, y := Find(bit0)
	fmt.Println("Find bitmap: ", x, y)

	arr := FindAll(bit0, bit, 0.1)
	fmt.Println("Find all bitmap:", arr)
	fmt.Println("find len: ", len(arr))
	// tt.Equal(t, 1, len(arr))

	c1 := robotgo.CHex(0xAADCDC)
	x, y = FindColor(c1)
	fmt.Println("Find color: ", x, y)
	arr = FindAllColor(c1)
	fmt.Println("Find all color: ", arr)

	img := robotgo.ToImage(bit)
	err := robotgo.SavePng(img, "robot_img.png")
	tt.Nil(t, err)

	bit1 := Open("robot_test.png")
	b := tt.TypeOf(bit, bit1)
	tt.True(t, b)
	tt.NotNil(t, bit1)
}
