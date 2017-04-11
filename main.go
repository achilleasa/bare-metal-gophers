package main

import (
	"reflect"
	"unsafe"
)

const (
	fbWidth            = 80
	fbHeight           = 25
	fbPhysAddr uintptr = 0xb8000
)

func main() {
	// Display a string to the top-left corner of the screen one character
	// at a time.
	buf := []byte{'H', 'e', 'l', 'l', 'o', ' ', 'G', 'o', 'l', 'a', 'n', 'g', 'U', 'K', '!'}
	attr := uint16(2<<4 | 0) // black text; green bg
	for i, b := range buf {
		*(*uint16)(unsafe.Pointer(fbPhysAddr + uintptr(i*2))) = attr<<8 | uint16(b)
	}

	delay(5000)

	// To setup our framebuffer slice we setup a fake reflect.SliceHeader that
	// points to the physical memory address of the mapped VGA text buffer
	var fb = *(*[]uint16)(unsafe.Pointer(&reflect.SliceHeader{
		Len:  fbWidth * fbHeight,
		Cap:  fbWidth * fbHeight,
		Data: fbPhysAddr,
	}))

	// Since both fb and logo are slices we can use the built-in copy funtion
	// to draw our logo.
	copy(fb, logo)
	transition(fb)

	// For the final part of the demo we will run a complex piece of code
	// that is designed so that no variables escape to the heap. This is
	// important as any calls to the runtime memory allocator will cause
	// the machine to triple-fault.
	bareMetalGophers(
		fb, // render target
		float64(fbWidth), float64(fbHeight),
		vec3{0, 5, 7.5}, // eye
		vec3{0, 0, 0},   // look-at
		vec3{0, 1, 0},   // camera up
		0.03,            // tick step
	)
}

// transition implements a slide transition using the current contents of the
// supplied framebuffer.
func transition(fb []uint16) {
	delay(5000)

	for i := 0; i < fbWidth; i++ {
		for y, off := 0, 0; y < fbHeight; y, off = y+1, off+fbWidth {
			// Even rows should slide one character to the left and
			// odd rows should slide one character to the right
			if y%2 == 0 {
				copy(fb[off:off+fbWidth], fb[off+1:off+fbWidth])
				fb[off+fbWidth-1] = ' '
			} else {
				copy(fb[off+1:off+fbWidth], fb[off:off+fbWidth-1])
				fb[off] = ' '
			}
		}
		delay(50)
	}
}

// delay implements a simple loop-based delay. The outer loop value is selected
// so that a reasonable delay is generated when running on virtualbox.
func delay(v int) {
	for i := 0; i < 684000; i++ {
		for j := 0; j < v; j++ {
		}
	}
}
