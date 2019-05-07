package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aykevl/go-llvm"
)

const (
	PAGESIZE = 65536
)

func main() {
	// Create a LLVM IR builder
	builder := llvm.NewBuilder()
	mod := llvm.NewModule("my_module")

	// Declare a type that returns an int32, and takes no parameters
	int32NoParams := llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false)

	// Create a function called "main"
	llvm.AddFunction(mod, "main", int32NoParams)

	//llvm.ExternalLinkage

	// Declare a type that returns a float, and takes 1 float as a parameter
	floatFloat := llvm.FunctionType(llvm.FloatType(), []llvm.Type{llvm.FloatType()}, false)

	// Import a global.  Trying with the external cos() function for now
	llvm.AddGlobal(mod, floatFloat, "cos")
	//cos := llvm.AddGlobal(mod, uInt32NoParams, "cos")

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

	// Verify the module is correct
	if ok := llvm.VerifyModule(mod, llvm.ReturnStatusAction); ok != nil {
		log.Fatal(ok)
	}

	// Write the IR for the module (text format) to stdout
	mod.Dump()

	// Compile and run the function
	engine, err := llvm.NewExecutionEngine(mod)
	if err != nil {
		log.Fatal(err)
	}
	funcResult := engine.RunFunction(mod.NamedFunction("main"), []llvm.GenericValue{})

	// Display the result of the function
	fmt.Printf("%d\n", funcResult.Int(false))

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
