package main

import "fmt"

// Represents the instruction types in RISC-V
type InstructionType int

// RISC-V instruction types
const (
	R_TYPE InstructionType = iota // Register (R-format) instructions
	I_TYPE                        // Immediate (I-format) instructions
	S_TYPE                        // Store (S-format) instructions
	B_TYPE                        // Branch (B-format) instructions
	U_TYPE                        // Upper immediate (U-format) instructions
	J_TYPE                        // Jump (J-format) instructions
)

// Represents an instruction in RISC-V
type Instruction struct {
	format InstructionType // The format of the instruction
	opcode uint8
	funct  uint8
}

// Instruction set
var instructionSet = map[string]Instruction{
	"add":     {R_TYPE, 0b000000, 0b100000},
	"addi":    {I_TYPE, 0b001000, 0},
	"and":     {R_TYPE, 0b000000, 0b100100},
	"andi":    {I_TYPE, 0b001100, 0},
	"nor":     {R_TYPE, 0b000000, 0b100111},
	"or":      {R_TYPE, 0b000000, 0b100101},
	"ori":     {I_TYPE, 0b001101, 0},
	"sub":     {R_TYPE, 0b000000, 0b100010},
	"xor":     {R_TYPE, 0b000000, 0b100110},
	"xori":    {I_TYPE, 0b001110, 0},
	"syscall": {R_TYPE, 0b000000, 0b001100},
}

// ABI registers map
var registerMap = map[string]int{
	"zero": 0, // Hardwired zero
	"ra":  1, // Return address
	"sp":  2, // Stack pointer
	"gp":  3, // Global pointer
	"tp":  4, // Thread pointer
	"t0":  5, // Temporary registers
	"t1":  6,
	"t2":  7,
	"s0":  8, // Saved register/Frame pointer
	"s1":  9, // Saved register
	"a0":  10, // Function arguments/return values
	"a1":  11,
	"a2":  12,
	"a3":  13,
	"a4":  14,
	"a5":  15,
	"a6":  16,
	"a7":  17,
	"s2":  18, // Saved registers
	"s3":  19,
	"s4":  20,
	"s5":  21,
	"s6":  22,
	"s7":  23,
	"s8":  24,
	"s9":  25,
	"s10": 26,
	"s11": 27,
	"t3":  28, // Temporary registers
	"t4":  29,
	"t5":  30,
	"t6":  31,
}

// Represents a RISC-V assembly instruction
type AssemblyInstruction struct {
	name  string            // The instruction name
	rType *RTypeInstruction // The R-type instruction
	iType *ITypeInstruction // The I-type instruction
	jType *JTypeInstruction // The J-type instruction
}

// Represents a R-type instruction
type RTypeInstruction struct {
	opcode uint8  // The opcode of the instruction
	rs     string // The source register
	rt     string // The target register
	rd     string // The destination register
	shamt  uint8  // The shift amount
	funct  uint8  // The function code
}

// Represents an I-type instruction
type ITypeInstruction struct {
	opcode uint8  // The opcode of the instruction
	rs     string // The source register
	rt     string // The target register
	imm    int16  // The immediate value
}

// Represents a J-type instruction
type JTypeInstruction struct {
	opcode uint8 // The opcode of the instruction
	addr   int32 // The address value for jump instructions
}

// Synthesize an R-type instruction
func synthesizeRType(asm RTypeInstruction, instruction Instruction) (uint32, error) {
	source, ok := registerMap[asm.rs]
	if !ok {
		return 0, fmt.Errorf("invalid source register: %s", asm.rs)
	}

	target, ok := registerMap[asm.rt]
	if !ok {
		return 0, fmt.Errorf("invalid target register: %s", asm.rt)
	}

	destination, ok := registerMap[asm.rd]
	if !ok {
		return 0, fmt.Errorf("invalid destination register: %s", asm.rd)
	}
	encodedInstruction := uint32(instruction.opcode) << 26
	encodedInstruction |= uint32(source) << 21
	encodedInstruction |= uint32(target) << 16
	encodedInstruction |= uint32(destination) << 11
	encodedInstruction |= uint32(asm.shamt) << 6
	encodedInstruction |= uint32(asm.funct)
	return encodedInstruction, nil
}

// Synthesize an I-type instruction
func synthesizeIType(asm ITypeInstruction, instruction Instruction) (uint32, error) {
	source, ok := registerMap[asm.rs]
	if !ok {
		return 0, fmt.Errorf("invalid source register: %s", asm.rs)
	}

	target, ok := registerMap[asm.rt]
	if !ok {
		return 0, fmt.Errorf("invalid target register: %s", asm.rt)
	}
	encodedInstruction := uint32(instruction.opcode) << 26
	encodedInstruction |= uint32(source) << 21
	encodedInstruction |= uint32(target) << 16
	encodedInstruction |= (uint32(asm.imm) & 0xFFFF)
	return encodedInstruction, nil
}

// Synthesize a J-type instruction
func synthesizeJType(asm JTypeInstruction, instruction Instruction) (uint32, error) {
	encodedInstruction := uint32(instruction.opcode) << 26
	encodedInstruction |= uint32(asm.addr) & 0x3FFFFFF
	return encodedInstruction, nil
}

func synthesize(asm AssemblyInstruction) (uint32, error) {
	instruction, ok := instructionSet[asm.name]
	if !ok {
		return 0, fmt.Errorf("invalid instruction name: %s", asm.name)
	}

	switch instruction.format {
	case R_TYPE:
		encodedInstruction, err := synthesizeRType(*asm.rType, instruction)
		if err != nil {
			return 0, fmt.Errorf("error synthesizing R-type instruction: %s", err)
		} else {
			return encodedInstruction, nil
		}
	case I_TYPE:
		encodedInstruction, err := synthesizeIType(*asm.iType, instruction)
		if err != nil {
			return 0, fmt.Errorf("error synthesizing I-type instruction: %s", err)
		} else {
			return encodedInstruction, nil
		}
	case J_TYPE:
		encodedInstruction, err := synthesizeJType(*asm.jType, instruction)
		if err != nil {
			return 0, fmt.Errorf("error synthesizing J-type instruction: %s", err)
		} else {
			return encodedInstruction, nil
		}
	default:
		return 0, fmt.Errorf("unknown instruction type: %d", instruction.format)
	}
}
