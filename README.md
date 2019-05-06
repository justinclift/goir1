# Go IR

This is a just a simple repo playing around with directly
generating LLVM IR from Go, just for learning purposes.

When this runs, it generates an LLVM IR bitcode file `goir1.bc`.

To convert that to a text representation, `llvm-link -S` seems
to work ok:

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

To convert the bitcode to LLVM SSA form, `llc` seems to work
ok:

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
