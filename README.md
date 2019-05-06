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

Using clang to generate WebAssembly seems like it works,
though even with `-Os` the resulting file seems to have
several text strings included.

Unsure yet if this is really a valid wasm file, and if those
text strings could be omitted or reduced.

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
$ hexdump -C goir1.wasm
00000000  00 61 73 6d 01 00 00 00  01 8b 80 80 80 00 02 60  |.asm...........`|
00000010  00 01 7f 60 02 7f 7f 01  7f 02 ba 80 80 80 00 02  |...`............|
00000020  03 65 6e 76 0f 5f 5f 6c  69 6e 65 61 72 5f 6d 65  |.env.__linear_me|
00000030  6d 6f 72 79 02 00 00 03  65 6e 76 19 5f 5f 69 6e  |mory....env.__in|
00000040  64 69 72 65 63 74 5f 66  75 6e 63 74 69 6f 6e 5f  |direct_function_|
00000050  74 61 62 6c 65 01 70 00  00 03 83 80 80 80 00 02  |table.p.........|
00000060  00 01 0a 8f 80 80 80 00  02 04 00 41 30 0b 08 00  |...........A0...|
00000070  10 80 80 80 80 00 0b 00  ab 80 80 80 00 07 6c 69  |..............li|
00000080  6e 6b 69 6e 67 02 08 9c  80 80 80 00 02 00 00 00  |nking...........|
00000090  0f 5f 5f 6f 72 69 67 69  6e 61 6c 5f 6d 61 69 6e  |.__original_main|
000000a0  00 00 01 04 6d 61 69 6e  00 90 80 80 80 00 0a 72  |....main.......r|
000000b0  65 6c 6f 63 2e 43 4f 44  45 03 01 00 09 00        |eloc.CODE.....|
000000be
```