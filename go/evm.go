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
			stack = pop(stack)
			stack = pop(stack)
			stack = push(stack, result)
		case 0x04: //DIVISION 
			n := len(stack)
			a,b := stack[n-1], stack[n-2]
			if b.Cmp(big.NewInt(0)) == 0 {
				stack = pop(stack)
				stack = pop(stack)
				stack = push(stack, big.NewInt(0))
				continue 
			}
			result := wrap(new(big.Int).Div(a,b))
			stack = pop(stack)
			stack = pop(stack)
			stack = push(stack,result)

	}}

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
