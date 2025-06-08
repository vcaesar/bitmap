// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>  
//
// This file may not be copied, modified, or distributed
// except according to those terms.

#include "bitmap_class.h"
#include "bitmap_find_c.h"
#include "../base/bitmap/MMBitmap_c.h"
#include "../base/bitmap/color_find_c.h"
//
#include "../base/bitmap/file_io_c.h"
#include "../base/bitmap/pasteboard_c.h"
#include "../base/bitmap/str_io_c.h"
//
#include <assert.h>
#include <stdio.h>

/* Returns false and sets error if |bitmap| is NULL. */
bool bitmap_ready(MMBitmapRef bitmap) {
	if (bitmap == NULL || bitmap->imageBuffer == NULL) {
		return false;
	}
	return true;
}

// void bitmap_dealloc(MMBitmapRef bitmap){
// 	if (bitmap != NULL) {
// 		destroyMMBitmap(bitmap);
// 		bitmap = NULL;
// 	}
// }

bool bitmap_copy_to_pboard(MMBitmapRef bitmap) {
	MMPasteError err;
	if (!bitmap_ready(bitmap)) { return false; }

	if ((err = copyMMBitmapToPasteboard(bitmap)) != kMMPasteNoError) {
		return false;
	}

	return true;
}

MMBitmapRef bitmap_deepcopy(MMBitmapRef bitmap) {
	return bitmap == NULL ? NULL : copyMMBitmap(bitmap);
}

MMPoint find_bitmap(MMBitmapRef bitmap, MMBitmapRef sbit, float tolerance) {
	MMPoint point = {-1, -1};
	if (!bitmap_ready(sbit) || !bitmap_ready(bitmap)) {
		printf("bitmap is not ready yet!\n");
		return point;
	}

	MMRect rect = MMBitmapGetBounds(sbit);
	if (findBitmapInRect(bitmap, sbit, &point, rect, tolerance) == 0) {
		return point;
	}

	return point;
}

MMPointArrayRef find_every_bitmap(MMBitmapRef bitmap, MMBitmapRef sbit, float tolerance, MMPoint *list) {
	if (!bitmap_ready(bitmap) || !bitmap_ready(sbit)) { return NULL; }

	MMPoint point;
	MMPointArrayRef pointArray;
	MMRect rect = MMBitmapGetBounds(sbit);

	if (findBitmapInRect(bitmap, sbit, &point, rect, tolerance) == 0) {
		// return NULL;
	}

	pointArray = findAllBitmapInRect(bitmap, sbit, rect, tolerance);
	return pointArray;
}

int count_of_bitmap(MMBitmapRef bitmap, MMBitmapRef sbit, float tolerance) {
	if (!bitmap_ready(bitmap) || !bitmap_ready(sbit)) { return 0; }

	MMRect rect = MMBitmapGetBounds(bitmap);
	return countOfBitmapInRect(bitmap, sbit, rect, tolerance);
}

bool point_in_bounds(MMBitmapRef bitmap, MMPoint point) {
	if (!bitmap_ready(bitmap)) {
		return NULL;
	}

	if (MMBitmapPointInBounds(bitmap, point)) {
		return true;
	}

	return false;
}

MMBitmapRef bitmap_open(char *path, uint16_t ttype) {
	MMBitmapRef bitmap;
	MMIOError err;

	bitmap = newMMBitmapFromFile(path, ttype, &err);
	return bitmap;
}

MMBitmapRef bitmap_from_string(const char *str) {
	size_t len = strlen(str);
	MMBitmapRef bitmap;
	MMBMPStringError err;

	if ((bitmap = createMMBitmapFromString((unsigned char*)str, len, &err)) == NULL) {
		return NULL;
	}

	return bitmap;
}

char *bitmap_save(MMBitmapRef bitmap, char *path, uint16_t type) {
	if (saveMMBitmapToFile(bitmap, path, (MMImageType) type) != 0) {
		return "Could not save image to file.";
	}
	// destroyMMBitmap(bitmap);
	return "";
}

char *tostring_bitmap(MMBitmapRef bitmap) {
	char *buf = NULL;
	MMBMPStringError err;

	buf = (char *)createStringFromMMBitmap(bitmap, &err);
	return buf;
}

// char out size 200 is enough
bool bitmap_str(MMBitmapRef bitmap, char *out) {
	if (!bitmap_ready(bitmap)) { return false; }
	sprintf(out, "<Bitmap with resolution %lu%lu, \
	                    %u bits per pixel, and %u bytes per pixel>",
	                    (unsigned long)bitmap->width,
	                    (unsigned long)bitmap->height,
	                    bitmap->bitsPerPixel,
	                    bitmap->bytesPerPixel);

	return true;
}

MMBitmapRef get_portion(MMBitmapRef bit_map, MMRect rect) {
	// MMRect rect;
	MMBitmapRef portion = NULL;

	portion = copyMMBitmapFromPortion(bit_map, rect);
	return portion;
}

MMRGBHex bitmap_get_color(MMBitmapRef bitmap, size_t x, size_t y) {
	if (!bitmap_ready(bitmap)) { return 0; }

	MMPoint point = MMPointMake(x, y);
	if (!MMBitmapPointInBounds(bitmap, point)) {
		return 0;
	}

	return MMRGBHexAtPoint(bitmap, point.x, point.y);
}

MMPoint bitmap_find_color(MMBitmapRef bitmap, MMRGBHex color, float tolerance) {
	MMRect rect = MMBitmapGetBounds(bitmap);
	MMPoint point = {-1, -1};

	if (findColorInRect(bitmap, color, &point, rect, tolerance) == 0) {
		return point;
	}

	return point;
}

MMPointArrayRef bitmap_find_every_color(MMBitmapRef bitmap, MMRGBHex color, float tolerance, MMPoint *list) {
	if (!bitmap_ready(bitmap)) { return NULL; }
	MMRect rect = MMBitmapGetBounds(bitmap);
	MMPointArrayRef pointArray;

	pointArray = findAllColorInRect(bitmap, color, rect, tolerance);
	return pointArray;
}

int bitmap_count_of_color(MMBitmapRef bitmap, MMRGBHex color, float tolerance) {
	if (!bitmap_ready(bitmap)) { return 0; }
	MMRect rect = MMBitmapGetBounds(bitmap);

	return countOfColorsInRect(bitmap, color, rect, tolerance);
}
