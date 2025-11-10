package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

// PrimeFactorization returns the prime factors of n.
func PrimeFactorization(n int) []int {
	var factors []int
	// Handle 2 separately
	for n%2 == 0 {
		factors = append(factors, 2)
		n /= 2
	}
	// Check odd divisors up to sqrt(n)
	for i := 3; i <= int(math.Sqrt(float64(n))); i += 2 {
		for n%i == 0 {
			factors = append(factors, i)
			n /= i
		}
	}
	// If remaining n is a prime > 2
	if n > 2 {
		factors = append(factors, n)
	}
	return factors
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run fizzbuzz.go <upper_limit> <value_for_factorization>")
		return
	}

	upperLimit, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Upper limit must be an integer.")
		return
	}

	value, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Value for factorization must be an integer.")
		return
	}

	factors := PrimeFactorization(value)
	if len(factors) == 0 {
		fmt.Println("No prime factors found.")
		return
	}

	lowest := factors[0]
	highest := factors[len(factors)-1]

	for i := 1; i <= upperLimit; i++ {
		output := ""
		if i%lowest == 0 {
			output += "fizz"
		}
		if i%highest == 0 {
			output += "buzz"
		}
		if output != "" {
			fmt.Printf("%d %s\n", i, output)
		}
	}

	fmt.Println("Lowest prime:", lowest, " Highest prime:", highest)
}
