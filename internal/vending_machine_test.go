package internal

import (
	"errors"
	"fmt"
	"testing"
)

func TestScanBucketSlice(t *testing.T) {
	data := []struct {
		scenario      string
		bucket        []int
		products      []int
		expectedSlice []int
	}{
		{
			scenario:      "Only first",
			bucket:        []int{2, 4, 7, 2, 9},
			products:      []int{5, 2, 1},
			expectedSlice: []int{2},
		},
		{
			scenario:      "Only first two",
			bucket:        []int{1, 5, 6, 2, 9},
			products:      []int{5, 2, 1},
			expectedSlice: []int{1, 5},
		},
		{
			scenario:      "Cant be reached v1",
			bucket:        []int{6, 5, 2, 2, 9},
			products:      []int{5, 2, 1},
			expectedSlice: []int{},
		},
		{
			scenario:      "Cant be reached v2",
			bucket:        []int{6, 4, 5, 2, 9},
			products:      []int{5, 2, 1},
			expectedSlice: []int{},
		},
		{
			scenario:      "Found but not in same order v1",
			bucket:        []int{2, 1, 5, 2, 9},
			products:      []int{5, 2, 1},
			expectedSlice: []int{2, 1, 5},
		},
		{
			scenario:      "Found but not in same order v2",
			bucket:        []int{5, 1, 2, 2, 9},
			products:      []int{5, 2, 1},
			expectedSlice: []int{5, 1, 2},
		},
		{
			scenario:      "Found with repeated elements v1",
			bucket:        []int{5, 1, 1, 2, 9},
			products:      []int{5, 2, 1},
			expectedSlice: []int{5, 1},
		},
		{
			scenario:      "Found with repeated elements v2",
			bucket:        []int{5, 1, 1, 1, 2, 9},
			products:      []int{5, 2, 1, 1},
			expectedSlice: []int{5, 1, 1},
		},
		{
			scenario:      "Found with repeated elements v3",
			bucket:        []int{5, 1, 2, 1, 2, 9},
			products:      []int{5, 2, 1, 1},
			expectedSlice: []int{5, 1, 2, 1},
		},
	}

	for _, d := range data {
		t.Run(d.scenario, func(t *testing.T) {
			resultSlice := ScanBucketSlice(d.bucket, &d.products)
			if len(*resultSlice) != len(d.expectedSlice) {
				t.Fatalf("Invalid size of arrays, expected size %d got %d", len(d.expectedSlice), len(*resultSlice))
			}
			for i := 0; i < len(*resultSlice); i++ {
				if (*resultSlice)[i] != d.expectedSlice[i] {
					t.Fatalf("Invalid value at %d: got %d expected %d", i, (*resultSlice)[i], d.expectedSlice[i])
				}
			}
		})
	}
}

func TestFindFirstPattern(t *testing.T) {
	data := []struct {
		scenario          string
		possibleSlice     *[]*PossibleBucketSlice
		possibleSliceCopy *[]*PossibleBucketSlice
		products          *[]int
		expectedPatterns  *[]*PopPattern
	}{
		{
			scenario: "Simple",
			possibleSlice: &[]*PossibleBucketSlice{
				{
					Index:  4,
					Values: []int{1, 2},
				},
				{
					Index:  0,
					Values: []int{2},
				},
				{
					Index:  1,
					Values: []int{1, 2},
				},
				{
					Index:  2,
					Values: []int{5, 2},
				},
			},
			possibleSliceCopy: &[]*PossibleBucketSlice{
				{
					Index:  4,
					Values: []int{1, 2},
				},
				{
					Index:  0,
					Values: []int{2},
				},
				{
					Index:  1,
					Values: []int{1, 2},
				},
				{
					Index:  2,
					Values: []int{5, 2},
				},
			},
			products: &[]int{5, 2, 1},
			expectedPatterns: &[]*PopPattern{
				{
					Index:        3,
					NumberPopped: 2,
				},
				{
					Index:        0,
					NumberPopped: 1,
				},
			},
		},
		{
			scenario: "Simple v2",
			possibleSlice: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{2},
				},
				{
					Index:  1,
					Values: []int{5, 1, 2},
				},
			},
			possibleSliceCopy: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{2},
				},
				{
					Index:  1,
					Values: []int{5, 1, 2},
				},
			},
			products: &[]int{5, 2, 1},
			expectedPatterns: &[]*PopPattern{
				{
					Index:        1,
					NumberPopped: 1,
				},
				{
					Index:        0,
					NumberPopped: 1,
				},
				{
					Index:        1,
					NumberPopped: 1,
				},
			},
		},
		{
			scenario: "Impossible",
			possibleSlice: &[]*PossibleBucketSlice{
				{
					Index:  4,
					Values: []int{1, 2},
				},
				{
					Index:  0,
					Values: []int{2},
				},
				{
					Index:  1,
					Values: []int{1, 2},
				},
				{
					Index:  2,
					Values: []int{1, 5},
				},
			},
			possibleSliceCopy: &[]*PossibleBucketSlice{
				{
					Index:  4,
					Values: []int{1, 2},
				},
				{
					Index:  0,
					Values: []int{2},
				},
				{
					Index:  1,
					Values: []int{1, 2},
				},
				{
					Index:  2,
					Values: []int{1, 5},
				},
			},
			products:         &[]int{5, 2, 1},
			expectedPatterns: nil,
		},
	}

	for _, d := range data {
		t.Run(d.scenario, func(t *testing.T) {
			patterns := FindFirstPattern(d.possibleSlice, d.products)
			if patterns == nil && d.expectedPatterns != nil {
				t.Fatal("Got nil patterns\n")
			}

			if err := assertEqualPatterns(patterns, d.expectedPatterns); err != nil {
				t.Fatal(err)
			}

			for i, originalSlice := range *d.possibleSliceCopy {
				val := *d.possibleSlice
				if areEqualInt(originalSlice.Values, val[i].Values) == false {
					t.Fatalf("Arrays have mutated\n")
				}
			}
		})
	}
}

func TestFindFirstNoOrderPattern(t *testing.T) {
	data := []struct {
		scenario          string
		possibleSlice     *[]*PossibleBucketSlice
		possibleSliceCopy *[]*PossibleBucketSlice
		products          *[]int
		expectedPatterns  *[]*PopPattern
	}{
		{
			scenario: "First row instant",
			possibleSlice: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{1, 3, 5, 2, 4},
				},
				{
					Index:  1,
					Values: []int{2, 5, 4, 3, 1},
				},
			},
			possibleSliceCopy: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{1, 3, 5, 2, 4},
				},
				{
					Index:  1,
					Values: []int{2, 5, 4, 3, 1},
				},
			},
			products: &[]int{1, 2, 3, 4, 5},
			expectedPatterns: &[]*PopPattern{
				{
					Index:        0,
					NumberPopped: 5,
				},
			},
		},
		{
			scenario: "Simple",
			possibleSlice: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{1, 2, 3, 5, 5},
				},
				{
					Index:  1,
					Values: []int{2, 5, 4, 3, 1},
				},
			},
			possibleSliceCopy: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{1, 2, 3, 5, 5},
				},
				{
					Index:  1,
					Values: []int{2, 5, 4, 3, 1},
				},
			},
			products: &[]int{1, 2, 3, 4, 5},
			expectedPatterns: &[]*PopPattern{
				{
					Index:        0,
					NumberPopped: 1,
				},
				{
					Index:        1,
					NumberPopped: 4,
				},
			},
		},
		{
			scenario: "first and third bucket",
			possibleSlice: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{1, 2, 3, 5, 5},
				},
				{
					Index:  1,
					Values: []int{3, 5, 4, 3, 1},
				},
				{
					Index:  1,
					Values: []int{5, 4, 2, 1, 1},
				},
			},
			possibleSliceCopy: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{1, 2, 3, 5, 5},
				},
				{
					Index:  1,
					Values: []int{3, 5, 4, 3, 1},
				},
				{
					Index:  1,
					Values: []int{5, 4, 2, 1, 1},
				},
			},
			products: &[]int{1, 2, 3, 4, 5},
			expectedPatterns: &[]*PopPattern{
				{
					Index:        0,
					NumberPopped: 3,
				},
				{
					Index:        2,
					NumberPopped: 2,
				},
			},
		},
		{
			scenario: "Second and third",
			possibleSlice: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{2, 1, 2, 2, 5},
				},
				{
					Index:  1,
					Values: []int{3, 5, 2, 3, 1},
				},
				{
					Index:  1,
					Values: []int{1, 2, 4, 1, 1},
				},
			},
			possibleSliceCopy: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{2, 1, 2, 2, 5},
				},
				{
					Index:  1,
					Values: []int{3, 5, 2, 3, 1},
				},
				{
					Index:  1,
					Values: []int{1, 2, 4, 1, 1},
				},
			},
			products: &[]int{1, 2, 3, 4, 5},
			expectedPatterns: &[]*PopPattern{
				{
					Index:        1,
					NumberPopped: 2,
				},
				{
					Index:        2,
					NumberPopped: 3,
				},
			},
		},
		{
			scenario: "Impossible scenario",
			possibleSlice: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{2, 1, 2, 2, 5},
				},
				{
					Index:  1,
					Values: []int{3, 5, 2, 3, 1},
				},
				{
					Index:  1,
					Values: []int{1, 2, 1, 4, 1},
				},
			},
			possibleSliceCopy: &[]*PossibleBucketSlice{
				{
					Index:  0,
					Values: []int{2, 1, 2, 2, 5},
				},
				{
					Index:  1,
					Values: []int{3, 5, 2, 3, 1},
				},
				{
					Index:  1,
					Values: []int{1, 2, 1, 4, 1},
				},
			},
			products:         &[]int{1, 2, 3, 4, 5},
			expectedPatterns: nil,
		},
	}

	for _, d := range data {
		t.Run(d.scenario, func(t *testing.T) {
			patterns := FindFirstNoOrderPattern(d.possibleSlice, d.products)
			if patterns == nil && d.expectedPatterns != nil {
				t.Fatal("Got nil patterns\n")
			}

			if err := assertEqualPatterns(patterns, d.expectedPatterns); err != nil {
				t.Fatal(err)
			}

			for i, originalSlice := range *d.possibleSliceCopy {
				val := *d.possibleSlice
				if areEqualInt(originalSlice.Values, val[i].Values) == false {
					t.Fatalf("Arrays have mutated\n")
				}
			}
		})
	}
}

func TestFindCumulativePopPattern_NoOrder(t *testing.T) {
	data := []struct {
		vendingMachine  [][]int
		products        []int
		expectedPattern []*PopPattern
	}{
		{
			vendingMachine: [][]int{
				{1, 2, 3, 5, 5},
				{2, 5, 4, 3, 1},
				{3, 5, 4, 1, 1},
				{5, 1, 1, 1, 1},
			},
			products: []int{1, 2, 3, 4, 5},
			expectedPattern: []*PopPattern{
				{
					Index:        0,
					NumberPopped: 1,
				},
				{
					Index:        1,
					NumberPopped: 4,
				},
			},
		},
	}

	for _, d := range data {
		pattern, err := FindCumulativePopPattern(&d.vendingMachine, &d.products, FindFirstNoOrderPattern)
		if err != nil {
			t.Fatal(err)
		}
		if err := assertEqualPatterns(&d.expectedPattern, pattern); err != nil {
			t.Fatal(err)
		}
	}
}

func areEqualInt(arr1 []int, arr2 []int) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}

func assertEqualPatterns(pattern *[]*PopPattern, other *[]*PopPattern) error {
	if pattern == nil && other == nil {
		return nil
	}
	if pattern == nil || other == nil {
		return errors.New("one pattern is nil")
	}
	if len(*pattern) != len(*other) {
		return errors.New(fmt.Sprintf("expected length of %d got %d", len(*other), len(*pattern)))
	}

	for i := 0; i < len(*pattern); i++ {
		if (*pattern)[i].NumberPopped != (*other)[i].NumberPopped ||
			(*pattern)[i].Index != (*other)[i].Index {
			return errors.New(fmt.Sprintf(
				"expected equality: %s = %s", (*pattern)[i].toString(), (*other)[i].toString(),
			))
		}
	}

	return nil
}
