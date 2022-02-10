package infiniteBitmask

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testValueDatum struct {
	suite.Suite
	Generator *Generator
	Value     *Value
}

func (suite *testValueDatum) SetupTest() {
	suite.Generator = NewGenerator()
	suite.Value = suite.Generator.newValue(bigZero)
}

func (suite *testValueDatum) TestString() {
	assert.Equal(suite.T(), "0", suite.Value.String())

	suite.Value.inner.number.SetUint64(32)

	assert.Equal(suite.T(), "100000", suite.Value.String())
}

func (suite *testValueDatum) TestClear() {
	suite.Value.inner.number.SetUint64(1e6)
	suite.Value.Clear()

	assert.Equal(suite.T(), uint64(0), suite.Value.inner.number.Uint64())
}

func (suite *testValueDatum) TestIsEmpty() {
	assert.Equal(suite.T(), true, suite.Value.IsEmpty())
}

func (suite *testValueDatum) TestClone() {
	valueCloned := suite.Value.Clone()

	assert.Equal(suite.T(), true, suite.Value.Equal(valueCloned))
	assert.Equal(suite.T(), suite.Value.String(), valueCloned.String())
}

func TestValue(t *testing.T) {
	suite.Run(t, new(testValueDatum))
}
