package auth

import (
	"app/myproj/pkg/workflow"
	"errors"
	"time"
)

func NewAuthWorkflow(authService *AuthService) *workflow.WorkflowDefinition {
	return &workflow.WorkflowDefinition{
		ID:           "auth_workflow",
		InitialState: "phone_input",
		Timeout:      15 * time.Minute,
		States: map[workflow.State][]workflow.StateTransition{
			"phone_input": {
				{
					Event:   "submit_phone",
					ToState: "otp_verification",
					Action: func(data map[string]interface{}) error {
						phone := data["phone_number"].(string)
						otp, err := authService.GenerateOTP()
						if err != nil {
							return err
						}
						data["otp"] = otp
						return authService.SendOTP(phone, otp)
					},
				},
			},
			"otp_verification": {
				{
					Event:   "verify_otp",
					ToState: "pin_setup",
					Action: func(data map[string]interface{}) error {
						phone := data["phone_number"].(string)
						otp := data["submitted_otp"].(string)
						if !authService.VerifyOTP(phone, otp) {
							return errors.New("invalid OTP")
						}
						return nil
					},
					Guard: func(data map[string]interface{}) bool {
						phone := data["phone_number"].(string)
						var user User
						return errors.Is(authService.db.Where("phone_number = ?", phone).
							First(&user).Error, authService.db.Where("phone_number = ?", phone).
							First(&user).Error)
					},
				},
			},
			"pin_setup": {
				{
					Event:   "set_pin",
					ToState: "completed",
					Action: func(data map[string]interface{}) error {
						phone := data["phone_number"].(string)
						pin := data["pin"].(string)
						return authService.SetPIN(phone, pin)
					},
				},
			},
		},
	}
}
