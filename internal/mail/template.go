package mail

import "fmt"

func WelcomeTemplate(name string) string {
	return fmt.Sprintf(`
		<h1>Welcome %s</h1>
		<p>Your account has been created successfully.</p>
	`, name)
}
