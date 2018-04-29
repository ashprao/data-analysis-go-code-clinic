/**
 * EnvData implements the Lake Pend Oreille data analysis challenge as outlined
 * in the code clinics at LinkedIn Learning at
 * https://www.linkedin.com/learning/code-clinic-java
 *
 *
 * This is just this authors implementation with some additional data analysis
 * over and above the mean and median computed in the LinkedIn version. The most
 * frequent reading along with the frequency and the highest reading is also
 * analyzed and reported in this version.
 *
 * @author Ashwin Rao
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type valueCountMap map[float64]int
type countValueMap map[int][]float64

func (m *valueCountMap) getAverage() (avg float64) {
	sum := 0.0
	list := m.getKeys()
	for _, v := range list {
		sum += v
	}
	return sum / float64(len(list))
}

func (m *valueCountMap) getMedian() (median float64) {
	// Sort the array of keys
	var sortedKeys []float64
	list := m.getKeys()
	for _, v := range list {
		sortedKeys = append(sortedKeys, v)
	}
	sort.Float64s(sortedKeys)

	// Check if array contains even set of keys
	keysLength, mid := len(sortedKeys), len(sortedKeys)/2
	if keysLength%2 == 0 {
		// median is the mean of the middle two numbers in array
		median = (sortedKeys[mid-1] + sortedKeys[mid]) / 2.0
		return
	}
	// Otherwise median is the mid point of the array
	median = sortedKeys[mid]
	return
}

func (m *valueCountMap) getKeys() (keys []float64) {
	// Extract all keys into an array. The count represents the number of times the
	// respective keys should be in the list
	for k, c := range *m {
		// Add the key for as many times as the count
		for i := 0; i < c; i++ {
			keys = append(keys, k)
		}
	}
	return
}

func (m *countValueMap) getKeys() (keys []int) {
	// Extract all keys into an array
	for k := range *m {
		keys = append(keys, k)
	}
	return
}

func (m *valueCountMap) getCounts() (lowFreq int, highFreq int, lowVals []float64, highVals []float64) {
	// Create a map of counts as keys and list of floats as values
	countsMap := make(countValueMap)
	for k, v := range *m {
		countsMap[v] = append(countsMap[v], k)
	}

	// Get all count keys
	keys := countsMap.getKeys()
	var sortedKeys []int
	for _, v := range keys {
		sortedKeys = append(sortedKeys, v)
	}
	sort.Ints(sortedKeys)

	return sortedKeys[0], sortedKeys[len(sortedKeys)-1], countsMap[sortedKeys[0]], countsMap[sortedKeys[len(sortedKeys)-1]]
}

func (m *valueCountMap) getMetrics() (avg float64, median float64, low float64, high float64, lowVals []float64, highVals []float64, lowFreq int, highFreq int) {
	keys := m.getKeys()
	var sortedKeys []float64
	for _, v := range keys {
		sortedKeys = append(sortedKeys, v)
	}
	sort.Float64s(sortedKeys)

	// Get high and low counts and corresponding values with those counts
	lowFreq, highFreq, lowVals, highVals = m.getCounts()

	sort.Float64s(lowVals)
	sort.Float64s(highVals)

	return m.getAverage(), m.getMedian(), sortedKeys[0], sortedKeys[len(keys)-1], lowVals, highVals, lowFreq, highFreq
}

func main() {
	start := time.Now()

	// Data construct to store all wind speed readings
	windSpeedsMap := make(valueCountMap)
	// Data construct to store all air temperature readings
	airTempMap := make(valueCountMap)
	// Data construct to store all barometric pressure readings
	barPressureMap := make(valueCountMap)

	// Open file
	f, err := os.Open(os.Args[1])
	check(err)
	defer f.Close()

	// Get a line scanner
	scanner := bufio.NewScanner(f)

	// Skip past the header line
	scanner.Scan()

	lineCount := 0

	// Get the data in the line
	for scanner.Scan() {
		lineCount++

		// Get a scanner for words
		wordScanner := bufio.NewScanner(strings.NewReader(scanner.Text()))

		// Set the split function for the scanning operation
		wordScanner.Split(bufio.ScanWords)

		// Parse the words in the line
		for column := 1; wordScanner.Scan(); column++ {
			// Get the current word from the sentence
			f, _ := strconv.ParseFloat(wordScanner.Text(), 64)
			switch column {
			case 3: // Air temperature is in the 3rd column
				airTempMap[f] = airTempMap[f] + 1
			case 4: // Barometric pressure is in the 4th column
				barPressureMap[f] = barPressureMap[f] + 1
			case 9: // Wind speed is in the 9th column
				windSpeedsMap[f] = windSpeedsMap[f] + 1
			}
		}
		if lineErr := scanner.Err(); lineErr != nil {
			fmt.Fprintln(os.Stderr, "reading record: ", lineErr)
		}
		if wordErr := wordScanner.Err(); wordErr != nil {
			fmt.Fprintln(os.Stderr, "reading record element: ", wordErr)
		}
	}

	fmt.Printf("Total readings: %d\n", lineCount)

	// Get temperature average and median
	avg, median, low, high, _, highVals, _, highFreq := airTempMap.getMetrics()
	fmt.Printf("Air temperature: Average %.2f, Median %.2f, Low %.2f, High %.2f, Most frequent %v, Frequency %d\n", avg, median, low, high, highVals, highFreq)

	// Get barometric pressure average and median
	avg, median, low, high, _, highVals, _, highFreq = barPressureMap.getMetrics()
	fmt.Printf("Barometric pressure: Average %.2f Median %.2f, Low %.2f, High %.2f, Most frequent %v, Frequency %d\n", avg, median, low, high, highVals, highFreq)

	// Get wind speed average and median
	avg, median, low, high, _, highVals, _, highFreq = windSpeedsMap.getMetrics()
	fmt.Printf("Wind speed: Average %.2f Median %.2f, Low %.2f, High %.2f, Most frequent %v, Frequency %d\n\n", avg, median, low, high, highVals, highFreq)

	elapsed := time.Since(start)
	log.Printf("Total time: %.2f", elapsed.Seconds())
}
