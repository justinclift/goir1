package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aykevl/go-llvm"
)

func main() {
	// Create an overall context
	ctx := llvm.NewContext()

	// Create a LLVM IR builder
	mod := ctx.NewModule("my_module")
	builder := ctx.NewBuilder()

	// Create the "nocapture" attribute
	// From: https://github.com/tinygo-org/tinygo/blob/fa5df4f524c3b4f15360c2da964932692e3ab4af/compiler/compiler.go#L374-L378
	getAttr := func(attrName string) llvm.Attribute {
		attrKind := llvm.AttributeKindID(attrName)
		return ctx.CreateEnumAttribute(attrKind, 0)
	}
	nocapture := getAttr("nocapture")
fmt.Printf("%v\n", nocapture)

	// Declare a type that returns an int32, and takes no parameters
	i32NoParams := llvm.FunctionType(ctx.Int32Type(), []llvm.Type{}, false)

	// Define the function type for the extern "puts" function
	puts1 := llvm.PointerType(ctx.Int8Type(), 0)
	putsType := llvm.FunctionType(ctx.Int32Type(), []llvm.Type{puts1}, false)

	// Create a function called "main"
	llvm.AddFunction(mod, "main", i32NoParams)

	// Add a global for the external puts() function
	llvm.AddGlobal(mod, putsType, "puts")
	//mod.NamedFunction("puts").AddAttributeAtIndex(0, nocapture)
	//putsGlobal := llvm.AddGlobal(mod, putsType, "puts")
	//putsGlobal.AddAttributeAtIndex(0, nocapture)


	// Create a basic block
	block := ctx.AddBasicBlock(mod.NamedFunction("main"), "entry")

	// Set the instruction insert point
	builder.SetInsertPoint(block, block.FirstInstruction())

	// Add "hello world" string
	builder.CreateGlobalString("hello world\n", ".str")

	// int a = 32
	a := builder.CreateAlloca(ctx.Int32Type(), "a")
	builder.CreateStore(llvm.ConstInt(ctx.Int32Type(), 32, false), a)

	// int b = 16
	b := builder.CreateAlloca(ctx.Int32Type(), "b")
	builder.CreateStore(llvm.ConstInt(ctx.Int32Type(), 16, false), b)

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
