// Copyright 2014 Elliott Stoneham and The TARDIS Go Authors
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package pogo

import (
	"code.google.com/p/go.tools/go/exact"
	"code.google.com/p/go.tools/go/ssa"
	"fmt"
	"go/token"
)

// emit the constant declarations
func emitNamedConstants() {
	allPack := rootProgram.AllPackages()
	for pkgIdx := range allPack {
		pkg := allPack[pkgIdx]
		for mName, mem := range pkg.Members {
			if mem.Token() == token.CONST {
				lit := mem.(*ssa.NamedConst).Value
				posStr := CodePosition(lit.Pos())
				pName := mem.(*ssa.NamedConst).Object().Pkg().Name()
				switch lit.Value.Kind() { // non language specific validation
				case exact.Bool, exact.String, exact.Float, exact.Int, exact.Complex: //OK
					isPublic := mem.Object().Exported()
					if isPublic { // constants will be inserted inline, these declarations of public constants are for exteral use in target language
						l := TargetLang
						_, _, isOverloaded := LanguageList[l].PackageOverloaded(pName)
						if !isOverloaded { // only emit constants from non-overloaded packages
							fmt.Fprintln(&LanguageList[l].buffer, LanguageList[l].NamedConst(pName, mName, *lit, posStr))
						}
					}
				default:
					LogError(posStr, "pogo", fmt.Errorf("%s.%s : emitConstants() internal error, unrecognised constant type: %v",
						pName, mName, lit.Value.Kind()))
				}
			}
		}
	}
}

// Float64Val is a utility function returns a string constant value from an exact.Value.
func Float64Val(eVal exact.Value, posStr string) string {
	fVal, isExact := exact.Float64Val(eVal)
	if !isExact {
		LogWarning(posStr, "inexact", fmt.Errorf("constant value %g cannot be accurately represented in float64", fVal))
	}
	if fVal < 0.0 {
		return fmt.Sprintf("(%g)", fVal)
	}
	return fmt.Sprintf("%g", fVal)
}

// IntVal is a utility function returns an int64 constant value from an exact.Value, split into high and low int32.
func IntVal(eVal exact.Value, posStr string) (high, low int32) {
	iVal, isExact := exact.Int64Val(eVal)
	if !isExact {
		LogWarning(posStr, "inexact", fmt.Errorf("constant value %d cannot be accurately represented in int64", iVal))
	}
	return int32(iVal >> 32), int32(iVal & 0xFFFFFFFF)
}
