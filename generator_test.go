package infiniteBitmask

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testGeneratorDatum struct {
	suite.Suite
	Generator *Generator
}

func (suite *testGeneratorDatum) SetupTest() {
	suite.Generator = NewGenerator()
}

func (suite *testGeneratorDatum) TestValueFromNameString() {
	assert.Equal(suite.T(), "1", suite.Generator.ValueFromName("start").String())
}

func (suite *testGeneratorDatum) TestValueFromNamesString() {
	const iterations = 32

	var name string
	allNames := make([]string, iterations)

	for i := 0; i < iterations; i++ {
		name += "-"
		allNames[i] = name
	}

	assert.Equal(
		suite.T(),
		strings.Repeat("1", iterations),
		suite.Generator.ValueFromNames(allNames...).String(),
	)
}

func (suite *testGeneratorDatum) TestValueFromNamesWithValueFromName() {
	one := suite.Generator.ValueFromName("one")
	two := suite.Generator.ValueFromName("two")
	oneAndTwo := suite.Generator.ValueFromNames("one", "two")

	assert.Equal(suite.T(), oneAndTwo.Number().Uint64(), one.Number().Uint64()|two.Number().Uint64())
}

func (suite *testGeneratorDatum) TestValueFromNamesEqualWithCloneAndCombine() {
	one := suite.Generator.ValueFromName("one")
	two := suite.Generator.ValueFromName("two")
	oneAndTwo := suite.Generator.ValueFromNames("one", "two")

	oneAndTwoCloned := one.Clone()
	oneAndTwoCloned.Combine(two)

	assert.Equal(suite.T(), true, oneAndTwo.Equal(oneAndTwoCloned))
	assert.Equal(suite.T(), oneAndTwo.String(), oneAndTwoCloned.String())
}

func (suite *testGeneratorDatum) TestGeneratorEqualWithClone() {
	suite.Generator.ValueFromName("one")
	suite.Generator.ValueFromName("two")
	suite.Generator.ValueFromName("three")

	generatorCloned := suite.Generator.Clone()

	assert.Equal(suite.T(), true, suite.Generator.Equal(generatorCloned))
	assert.Equal(suite.T(), suite.Generator.String(), generatorCloned.String())
}

func (suite *testGeneratorDatum) TestValueFromAllNames() {
	suite.Generator.ValueFromName("one")
	suite.Generator.ValueFromName("two")
	suite.Generator.ValueFromName("three")
	suite.Generator.ValueFromName("four")
	suite.Generator.ValueFromName("five")

	assert.Equal(suite.T(), uint64((1<<5)-1), suite.Generator.ValueFromAllNames().Number().Uint64())
}

func TestGenerator(t *testing.T) {
	suite.Run(t, new(testGeneratorDatum))
}
