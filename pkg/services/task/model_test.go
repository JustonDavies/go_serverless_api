//-- Package Declaration -----------------------------------------------------------------------------------------------
package task

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//-- Local Constants ---------------------------------------------------------------------------------------------------

//-- Decorators --------------------------------------------------------------------------------------------------------

//-- Helpers -----------------------------------------------------------------------------------------------------------
func newValidTask() *Task {
	var id uint = 0
	var name = `Test valid Task`
	var details = `The default valid Task structure used in testing`

	return &Task{
		ID:      id,
		Name:    name,
		Details: &details,
	}
}

//-- Tests -------------------------------------------------------------------------------------------------------------
func TestModelString(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var result string

	//-- Test Parameters ----------
	var name = `Test string`
	var details = `Test the string method with a clean model`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	//-- Action ----------
	result = model.String()

	//-- Post-conditions ----------
	assert.True(test, len(result) > 0)
}

func TestModelStringWithNils(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var result string

	//-- Test Parameters ----------
	var name = `Test string with nils`

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name

	//-- Action ----------
	result = model.String()

	//-- Post-conditions ----------
	assert.True(test, len(result) > 0)
}

func TestModelCompare(test *testing.T) {
	//-- Shared Variables ----------
	var model, other *Task
	var result bool

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test compare`
	var details = `Test the compare method with identical models`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	other = newValidTask()
	other.ID = id
	other.Name = name
	other.Details = &details
	other.ResolvedAt = &resolvedAt
	other.CreatedAt = createdAt
	other.UpdatedAt = &updatedAt

	//-- Action ----------
	result = model.compare(*other)

	//-- Post-conditions ----------
	assert.True(test, result)
}

func TestModelCompareWithNils(test *testing.T) {
	//-- Shared Variables ----------
	var model, other *Task
	var result bool

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test compare with nils`
	var createdAt = time.Now()

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = nil
	model.ResolvedAt = nil
	model.CreatedAt = createdAt
	model.UpdatedAt = nil

	other = newValidTask()
	other.ID = id
	other.Name = name
	other.Details = nil
	other.ResolvedAt = nil
	other.CreatedAt = createdAt
	other.UpdatedAt = nil

	//-- Action ----------
	result = model.compare(*other)

	//-- Post-conditions ----------
	assert.True(test, result)
}

func TestModelCompareDifferentID(test *testing.T) {
	//-- Shared Variables ----------
	var model, other *Task
	var result bool

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test id compare`
	var details = `Test the compare method with non-identical models`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	var otherAttr uint = 2

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	other = newValidTask()
	other.ID = otherAttr
	other.Name = name
	other.Details = &details
	other.ResolvedAt = &resolvedAt
	other.CreatedAt = createdAt
	other.UpdatedAt = &updatedAt

	//-- Action ----------
	result = model.compare(*other)

	//-- Post-conditions ----------
	assert.False(test, result)
}

func TestModelCompareDifferentName(test *testing.T) {
	//-- Shared Variables ----------
	var model, other *Task
	var result bool

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test name compare`
	var details = `Test the compare method with non-identical models`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	var otherAttr = `Test different name compare`

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	other = newValidTask()
	other.ID = id
	other.Name = otherAttr
	other.Details = &details
	other.ResolvedAt = &resolvedAt
	other.CreatedAt = createdAt
	other.UpdatedAt = &updatedAt

	//-- Action ----------
	result = model.compare(*other)

	//-- Post-conditions ----------
	assert.False(test, result)
}

func TestModelCompareDifferentDetails(test *testing.T) {
	//-- Shared Variables ----------
	var model, other *Task
	var result bool

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test details compare`
	var details = `Test the compare method with non-identical models`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	var otherAttr = `Test different details compare`

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	other = newValidTask()
	other.ID = id
	other.Name = name
	other.Details = &otherAttr
	other.ResolvedAt = &resolvedAt
	other.CreatedAt = createdAt
	other.UpdatedAt = &updatedAt

	//-- Action ----------
	result = model.compare(*other)

	//-- Post-conditions ----------
	assert.False(test, result)
}

func TestModelCompareNilDetails(test *testing.T) {
	//-- Shared Variables ----------
	var model, other *Task
	var result bool

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test nil details compare`
	var details = `Test the compare method with non-identical models`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	var otherAttr *string = nil

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	other = newValidTask()
	other.ID = id
	other.Name = name
	other.Details = otherAttr
	other.ResolvedAt = &resolvedAt
	other.CreatedAt = createdAt
	other.UpdatedAt = &updatedAt

	//-- Action ----------
	result = model.compare(*other)

	//-- Post-conditions ----------
	assert.False(test, result)
}

func TestModelCompareDifferentResolvedAt(test *testing.T) {
	//-- Shared Variables ----------
	var model, other *Task
	var result bool

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test resolved at compare`
	var details = `Test the compare method with non-identical models`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	var otherAttr = time.Unix(0, 0)

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	other = newValidTask()
	other.ID = id
	other.Name = name
	other.Details = &details
	other.ResolvedAt = &otherAttr
	other.CreatedAt = createdAt
	other.UpdatedAt = &updatedAt

	//-- Action ----------
	result = model.compare(*other)

	//-- Post-conditions ----------
	assert.False(test, result)
}

func TestModelCompareNilResolvedAt(test *testing.T) {
	//-- Shared Variables ----------
	var model, other *Task
	var result bool

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test nil resolved at compare`
	var details = `Test the compare method with non-identical models`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	var otherAttr *time.Time = nil

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	other = newValidTask()
	other.ID = id
	other.Name = name
	other.Details = &details
	other.ResolvedAt = otherAttr
	other.CreatedAt = createdAt
	other.UpdatedAt = &updatedAt

	//-- Action ----------
	result = model.compare(*other)

	//-- Post-conditions ----------
	assert.False(test, result)
}

func TestModelCompareDifferentCreatedAt(test *testing.T) {
	//-- Shared Variables ----------
	var model, other *Task
	var result bool

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test created at compare`
	var details = `Test the compare method with non-identical models`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	var otherAttr = time.Unix(0, 0)

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	other = newValidTask()
	other.ID = id
	other.Name = name
	other.Details = &details
	other.ResolvedAt = &resolvedAt
	other.CreatedAt = otherAttr
	other.UpdatedAt = &updatedAt

	//-- Action ----------
	result = model.compare(*other)

	//-- Post-conditions ----------
	assert.False(test, result)
}

func TestModelCompareDifferentUpdatedAt(test *testing.T) {
	//-- Shared Variables ----------
	var model, other *Task
	var result bool

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test updated at compare`
	var details = `Test the compare method with non-identical models`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	var otherAttr = time.Unix(0, 0)

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	other = newValidTask()
	other.ID = id
	other.Name = name
	other.Details = &details
	other.ResolvedAt = &resolvedAt
	other.CreatedAt = createdAt
	other.UpdatedAt = &otherAttr

	//-- Action ----------
	result = model.compare(*other)

	//-- Post-conditions ----------
	assert.False(test, result)
}

func TestModelCompareNilUpdatedAt(test *testing.T) {
	//-- Shared Variables ----------
	var model, other *Task
	var result bool

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test nil updated at compare`
	var details = `Test the compare method with non-identical models`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	var otherAttr *time.Time = nil

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	other = newValidTask()
	other.ID = id
	other.Name = name
	other.Details = &details
	other.ResolvedAt = &resolvedAt
	other.CreatedAt = createdAt
	other.UpdatedAt = otherAttr

	//-- Action ----------
	result = model.compare(*other)

	//-- Post-conditions ----------
	assert.False(test, result)
}

func TestModelSanitize(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var sanitizeErr error

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test sanitize`
	var details = `Test the sanitize method with a clean model`
	var resolvedAt = time.Now()

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.Details = &details
	model.ResolvedAt = &resolvedAt

	//-- Action ----------
	sanitizeErr = model.sanitize()

	//-- Post-conditions ----------
	assert.Nil(test, sanitizeErr)
	assert.Equal(test, id, model.ID)
	assert.Equal(test, name, model.Name)
	assert.Equal(test, details, *model.Details)
	assert.Equal(test, resolvedAt, *model.ResolvedAt)
}

func TestModelSanitizeDetailsNil(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var sanitizeErr error

	//-- Test Parameters ----------
	var name = `Test sanitize of zero-value details`
	var details *string = nil

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name
	model.Details = details

	//-- Action ----------
	sanitizeErr = model.sanitize()

	//-- Post-conditions ----------
	assert.Nil(test, sanitizeErr)
	assert.Equal(test, name, model.Name)
	assert.Equal(test, details, model.Details)
}

func TestModelSanitizeDetailsZeroValue(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var sanitizeErr error

	//-- Test Parameters ----------
	var name = `Test sanitize of zero-value details`
	var details = ``

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name
	model.Details = &details

	//-- Action ----------
	sanitizeErr = model.sanitize()

	//-- Post-conditions ----------
	assert.Nil(test, sanitizeErr)
	assert.Equal(test, name, model.Name)
	assert.Nil(test, model.Details)
}

func TestModelSanitizeTimezones(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var sanitizeErr error

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Test sanitize of timezones`
	var resolvedAt = time.Now()
	var createdAt = time.Now()
	var updatedAt = time.Now()

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.Name = name
	model.ResolvedAt = &resolvedAt
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	assert.Equal(test, `Local`, model.ResolvedAt.Location().String())
	assert.Equal(test, `Local`, model.CreatedAt.Location().String())
	assert.Equal(test, `Local`, model.UpdatedAt.Location().String())

	//-- Action ----------
	sanitizeErr = model.sanitize()

	//-- Post-conditions ----------
	assert.Nil(test, sanitizeErr)
	assert.Equal(test, name, model.Name)
	assert.Equal(test, `UTC`, model.ResolvedAt.Location().String())
	assert.Equal(test, `UTC`, model.CreatedAt.Location().String())
	assert.Equal(test, `UTC`, model.UpdatedAt.Location().String())

	assert.Equal(test, resolvedAt.Unix(), model.ResolvedAt.Unix())
	assert.Equal(test, createdAt.Unix(), model.CreatedAt.Unix())
	assert.Equal(test, updatedAt.Unix(), model.UpdatedAt.Unix())

}

func TestModelSanitizeNilTimezones(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var sanitizeErr error

	//-- Test Parameters ----------
	var name = `Test sanitize of timezones with nils`

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name

	//-- Action ----------
	sanitizeErr = model.sanitize()

	//-- Post-conditions ----------
	assert.Nil(test, sanitizeErr)
	assert.Equal(test, name, model.Name)
	assert.Nil(test, model.ResolvedAt)
	assert.Nil(test, model.UpdatedAt)

}

func TestModelValidateValid(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var name = `Test validate valid`

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name

	//-- Action ----------
	validationErr = model.validate()

	//-- Post-conditions ----------
	assert.Nil(test, validationErr)
}

func TestModelValidateNotValidName(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var name = `Test invalid Name ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name

	//-- Action ----------
	validationErr = model.validate()

	//-- Post-conditions ----------
	assert.NotNil(test, validationErr)
}

func TestModelValidateNotValidDetails(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var name = `Test invalid Details`
	var details = `Test the validation method with an invalid Details attribute ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name
	model.Details = &details

	//-- Action ----------
	validationErr = model.validate()

	//-- Post-conditions ----------
	assert.NotNil(test, validationErr)
}

func TestModelValidateNotValidUpdatedAt(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var name = `Test invalid UpdatedAt`
	var createdAt = time.Now()
	var updatedAt = time.Unix(0, 0)

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	//-- Action ----------
	validationErr = model.validate()

	//-- Post-conditions ----------
	assert.NotNil(test, validationErr)
}

func TestModelValidateNameValid(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var name = `Test valid name`

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name

	//-- Action ----------
	validationErr = model.validateName()

	//-- Post-conditions ----------
	assert.Nil(test, validationErr)
}

func TestModelValidateNameZeroValue(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var name = ``

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name

	//-- Action ----------
	validationErr = model.validateName()

	//-- Post-conditions ----------
	assert.NotNil(test, validationErr)
}

func TestModelValidateNameNotValid(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var name = `Test invalid name ~~~`

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Name = name

	//-- Action ----------
	validationErr = model.validateName()

	//-- Post-conditions ----------
	assert.NotNil(test, validationErr)
}

func TestModelValidateDetailsValid(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var details = `Test the details validation method with a valid attribute value`

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Details = &details

	//-- Action ----------
	validationErr = model.validateDetails()

	//-- Post-conditions ----------
	assert.Nil(test, validationErr)
}

func TestModelValidateDetailsZeroValue(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var details = ``

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Details = &details

	//-- Action ----------
	validationErr = model.validateDetails()

	//-- Post-conditions ----------
	assert.Nil(test, validationErr)
}

func TestModelValidateDetailsNil(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var details *string = nil

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Details = details

	//-- Action ----------
	validationErr = model.validateDetails()

	//-- Post-conditions ----------
	assert.Nil(test, validationErr)
}

func TestModelValidateDetailsNotValid(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var details = `Test the details validation method with an invalid attribute value ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`

	//-- Pre-conditions ----------
	model = newValidTask()
	model.Details = &details

	//-- Action ----------
	validationErr = model.validateDetails()

	//-- Post-conditions ----------
	assert.NotNil(test, validationErr)
}

func TestModelValidateUpdatedAtValid(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var id uint = 1
	var createdAt = time.Now()
	var updatedAt = time.Now()

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	//-- Action ----------
	validationErr = model.validateUpdatedAt()

	//-- Post-conditions ----------
	assert.Nil(test, validationErr)
}

func TestModelValidateUpdatedAtWithoutID(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var id uint = 0
	var updatedAt = time.Now()

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.UpdatedAt = &updatedAt

	//-- Action ----------
	validationErr = model.validateUpdatedAt()

	//-- Post-conditions ----------
	assert.NotNil(test, validationErr)
}

func TestModelValidateUpdatedAtPreDated(test *testing.T) {
	//-- Shared Variables ----------
	var model *Task
	var validationErr error

	//-- Test Parameters ----------
	var id uint = 1
	var createdAt = time.Now()
	var updatedAt = time.Unix(0, 0)

	//-- Pre-conditions ----------
	model = newValidTask()
	model.ID = id
	model.CreatedAt = createdAt
	model.UpdatedAt = &updatedAt

	//-- Action ----------
	validationErr = model.validateUpdatedAt()

	//-- Post-conditions ----------
	assert.NotNil(test, validationErr)
}
