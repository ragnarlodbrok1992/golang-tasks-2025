package main

import (
	"fmt"
	"os"
	// "encoding/hex"
	// "strings"
)

const ASCII_SPACE = 0x20
const ASCII_NEWLINE = 0x0A

type ErrorCode int
type Token int
type ParseControl int

const (
	_ ErrorCode = iota // Skipping 0
	ErrNoSourceFile
	ErrWrongSourceFile
	ErrSourceFileStats
	ErrLoadingBuffer
)

const (
	TokenMain Token = iota // Special token - entry function
	TokenInt
	TokenFunc
	TokenOp
	TokenComment
	TokenEnd // Stop keyword
	TokenEOF // Keep last
)

const (
	ParseControlNoMoreTokens ParseControl = -1
)

var asm_empty_program = `
section .text
global _start

_start:
	mov rax, 60   ; syscall number for exit (60)
	xor rdi, rdi  ; exit status 0
	syscall
`

// Creating a closure that will keep the source file data
// and return tokens until it is empty
func NewTokenParser(source_code_buffer []byte) func() (Token, string) { // Probably will have to return pair of Token and value as string
	parser_buffer := make([]byte, len(source_code_buffer))
	copy(parser_buffer, source_code_buffer)
	current_index := 0 // Start at beginning

	return func() (Token, string) {
		var token_buffer []byte

		// // Right now returning source file bytes - characters
		// cur_char := parser_buffer[current_index]
		// current_index++
		// // fmt.Println(cur_char)
		// return cur_char 

		for {
			if current_index + 1 > len(parser_buffer) {
				return TokenEOF, "EOF"
			}

			cur_byte := parser_buffer[current_index];
			// fmt.Println("cur_byte --> 0x", hex.EncodeToString([]byte{cur_byte}))

			if cur_byte == ASCII_SPACE || cur_byte == ASCII_NEWLINE {
				// fmt.Println("Encountered ASCII_SPACE or ASCII_NEWLINE!")
				ret_buffer := token_buffer
				token_buffer = token_buffer[:0] // Idk if this is how you clean the byte buffer
				current_index++
				return TokenOp, string(ret_buffer) // For now - it's getting quite complicated real fast
			} else {
				token_buffer = append(token_buffer, cur_byte)
				current_index++
			}
		}
	}
}

func main() {
	fmt.Println("Hello, gorth!")
	defer fmt.Println("Goodbye...")

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

					// Allocate array of chars
					source_file_data := make([]byte, file_information.Size())

					n, err := file.Read(source_file_data)
					if err != nil {
						fmt.Println("Couldn't populate src file buffer!")
						os.Exit(int(ErrLoadingBuffer))
					}
					fmt.Println("Source file is ", n, " bytes")

					// Move to code parsing function
					parser := NewTokenParser(source_file_data) // We return function here

					for {
						next_token, token_value := parser()
						if next_token == TokenEOF {
							break
						}
						// Write out char of source file
						fmt.Println(string(next_token), token_value)
					}

					defer file.Close()
				} else {
					fmt.Println("Src file not provided, exiting!")
					os.Exit(int(ErrNoSourceFile))
				}
			}
		}
	}
}
