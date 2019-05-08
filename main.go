package main

import (
	"log"
	"os"

	"github.com/aykevl/go-llvm"
)

var (
	builder llvm.Builder
	ctx     llvm.Context
	mod     llvm.Module
)

func init() {
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()
}

func main() {
	// Set up target machine
	cfg := llvm.DefaultTargetTriple()
	tgt, err := llvm.GetTargetFromTriple(cfg)
	if err != nil {
		log.Fatal(err)
	}
	machine := tgt.CreateTargetMachine(cfg, "generic", "", llvm.CodeGenLevelDefault, llvm.RelocStatic, llvm.CodeModelDefault)
	targetData := machine.CreateTargetData()

	// Create an overall context
	ctx = llvm.NewContext()

	// Create an LLVM IR builder
	mod = ctx.NewModule("")
	mod.SetTarget(cfg)
	mod.SetDataLayout(targetData.String())
	builder = ctx.NewBuilder()

	// Create the "nocapture" attribute
	// From: https://github.com/tinygo-org/tinygo/blob/fa5df4f524c3b4f15360c2da964932692e3ab4af/compiler/compiler.go#L374-L378
	getAttr := func(attrName string) llvm.Attribute {
		attrKind := llvm.AttributeKindID(attrName)
		return ctx.CreateEnumAttribute(attrKind, 0)
	}
	nocapture := getAttr("nocapture")

	// Define the function type for the extern "puts" function
	puts1 := llvm.PointerType(ctx.Int8Type(), 0)
	putsType := llvm.FunctionType(ctx.Int32Type(), []llvm.Type{puts1}, false)

	// Add a global for the external puts() function
	llvm.AddFunction(mod, "puts", putsType)
	mod.NamedFunction("puts").AddAttributeAtIndex(1, nocapture)

	// Declare a type that returns an int32, and takes no parameters
	i32NoParams := llvm.FunctionType(ctx.Int32Type(), []llvm.Type{}, false)

	// Create a function called "main"
	llvm.AddFunction(mod, "main", i32NoParams)

	// Create a basic block and set the instruction insert point
	block := ctx.AddBasicBlock(mod.NamedFunction("main"), "")
	builder.SetInsertPointAtEnd(block)

	// Print "hello world" strings
	printString("hello world")
	printString("hello world2")

	// Return 0 from the main function
	builder.CreateRet(llvm.ConstInt(ctx.Int32Type(), 0, false))

	// Verify the module is correct
	if ok := llvm.VerifyModule(mod, llvm.ReturnStatusAction); ok != nil {
		log.Fatal(ok)
	}

	// Write the IR for the module (text format) to stdout
	mod.Dump()

	// Write out the IR as bitcode
	fName := "goir1.bc"
	outFile, err := os.Create(fName)
	if err != nil {
		log.Fatal(err)
	}
	err = llvm.WriteBitcodeToFile(mod, outFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("********")
	log.Printf("Bitcode written to %s\n", fName)

	// Compile and run the function
	engine, err := llvm.NewExecutionEngine(mod)
	if err != nil {
		log.Fatal(err)
	}
	funcResult := engine.RunFunction(mod.NamedFunction("main"), []llvm.GenericValue{})

	// Display the result of the function
	log.Printf("Program compiled and executed, returning exit code: %d\n", funcResult.Int(false))
	log.Println("********")
}

func printString(text string) {
	// Add the string
	str := builder.CreateGlobalString(text, "")

	// Call the external puts function
	strPtr := builder.CreatePointerCast(str, llvm.PointerType(ctx.Int8Type(), 0), "")
	builder.CreateCall(mod.NamedFunction("puts"), []llvm.Value{strPtr}, "")
}