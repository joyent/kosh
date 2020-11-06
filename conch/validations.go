package conch

import "github.com/joyent/kosh/conch/types"

// GetAllValidationPlans ( GET /validation_plan ) returns a list of all
// validations plans. See also
// https://joyent.github.io/conch-api/modules/Conch%3A%3ARoute%3A%3AValidation#get-validation_plans
func (c *Client) GetAllValidationPlans() (plans types.ValidationPlans, e error) {
	_, e = c.ValidationPlan().Receive(&plans)
	return
}

// GetValidationPlanByName (GET /validation_plan/:validation_plan_id_or_name)
// retrieves a single validation plan with the given name
func (c *Client) GetValidationPlanByName(name string) (plan types.ValidationPlan, e error) {
	_, e = c.ValidationPlan(name).Receive(&plan)
	return
}

// GetValidationPlanByID (GET /validation_plan/:validation_plan_id_or_name)
// retrieves a single validation plan with the given UUID
func (c *Client) GetValidationPlanByID(id types.UUID) (plan types.ValidationPlan, e error) {
	_, e = c.ValidationPlan(id.String()).Receive(&plan)
	return
}

// GetValidationPlanValidations (GET /validation_plan/:validation_plan_id_or_name/validation)
// get the list of validations for the named validation plan
func (c *Client) GetValidationPlanValidations(name string) (validations types.Validations, e error) {
	_, e = c.ValidationPlan(name).Validation().Receive(&validations)
	return
}

// GetValidationStateByName (GET /validation_state/:validation_state_id_or_name)
// retrieves the validation state with the given name
func (c *Client) GetValidationStateByName(name string) (state types.ValidationStateWithResults, e error) {
	_, e = c.ValidationState(name).Receive(&state)
	return
}

// GetValidationStateByID (GET /validation_state/:validation_state_id_or_name)
// retrieves the validation state with the given UUID
func (c *Client) GetValidationStateByID(id types.UUID) (state types.ValidationStateWithResults, e error) {
	_, e = c.ValidationState(id.String()).Receive(&state)
	return
}
