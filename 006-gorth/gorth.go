package main

import (
	"fmt"
	"os"
	"strconv"
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
	TokenFuncParam
	TokenOp
	TokenComment // TODO: do we even need a token for comment? maybe for some meta-tuning or some stuff like that
	TokenEnd // Stop keyword
	TokenEOF // Keep last
)

const (
	ParseControlNoMoreTokens ParseControl = -1
)

var asm_empty_program_start = `
section .text
global _start

_start:
`

var asm_empty_program_end = `
	mov rax, 60   ; syscall number for exit (60)
	xor rdi, rdi  ; exit status 0
	syscall
`

var asm_empty_program = asm_empty_program_start + asm_empty_program_end

func (tkn Token) ToString() string {
	switch tkn {
		case TokenMain:
			return "TokenMain"
		case TokenInt:
			return "TokenInt"
		case TokenFunc:
			return "TokenFunc"
		case TokenOp:
			return "TokenOp"
		case TokenComment:
			return "TokenComment"
		case TokenEnd:
			return "TokenEnd"
		case TokenEOF:
			return "TokenEOF"
		default:
			panic("Unknown token to string conversion...")
	}
}

func canCastToInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// TODO: evaluate token based on it's contents
func EvaluateToken(str []byte) Token {
	// fmt.Println("Passing token --> ", string(str))
	tkn_str := string(str)
	fmt.Printf("Passing token --> %sH\n", tkn_str)

	if str[0] == ';' {
		return TokenComment
	}


	switch tkn_str {
	case "main":
		return TokenMain
	case "stop":
		return TokenEnd
	case "print": // TODO: Move this case to something more sophisticated, built-in function will land in some kind of hashmap
		return TokenFunc
	}

	// Check if token - value can be cast to integer
	if canCastToInt(tkn_str) {
		return TokenInt
	}

	return TokenOp
}

// Creating a closure that will keep the source file data
// and return tokens until it is empty
func NewTokenParser(source_code_buffer []byte) func() (Token, string) { // Probably will have to return pair of Token and value as string
	parser_buffer := make([]byte, len(source_code_buffer))
	copy(parser_buffer, source_code_buffer)
	current_index := 0 // Start at beginning

	return func() (Token, string) {
		var token_buffer []byte

		for {
			if current_index + 1 > len(parser_buffer) {
				return TokenEOF, "EOF"
			}

			cur_byte := parser_buffer[current_index];
			// fmt.Println("cur_byte --> 0x", hex.EncodeToString([]byte{cur_byte}))

			// TODO: a lot of reused code can be factored out here
			// What about multiple spaces? TODO: check if token_buffer is empty and do not return it
			if cur_byte == ASCII_SPACE || cur_byte == ASCII_NEWLINE {
				// fmt.Println("Encountered ASCII_SPACE or ASCII_NEWLINE!")
				ret_buffer := token_buffer
				token_buffer = token_buffer[:0] // Idk if this is how you clean the byte buffer
				current_index++

				ret_token := EvaluateToken(ret_buffer)

				if ret_token == TokenComment {
					// We keep the ';'
					token_buffer = append(token_buffer, ';')
					token_buffer = append(token_buffer, ' ')

					for {
						cur_byte = parser_buffer[current_index];
						if cur_byte == ASCII_NEWLINE {
							ret_buffer = token_buffer;
							token_buffer = token_buffer[:0]
							current_index++

							return ret_token, string(ret_buffer);
						}

						token_buffer = append(token_buffer, cur_byte)
						current_index++
					}
				}

				// TODO: there will be a special case for TokenComment...
				return ret_token, string(ret_buffer) // For now - it's getting quite complicated real fast
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
						fmt.Println(next_token.ToString(), token_value)
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
