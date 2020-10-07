package conch

import "github.com/joyent/kosh/conch/types"

// GET /validation
func (c *Client) GetValidations() (validations types.Validations) {
	c.Validation().Receive(validations)
	return
}

// GET /validation/:validation_id_or_name
func (c *Client) GetValidationByID(id string) (validation types.Validation) {
	c.Validation(id).Receive(validation)
	return
}

// GET /validation_plan
func (c *Client) GetValidationPlans() (plans types.ValidationPlans) {
	c.ValidationPlan().Receive(plans)
	return
}

// GET /validation_plan/:validation_plan_id_or_name
func (c *Client) GetValidationPlanByID(id string) (plan types.ValidationPlan) {
	c.ValidationPlan(id).Receive(plan)
	return
}

// GET /validation_plan/:validation_plan_id_or_name/validation
func (c *Client) GetValidationPlanValidations(id string) (validations types.Validations) {
	c.ValidationPlan(id).Validation().Receive(validations)
	return
}

// GET /validation_state/:validation_state_id_or_name
func (c *Client) GetValidationStateByID(id string) (state types.ValidationStateWithResults) {
	c.ValidationState(id).Receive(state)
	return
}
