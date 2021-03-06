package generator

import (
	"fmt"
)

// IntegerGenerator is a generator capable of generating integers and strings.
type IntegerGenerator interface {
	Generator
	// NextInt returns the next value as an int. When overriding this method,
	// be sure to call setLastString() properly, or the LastString() call
	// won't work.
	NextInt() int64

	// LastInt returns the previous int generated by the distribution.
	// This call is unique to IntegerGenerator implementation struct, and
	// assumes all implementation of this interface always return ints for
	// NextInt (e.g. not arbitrary values).
	LastInt() int64

	// Mean returns the expected value(mean) of the values this generator will
	// return.
	Mean() float64
}

// IntegerGeneratorBase is a parent class for all IntegerGenerator subclasses.
type IntegerGeneratorBase struct {
	lastInt int64
}

func NewIntegerGeneratorBase(last int64) *IntegerGeneratorBase {
	return &IntegerGeneratorBase{
		lastInt: last,
	}
}

// SetLastInt sets the last value to be generated.
// IntegerGenerator subclasses must use this call to properly set the last
// int value, or the LastString() and LastInt() calls won't work.
func (self *IntegerGeneratorBase) SetLastInt(value int64) {
	self.lastInt = value
}

// NextString generates the next string in the distribution.
func (self *IntegerGeneratorBase) NextString(g IntegerGenerator) string {
	return fmt.Sprintf("%d", g.NextInt())
}

func (self *IntegerGeneratorBase) LastInt() int64 {
	return self.lastInt
}

func (self *IntegerGeneratorBase) lastStringFrom(g IntegerGenerator) string {
	return fmt.Sprintf("%d", g.LastInt())
}

func (self *IntegerGeneratorBase) LastString() string {
	return fmt.Sprintf("%d", self.LastInt())
}

// ConstantIntegerGenerator is a trivial integer generator that always returns
// the same value.
type ConstantIntegerGenerator struct {
	*IntegerGeneratorBase
	value int64
}

func NewConstantIntegerGenerator(i int64) *ConstantIntegerGenerator {
	return &ConstantIntegerGenerator{
		IntegerGeneratorBase: NewIntegerGeneratorBase(i - 1),
		value:                i,
	}
}

func (self *ConstantIntegerGenerator) NextInt() int64 {
	return self.value
}

func (self *ConstantIntegerGenerator) NextString() string {
	return self.IntegerGeneratorBase.NextString(self)
}

func (self *ConstantIntegerGenerator) Mean() float64 {
	return float64(self.NextInt())
}

// Generate a popularity distribution of items, skewed to favor recent items
// significantly more than older items.
type SkewedLatestGenerator struct {
	*IntegerGeneratorBase
	basis   *CounterGenerator
	zipfian *ZipfianGenerator
}

func NewSkewedLatestGenerator(basis *CounterGenerator) *SkewedLatestGenerator {
	zipfian := NewZipfianGeneratorByInterval(0, basis.LastInt()-1)
	object := &SkewedLatestGenerator{
		IntegerGeneratorBase: NewIntegerGeneratorBase(0),
		basis:                basis,
		zipfian:              zipfian,
	}
	object.NextInt()
	return object
}

// Generate the next value in the distribution, skewed Zipfian favoring
// the items most recently returned by the basis generator.
func (self *SkewedLatestGenerator) NextInt() int64 {
	max := self.basis.LastInt()
	nextInt := max - self.zipfian.Next(max)
	self.SetLastInt(nextInt)
	return nextInt
}

func (self *SkewedLatestGenerator) NextString() string {
	return self.IntegerGeneratorBase.NextString(self)
}

func (self *SkewedLatestGenerator) Mean() float64 {
	panic("can't compute mean of non-stationary distribution")
}
