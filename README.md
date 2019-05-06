# Go IR

This is a just a simple repo playing around with directly
generating LLVM IR from Go, just for learning purposes.

When this runs, it generates an LLVM IR bitcode file `goir1.bc`.

## Displaying the text representation of IR

To convert LLVM bitcode to a text representation, `llvm-link -S`
seems to work ok:

```
$ llvm-link -S -v goir1.bc
Loading 'goir1.bc'
Linking in 'goir1.bc'
Writing bitcode...
; ModuleID = 'llvm-link'
source_filename = "llvm-link"

define i32 @main() {
entry:
  %a = alloca i32
  store i32 32, i32* %a
  %b = alloca i32
  store i32 16, i32* %b
  %a_val = load i32, i32* %a
  %b_val = load i32, i32* %b
  %ab_value = add i32 %a_val, %b_val
  ret i32 %ab_value
}

```

## Executing the LLVM bitcode

The generated bitcode can be run directly from the command
line using `lli`, without needing to turn it into an executable.

For example:

```
$ lli goir1.bc 
$ echo $?
48
```

The return code of 48 there is correct, as the bitcode in this
playground example returns a value of 48 to the caller from its
`main()`.

## Compiling the LLVM bitcode to an executable

`clang` can directly compile the bitcode to a runnable
executable:

```
$ clang -o goir1 goir1.bc
warning: overriding the module target triple with x86_64-unknown-linux-gnu [-Woverride-module]
1 warning generated.
$ ls -la goir1
-rwxrwxr-x. 1 jc jc 8288 May  6 21:28 goir1
$ ./goir1
$ echo $?
48
```

## Converting the LLVM bitcode to SSA form

To convert the bitcode to LLVM SSA form, `llc` seems to work ok:

```
$ llc goir1.bc
$ cat goir1.s 
        .text
        .file   "my_module"
        .globl  main                    # -- Begin function main
        .p2align        4, 0x90
        .type   main,@function
main:                                   # @main
        .cfi_startproc
# %bb.0:                                # %entry
        movl    $32, -4(%rsp)
        movl    $16, -8(%rsp)
        movl    $48, %eax
        retq
.Lfunc_end0:
        .size   main, .Lfunc_end0-main
        .cfi_endproc
                                        # -- End function

        .section        ".note.GNU-stack","",@progbits
```

## Compiling the bitcode to WebAssembly

Using clang to generate WebAssembly works, with the resulting wasm
able to be processed by [wabt](https://github.com/WebAssembly/wabt).

```
$ clang --compile -Os --target=wasm32-unknown-unknown-wasm -o goir1.wasm goir1.bc
warning: overriding the module target triple with wasm32-unknown-unknown-wasm [-Woverride-module]
1 warning generated.
```

```
$ ls -la goir1.wasm
-rw-rw-r--. 1 jc jc 190 May  6 22:04 goir1.wasm
```

```
$ wasm2wat --generate-names goir1.wasm
(module
  (type $t0 (func (result i32)))
  (type $t1 (func (param i32 i32) (result i32)))
  (import "env" "__linear_memory" (memory $env.__linear_memory 0))
  (import "env" "__indirect_function_table" (table $env.__indirect_function_table 0 funcref))
  (import "env" "__stack_pointer" (global $env.__stack_pointer (mut i32)))
  (func $f0 (type $t0) (result i32)
    (local $l0 i32) (local $l1 i32) (local $l2 i32) (local $l3 i32) (local $l4 i32) (local $l5 i32) (local $l6 i32) (local $l7 i32)
    global.get $env.__stack_pointer
    local.set $l0
    i32.const 16
    local.set $l1
    local.get $l0
    local.get $l1
    i32.sub
    local.set $l2
    i32.const 16
    local.set $l3
    i32.const 32
    local.set $l4
    local.get $l2
    local.get $l4
    i32.store offset=12
    local.get $l2
    local.get $l3
    i32.store offset=8
    local.get $l2
    i32.load offset=12
    local.set $l5
    local.get $l2
    i32.load offset=8
    local.set $l6
    local.get $l5
    local.get $l6
    i32.add
    local.set $l7
    local.get $l7
    return)
  (func $f1 (type $t1) (param $p0 i32) (param $p1 i32) (result i32)
    (local $l2 i32)
    call $f0
    local.set $l2
    local.get $l2
    return))
```