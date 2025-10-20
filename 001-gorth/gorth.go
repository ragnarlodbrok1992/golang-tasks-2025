package main

import (
	"fmt"
	"os"
	// "strings"
)

type ErrorCode int

const (
	_ ErrorCode = iota // Skipping 0
	ErrNoSourceFile
	ErrWrongSourceFile
	ErrSourceFileStats
	ErrLoadingBuffer
)

var asm_empty_program = `
section .text
global _start

_start:
	mov rax, 60   ; syscall number for exit (60)
	xor rdi, rdi  ; exit status 0
	syscall
`

func main() {
	fmt.Println("Hello, gorth!")
	defer fmt.Println("Goodbye...")

	// Test prints
	fmt.Println("Some test prints:")
	fmt.Println(asm_empty_program)

	// Opening input file
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		// Test print
		// fmt.Printf("arg: %s\n", arg)

		switch arg {
			case "-src": {
				// fmt.Println("Found src value!")

				if i + 1 < len(os.Args) {
					src_file := os.Args[i + 1]
					fmt.Println("src_file: ", src_file)

					// Trying to open the source file
					file, err := os.Open(src_file)
					if err != nil {
						fmt.Println("Src file is wrong, exiting!")
						os.Exit(int(ErrWrongSourceFile))
					}
					file_information, err := file.Stat()
					if err != nil {
						fmt.Println("Couldn't get file stats!")
						os.Exit(int(ErrSourceFileStats))
					}
					// File opened correctly - do something with it
					// Write out information about the file
					// Test code
					fmt.Printf("Name: %s\n", file.Name()) 
					fmt.Printf("Size: %d bytes\n", file_information.Size()) 

					// Allocate array of chars
					source_file_data := make([]byte, file_information.Size())

					// Populate allocated memory
					n, err := file.Read(source_file_data)
					if err != nil {
						fmt.Println("Couldn't populate src file buffer!")
						os.Exit(int(ErrLoadingBuffer))
					}
					fmt.Println("Read up to ", n)

					// Print out source file
					fmt.Println("--- source file ---")
					fmt.Println(string(source_file_data))
					fmt.Println("--- end of source file ---")

					defer file.Close()
				} else {
					fmt.Println("Src file not provided, exiting!")
					os.Exit(int(ErrNoSourceFile))
				}
			}
		}
	}
}
