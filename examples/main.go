// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package main

import (
	"fmt"
	"log"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/bitmap"
	"github.com/vcaesar/imgo"
)

func toBitmap(bmp robotgo.CBitmap) {
	img := robotgo.ToImage(bmp)
	fmt.Println("img: ", img)
	imgo.SaveToPNG("test_IMG.png", img)

	gbit := robotgo.ToBitmap(bmp)
	fmt.Println("go bitmap", gbit, gbit.Width)

	cbit := robotgo.ToCBitmap(gbit)
	// defer robotgo.FreeBitmap(cbit)
	log.Println("cbit == bitmap: ", cbit == bmp)
}

// find color
func findColor(bmp robotgo.CBitmap) {
	// find the color in bitmap
	color := bitmap.GetColor(bmp, 1, 2)
	fmt.Println("color...", color)

	cx, cy := bitmap.FindColor(robotgo.CHex(color), bmp, 1.0)
	fmt.Println("pos...", cx, cy)
	cx, cy = bitmap.FindColor(robotgo.CHex(color))
	fmt.Println("pos...", cx, cy)

	cx, cy = bitmap.FindColor(0xAADCDC, bmp)
	fmt.Println("pos...", cx, cy)

	cx, cy = bitmap.FindColor(0xAADCDC, nil, 0.1)
	fmt.Println("pos...", cx, cy)

	cx, cy = bitmap.FindColorCS(0xAADCDC, 388, 179, 300, 300)
	fmt.Println("pos...", cx, cy)

	cnt := bitmap.CountColor(0xAADCDC, bmp)
	fmt.Println("count...", cnt)

	cnt1 := bitmap.CountColorCS(0xAADCDC, 10, 20, 30, 40)
	fmt.Println("count...", cnt1)

	arr := bitmap.FindAllColor(0xAADCDC)
	fmt.Println("find all color: ", arr)
	for i := 0; i < len(arr); i++ {
		fmt.Println("pos is: ", arr[i].X, arr[i].Y)
	}
}

func bitmapString(bmp robotgo.CBitmap) {
	// creates bitmap from string by bitmap
	bitstr := bitmap.Tostring(bmp)
	fmt.Println("bitstr...", bitstr)

	sbitmap := bitmap.FromStr(bitstr)
	fmt.Println("bitmap str...", sbitmap)

	bitmap.Save(sbitmap, "teststr.png")
}

func bitmapTool(bmp robotgo.CBitmap) {
	abool := bitmap.PointInBounds(bmp, 1, 2)
	fmt.Println("point in bounds...", abool)

	// returns new bitmap object created from a portion of another
	bitpos := bitmap.GetPortion(bmp, 10, 10, 11, 10)
	fmt.Println(bitpos)

	// saves image to absolute filepath in the given format
	bitmap.Save(bmp, "test.png")
}

func decode() {
	img, name, err := robotgo.DecodeImg("test.png")
	if err != nil {
		log.Println("decode image ", err)
	}
	fmt.Println("decode test.png", img, name)

	byt, _ := robotgo.OpenImg("test.png")
	imgo.SaveByte("test2.png", byt)

	w, h := bitmap.GetSize("test.png")
	fmt.Println("image width and hight ", w, h)
	w, h, _ = imgo.GetSize("test.png")
	fmt.Println("image width and hight ", w, h)

	// convert image
	bitmap.Convert("test.png", "test.tif")
}

func bitmapTest(bmp robotgo.CBitmap) {
	bit := robotgo.CaptureScreen(1, 2, 40, 40)
	defer robotgo.FreeBitmap(bit)
	fmt.Println("CaptureScreen...", bit)

	// searches for needle in bitmap
	fx, fy := bitmap.Find(bit, bmp)
	fmt.Println("FindBitmap------", fx, fy)

	fx, fy = bitmap.Find(bit)
	fmt.Println("FindBitmap------", fx, fy)

	fx, fy = bitmap.Find(bit, nil, 0.2)
	fmt.Println("find bitmap: ", fx, fy)

	fx, fy = bitmap.Find(bit, bmp, 0.3)
	fmt.Println("find bitmap: ", fx, fy)
}

func findBitmap(bmp robotgo.CBitmap) {
	fx, fy := bitmap.Find(bmp)
	fmt.Println("findBitmap: ", fx, fy)

	// open image bitmap
	openbit := bitmap.Open("test.tif")
	fmt.Println("openBitmap...", openbit)

	fx, fy = bitmap.Find(openbit)
	fmt.Println("FindBitmap------", fx, fy)

	fx, fy = bitmap.FindPic("test.tif")
	fmt.Println("FindPic------", fx, fy)

	arr := bitmap.FindAll(openbit)
	fmt.Println("find all bitmap: ", arr)
	for i := 0; i < len(arr); i++ {
		fmt.Println("pos is: ", arr[i].X, arr[i].Y)
	}
}

func bitmap1() {
	////////////////////////////////////////////////////////////////////////////////
	// Bitmap
	////////////////////////////////////////////////////////////////////////////////

	// gets all of the screen
	abitMap := robotgo.CaptureScreen()
	fmt.Println("abitMap...", abitMap)

	// gets part of the screen
	cbit := robotgo.CaptureScreen(100, 200, 30, 30)
	defer robotgo.FreeBitmap(cbit)
	fmt.Println("CaptureScreen...", cbit)

	toBitmap(cbit)

	findColor(cbit)

	count := bitmap.Count(abitMap, cbit)
	fmt.Println("count...", count)

	bitmapTest(cbit)
	findBitmap(cbit)

	bitmapString(cbit)
	bitmapTool(cbit)

	decode()

	// free the bitmap
	robotgo.FreeBitmap(abitMap)
}

func main() {
	bitmap1()
}
