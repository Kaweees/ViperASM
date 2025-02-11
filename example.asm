#-------------------------------------------------------------------------------
# author       : Miguel Villa Floran
# date         : 2025.02.06
# description  : RISC-V - Hello World
#-------------------------------------------------------------------------------

  .globl main	# declare global symbols

  .data # start of the data section
        hello_str: .asciiz "Hello, World!\n"
  .text # start of the code section

main:
  li $v0, 4 # syscall code for printing a string
  la $a0, hello_str # Load address of hello_str into $a0
  syscall # Make the syscall to print the string

  li $v0, 10 # syscall code for exiting the program
  ecall # execute the syscall (exit the program)
