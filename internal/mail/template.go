package mail

import "fmt"

func OTPTemplate(appName, otp string, expiryMinutes int) string {
	return fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 24px; border: 1px solid #eaeaea; border-radius: 8px;">
			
			<h2 style="color: #2c3e50;">
				Hello from <span style="color: #0070f3;">%s</span> 👋
			</h2>

			<p style="font-size: 16px; color: #333;">
				To keep your account secure, we’ve generated a One-Time Password (OTP) for you.
			</p>

			<p style="font-size: 16px; color: #333; margin-bottom: 8px;">
				Please use the following OTP to complete your verification:
			</p>

			<div style="font-size: 28px; font-weight: bold; color: #0070f3; padding: 12px 0;">
				%s
			</div>

			<p style="font-size: 14px; color: #666;">
				This OTP is valid for <strong>%d minutes</strong>. Do not share this code with anyone—even if they claim to be from %s.
			</p>

			<hr style="margin: 24px 0; border: none; border-top: 1px solid #eee;" />

			<p style="font-size: 12px; color: #999;">
				If you didn’t request this OTP, please ignore this email or contact our support team.
			</p>

			<p style="font-size: 12px; color: #999;">
				— The %s Team
			</p>

		</div>
	`, appName, otp, expiryMinutes, appName, appName)
}
