package main

import (
	"log"
	"os"

	"github.com/aykevl/go-llvm"
)

func main() {
	// Create a LLVM IR builder
	builder := llvm.NewBuilder()
	mod := llvm.NewModule("my_module")

	// Create a function called "main"
	main := llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false)
	llvm.AddFunction(mod, "main", main)

	// Create a basic block
	block := llvm.AddBasicBlock(mod.NamedFunction("main"), "entry")

	// Set the instruction insert point
	builder.SetInsertPoint(block, block.FirstInstruction())

	// int a = 32
	a := builder.CreateAlloca(llvm.Int32Type(), "a")
	builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 32, false), a)

	// int b = 16
	b := builder.CreateAlloca(llvm.Int32Type(), "b")
	builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 16, false), b)

	// a + b
	aVal := builder.CreateLoad(a, "a_val")
	bVal := builder.CreateLoad(b, "b_val")
	result := builder.CreateAdd(aVal, bVal, "ab_value")

	// Return
	builder.CreateRet(result)

	// Write out the IR as bitcode
	outFile, err := os.Create("goir1.bc")
	if err != nil {
		log.Fatal(err)
	}
	err = llvm.WriteBitcodeToFile(mod, outFile)
	if err != nil {
		log.Fatal(err)
	}
}
