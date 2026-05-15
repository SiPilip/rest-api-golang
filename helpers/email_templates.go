package helpers

import "fmt"

func WelcomeEmailBody(email string) string {
	return fmt.Sprintf(`
    <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #333;">Welcome to Event API! 🎉</h1>
        <p>Hi <strong>%s</strong>,</p>
        <p>Your account has been created successfully.</p>
        <p>You can now create events, register for events, and more.</p>
        <hr style="border: 1px solid #eee;">
        <p style="color: #999; font-size: 12px;">This is an automated email. Please do not reply.</p>
    </div>
    `, email)
}

func EventRegistrationEmailBody(email, eventName string) string {
    return fmt.Sprintf(`
    <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #333;">Registration Confirmed! ✅</h1>
        <p>Hi <strong>%s</strong>,</p>
        <p>You have successfully registered for:</p>
        <div style="background: #f5f5f5; padding: 15px; border-radius: 8px; margin: 15px 0;">
            <h2 style="margin: 0; color: #2196F3;">%s</h2>
        </div>
        <p>See you there!</p>
        <hr style="border: 1px solid #eee;">
        <p style="color: #999; font-size: 12px;">This is an automated email. Please do not reply.</p>
    </div>
    `, email, eventName)
}