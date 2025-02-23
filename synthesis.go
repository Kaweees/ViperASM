package main

import "fmt"

// Represents the instruction types in RISC-V
type InstructionType int

// RISC-V instruction types
const (
	R_TYPE InstructionType = iota // Register-Register (R-format) instructions
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
	funct3 uint8
	funct7 uint16
}

// Instruction set
var instructionSet = map[string]Instruction{
	"add":  {R_TYPE, 0b0110011, 0x0, 0x00}, // rd = rs1 + rs2
	"sub":  {R_TYPE, 0b0110011, 0x0, 0x20}, // rd = rs1 - rs2
	"xor":  {R_TYPE, 0b0110011, 0x4, 0x00}, // rd = rs1 ^ rs2
	"or":   {R_TYPE, 0b0110011, 0x6, 0x00}, // rd = rs1 | rs2
	"and":  {R_TYPE, 0b0110011, 0x7, 0x00}, // rd = rs1 & rs2
	"sll":  {R_TYPE, 0b0110011, 0x1, 0x00}, // rd = rs1 << rs2
	"srl":  {R_TYPE, 0b0110011, 0x5, 0x00}, // rd = rs1 >> rs2 (logical)
	"sra":  {R_TYPE, 0b0110011, 0x5, 0x20}, // rd = rs1 >> rs2 (arithmetic)
	"slt":  {R_TYPE, 0b0110011, 0x2, 0x00}, // rd = (rs1 < rs2) ? 1 : 0
	"sltu": {R_TYPE, 0b0110011, 0x3, 0x00}, // rd = (rs1 < rs2) ? 1 : 0

	"addi":  {I_TYPE, 0b0010011, 0x0, 0x00}, // rd = rs1 + imm
	"xori":  {I_TYPE, 0b0010011, 0x4, 0x00}, // rd = rs1 ^ imm
	"ori":   {I_TYPE, 0b0010011, 0x6, 0x00}, // rd = rs1 | imm
	"andi":  {I_TYPE, 0b0010011, 0x7, 0x00}, // rd = rs1 & imm
	"slli":  {I_TYPE, 0b0010011, 0x1, 0x00}, // rd = rs1 << imm[0:4]
	"srli":  {I_TYPE, 0b0010011, 0x5, 0x00}, // rd = rs1 >> imm[0:4] (logical)
	"srai":  {I_TYPE, 0b0010011, 0x5, 0x20}, // rd = rs1 >> imm[0:4] (arithmetic)
	"slti":  {I_TYPE, 0b0010011, 0x2, 0x00}, // rd = (rs1 < imm) ? 1 : 0
	"sltiu": {I_TYPE, 0b0010011, 0x3, 0x00}, // rd = (rs1 < imm) ? 1 : 0

	"lb":  {I_TYPE, 0b0000011, 0x0, 0x00}, // rd = M[rs1+imm][0:7]
	"lh":  {I_TYPE, 0b0000011, 0x1, 0x00}, // rd = M[rs1+imm][0:15]
	"lw":  {I_TYPE, 0b0000011, 0x2, 0x00}, // rd = M[rs1+imm][0:31]
	"lbu": {I_TYPE, 0b0000011, 0x4, 0x00}, // rd = M[rs1+imm][0:7] zero-extends
	"lhu": {I_TYPE, 0b0000011, 0x5, 0x00}, // rd = M[rs1+imm][0:15] zero-extends

	"sb": {S_TYPE, 0b0100011, 0x0, 0x00}, // M[rs1+imm][0:7] = rs2[0:7]
	"sh": {S_TYPE, 0b0100011, 0x1, 0x00}, // M[rs1+imm][0:15] = rs2[0:15]
	"sw": {S_TYPE, 0b0100011, 0x2, 0x00}, // M[rs1+imm][0:31] = rs2[0:31]

	"beq":  {B_TYPE, 0b1100011, 0x0, 0x00}, // if(rs1 == rs2) PC += imm
	"bne":  {B_TYPE, 0b1100011, 0x1, 0x00}, // if(rs1 != rs2) PC += imm
	"blt":  {B_TYPE, 0b1100011, 0x4, 0x00}, // if(rs1 < rs2) PC += imm
	"bge":  {B_TYPE, 0b1100011, 0x5, 0x00}, // if(rs1 >= rs2) PC += imm
	"bltu": {B_TYPE, 0b1100011, 0x6, 0x00}, // if(rs1 < rs2) PC += imm zero-extends
	"bgeu": {B_TYPE, 0b1100011, 0x7, 0x00}, // if(rs1 >= rs2) PC += imm zero-extends

	"jal":  {J_TYPE, 0b1101111, 0x0, 0x00}, // rd = PC+4; PC += imm
	"jalr": {I_TYPE, 0b1100111, 0x0, 0x00}, // rd = PC+4; PC = rs1 + imm

	"lui":   {U_TYPE, 0b0110111, 0x0, 0x00}, // rd = imm << 12
	"auipc": {U_TYPE, 0b0010111, 0x0, 0x00}, // rd = PC + (imm << 12)

	"ecall": {I_TYPE, 0b1110011, 0x0, 0x00}, // imm=0x0 Transfer control to OS

	"ebreak": {I_TYPE, 0b1110011, 0x0, 0x01}, // imm=0x1 Transfer control to debugger
}

// Pseudo instructions expansions
var pseudoInstructionSet = map[string][]AssemblyInstruction{
	"nop":  {{name: "addi", iType: &ITypeInstruction{rd: "x0", rs1: "x0", imm: 0}}},
	"mv":   {{name: "addi", iType: &ITypeInstruction{rd: "", rs1: "", imm: 0}}}, // Will be populated during parsing
	"not":  {{name: "xori", iType: &ITypeInstruction{rd: "", rs1: "", imm: -1}}},
	"j":    {{name: "jal", jType: &JTypeInstruction{rd: "x0", imm: 0}}}, // imm filled during parsing
	"ret":  {{name: "jalr", iType: &ITypeInstruction{rd: "x0", rs1: "x1", imm: 0}}},
	"li":   {{name: "lui", uType: &UTypeInstruction{rd: "", imm: 0}}, {name: "addi", iType: &ITypeInstruction{rd: "", rs1: "", imm: 0}}},
	"call": {{name: "auipc", uType: &UTypeInstruction{rd: "x1", imm: 0}}, {name: "jalr", iType: &ITypeInstruction{rd: "x1", rs1: "x1", imm: 0}}},
}

func splitImmediate(imm int32) (int32, int32) {
	high20 := (imm + 0x800) >> 12
	low12 := imm - (high20 << 12)
	return high20, low12
}

// ABI registers map
var registerMap = map[string]int{
	"zero": 0, // Hardwired zero
	"ra":   1, // Return address
	"sp":   2, // Stack pointer
	"gp":   3, // Global pointer
	"tp":   4, // Thread pointer
	"t0":   5, // Temporary registers
	"t1":   6,
	"t2":   7,
	"s0":   8,  // Saved register/Frame pointer
	"s1":   9,  // Saved register
	"a0":   10, // Function arguments/return values
	"a1":   11,
	"a2":   12,
	"a3":   13,
	"a4":   14,
	"a5":   15,
	"a6":   16,
	"a7":   17,
	"s2":   18, // Saved registers
	"s3":   19,
	"s4":   20,
	"s5":   21,
	"s6":   22,
	"s7":   23,
	"s8":   24,
	"s9":   25,
	"s10":  26,
	"s11":  27,
	"t3":   28, // Temporary registers
	"t4":   29,
	"t5":   30,
	"t6":   31,
}

// Represents a RISC-V assembly instruction
type AssemblyInstruction struct {
	name  string            // The instruction name
	rType *RTypeInstruction // The R-type instruction
	iType *ITypeInstruction // The I-type instruction
	sType *STypeInstruction // The S-type instruction
	bType *BTypeInstruction // The B-type instruction
	uType *UTypeInstruction // The U-type instruction
	jType *JTypeInstruction // The J-type instruction
}

// Represents a R-type instruction
type RTypeInstruction struct {
	opcode uint8  // The opcode of the instruction
	rs1    string // The source register 1
	rs2    string // The source register 2
	rd     string // The destination register
	funct3 uint8  // The function code 3
	funct7 uint8  // The function code 7
}

// Represents an I-type instruction
type ITypeInstruction struct {
	opcode uint8  // The opcode of the instruction
	rs1    string // The source register 1
	rd     string // The destination register
	funct3 uint8  // The function code 3
	imm    int16  // The immediate value
}

// Represents a S-type instruction
type STypeInstruction struct {
	opcode uint8  // The opcode of the instruction
	rs1    string // The source register 1
	rs2    string // The source register 2
	funct3 uint8  // The function code 3
	imm    int16  // The immediate value
}

// Represents a B-type instruction
type BTypeInstruction struct {
	opcode uint8  // The opcode of the instruction
	rs1    string // The source register 1
	rs2    string // The source register 2
	funct3 uint8  // The function code 3
	imm    int16  // The immediate value
}

// Represents a U-type instruction
type UTypeInstruction struct {
	opcode uint8  // The opcode of the instruction
	rd     string // The destination register
	imm    int32  // The immediate value
}

// Represents a J-type instruction
type JTypeInstruction struct {
	opcode uint8  // The opcode of the instruction
	rd     string // The destination register
	imm    int32  // The immediate value
}

// Synthesize an R-type instruction
func synthesizeRType(asm RTypeInstruction, instruction Instruction) (uint32, error) {
	source, ok := registerMap[asm.rs1]
	if !ok {
		return 0, fmt.Errorf("invalid source register: %s", asm.rs1)
	}

	target, ok := registerMap[asm.rs2]
	if !ok {
		return 0, fmt.Errorf("invalid target register: %s", asm.rs2)
	}

	destination, ok := registerMap[asm.rd]
	if !ok {
		return 0, fmt.Errorf("invalid destination register: %s", asm.rd)
	}
	encodedInstruction := uint32(instruction.opcode) << 26
	encodedInstruction |= uint32(source) << 21
	encodedInstruction |= uint32(target) << 16
	encodedInstruction |= uint32(destination) << 11
	encodedInstruction |= uint32(asm.funct3) << 12
	encodedInstruction |= uint32(asm.funct7) << 25
	return encodedInstruction, nil
}

// Synthesize an I-type instruction
func synthesizeIType(asm ITypeInstruction, instruction Instruction) (uint32, error) {
	source, ok := registerMap[asm.rs1]
	if !ok {
		return 0, fmt.Errorf("invalid source register: %s", asm.rs1)
	}

	target, ok := registerMap[asm.rd]
	if !ok {
		return 0, fmt.Errorf("invalid target register: %s", asm.rd)
	}
	encodedInstruction := uint32(instruction.opcode) << 26
	encodedInstruction |= uint32(source) << 21
	encodedInstruction |= uint32(target) << 16
	encodedInstruction |= (uint32(asm.imm) & 0xFFFF)
	return encodedInstruction, nil
}

// Synthesize a J-type instruction
func synthesizeJType(asm JTypeInstruction, instruction Instruction) (uint32, error) {
	destination, ok := registerMap[asm.rd]
	if !ok {
		return 0, fmt.Errorf("invalid destination register: %s", asm.rd)
	}
	encodedInstruction := uint32(instruction.opcode) << 26
	encodedInstruction |= uint32(destination) << 16
	encodedInstruction |= uint32(asm.imm) & 0x3FFFFFF
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
