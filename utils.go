package main

import "math/rand"

// RandomSample uses modified Fisher-Yates algorithm to return k randomly selected nums from a range of n length.
func RandomSample(n, k int) []int {
	population := make([]int, n)

	for i := range population {
		population[i] = i
	}

	for i := 0; i < k; i++ {
		j := rand.Intn(n-i) + i
		population[i], population[j] = population[j], population[i]
	}

	return population[:k]
}
