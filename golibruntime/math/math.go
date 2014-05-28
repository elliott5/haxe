// Copyright 2014 Elliott Stoneham and The tardisgo Authors
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Bitcast elements adapted from https://github.com/gopherjs/gopherjs/blob/master/bitcasts/bitcasts.go
/*
Copyright (c) 2013 Richard Musiol. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Package math provides implementation and documentation of math functions overloaded by TARDIS Go->Haxe transpiler
package math

import "math"

// Code below adapted from https://github.com/gopherjs/gopherjs/blob/master/bitcasts/bitcasts.go

//commented out code from math package unsafe.go, which is overloaded with the code below
//import "unsafe"

// Float32bits returns the IEEE 754 binary representation of f.
//func Float32bits(f float32) uint32 { return *(*uint32)(unsafe.Pointer(&f)) }
func glrFloat32bits(f float32) uint32 {
	var Zero = 0.0
	var NegZero = -Zero
	//var NaN = Zero / Zero
	if f == 0 {
		if f == 0 && 1/f == float32(1/NegZero) {
			return 1 << 31
		}
		return 0
	}
	if f != f { // NaN
		return 2143289344
	}

	s := uint32(0)
	if f < 0 {
		s = 1 << 31
		f = -f
	}

	e := uint32(127 + 23)
	for f >= 1<<24 {
		f /= 2
		if e == (1<<8)-1 {
			break
		}
		e++
	}
	for f < 1<<23 {
		e--
		if e == 0 {
			break
		}
		f *= 2
	}

	r := math.Mod(float64(f), 2)
	if (r > 0.5 && r < 1) || r >= 1.5 { // round to nearest even
		f++
	}

	return s | uint32(e)<<23 | (uint32(f) &^ (1 << 23))
}

// Float32frombits returns the floating point number corresponding
// to the IEEE 754 binary representation b.
// func Float32frombits(b uint32) float32 { return *(*float32)(unsafe.Pointer(&b)) }
func glrFloat32frombits(b uint32) float32 {
	var Zero = 0.0
	//var NegZero = -Zero
	var NaN = Zero / Zero
	s := float32(+1)
	if b&(1<<31) != 0 {
		s = -1
	}
	e := (b >> 23) & (1<<8 - 1)
	m := b & (1<<23 - 1)

	if e == (1<<8)-1 {
		if m == 0 {
			return s / 0 // Inf
		}
		return float32(NaN)
	}
	if e != 0 {
		m += 1 << 23
	}
	if e == 0 {
		e = 1
	}

	return float32(math.Ldexp(float64(m), int(e)-127-23)) * s
}

// Float64bits returns the IEEE 754 binary representation of f.
//func Float64bits(f float64) uint64 { return *(*uint64)(unsafe.Pointer(&f)) }
func glrFloat64bits(f float64) uint64 {
	var Zero = 0.0
	var NegZero = -Zero
	//var NaN = Zero / Zero
	if f == 0 {
		if f == 0 && 1/f == 1/NegZero {
			return 1 << 63
		}
		return 0
	}
	if f != f { // NaN
		return 9221120237041090561
	}

	s := uint64(0)
	if f < 0 {
		s = 1 << 63
		f = -f
	}

	e := uint32(1023 + 52)
	for f >= 1<<53 {
		f /= 2
		if e == (1<<11)-1 {
			break
		}
		e++
	}
	for f < 1<<52 {
		e--
		if e == 0 {
			break
		}
		f *= 2
	}

	return s | uint64(e)<<52 | (uint64(f) &^ (1 << 52))
}

// Float64frombits returns the floating point number corresponding
// the IEEE 754 binary representation b.
//func Float64frombits(b uint64) float64 { return *(*float64)(unsafe.Pointer(&b)) }
func glrFloat64frombits(b uint64) float64 {
	var Zero = 0.0
	//var NegZero = -Zero
	var NaN = Zero / Zero
	s := float64(+1)
	if b&(1<<63) != 0 {
		s = -1
	}
	e := (b >> 52) & (1<<11 - 1)
	m := b & (1<<52 - 1)

	if e == (1<<11)-1 {
		if m == 0 {
			return s / 0
		}
		return NaN
	}
	if e != 0 {
		m += 1 << 52
	}
	if e == 0 {
		e = 1
	}

	return math.Ldexp(float64(m), int(e)-1023-52) * s
}

// end adapted code

// the entries below for documentation purposes only

//emulated in Go standard math package:
/*
func Atan2(y, x float64) float64                 { return atan2(x, y) }
func Frexp(f float64) (frac float64, exp int)    { return frexp(f) }
func Hypot(p, q float64) float64                 { return hypot(p, q) }
func Ldexp(frac float64, exp int) float64        { return ldexp(frac, exp) }
func Log1p(x float64) float64                    { return log1p(x) }
func Mod(x, y float64) float64                   { return mod(x, y) }
func Modf(f float64) (int float64, frac float64) { return modf(f) }
func Sincos(x float64) (sin, cos float64)        { return sincos(x) }
*/

//Go Math functions
//mapped into Haxe:
/*
	"math_Abs":   "Math.abs",
	"math_Acos":  "Math.acos",
	"math_Asin":  "Math.asin",
	"math_Atan":  "Math.atan",
	"math_Ceil":  "Math.fceil",
	"math_Cos":   "Math.cos",
	"math_Exp":   "Math.exp",
	"math_Floor": "Math.ffloor",
	"math_Log":   "Math.log",
	"math_Sin":   "Math.sin",
	"math_Sqrt":  "Math.sqrt",
	"math_Tan":   "Math.tan",
*/
