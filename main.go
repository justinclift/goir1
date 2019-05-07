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
	i32NoParams := llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false)

	// Define the function type for the extern "puts" function
	putsType := llvm.FunctionType(llvm.Int32Type(), []llvm.Type{llvm.PointerType(llvm.Int8Type(), 0)}, false)

	// Declare a type that returns a double, and takes no parameters
	//DoubleNoParams := llvm.FunctionType(llvm.DoubleType(), []llvm.Type{}, false)

	// Create a function called "main"
	llvm.AddFunction(mod, "main", i32NoParams)

	//llvm.ExternalLinkage

	// Declare a type that returns a double, and takes 1 double as a parameter
	//doubleDouble := llvm.FunctionType(llvm.DoubleType(), []llvm.Type{llvm.DoubleType()}, false)

	// Import a global.  Trying with the external cos() function for now
	//llvm.AddGlobal(mod, doubleDouble, "cos")
	//cos := llvm.AddGlobal(mod, uInt32NoParams, "cos")

	// Add a global for the external puts() function
	llvm.AddGlobal(mod, putsType, "puts")

	// Create a basic block
	block := llvm.AddBasicBlock(mod.NamedFunction("main"), "entry")

	// Set the instruction insert point
	builder.SetInsertPoint(block, block.FirstInstruction())

	// Add "hello world" string
	builder.CreateGlobalString("hello world\n", ".str")

	// int a = 32
	a := builder.CreateAlloca(llvm.Int32Type(), "a")
	builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 32, false), a)

	// int b = 16
	b := builder.CreateAlloca(llvm.Int32Type(), "b")
	builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 16, false), b)

	// a + b
	aVal := builder.CreateLoad(a, "a_val")
	bVal := builder.CreateLoad(b, "b_val")

	// cos (a + b)
	//c := builder.CreateAlloca(llvm.DoubleType(), "cosresult")
	//c := builder.CreateAdd(aVal, bVal, "ab_value")
	result := builder.CreateAdd(aVal, bVal, "ab_value")

	//foo := llvm.UIToFP
	builder.CreateUIToFP(result, llvm.DoubleType(), "convertint")
	//c := builder.CreateUIToFP(result, llvm.DoubleType(), "convertint")

	//result := builder.

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
