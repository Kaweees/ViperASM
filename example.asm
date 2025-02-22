#-------------------------------------------------------------------------------
# author       : Miguel Villa Floran
# date         : 2025.01.31
# description  : RISC-V - Hello World
#-------------------------------------------------------------------------------

  .globl main	# declare global symbols

  .data # start of the data section
    hello_str: .asciiz "Hello, World!\n"
    newline: .byte '\n'    # Newline character
    tab:     .byte '\t'    # Tab character
    char:    .byte 'A'     # Regular character
  .text # start of the code section

main:
  li a7, 4 # syscall code for printing a string
  la a0, hello_str # Load address of hello_str into $a0
  ecall # Make the syscall to print the string

  li a7, 10 # syscall code for exiting the program
  ecall # execute the syscall (exit the program)
