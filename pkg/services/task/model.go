//-- Package Declaration -----------------------------------------------------------------------------------------------
package task

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

//-- Constants ---------------------------------------------------------------------------------------------------------

//-- Structs -----------------------------------------------------------------------------------------------------------
type Task struct {
	//-- Primary Key ----------
	ID uint

	//-- User Variables ----------
	Name       string
	Details    *string
	ResolvedAt *time.Time

	//-- System Variables ----------

	//-- Relations ----------

	//-- Automated fields (Timestamps) ----------
	CreatedAt time.Time
	UpdatedAt *time.Time
}

//-- Exported Functions ------------------------------------------------------------------------------------------------
func (task Task) String() string {
	var details, resolvedAt, updatedAt = `<nil>`, `<nil>`, `<nil>`

	if task.Details != nil {
		details = *task.Details
	}
	if task.ResolvedAt != nil {
		resolvedAt = task.ResolvedAt.String()
	}
	if task.UpdatedAt != nil {
		updatedAt = task.UpdatedAt.String()
	}

	return fmt.Sprintf(`{ID: %d, Name: %s, Details: %s, ResolvedAt: %s, CreatedAt: %s, UpdatedAt: %s}`, task.ID, task.Name, details, resolvedAt, task.CreatedAt, updatedAt)
}

//-- Store Functions ---------------------------------------------------------------------------------------------------
func (task Task) compare(other Task) bool {
	if task.ID != other.ID {
		return false
	}

	if task.Name != other.Name {
		return false
	}

	if (task.Details == nil && other.Details != nil) || (task.Details != nil && other.Details == nil) {
		return false
	} else if task.Details != nil && other.Details != nil && *task.Details != *other.Details {
		return false
	}

	if (task.ResolvedAt == nil && other.ResolvedAt != nil) || (task.ResolvedAt != nil && other.ResolvedAt == nil) {
		return false
	} else if task.ResolvedAt != nil && other.ResolvedAt != nil && task.ResolvedAt.Unix() != other.ResolvedAt.Unix() {
		return false
	}

	if task.CreatedAt.Unix() != other.CreatedAt.Unix() {
		return false
	}

	if (task.UpdatedAt == nil && other.UpdatedAt != nil) || (task.UpdatedAt != nil && other.UpdatedAt == nil) {
		return false
	} else if task.UpdatedAt != nil && other.UpdatedAt != nil && task.UpdatedAt.Unix() != other.UpdatedAt.Unix() {
		return false
	}

	return true
}

func (task *Task) sanitize() error {
	if task.ID == 0 {
		task.UpdatedAt = nil
	}

	if task.Details != nil && len(*task.Details) == 0 {
		task.Details = nil
	}

	if task.ResolvedAt != nil {
		*task.ResolvedAt = task.ResolvedAt.UTC()
	}

	task.CreatedAt = task.CreatedAt.UTC()

	if task.UpdatedAt != nil {
		*task.UpdatedAt = task.UpdatedAt.UTC()
	}

	return nil
}

func (task Task) validate() error {
	if err := task.validateName(); err != nil {
		return err
	}

	if err := task.validateDetails(); err != nil {
		return err
	}

	if err := task.validateUpdatedAt(); err != nil {
		return err
	}

	return nil
}

//-- Internal Functions ------------------------------------------------------------------------------------------------
func (task Task) validateName() error {
	//-- Common variables ----------
	var validPattern = regexp.MustCompile(`\A[a-zA-Z0-9 \-:]{1,50}\z`)

	//-- Check for pattern adherence ----------
	if !validPattern.MatchString(task.Name) {
		return errors.New(fmt.Sprintf(`validation - Name '%s' must be comprised only of letters, numbers, spaces and hyphens/colons and may not be empty and may not exceed 50 characters`, task.Name))
	}

	return nil
}

func (task Task) validateDetails() error {
	//-- Common variables ----------
	var validPattern = regexp.MustCompile(`\A[a-zA-Z0-9 \-:]{0,512}\z`)

	//-- Check for pattern adherence ----------
	if task.Details != nil && !validPattern.MatchString(*task.Details) {
		return errors.New(fmt.Sprintf(`validation - Details '%s' must be comprised only of letters, numbers, spaces and hyphens/colons and may not exceed 512 characters`, *task.Details))
	}

	return nil
}

func (task Task) validateUpdatedAt() error {
	//-- Check for non-sensical value ----------
	if task.ID == 0 && task.UpdatedAt != nil {
		return errors.New(fmt.Sprintf(`validation - UpdatedAt '%s' should not be present when there is no pre-assigned ID`, *task.UpdatedAt))
	}

	if task.UpdatedAt != nil && task.CreatedAt.Unix() > task.UpdatedAt.Unix() {
		return errors.New(fmt.Sprintf(`validation - UpdatedAt '%s' should not occure before the CreatedAt timestamp`, *task.UpdatedAt))
	}

	return nil
}
