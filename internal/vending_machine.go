package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// input
// 5 2 1

// 2 4 7 2 9
// 1 2 1 3
// 5 2 3 3 3
// 4 3 2 6
// 1 2
// 7 8 2 5 7
// 4
// 7 2 5 1

// Best scenario seems to be to go through the buckets we already went up to the length of input,
// doing this we would create a new matrix of buckets that contain the products we need, so that
// we can use them in our combination. We will always take a slice of the bucket that is no more
// than the max length of input, because it can't go any deeper than that. We must keep index of
// the original bucket in matrix

type PopPattern struct {
	Index        int
	NumberPopped int
}

func (pp *PopPattern) Print() {
	fmt.Printf("[PopPattern]: Index: %d NumberPopped %d\n", pp.Index, pp.NumberPopped)
}
func (pp *PopPattern) toString() string {
	return fmt.Sprintf("[PopPattern]: Index %d, NumberPopped %d", pp.Index, pp.NumberPopped)
}

type PossibleBucketSlice struct {
	Index  int
	Values []int
}

func (pb *PossibleBucketSlice) Copy() *PossibleBucketSlice {
	return &PossibleBucketSlice{
		Index:  pb.Index,
		Values: pb.Values,
	}
}

func (pb *PossibleBucketSlice) Print() {
	fmt.Printf("[PossibleBucketSlice]: index %d values %+v\n", pb.Index, pb.Values)
}

type PatternFunc func(possibleSlice *[]*PossibleBucketSlice, products *[]int) *[]*PopPattern

var ImpossibleErr = errors.New("IMPOSSIBLE")
var InvalidArgument = errors.New("invalid argument")

func PrintPretty(vendingMachine *[][]int) {
	fmt.Println("Vending machine")
	for _, bucket := range *vendingMachine {
		fmt.Printf("\t%+v\n", bucket)
	}
}

func CreateFromString(str string) (*[][]int, error) {
	matrix := [][]int{}

	if len(str) == 0 {
		return &matrix, nil
	}

	buckets := strings.Split(str, ";")

	for _, bucket := range buckets {
		bucketProducts := strings.Split(bucket, ",")
		parsedBucketProducts := []int{}

		for _, product := range bucketProducts {
			parsedProduct, err := strconv.Atoi(product)
			if err != nil {
				return nil, InvalidArgument
			}
			parsedBucketProducts = append(parsedBucketProducts, parsedProduct)
		}
		matrix = append(matrix, parsedBucketProducts)
	}

	return &matrix, nil
}

func ParseInput(input string) (*[]int, error) {
	parsedInput := strings.Split(input, ",")

	products := []int{}
	for _, p := range parsedInput {
		parsedProduct, err := strconv.Atoi(p)
		if err != nil {
			return nil, InvalidArgument
		}
		products = append(products, parsedProduct)
	}

	return &products, nil
}

// Scan products from bucket to see if all match by returning array of the same
// length as the products array. If returned array has smaller length than products
// it can be used later for repeated pattern scans.
// Found slice maybe equal to products, in that case it is found.
func ScanBucketSlice(bucket []int, products *[]int) *[]int {
	slice := []int{}

	// We need to avoid duplicates, hence every found is removed from the stack
	productsStack := make([]int, len(*products))
	copy(productsStack, *products)

	var length int
	if len(bucket) < len(*products) {
		length = len(bucket)
	} else {
		length = len(*products)
	}

	for i := 0; i < length; i++ {
		index := getIndexFromArray(bucket[i], productsStack)
		if index == -1 {
			break
		}
		slice = append(slice, bucket[i])
		// remove found product (cut)
		cutIntFromSlice(&productsStack, index)
	}

	return &slice
}

func getIndexFromArray(item int, items []int) int {
	for i, p := range items {
		if item == p {
			return i
		}
	}

	return -1
}

func PopByPattern(vendingMachine *[][]int, patterns *[]*PopPattern) {
	// Order matters
	for _, pattern := range *patterns {
		(*vendingMachine)[pattern.Index] = (*vendingMachine)[pattern.Index][pattern.NumberPopped:]
	}
}

func cutIntFromSlice(slice *[]int, index int) {
	*slice = append((*slice)[:index], (*slice)[index+1:]...)
}

// We try stacking up to the maximum product length, from first bucket to the last.
// Algorithm will focus on 1 bucket search, then 2 bucket search, then 3 bucket search, etc
// so that we avoid looping over an entire array, because there is a greater chance that
// the next bucket or two will be enough for constructing a pattern and have less items
// to loop through.
// Example buckets:
//
// [1 2 3 5 5]
// [2 5 4 3 1]
// [3 5 4 1 1]
// [5 1 1 1 1]
//
// We are looking for 1 2 3 4 5, order does not matter
//
// If the first bucket does NOT stack 5, but stacks 4, we try the second bucket for what's left.
// If the second bucket does not stack, then we remove one item from the stack and try again,
// this time checking two values, if it fails we remove from the stack and continue.
// There are two cases from here:
//		- Stack can be empty, then we move to the next bucket and start all over
//		- Bucket we are stacking is empty, then we just move to the next bucket and repeat
func FindFirstNoOrderPattern(possibleSlice *[]*PossibleBucketSlice, products *[]int) *[]*PopPattern {
	bucketsCopy := make([]*PossibleBucketSlice, len(*possibleSlice))
	copy(bucketsCopy, *possibleSlice)

	productsStack := make([]int, len(*products))
	copy(productsStack, *products)
	var pops []*PopPattern

	// Used for tracking of current pattern
	// For example the stack is on first 4 products looking for 5th
	// if not found, the currentPopPattern is decremented by 1 to
	// look for remaining 2, etc until it's down to 0 count
	var currentPopPattern *PopPattern

	// Loop through buckets
	for i := 0; i < len(bucketsCopy) && len(productsStack) > 0; i++ {

		var pop *PopPattern

		// Loop through products
		for j := 0; j < len(bucketsCopy[i].Values) && len(productsStack) > 0; j++ {
			product := bucketsCopy[i].Values[j]
			index := getIndexFromArray(product, productsStack)
			if index == -1 {
				// Not in the product list, go to the next bucket
				break
			}
			cutIntFromSlice(&productsStack, index)

			// Update tracking pop pattern
			if currentPopPattern == nil {
				currentPopPattern = &PopPattern{
					Index:        i,
					NumberPopped: 1,
				}
			} else {
				if i == currentPopPattern.Index {
					currentPopPattern.NumberPopped++
				}
			}

			// Add pop pattern
			if pop == nil {
				pop = &PopPattern{
					Index:        i,
					NumberPopped: 1,
				}
			} else {
				pop.NumberPopped = pop.NumberPopped + 1
			}
		}

		if pop != nil {
			// Avoid adding current pop pattern as it will be added at the end
			if currentPopPattern == nil || pop.Index != currentPopPattern.Index {
				pops = append(pops, pop)
			}
		}

		// If last element of the loop and we still don't have the products
		// decrement pattern count and retry
		if i == len(bucketsCopy)-1 && len(productsStack) > 0 {
			// If the bucket products didn't match, we move to the next bucket
			if currentPopPattern != nil && currentPopPattern.NumberPopped > 0 {
				// break if there are no more buckets
				if len(bucketsCopy) == currentPopPattern.Index {
					break
				}
				currBucket := bucketsCopy[currentPopPattern.Index]
				lastPop := currBucket.Values[currentPopPattern.NumberPopped-1]
				productsStack = append(productsStack, lastPop) // add one back to the stack
				currentPopPattern.NumberPopped--
				i = currentPopPattern.Index // i will get incremented in next for each

				// Revert products stack from all pops
				for _, p := range pops {
					bucket := bucketsCopy[p.Index]
					// get popped products
					productsPopped := (*bucket).Values[:p.NumberPopped]
					productsStack = append(productsStack, productsPopped...)
				}
				// clear pop patterns as the current tracking pop pattern didn't work
				pops = nil

				// Reset the current pop pattern so that a new one can be created
				if currentPopPattern.NumberPopped == 0 {
					currentPopPattern = nil
				}
			}
		}
	}

	if currentPopPattern != nil && currentPopPattern.NumberPopped == 0 {
		return nil
	}
	appended := append([]*PopPattern{currentPopPattern}, pops...)

	patternPopSum := 0
	for _, ap := range appended {
		patternPopSum += ap.NumberPopped
	}

	if patternPopSum != len(*products) {
		return nil
	}

	return &appended
}

func FindFirstPattern(possibleSlice *[]*PossibleBucketSlice, products *[]int) *[]*PopPattern {
	bucketsCopy := make([]*PossibleBucketSlice, len(*possibleSlice))
	copy(bucketsCopy, *possibleSlice)

	var pops []*PopPattern
	currProductIndex := 0

	lastLoopProductsFound := 0

	// Loop through sliced buckets
	for i := 0; i < len(bucketsCopy) && currProductIndex < len(*products); i++ {
		pop := &PopPattern{
			Index:        i,
			NumberPopped: 0,
		}

		// Check bucket products by sequence
		for currProductIndex < len(*products) {
			if len(bucketsCopy[i].Values) == 0 {
				break
			}
			firstProduct := bucketsCopy[i].Values[0]

			if firstProduct != (*products)[currProductIndex] {
				break
			}

			currProductIndex++
			pop.NumberPopped++
			// Remove popped, ensure not to mutate original
			sliced := bucketsCopy[i].Values[1:]
			cloned := bucketsCopy[i].Copy()
			cloned.Values = sliced
			bucketsCopy[i] = cloned
		}

		if pop.NumberPopped > 0 {
			pops = append(pops, pop)
		}

		// Start again from the start, but only if you found a new product this loops
		if i+1 == len(bucketsCopy) {
			if lastLoopProductsFound != len(pops) {
				// -1 because at the end of the loop it will get increased by 1
				i = -1
				lastLoopProductsFound = len(pops)
			}
		}
	}

	// If we don't have all the pops, the pattern is not valid
	sumPops := 0
	for _, p := range pops {
		sumPops += p.NumberPopped
	}
	if sumPops != len(*products) {
		return nil
	}

	return &pops
}

func FindCumulativePopPattern(vendingMachine *[][]int, products *[]int, fn PatternFunc) (*[]*PopPattern, error) {
	var possibleSlices []*PossibleBucketSlice

	for i, bucket := range *vendingMachine {
		// n row times
		slice := ScanBucketSlice(bucket, products)
		if len(*slice) == 0 {
			continue
		}

		possibleSlices = append(possibleSlices, &PossibleBucketSlice{
			Index:  i,
			Values: *slice,
		})

		patterns := fn(&possibleSlices, products)
		if patterns != nil {
			return patterns, nil
		}
	}

	return nil, ImpossibleErr
}

func FindAndPopByOrder(vendingMachine *[][]int, products *[]int, fn PatternFunc) error {
	patterns, err := FindCumulativePopPattern(vendingMachine, products, fn)
	if err != nil {
		return err
	}
	if patterns == nil {
		return ImpossibleErr
	}

	PopByPattern(vendingMachine, patterns)

	return nil
}
