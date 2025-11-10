package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var stages = []string{
    `
  +---+
  |   |
      |
      |
      |
      |
=========
    `,
    `
  +---+
  |   |
  O   |
      |
      |
      |
=========
    `,
    `
  +---+
  |   |
  O   |
  |   |
      |
      |
=========
    `,
    `
  +---+
  |   |
  O   |
 /|   |
      |
      |
=========
    `,
    `
  +---+
  |   |
  O   |
 /|\  |
      |
      |
=========
    `,
    `
  +---+
  |   |
  O   |
 /|\  |
 /    |
      |
=========
    `,
    `
  +---+
  |   |
  O   |
 /|\  |
 / \  |
      |
=========
    `,
}

var word_list = []string {
	"gopher",
	"programming",
	"hangman",
	"computer",
	"keyboard",
	"developer",
	"algorithm",
}

func main() {
	rand.Seed(time.Now().UnixNano())
	word := word_list[rand.Intn(len(word_list))]
	guessed := make(map[rune]bool)
	incorrectGuesses := 0
	maxIncorrectGuesses := len(stages) - 1

	// Initialize the display with underscores
	display := make([]rune, len(word))
	for i := range display {
		display[i] = '_'
	}

	scanner := bufio.NewScanner(os.Stdin)

	for incorrectGuesses < maxIncorrectGuesses {
		fmt.Print(stages[incorrectGuesses])
		fmt.Printf("\nWord: %s\n", string(display))
		fmt.Printf("Guessed letters: ")
		for letter := range guessed {
			fmt.Printf("%c ", letter)
		}
		fmt.Print("\nEnter a letter: ")

		scanner.Scan()
		input := strings.ToLower(scanner.Text())

		if len(input) != 1 {
			fmt.Println("Please enter a single letter.")
			continue
		}

		letter := rune(input[0])
		if guessed[letter] {
			fmt.Println("You already guessed that letter.")
			continue
		}

		guessed[letter] = true
		correct := false

		for i, char := range word {
			if char == letter {
				display[i] = letter
				correct = true
			}
		}

		if !correct {
			incorrectGuesses++
			fmt.Println("Incorrect guess!")
		}

		if string(display) == word {
			fmt.Printf("\nCongratulations! You guessed the word: %s\n", word)
			return
		}
	}

	fmt.Print(stages[maxIncorrectGuesses])
	fmt.Printf("\nGame over! The word was: %s\n", word)
}
