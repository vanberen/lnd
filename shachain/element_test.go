package shachain

import (
	"reflect"
	"testing"

	"github.com/go-errors/errors"
)

// bitsToIndex is a helper function which takes 'n' last bits as input and
// create shachain index.
// Example:
//  Input: 0,1,1,0,0
//  Output: 0b000000000000000000000000000000000000000[01100] == 12
func bitsToIndex(bs ...uint64) (index, error) {
	if len(bs) > 64 {
		return 0, errors.New("number of elements should be lower then" +
			" 64")
	}

	var res uint64
	for i, e := range bs {
		if e != 1 && e != 0 {
			return 0, errors.New("wrong element, should be '0' or" +
				" '1'")
		}

		res += e * 1 << uint(len(bs)-i-1)
	}

	return index(res), nil
}

type deriveTest struct {
	name       string
	from       index
	to         index
	position   []uint8
	shouldFail bool
}

func generateTests(t *testing.T) []deriveTest {
	var (
		tests []deriveTest
		from  index
		to    index
		err   error
	)

	from, err = bitsToIndex(0)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	to, err = bitsToIndex(0)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	tests = append(tests, deriveTest{
		name:       "zero 'from' 'to'",
		from:       from,
		to:         to,
		position:   nil,
		shouldFail: false,
	})

	from, err = bitsToIndex(0, 1, 0, 0)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	to, err = bitsToIndex(0, 1, 0, 0)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	tests = append(tests, deriveTest{
		name:       "same indexes #1",
		from:       from,
		to:         to,
		position:   nil,
		shouldFail: false,
	})

	from, err = bitsToIndex(1)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	to, err = bitsToIndex(0)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	tests = append(tests, deriveTest{
		name:       "same indexes #2",
		from:       from,
		to:         to,
		shouldFail: true,
	})

	from, err = bitsToIndex(0, 0, 0, 0)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	to, err = bitsToIndex(0, 0, 1, 0)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	tests = append(tests, deriveTest{
		name:       "test seed 'from'",
		from:       from,
		to:         to,
		position:   []uint8{1},
		shouldFail: false,
	})

	from, err = bitsToIndex(1, 1, 0, 0)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	to, err = bitsToIndex(0, 1, 0, 0)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	tests = append(tests, deriveTest{
		name:       "not the same indexes",
		from:       from,
		to:         to,
		shouldFail: true,
	})

	from, err = bitsToIndex(1, 0, 1, 0)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	to, err = bitsToIndex(1, 0, 0, 0)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	tests = append(tests, deriveTest{
		name:       "'from' index greater then 'to' index",
		from:       from,
		to:         to,
		shouldFail: true,
	})

	from, err = bitsToIndex(1)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	to, err = bitsToIndex(1)
	if err != nil {
		t.Fatalf("can't generate from index: %v", err)
	}
	tests = append(tests, deriveTest{
		name:       "zero number trailing zeros",
		from:       from,
		to:         to,
		position:   nil,
		shouldFail: false,
	})

	return tests
}

// TestDeriveIndex check the correctness of index derive function by testing
// the index corner cases.
func TestDeriveIndex(t *testing.T) {
	for _, test := range generateTests(t) {
		pos, err := test.from.deriveBitTransformations(test.to)
		if err != nil {
			if !test.shouldFail {
				t.Fatalf("Failed (%v): %v", test.name, err)
			}
		} else {
			if test.shouldFail {
				t.Fatalf("Failed (%v): test should failed "+
					"but it's not", test.name)
			}

			if !reflect.DeepEqual(pos, test.position) {
				t.Fatalf("Failed(%v): position is wrong real:"+
					"%v expected:%v", test.name, pos, test.position)
			}
		}

		t.Logf("Passed: %v", test.name)

	}
}

var deriveElementTests = []struct {
	name       string
	index      index
	output     string
	seed       string
	shouldFail bool
}{
	{
		name:       "generate_from_seed FF alternate bits 1",
		index:      0xaaaaaaaaaaa,
		output:     "56f4008fb007ca9acf0e15b054d5c9fd12ee06cea347914ddbaed70d1c13a528",
		seed:       "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		shouldFail: false,
	},
	{
		name:       "generate_from_seed FF alternate bits 2",
		index:      0x555555555555,
		output:     "9015daaeb06dba4ccc05b91b2f73bd54405f2be9f217fbacd3c5ac2e62327d31",
		seed:       "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		shouldFail: false,
	},
	{
		name:       "generate_from_seed 01 last nontrivial node",
		index:      1,
		output:     "915c75942a26bb3a433a8ce2cb0427c29ec6c1775cfc78328b57f6ba7bfeaa9c",
		seed:       "0101010101010101010101010101010101010101010101010101010101010101",
		shouldFail: false,
	},
}

// TestSpecificationDeriveElement is used to check the consistency with
// specification hash derivation function.
func TestSpecificationDeriveElement(t *testing.T) {
	for _, test := range deriveElementTests {
		// Generate seed element.
		element, err := newElementFromStr(test.seed, rootIndex)
		if err != nil {
			t.Fatal(err)
		}

		// Derive element by index.
		result, err := element.derive(test.index)
		if err != nil {
			if !test.shouldFail {
				t.Fatalf("Failed (%v): %v", test.name, err)
			}
		} else {
			if test.shouldFail {
				t.Fatalf("Failed (%v): test should failed "+
					"but it's not", test.name)
			}

			// Generate element which we should get after deriviation.
			output, err := newElementFromStr(test.output, test.index)
			if err != nil {
				t.Fatal(err)
			}

			// Check that they are equal.
			if !result.isEqual(output) {
				t.Fatalf("Failed (%v): hash is wrong, real:"+
					"%v expected:%v", test.name,
					result.hash.String(), output.hash.String())
			}
		}

		t.Logf("Passed (%v)", test.name)
	}
}