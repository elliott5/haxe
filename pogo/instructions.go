// Copyright 2014 Elliott Stoneham and The TARDIS Go Authors
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package pogo

import (
	"fmt"
	"reflect"

	"code.google.com/p/go.tools/go/ssa"
	"code.google.com/p/go.tools/go/types"
)

// RegisterName returns the name of an ssa.Value, a utility function in case it needs to be altered.
func RegisterName(val ssa.Value) string {
	//NOTE the SSA code says that name() should not be relied on, so this code may need to alter
	return "_" + val.Name()
}

var previousErrorInfo string // used to give some indication of the error's location, even if it is not given

// Handle an individual instruction.
func emitInstruction(instruction interface{}, operands []*ssa.Value) (emitPhiFlag bool) {
	l := TargetLang
	emitPhiFlag = true
	prev := LatestValidPosHash
	MakePosHash(instruction.(ssa.Instruction).Pos()) // this so that we log the nearby position info
	if prev != LatestValidPosHash {                  // new info, so put out an update
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].SetPosHash())
	}
	errorInfo := CodePosition(instruction.(ssa.Instruction).Pos())
	if errorInfo == "" {
		errorInfo = previousErrorInfo
	} else {
		previousErrorInfo = "near " + errorInfo
		errorInfo = "@ " + errorInfo
	}
	errorInfo = reflect.TypeOf(instruction).String() + " " + errorInfo //TODO consider removing as for DEBUG only
	instrVal, hasVal := instruction.(ssa.Value)
	register := ""
	comment := ""
	if hasVal {
		register = RegisterName(instrVal)
		comment = fmt.Sprintf("%s = %+v %s", register, instruction, errorInfo)
		//emitComment(comment)
		switch len(*instruction.(ssa.Value).Referrers()) {
		case 0: // no other instruction uses the result of this one
			comment += " [REGISTER VALUE UN-USED]"
			register = ""
		case 1: // only 1 other use of the register
		// TODO register optimisation currently disabled, consider reimplimentation
		//	user := (*instruction.(ssa.Value).Referrers())[0]
		//	if user.Block() == instruction.(ssa.Instruction).Block() {
		//		comment += " [REGISTER MAY BE OPTIMIZABLE]"
		//	}
		default: //multiple usage of the register
		}
		if len(register) > 0 {
			if LanguageList[TargetLang].LangType(instruction.(ssa.Value).Type(), false, errorInfo) == "" { // NOTE an empty type def makes a register useless too
				register = ""
			}
		}
	} else {
		comment = fmt.Sprintf("%+v %s", instruction, errorInfo)
		//emitComment(comment)
	}
	switch instruction.(type) {
	case *ssa.Jump:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].Jump(instruction.(*ssa.Jump).Block().Succs[0].Index)+LanguageList[l].Comment(comment))

	case *ssa.If:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].If(*operands[0],
				instruction.(*ssa.If).Block().Succs[0].Index,
				instruction.(*ssa.If).Block().Succs[1].Index,
				errorInfo)+LanguageList[l].Comment(comment))

	case *ssa.Phi:
		text := LanguageList[l].PhiStart(register,
			LanguageList[l].LangType(instrVal.Type(), false, errorInfo),
			LanguageList[l].LangType(instrVal.Type(), true, errorInfo))
		for o := range operands {
			text += LanguageList[l].PhiEntry(register, instruction.(*ssa.Phi).Block().Preds[o].Index, *operands[o], errorInfo)
		}
		text += LanguageList[l].PhiEnd(LanguageList[l].LangType(instrVal.Type(), true, errorInfo))
		fmt.Fprintln(&LanguageList[l].buffer, text+LanguageList[l].Comment(comment))

	case *ssa.Call:
		if instruction.(*ssa.Call).Call.IsInvoke() {
			fmt.Fprintln(&LanguageList[l].buffer,
				LanguageList[l].EmitInvoke(register, false, false, instruction.(*ssa.Call).Call, errorInfo)+
					LanguageList[l].Comment(comment))
		} else {
			switch instruction.(*ssa.Call).Call.Value.(type) {
			case *ssa.Builtin:
				emitCall(true, false, false, register, instruction.(*ssa.Call).Call, errorInfo, comment)
			default:
				emitCall(false, false, false, register, instruction.(*ssa.Call).Call, errorInfo, comment)
			}
		}

	case *ssa.Go:
		if instruction.(*ssa.Go).Call.IsInvoke() {
			fmt.Fprintln(&LanguageList[l].buffer,
				LanguageList[l].EmitInvoke(register, true, false, instruction.(*ssa.Go).Call, errorInfo)+
					LanguageList[l].Comment(comment))
		} else {
			switch instruction.(*ssa.Go).Call.Value.(type) {
			case *ssa.Builtin: // no builtin functions can be go'ed
				LogError(errorInfo, "pogo", fmt.Errorf("builtin functions cannot be go'ed"))
			default:
				emitCall(false, true, false, register, instruction.(*ssa.Go).Call, errorInfo, comment)
			}
		}

	case *ssa.Defer:
		if instruction.(*ssa.Defer).Call.IsInvoke() {
			fmt.Fprintln(&LanguageList[l].buffer,
				LanguageList[l].EmitInvoke(register, true, true, instruction.(*ssa.Defer).Call, errorInfo)+
					LanguageList[l].Comment(comment))
		} else {
			switch instruction.(*ssa.Defer).Call.Value.(type) {
			case *ssa.Builtin: // no builtin functions can be defer'ed - TODO: the spec does allow this in some circumstances
				LogError(errorInfo, "pogo", fmt.Errorf("builtin functions cannot be defer'ed"))
			default:
				emitCall(false, false, true, register, instruction.(*ssa.Defer).Call, errorInfo, comment)
			}
		}

	case *ssa.Return:
		emitPhiFlag = false
		r := ""
		switch len(operands) {
		case 0:
			r = LanguageList[l].Ret0()
		case 1:
			r = LanguageList[l].Ret1(*operands[0], errorInfo)
		default:
			r = LanguageList[l].RetN(operands, errorInfo)
		}
		fmt.Fprintln(&LanguageList[l].buffer, r+LanguageList[l].Comment(comment))

	case *ssa.Panic:
		emitPhiFlag = false
		fmt.Fprintln(&LanguageList[l].buffer, LanguageList[l].Panic(*operands[0], errorInfo)+LanguageList[l].Comment(comment))

	case *ssa.UnOp:
		if register == "" && instruction.(*ssa.UnOp).Op.String() != "<-" {
			emitComment(comment)
		} else {
			fmt.Fprintln(&LanguageList[l].buffer,
				LanguageList[l].UnOp(register, instruction.(*ssa.UnOp).Op.String(), *operands[0],
					instruction.(*ssa.UnOp).CommaOk, errorInfo)+
					LanguageList[l].Comment(comment))
		}

	case *ssa.BinOp:
		if register == "" {
			emitComment(comment)
		} else {
			op := instruction.(*ssa.BinOp).Op.String()
			fmt.Fprintln(&LanguageList[l].buffer,
				LanguageList[l].BinOp(register, op, *operands[0], *operands[1], errorInfo)+
					LanguageList[l].Comment(comment))
		}

	case *ssa.Store:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].Store(*operands[0], *operands[1], errorInfo)+LanguageList[l].Comment(comment))

	case *ssa.Send:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].Send(*operands[0], *operands[1], errorInfo)+LanguageList[l].Comment(comment))

	case *ssa.Convert:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].Convert(register, LanguageList[l].LangType(instrVal.Type(), false, errorInfo), instrVal.Type(), *operands[0], errorInfo)+
				LanguageList[l].Comment(comment))

	case *ssa.ChangeType:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].ChangeType(register, instruction.(ssa.Value).Type(), *operands[0], errorInfo)+
				LanguageList[l].Comment(comment))

	case *ssa.MakeInterface:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].MakeInterface(register, instruction.(ssa.Value).Type(), *operands[0], errorInfo)+
				LanguageList[l].Comment(comment))

	case *ssa.ChangeInterface:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].ChangeInterface(register, instruction.(ssa.Value).Type(), *operands[0], errorInfo)+
				LanguageList[l].Comment(comment))

	case *ssa.TypeAssert:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].TypeAssert(register, instruction.(*ssa.TypeAssert).X,
				instruction.(*ssa.TypeAssert).AssertedType, instruction.(*ssa.TypeAssert).CommaOk, errorInfo)+
				LanguageList[l].Comment(comment))

	case *ssa.RunDefers:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].RunDefers()+
				LanguageList[l].Comment(comment))

	case *ssa.Alloc:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].Alloc(register,
				instruction.(*ssa.Alloc).Type() /*was+: .(*types.Pointer).Elem() */, errorInfo)+
				LanguageList[l].Comment(comment))

	case *ssa.MakeClosure:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].MakeClosure(register,
				instruction,
				errorInfo)+
				LanguageList[l].Comment(comment))

	case *ssa.MakeSlice:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].MakeSlice(register,
				instruction,
				errorInfo)+
				LanguageList[l].Comment(comment))

	case *ssa.MakeChan:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].MakeChan(register,
				instruction,
				errorInfo)+
				LanguageList[l].Comment(comment))

	case *ssa.MakeMap:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].MakeMap(register,
				instruction,
				errorInfo)+
				LanguageList[l].Comment(comment))

	case *ssa.MapUpdate:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].MapUpdate(*operands[0], *operands[1], *operands[2], errorInfo)+LanguageList[l].Comment(comment))

	case *ssa.Range:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].Range(register, *operands[0], errorInfo)+LanguageList[l].Comment(comment))

	case *ssa.Next:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].Next(register, *operands[0], instruction.(*ssa.Next).IsString,
				errorInfo)+LanguageList[l].Comment(comment))

	case *ssa.Lookup:
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].Lookup(register, *operands[0], *operands[1], instruction.(*ssa.Lookup).CommaOk, errorInfo)+
				LanguageList[l].Comment(comment))

	case *ssa.Extract:
		if register == "" { // rquired here because of a "feature" in the generated SSA form
			emitComment(comment)
		} else {
			fmt.Fprintln(&LanguageList[l].buffer,
				LanguageList[l].Extract(register, *operands[0], instruction.(*ssa.Extract).Index, errorInfo)+
					LanguageList[l].Comment(comment))
		}

	case *ssa.Slice:
		// TODO see http://tip.golang.org/doc/go1.2#three_index
		// TODO add third parameter when SSA code provides it to enable slice instructions to specify a capacity
		if register == "" {
			emitComment(comment)
		} else {
			fmt.Fprintln(&LanguageList[l].buffer,
				LanguageList[l].Slice(register, instruction.(*ssa.Slice).X,
					instruction.(*ssa.Slice).Low, instruction.(*ssa.Slice).High, errorInfo)+
					LanguageList[l].Comment(comment))

		}

	case *ssa.Index:
		if register == "" {
			emitComment(comment)
		} else {
			doRangeCheck := true
			aLen := 0
			switch instruction.(*ssa.Index).X.Type().(type) {
			case *types.Array:
				aLen = int(instruction.(*ssa.Index).X.Type().(*types.Array).Len())
			case *types.Pointer:
				switch instruction.(*ssa.Index).X.Type().(*types.Pointer).Elem().(type) {
				case *types.Array:
					aLen = int(instruction.(*ssa.Index).X.Type().(*types.Pointer).Elem().(*types.Array).Len())
				}
			}
			if aLen > 0 {
				_, indexIsConst := instruction.(*ssa.Index).Index.(*ssa.Const)
				if indexIsConst {
					// this error handling is defensive, as the Go SSA code catches this error
					index := instruction.(*ssa.Index).Index.(*ssa.Const).Int64()
					if (index < 0) || (index >= int64(aLen)) {
						LogError(errorInfo, "pogo", fmt.Errorf("index [%d] out of range: 0 to %d", index, aLen-1))
					}
					doRangeCheck = false
				}
			}
			if doRangeCheck {
				fmt.Fprintln(&LanguageList[l].buffer,
					LanguageList[l].RangeCheck(instruction.(*ssa.Index).X, instruction.(*ssa.Index).Index, aLen, errorInfo))
			}
			fmt.Fprintln(&LanguageList[l].buffer,
				LanguageList[l].Index(register, *operands[0], *operands[1], errorInfo)+
					LanguageList[l].Comment(comment))
		}

	case *ssa.IndexAddr:
		if register == "" {
			emitComment(comment)
		} else {
			doRangeCheck := true
			aLen := 0
			switch instruction.(*ssa.IndexAddr).X.Type().(type) {
			case *types.Array:
				aLen = int(instruction.(*ssa.IndexAddr).X.Type().(*types.Array).Len())
			case *types.Pointer:
				switch instruction.(*ssa.IndexAddr).X.Type().(*types.Pointer).Elem().(type) {
				case *types.Array:
					aLen = int(instruction.(*ssa.IndexAddr).X.Type().(*types.Pointer).Elem().(*types.Array).Len())
				}
			}
			if aLen > 0 {
				_, indexIsConst := instruction.(*ssa.IndexAddr).Index.(*ssa.Const)
				if indexIsConst {
					index := instruction.(*ssa.IndexAddr).Index.(*ssa.Const).Int64()
					if (index < 0) || (index >= int64(aLen)) {
						LogError(errorInfo, "pogo", fmt.Errorf("index [%d] out of range: 0 to %d", index, aLen-1))
					}
					doRangeCheck = false
				}
			}
			if doRangeCheck {
				fmt.Fprintln(&LanguageList[l].buffer,
					LanguageList[l].RangeCheck(instruction.(*ssa.IndexAddr).X, instruction.(*ssa.IndexAddr).Index, aLen, errorInfo)+
						LanguageList[l].Comment(comment+" [POINTER]"))
			}
			fmt.Fprintln(&LanguageList[l].buffer, LanguageList[l].IndexAddr(register, instruction, errorInfo),
				LanguageList[l].Comment(comment+" [POINTER]"))

		}

	case *ssa.FieldAddr:
		fmt.Fprintln(&LanguageList[l].buffer, LanguageList[l].FieldAddr(register, instruction, errorInfo),
			LanguageList[l].Comment(comment+" [POINTER]"))

	case *ssa.Field:
		if register == "" {
			emitComment(comment)
		} else { // TODO review if Haxe stops using Array<Dynamic> for struct
			st := instruction.(*ssa.Field).X.Type().Underlying().(*types.Struct)
			fName := MakeID(st.Field(instruction.(*ssa.Field).Field).Name())
			l := TargetLang
			fmt.Fprintln(&LanguageList[l].buffer,
				LanguageList[l].Field(register, instruction.(*ssa.Field).X,
					instruction.(*ssa.Field).Field, fName, errorInfo, false)+
					LanguageList[l].Comment(comment))
		}

	case *ssa.DebugRef: // TODO just generates a comment at the moment, short term the comment could include the actual Go code, long term it needs some way to link to a debugger
		fmt.Fprintln(&LanguageList[l].buffer,
			LanguageList[l].Comment(comment))

	case *ssa.Select:
		text := LanguageList[l].Select(true, register, instruction, false, errorInfo)
		fmt.Fprintln(&LanguageList[l].buffer, text+LanguageList[l].Comment(comment))

	default:
		emitComment(comment + " [NO CODE GENERATED]")
		LogError(errorInfo, "pogo", fmt.Errorf("SSA instruction not implemented: %v", reflect.TypeOf(instruction)))
	}
	if false { //TODO add instruction detail DEBUG FLAG
		for o := range operands { // this loop for the creation of comments to show what is in the instructions
			val := *operands[o]
			vip := valIsPointer(val)
			if vip {
				vipOut := showIndirectValue(val)
				emitComment(fmt.Sprintf("Op[%d].VIP: %+v", o, vipOut))
			} else {
				var ic interface{} = *operands[o]
				constVal, isConst := ic.(*ssa.Const)
				if isConst {
					emitComment(fmt.Sprintf("Op[%d]: Constant= %+v", o, constVal))
				} else {
					emitComment(fmt.Sprintf("Op[%d]: %v = %+v", o, (*operands[o]), val))
				}
			}
			// l := TargetLang
			// fmt.Fprintln(&LanguageList[l].buffer, LanguageList[l].Value(*operands[o], "TEST"))
		}
	}
	return // return value is named and set in the code above
}
