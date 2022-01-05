// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.
//

package bitmap

/*
#cgo darwin CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo darwin LDFLAGS: -framework Cocoa
//
#cgo darwin,amd64 LDFLAGS:-L${SRCDIR}/cdeps/mac/amd -lpng -lz
#cgo darwin,arm64 LDFLAGS:-L${SRCDIR}/cdeps/mac/m1 -lpng -lz
//
#cgo linux LDFLAGS: -L/usr/src -lpng -lz
//
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/cdeps/win/amd/win64 -lpng -lz
#cgo windows,386 LDFLAGS: -L${SRCDIR}/cdeps/win/amd/win32 -lpng -lz
#cgo windows,arm64 LDFLAGS:-L${SRCDIR}/cdeps/win/arm -lpng -lz
//
#include "c/goBitmap.h"
*/
import "C"

import (
	"reflect"
	"unsafe"

	"github.com/go-vgo/robotgo"
	// "github.com/vcaesar/tt"
)

/*
.______    __  .___________..___  ___.      ___      .______
|   _  \  |  | |           ||   \/   |     /   \     |   _  \
|  |_)  | |  | `---|  |----`|  \  /  |    /  ^  \    |  |_)  |
|   _  <  |  |     |  |     |  |\/|  |   /  /_\  \   |   ___/
|  |_)  | |  |     |  |     |  |  |  |  /  _____  \  |  |
|______/  |__|     |__|     |__|  |__| /__/     \__\ | _|
*/

// SaveCapture capture the screen and save to path
func SaveCapture(spath string, args ...int) string {
	bit := robotgo.CaptureScreen(args...)

	err := Save(bit, spath)
	robotgo.FreeBitmap(bit)
	return err
}

// ToThis trans robotgo.CBitmap to C.MMBitmapRef
func ToThis(bit robotgo.CBitmap) C.MMBitmapRef {
	return (C.MMBitmapRef)(unsafe.Pointer(reflect.ValueOf(bit).Pointer()))
}

// ToRobot trans C.MMBitmapRef to robotgo.CBitmap
func ToRobot(bit C.MMBitmapRef) robotgo.CBitmap {
	return (robotgo.CBitmap)(unsafe.Pointer(reflect.ValueOf(bit).Pointer()))
}

// ToC trans robotgo.Bitmap to C.MMBitmapRef
func ToC(bit robotgo.Bitmap) C.MMBitmapRef {
	cbitmap := C.createMMBitmap(
		(*C.uint8_t)(bit.ImgBuf),
		C.size_t(bit.Width),
		C.size_t(bit.Height),
		C.size_t(bit.Bytewidth),
		C.uint8_t(bit.BitsPixel),
		C.uint8_t(bit.BytesPerPixel),
	)

	return cbitmap
}

// ToMMBitmapRef trans CBitmap to C.MMBitmapRef
func ToMMBitmapRef(bit robotgo.CBitmap) C.MMBitmapRef {
	return ToThis(bit)
}

// ToBytes saves robotgo.CBitmap to format in bytes
func ToBytes(bit robotgo.CBitmap) []byte {
	var len C.size_t
	ptr := C.saveMMBitmapAsBytes(ToThis(bit), &len)
	if int(len) < 0 {
		return nil
	}

	bs := C.GoBytes(unsafe.Pointer(ptr), C.int(len))
	C.free(unsafe.Pointer(ptr))
	return bs
}

// Tostring tostring bitmap to string
func Tostring(bit robotgo.CBitmap) string {
	strBit := C.tostring_bitmap(ToThis(bit))
	return C.GoString(strBit)
}

// Tochar tostring bitmap to C.char
func Tochar(bit robotgo.CBitmap) *C.char {
	strBit := C.tostring_bitmap(ToThis(bit))
	return strBit
}

func internalFind(bit, sbit robotgo.CBitmap, tolerance float64) (int, int) {
	pos := C.find_bitmap(ToThis(bit), ToThis(sbit), C.float(tolerance))
	return int(pos.x), int(pos.y)
}

// Find find the bitmap's pos in source bitmap
//
//	bitmap.Find(bitmap, source_bitmap robotgo.CBitmap, tolerance float64)
//
// 	|tolerance| should be in the range 0.0f - 1.0f, denoting how closely the
// 	colors in the bitmaps need to match, with 0 being exact and 1 being any.
//
// This method only automatically free the internal bitmap,
// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
func Find(bit robotgo.CBitmap, args ...interface{}) (int, int) {
	var (
		sbit      robotgo.CBitmap
		tolerance = 0.01
	)

	if len(args) > 0 && args[0] != nil {
		sbit = args[0].(robotgo.CBitmap)
	} else {
		sbit = robotgo.CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = args[1].(float64)
	}

	fx, fy := internalFind(bit, sbit, tolerance)
	// FreeBitmap(bit)
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
		robotgo.FreeBitmap(sbit)
	}

	return fx, fy
}

// FindPic finding the image by path
//
//	bitmap.FindPic(path string, source_bitmap robotgo.CBitmap, tolerance float64)
//
// This method only automatically free the internal bitmap,
// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
func FindPic(path string, args ...interface{}) (int, int) {
	var (
		sbit      robotgo.CBitmap
		tolerance = 0.01
	)

	openbit := Open(path)
	if len(args) > 0 && args[0] != nil {
		sbit = args[0].(robotgo.CBitmap)
	} else {
		sbit = robotgo.CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = args[1].(float64)
	}

	fx, fy := internalFind(openbit, sbit, tolerance)
	robotgo.FreeBitmap(openbit)
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
		robotgo.FreeBitmap(sbit)
	}

	return fx, fy
}

// FreeMMPointArr free MMPoint array
func FreeMMPointArr(pointArray C.MMPointArrayRef) {
	C.destroyMMPointArray(pointArray)
}

// FindAll find the all bitmap
func FindAll(bit robotgo.CBitmap, args ...interface{}) (posArr []robotgo.Point) {
	var (
		sbit      robotgo.CBitmap
		tolerance C.float = 0.01
		lpos      C.MMPoint
	)

	if len(args) > 0 && args[0] != nil {
		sbit = args[0].(robotgo.CBitmap)
	} else {
		sbit = robotgo.CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = C.float(args[1].(float64))
	}

	if len(args) > 2 {
		lpos.x = C.size_t(args[2].(int))
		lpos.y = 0
	} else {
		lpos.x = 0
		lpos.y = 0
	}

	if len(args) > 3 {
		lpos.x = C.size_t(args[2].(int))
		lpos.y = C.size_t(args[3].(int))
	}

	pos := C.find_every_bitmap(ToThis(bit), ToThis(sbit), tolerance, &lpos)
	// FreeBitmap(bit)
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
		robotgo.FreeBitmap(sbit)
	}
	if pos == nil {
		return
	}
	defer FreeMMPointArr(pos)

	cSize := pos.count
	cArray := pos.array
	gSlice := (*[(1 << 28) - 1]C.MMPoint)(unsafe.Pointer(cArray))[:cSize:cSize]
	for i := 0; i < len(gSlice); i++ {
		posArr = append(posArr, robotgo.Point{
			X: int(gSlice[i].x),
			Y: int(gSlice[i].y),
		})
	}

	return
}

// Count count of the bitmap
func Count(bit, sourceBitmap robotgo.CBitmap, args ...float32) int {
	var tolerance C.float = 0.01
	if len(args) > 0 {
		tolerance = C.float(args[0])
	}

	count := C.count_of_bitmap(ToThis(bit), ToThis(sourceBitmap), tolerance)
	return int(count)
}

// Click find the bitmap and click
func Click(bit robotgo.CBitmap, args ...interface{}) {
	x, y := Find(bit)
	robotgo.MovesClick(x, y, args...)
}

// PointInBounds bitmap point in bounds
func PointInBounds(bit robotgo.CBitmap, x, y int) bool {
	var point C.MMPoint
	point.x = C.size_t(x)
	point.y = C.size_t(y)
	cbool := C.point_in_bounds(ToThis(bit), point)

	return bool(cbool)
}

// OpenC open the bitmap return C.MMBitmapRef
//
// bitmap.Open(path string, type int)
func OpenC(gpath string, args ...int) C.MMBitmapRef {
	path := C.CString(gpath)
	var mtype C.uint16_t = 1

	if len(args) > 0 {
		mtype = C.uint16_t(args[0])
	}

	bit := C.bitmap_open(path, mtype)
	C.free(unsafe.Pointer(path))

	return bit
}

// Open open the bitmap image return robotgo.CBitmap
func Open(path string, args ...int) robotgo.CBitmap {
	return ToRobot(OpenC(path, args...))
}

// FromStr bitmap from string
func FromStr(str string) C.MMBitmapRef {
	cs := C.CString(str)
	bit := C.bitmap_from_string(cs)
	C.free(unsafe.Pointer(cs))

	return bit
}

// Save save the bitmap to image
//
// bitmap.Save(bitmap robotgo.CBitmap, path string, type int)
func Save(bit robotgo.CBitmap, gpath string, args ...int) string {
	var mtype C.uint16_t = 1
	if len(args) > 0 {
		mtype = C.uint16_t(args[0])
	}

	path := C.CString(gpath)
	saveBit := C.bitmap_save(ToThis(bit), path, mtype)
	C.free(unsafe.Pointer(path))

	return C.GoString(saveBit)
}

// GetPortion get bitmap portion
func GetPortion(bit robotgo.CBitmap, x, y, w, h int) robotgo.CBitmap {
	var rect C.MMRect
	rect.origin.x = C.size_t(x)
	rect.origin.y = C.size_t(y)
	rect.size.width = C.size_t(w)
	rect.size.height = C.size_t(h)

	pos := C.get_portion(ToThis(bit), rect)
	return ToRobot(pos)
}

// Convert convert the bitmap
//
// bitmap.Convert(opath, spath string, type int)
func Convert(opath, spath string, args ...int) string {
	var mtype = 1
	if len(args) > 0 {
		mtype = args[0]
	}

	bit := Open(opath)
	return Save(bit, spath, mtype)
}

// FreeBitmapArr free and dealloc the C bitmap array
func FreeArr(bit ...robotgo.CBitmap) {
	for i := 0; i < len(bit); i++ {
		robotgo.FreeBitmap(bit[i])
	}
}

// FreeArrC free and dealloc the C.MMBitmapRef bitmap array
func FreeArrC(bit ...C.MMBitmapRef) {
	for i := 0; i < len(bit); i++ {
		robotgo.FreeBitmap(ToRobot(bit[i]))
	}
}

// Read returns false and sets error if |bitmap| is NULL
func Read(bit robotgo.CBitmap) bool {
	abool := C.bitmap_ready(ToThis(bit))
	return bool(abool)
}

// CopyToPB copy bitmap to pasteboard
func CopyToPB(bit robotgo.CBitmap) bool {
	abool := C.bitmap_copy_to_pboard(ToThis(bit))
	return bool(abool)
}

// DeepCopyC deep copy bitmap to new bitamp
func DeepCopyC(bit C.MMBitmapRef) C.MMBitmapRef {
	bit1 := C.bitmap_deepcopy(bit)
	return bit1
}

// DeepCopy deep copy bitmap to new bitmap
func DeepCopy(bit robotgo.CBitmap) robotgo.CBitmap {
	return ToRobot(DeepCopyC(ToThis(bit)))
}

// GetColor get the bitmap color
func GetColor(bit robotgo.CBitmap, x, y int) C.MMRGBHex {
	color := C.bitmap_get_color(ToThis(bit), C.size_t(x), C.size_t(y))
	return color
}

// ToRHex trans C.MMRGBHex to robotgo.CHex
func ToRHex(hex C.MMRGBHex) robotgo.CHex {
	return robotgo.CHex(hex)
}

// GetColors get bitmap color retrun string
func GetColors(bit robotgo.CBitmap, x, y int) string {
	clo := GetColor(bit, x, y)
	return robotgo.PadHexs(ToRHex(clo))
}

// FindColor find bitmap color
//
// bitmap.FindColor(color CHex, bitmap robotgo.CBitmap, tolerance float)
func FindColor(color robotgo.CHex, args ...interface{}) (int, int) {
	var (
		tolerance C.float = 0.01
		bit       robotgo.CBitmap
	)

	if len(args) > 0 && args[0] != nil {
		bit = args[0].(robotgo.CBitmap)
	} else {
		bit = robotgo.CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = C.float(args[1].(float64))
	}

	pos := C.bitmap_find_color(ToThis(bit), C.MMRGBHex(color), tolerance)
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
		robotgo.FreeBitmap(bit)
	}

	x := int(pos.x)
	y := int(pos.y)

	return x, y
}

// FindColorCS findcolor by CaptureScreen
func FindColorCS(color robotgo.CHex, x, y, w, h int, args ...float64) (int, int) {
	var tolerance = 0.01

	if len(args) > 0 {
		tolerance = args[0]
	}

	bit := robotgo.CaptureScreen(x, y, w, h)
	rx, ry := FindColor(color, bit, tolerance)
	robotgo.FreeBitmap(bit)

	return rx, ry
}

// FindAllColor find the all color
func FindAllColor(color robotgo.CHex, args ...interface{}) (posArr []robotgo.Point) {
	var (
		bit robotgo.CBitmap
		// bitmap    C.MMBitmapRef
		tolerance C.float = 0.01
		lpos      C.MMPoint
	)

	if len(args) > 0 && args[0] != nil {
		bit = args[0].(robotgo.CBitmap)
	} else {
		bit = robotgo.CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = C.float(args[1].(float64))
	}

	if len(args) > 2 {
		lpos.x = C.size_t(args[2].(int))
		lpos.y = 0
	} else {
		lpos.x = 0
		lpos.y = 0
	}

	if len(args) > 3 {
		lpos.x = C.size_t(args[2].(int))
		lpos.y = C.size_t(args[3].(int))
	}

	pos := C.bitmap_find_every_color(ToThis(bit), C.MMRGBHex(color), tolerance, &lpos)
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
		robotgo.FreeBitmap(bit)
	}

	if pos == nil {
		return
	}
	defer FreeMMPointArr(pos)

	cSize := pos.count
	cArray := pos.array
	gSlice := (*[(1 << 28) - 1]C.MMPoint)(unsafe.Pointer(cArray))[:cSize:cSize]
	for i := 0; i < len(gSlice); i++ {
		posArr = append(posArr, robotgo.Point{
			X: int(gSlice[i].x),
			Y: int(gSlice[i].y),
		})
	}

	return
}

// CountColor count bitmap color
func CountColor(color robotgo.CHex, args ...interface{}) int {
	var (
		tolerance C.float = 0.01
		bit       robotgo.CBitmap
	)

	if len(args) > 0 && args[0] != nil {
		bit = args[0].(robotgo.CBitmap)
	} else {
		bit = robotgo.CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = C.float(args[1].(float64))
	}

	count := C.bitmap_count_of_color(ToThis(bit), C.MMRGBHex(color), tolerance)
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
		robotgo.FreeBitmap(bit)
	}

	return int(count)
}

// CountColorCS count bitmap color by CaptureScreen
func CountColorCS(color robotgo.CHex, x, y, w, h int, args ...float64) int {
	var tolerance = 0.01

	if len(args) > 0 {
		tolerance = args[0]
	}

	bit := robotgo.CaptureScreen(x, y, w, h)
	rx := CountColor(color, bit, tolerance)
	robotgo.FreeBitmap(bit)

	return rx
}

// GetSize get the image size
func GetSize(imgPath string) (int, int) {
	bit := Open(imgPath)
	gbit := robotgo.ToBitmap(bit)

	w := gbit.Width / 2
	h := gbit.Height / 2
	robotgo.FreeBitmap(bit)

	return w, h
}
