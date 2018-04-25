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

func getKeysFromValueMap(m map[float64]int) (keys []float64) {
	// Extract all keys into an array
	for k, v := range m {
		for i := 0; i < v; i++ {
			keys = append(keys, k)
		}
	}
	return
}

func getKeysFromCountMap(m map[int][]float64) (keys []int) {
	// Extract all keys into an array
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func getAverage(list []float64) (avg float64) {
	sum := 0.0
	for _, v := range list {
		sum += v
	}
	return sum / float64(len(list))
}

func getMedian(list []float64) (median float64) {
	// Sort the array of keys
	sort.Float64s(list)

	// Check if array contains even set of keys
	keysLength, mid := len(list), len(list)/2
	if keysLength%2 == 0 {
		// median is the mean of the middle two numbers in array
		median = (list[mid-1] + list[mid]) / 2.0
		return
	}
	// Otherwise median is the mid point of the array
	median = list[mid]
	return
}

func getCounts(m map[float64]int) (lowFreq int, highFreq int, lowVals []float64, highVals []float64) {
	// Create a map of counts as keys and list of floats as values
	countsMap := make(map[int][]float64)
	for k, v := range m {
		countsMap[v] = append(countsMap[v], k)
	}

	// Get all count keys
	keys := getKeysFromCountMap(countsMap)
	sort.Ints(keys)

	return keys[0], keys[len(keys)-1], countsMap[keys[0]], countsMap[keys[len(keys)-1]]
}

func getMetrics(m map[float64]int) (avg float64, median float64, low float64, high float64, lowVals []float64, highVals []float64, lowFreq int, highFreq int) {
	keys := getKeysFromValueMap(m)

	// Get high and low counts and corresponding values with those counts
	lowFreq, highFreq, lowVals, highVals = getCounts(m)

	sort.Float64s(lowVals)
	sort.Float64s(highVals)

	avg = getAverage(keys)
	median = getMedian(keys)
	return avg, median, keys[0], keys[len(keys)-1], lowVals, highVals, lowFreq, highFreq
}

func main() {
	start := time.Now()

	// Data construct to store all wind speed readings
	windSpeedsMap := make(map[float64]int)
	// Data construct to store all air temperature readings
	airTempMap := make(map[float64]int)
	// Data construct to store all barometric pressure readings
	barPressureMap := make(map[float64]int)

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
		for count := 1; wordScanner.Scan(); count++ {
			// Get the current word from the sentence
			f, _ := strconv.ParseFloat(wordScanner.Text(), 64)
			switch count {
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
	avg, median, low, high, _, highVals, _, highFreq := getMetrics(airTempMap)
	fmt.Printf("Air temperature: Average %.2f, Median %.2f, Low %.2f, High %.2f, Most frequent %v, Frequency %d\n", avg, median, low, high, highVals, highFreq)

	// Get barometric pressure average and median
	avg, median, low, high, _, highVals, _, highFreq = getMetrics(barPressureMap)
	fmt.Printf("Barometric pressure: Average %.2f Median %.2f, Low %.2f, High %.2f, Most frequent %v, Frequency %d\n", avg, median, low, high, highVals, highFreq)

	// Get wind speed average and median
	avg, median, low, high, _, highVals, _, highFreq = getMetrics(windSpeedsMap)
	fmt.Printf("Wind speed: Average %.2f Median %.2f, Low %.2f, High %.2f, Most frequent %v, Frequency %d\n\n", avg, median, low, high, highVals, highFreq)

	elapsed := time.Since(start)
	log.Printf("Total time: %.2f", elapsed.Seconds())
}
