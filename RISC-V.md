# RISC-V Assembly Language Reference Guide

## Core Instruction Formats
RISC-V uses fixed-length 32-bit instructions with these basic formats:
- R-type: Register-to-register operations
- I-type: Immediate and load operations
- S-type: Store operations
- B-type: Branch operations
- U-type: Upper immediate operations
- J-type: Jump operations

## Base Integer Instruction Set (RV32I)

### Arithmetic Operations
- `add rd, rs1, rs2`    # rd = rs1 + rs2
- `sub rd, rs1, rs2`    # rd = rs1 - rs2
- `addi rd, rs1, imm`   # rd = rs1 + immediate
- `lui rd, imm`         # rd = immediate << 12
- `auipc rd, imm`       # rd = pc + (immediate << 12)

### Logical Operations
- `and rd, rs1, rs2`    # Bitwise AND
- `or rd, rs1, rs2`     # Bitwise OR
- `xor rd, rs1, rs2`    # Bitwise XOR
- `andi rd, rs1, imm`   # Immediate AND
- `ori rd, rs1, imm`    # Immediate OR
- `xori rd, rs1, imm`   # Immediate XOR

### Shift Operations
- `sll rd, rs1, rs2`    # Logical left shift
- `srl rd, rs1, rs2`    # Logical right shift
- `sra rd, rs1, rs2`    # Arithmetic right shift
- `slli rd, rs1, imm`   # Immediate logical left shift
- `srli rd, rs1, imm`   # Immediate logical right shift
- `srai rd, rs1, imm`   # Immediate arithmetic right shift

### Comparison Operations
- `slt rd, rs1, rs2`    # Set if less than (signed)
- `sltu rd, rs1, rs2`   # Set if less than (unsigned)
- `slti rd, rs1, imm`   # Set if less than immediate (signed)
- `sltiu rd, rs1, imm`  # Set if less than immediate (unsigned)

### Memory Operations
- `lb rd, offset(rs1)`  # Load byte
- `lh rd, offset(rs1)`  # Load halfword
- `lw rd, offset(rs1)`  # Load word
- `lbu rd, offset(rs1)` # Load byte unsigned
- `lhu rd, offset(rs1)` # Load halfword unsigned
- `sb rs2, offset(rs1)` # Store byte
- `sh rs2, offset(rs1)` # Store halfword
- `sw rs2, offset(rs1)` # Store word

### Control Flow
- `beq rs1, rs2, offset`  # Branch if equal
- `bne rs1, rs2, offset`  # Branch if not equal
- `blt rs1, rs2, offset`  # Branch if less than
- `bge rs1, rs2, offset`  # Branch if greater or equal
- `bltu rs1, rs2, offset` # Branch if less than unsigned
- `bgeu rs1, rs2, offset` # Branch if greater or equal unsigned
- `jal rd, offset`        # Jump and link
- `jalr rd, offset(rs1)`  # Jump and link register

### System/Special Instructions
- `fence`               # Memory and I/O fence
- `ecall`              # Environment call
- `ebreak`             # Environment break

## Standard Extensions

### M Extension (Integer Multiplication and Division)
- `mul rd, rs1, rs2`     # Multiplication
- `mulh rd, rs1, rs2`    # Multiplication (high bits, signed)
- `mulhu rd, rs1, rs2`   # Multiplication (high bits, unsigned)
- `mulhsu rd, rs1, rs2`  # Multiplication (high bits, signed*unsigned)
- `div rd, rs1, rs2`     # Division
- `divu rd, rs1, rs2`    # Division unsigned
- `rem rd, rs1, rs2`     # Remainder
- `remu rd, rs1, rs2`    # Remainder unsigned

### F Extension (Single-Precision Floating-Point)
- `flw rd, offset(rs1)`  # Float load word
- `fsw rs2, offset(rs1)` # Float store word
- `fadd.s`               # Float add
- `fsub.s`               # Float subtract
- `fmul.s`               # Float multiply
- `fdiv.s`               # Float divide
- `fsqrt.s`              # Float square root
- `fmin.s`               # Float minimum
- `fmax.s`               # Float maximum
- `fcvt.w.s`            # Float convert to integer
- `fcvt.s.w`            # Integer convert to float

### D Extension (Double-Precision Floating-Point)
- `fld rd, offset(rs1)`  # Float load double
- `fsd rs2, offset(rs1)` # Float store double
- `fadd.d`               # Double add
- `fsub.d`               # Double subtract
- `fmul.d`               # Double multiply
- `fdiv.d`               # Double divide
- `fsqrt.d`              # Double square root

## Labels and Symbol References
### Local Labels
- Numeric labels (1:, 2:, etc.) for local reference
- Referenced with 1f (forward) or 1b (backward)
- Scope limited to current function
```
1:
    addi t0, t0, 1
    bne t0, t1, 1b    # Branch back to label 1
```

### Global Labels
- Named labels visible throughout the program
- Can be referenced by other files when made global
```
function_name:
    addi sp, sp, -16
    sw ra, 12(sp)
```

### Symbol References
- Can reference data section symbols
- PC-relative addressing for position-independent code
```
    la t0, data_symbol    # Load address of symbol
    lw t1, data_symbol    # Load value at symbol
```

## Macros
### Basic Macro Definition
```
.macro name arg1, arg2
    # Macro body using \arg1, \arg2
.endm
```

### Macro Examples
```
.macro push reg
    addi sp, sp, -4
    sw \reg, 0(sp)
.endm

.macro pop reg
    lw \reg, 0(sp)
    addi sp, sp, 4
.endm

.macro load_const reg, const
    li \reg, \const
.endm
```

### Advanced Macro Features
- Conditional assembly within macros
- Macro parameter default values
- Macro string concatenation
```
.macro conditional_move rd, rs, rt, cond
    b\cond 1f
    mv \rd, \rs
    j 2f
1:  mv \rd, \rt
2:
.endm
```

## Preprocessor Directives
### Basic Directives
- `.equ SYMBOL, expression`   # Define constant
- `.set SYMBOL, expression`   # Define/redefine symbol
- `.ifdef SYMBOL`            # Conditional assembly if symbol defined
- `.ifndef SYMBOL`           # Conditional assembly if symbol not defined
- `.endif`                   # End conditional block

### Include Directives
- `.include "filename"`      # Include source file
- `.incbin "filename"`      # Include binary file

### Conditional Assembly
```
.ifdef DEBUG
    # Debug code here
.else
    # Release code here
.endif
```

### Constants and Expressions
```
.equ MAX_SIZE, 1024
.equ OFFSET, MAX_SIZE - 4

.if MAX_SIZE > 2048
    .error "MAX_SIZE too large"
.endif
```

## Assembler Directives
- `.text`               # Code section
- `.data`               # Data section
- `.global symbol`      # Export symbol
- `.extern symbol`      # Import symbol
- `.byte`               # Define bytes
- `.half`               # Define halfwords
- `.word`               # Define words
- `.dword`              # Define doublewords
- `.string`             # Define string
- `.align n`            # Align to 2^n boundary
- `.section name`       # Define custom section
- `.file "filename"`    # Source file name for debugging
- `.loc file line`      # Source line number for debugging

## Registers
### Integer Registers (x0-x31)
- `x0/zero`: Hardwired zero
- `x1/ra`: Return address
- `x2/sp`: Stack pointer
- `x3/gp`: Global pointer
- `x4/tp`: Thread pointer
- `x5-x7/t0-t2`: Temporaries
- `x8-x9/s0-s1`: Saved registers
- `x10-x11/a0-a1`: Function arguments/returns
- `x12-x17/a2-a7`: Function arguments
- `x18-x27/s2-s11`: Saved registers
- `x28-x31/t3-t6`: Temporaries

### Floating-Point Registers (f0-f31)
- `f0-f7/ft0-ft7`: FP temporaries
- `f8-f9/fs0-fs1`: FP saved registers
- `f10-f11/fa0-fa1`: FP arguments/returns
- `f12-f17/fa2-fa7`: FP arguments
- `f18-f27/fs2-fs11`: FP saved registers
- `f28-f31/ft8-ft11`: FP temporaries
