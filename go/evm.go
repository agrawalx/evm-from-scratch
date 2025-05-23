// Package evm is an **incomplete** implementation of the Ethereum Virtual
// Machine for the "EVM From Scratch" course:
// https://github.com/w1nt3r-eth/evm-from-scratch
//
// To work on EVM From Scratch In Go:
//
// - Install Golang: https://golang.org/doc/install
// - Go to the `go` directory: `cd go`
// - Edit `evm.go` (this file!), see TODO below
// - Run `go test ./...` to run the tests
package evm

import (
	"math/big"
)

// Run runs the EVM code and returns the stack and a success indicator.
func Evm(code []byte) ([]*big.Int, bool) {
	var stack []*big.Int
	pc := 0

	for pc < len(code) {
		op := code[pc]
		pc++

		// TODO: Implement the EVM here!
		if op >= 0x60 && op <= 0x7f { //PUSH1 To PUSH32 	
			pushsize := int(op - 0x5f)
			valueBytes := code[pc:pc +pushsize]
			pc += pushsize
			value := new(big.Int).SetBytes(valueBytes)
			stack = push(stack, value)
			continue 
		}
		switch op {
		case 0x00: // STOP 
			return stack, true
		case 0x5f: // PUSH0 
			stack = push(stack, big.NewInt(0))
		case 0x50: //POP 
			stack = pop(stack)
		case 0x01: //ADD 
			stack = add(stack)
		case 0x02: //MULTIPLY 
			stack = mul(stack)
		case 0x03: //SUBTRACT
			n := len(stack)
			a,b := stack[n-1], stack[n-2]
			result := wrap(new(big.Int).Sub(a,b))
			stack = pop2AndPush(stack, result)
		case 0x04: //DIVISION 
			n := len(stack)
			a,b := stack[n-1], stack[n-2]
			if b.Cmp(big.NewInt(0)) == 0 {
				stack = pop2AndPush(stack, big.NewInt(0))
				continue 
			}
			result := wrap(new(big.Int).Div(a,b))
			stack = pop2AndPush(stack, result)
		case 0x06: //MOD 
			n := len(stack)
			a,modulo := stack[n-1], stack[n-2]
			result := mod(a, modulo)
			stack = pop2AndPush(stack, result)
		case 0x08: //ADDMOD
			n := len(stack)
			a , b , c := stack[n-1], stack[n-2], stack[n-3] 
			add := wrap(new(big.Int).Add(a,b))
			result := mod(add, c) 
			for i := 0;i < 3;i++ {
				stack = pop(stack)
			}
			stack = push(stack, result)
		case 0x09: //MULMOD 
			n := len(stack)
			a , b , c := stack[n-1], stack[n-2], stack[n-3] 
			mul := new(big.Int).Mul(a,b)
			result := wrap(mod(mul, c))
			for i := 0;i < 3;i++ {
				stack = pop(stack)
			}
			stack = push(stack, result)
		case 0x0a: //EXP 
			n := len(stack)
			a, power := stack[n-1] , stack[n-2]
			result := new(big.Int).Exp(a,power,nil)
			stack = pop2AndPush(stack, result)
		case 0x0b: //SIGNEXTEND
			n := len(stack)
			value , byteIndex := stack[n-2] , stack[n-1]
			idx := int(byteIndex.Uint64())
			padded := make([]byte, 32)
        	value.FillBytes(padded)
			signByte := padded[31-idx]
			signBit := (signByte >> 7) & 1
			mask := new(big.Int).Lsh(big.NewInt(1), uint((idx+1)*8))
			mask.Sub(mask, big.NewInt(1))

			var result *big.Int
			if signBit == 0 {
				// Positive number, clear upper bits
				result = new(big.Int).And(value, mask)
			} else {
				// Negative number, sign-extend upper bits
				negOne := new(big.Int).Sub(mask, big.NewInt(1))
				result = new(big.Int).Or(value, new(big.Int).Not(negOne))
			}

			// Mask result to 256 bits unsigned representation
			mask256 := new(big.Int).Lsh(big.NewInt(1), 256)
			mask256.Sub(mask256, big.NewInt(1))
			result.And(result, mask256)

			stack = pop2AndPush(stack, result)
			
		case 0x05: //SDIV 
			n := len(stack)
			num , den := stack[n-1] , stack[n-2]
			if den.Cmp(big.NewInt(0)) == 0 {
				stack = pop2AndPush(stack, big.NewInt(0))
				break
			}
			signedNum := toSigned(num)
			signedDen := toSigned(den)
			result := wrap(new(big.Int).Div(signedNum,signedDen))
			stack = pop2AndPush(stack,result)
		case 0x07: //SMOD 
			n := len(stack)
			num, mod := stack[n-1] , stack[n-2]
			if mod.Cmp(big.NewInt(0)) == 0 {
				stack = pop2AndPush(stack, big.NewInt(0))
				break
			}
			signedNum := toSigned(num)
			signedMod := toSigned(mod)
			absNum := new(big.Int).Abs(signedNum)
			absMod := new(big.Int).Abs(signedMod)

			remainder := new(big.Int).Mod(absNum, absMod)

			if signedNum.Sign() < 0 {
				remainder.Neg(remainder)
			}

			result := wrap(remainder)
			stack = pop2AndPush(stack, result)
		case 0x10: 
			n := len(stack)
			left , right := stack[n-1] , stack[n-2]
			result := new(big.Int)
			if left.Cmp(right) < 0 {
				result = big.NewInt(1) 
			} else {
				result = big.NewInt(0)
			}
			stack = pop2AndPush(stack, result)
		case 0x11: 
			n:= len(stack)
			left , right := stack[n-1] , stack[n-2]
			result := new(big.Int)
			if right.Cmp(left) < 0 {
				result = big.NewInt(1) 
			} else {
				result = big.NewInt(0)
			}
			stack = pop2AndPush(stack, result)
		case 0x12: 
			n:= len(stack)
			left , right := toSigned(stack[n-1]) , toSigned(stack[n-2])
			result := new(big.Int)
			if left.Cmp(right) < 0 {
				result = big.NewInt(1) 
			} else {
				result = big.NewInt(0)
			}
			stack = pop2AndPush(stack, result)
		case 0x13: 
			n:= len(stack)
			left , right := toSigned(stack[n-1]) , toSigned(stack[n-2])
			result := new(big.Int)
			if right.Cmp(left) < 0 {
				result = big.NewInt(1) 
			} else {
				result = big.NewInt(0)
			}
			stack = pop2AndPush(stack, result)
		case 0x14: 
			n:= len(stack)
			left , right := stack[n-1] , stack[n-2]
			result := new(big.Int)
			if right.Cmp(left) == 0 {
				result = big.NewInt(1) 
			} else {
				result = big.NewInt(0)
			}
			stack = pop2AndPush(stack, result)
		}
		
	}
	return reverse(stack), true
}

func push(stack []*big.Int, num *big.Int) ([]*big.Int) {
	stack = append(stack, num)
	return stack 
}
func reverse(stack []*big.Int) ([]*big.Int) {
	var reversed_stack []*big.Int 
	n := len(stack)
	for i := 0; i < n; i++ {
		reversed_stack = append(reversed_stack, stack[n-i-1])
	}
	return reversed_stack
}

func pop(stack []*big.Int) ([]*big.Int) {
	stack = stack[:len(stack)-1]
	return stack
}
func add(stack []*big.Int) ([]*big.Int) {
	n := len(stack) 
	result := wrap(new(big.Int).Add(stack[n-1],stack[n-2]))
	stack = pop(stack)
	stack = pop(stack)
	stack = push(stack, result)
	return stack 
}
func mul(stack []*big.Int) ([]*big.Int) {
	n := len(stack) 
	result := wrap(new(big.Int).Mul(stack[n-1],stack[n-2]))
	stack = pop(stack)
	stack = pop(stack)
	stack = push(stack, result)
	return stack 
}
func wrap(num *big.Int) (*big.Int) {
	modulo := new(big.Int).Lsh(big.NewInt(1), 256)
	result := num.Mod(num, modulo)
	return result 
}

func pop2AndPush(stack []*big.Int, num *big.Int) ([]*big.Int){
	stack = pop(stack)
	stack = pop(stack)
	stack = push(stack,num)
	return stack
}

func mod(a *big.Int, modulo *big.Int) *big.Int {
	if modulo.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0)
	}
	return new(big.Int).Mod(a, modulo)
}

func toSigned(x *big.Int) *big.Int {
	if x.Bit(255) == 1 {
		// If MSB is set, it's negative
		return new(big.Int).Sub(x, new(big.Int).Lsh(big.NewInt(1), 256))
	}
	return new(big.Int).Set(x)
}
