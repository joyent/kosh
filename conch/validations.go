package conch

import "github.com/joyent/kosh/conch/types"

// GET /validation
func (c *Client) GetValidations() (validations types.Validations) {
	c.Validation().Receive(&validations)
	return
}

// GET /validation/:validation_id_or_name
func (c *Client) GetValidationByName(name string) (validation types.Validation) {
	c.Validation(name).Receive(&validation)
	return
}

func (c *Client) GetValidationByID(id types.UUID) (validation types.Validation) {
	c.Validation(id.String()).Receive(&validation)
	return
}

// GET /validation_plan
func (c *Client) GetAllValidationPlans() (plans types.ValidationPlans) {
	c.ValidationPlan().Receive(&plans)
	return
}

// GET /validation_plan/:validation_plan_id_or_name
func (c *Client) GetValidationPlanByName(name string) (plan types.ValidationPlan) {
	c.ValidationPlan(name).Receive(&plan)
	return
}

func (c *Client) GetValidationPlanByID(id types.UUID) (plan types.ValidationPlan) {
	c.ValidationPlan(id.String()).Receive(&plan)
	return
}

// GET /validation_plan/:validation_plan_id_or_name/validation
func (c *Client) GetValidationPlanValidations(id string) (validations types.Validations) {
	c.ValidationPlan(id).Validation().Receive(&validations)
	return
}

// GET /validation_state/:validation_state_id_or_name
func (c *Client) GetValidationStateByName(name string) (state types.ValidationStateWithResults) {
	c.ValidationState(name).Receive(&state)
	return
}

func (c *Client) GetValidationStateByID(id types.UUID) (state types.ValidationStateWithResults) {
	c.ValidationState(id.String()).Receive(&state)
	return
}
