package auth

import (
	"fmt"
	"html"
)

func verifyOTPEmail(name, code string, minutes int) string {
	return fmt.Sprintf(`<!doctype html>
<html>
<body style="margin:0;padding:0;background:#f4f6f8;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6f8;padding:24px 0;">
    <tr><td align="center">
      <table role="presentation" width="480" cellpadding="0" cellspacing="0" style="max-width:480px;width:100%%;background:#ffffff;border-radius:14px;overflow:hidden;box-shadow:0 1px 4px rgba(0,0,0,.06);font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;">
        <tr><td style="background:#2f6feb;height:6px;font-size:0;line-height:0;">&nbsp;</td></tr>
        <tr><td style="padding:36px 40px 8px;">
          <h1 style="margin:0 0 10px;font-size:20px;color:#1c2230;">Verify your email</h1>
          <p style="margin:0 0 20px;font-size:14px;line-height:1.6;color:#5c6675;">Hi %s, enter this code to confirm your email address.</p>
        </td></tr>
        <tr><td align="center" style="padding:0 40px 8px;">
          <div style="display:inline-block;background:#f0f4ff;border:1px solid #dbe4ff;border-radius:12px;padding:16px 26px;font-size:34px;font-weight:700;letter-spacing:10px;color:#2f6feb;">%s</div>
        </td></tr>
        <tr><td style="padding:18px 40px 36px;">
          <p style="margin:0 0 6px;font-size:13px;color:#5c6675;">This code expires in %d minutes.</p>
          <p style="margin:0;font-size:12px;color:#9aa5b1;">Didn't create an account? You can safely ignore this email.</p>
        </td></tr>
      </table>
      <p style="margin:16px 0 0;font-size:11px;color:#9aa5b1;font-family:-apple-system,'Segoe UI',Roboto,Arial,sans-serif;">go-backend-boilerplate</p>
    </td></tr>
  </table>
</body>
</html>`, html.EscapeString(name), html.EscapeString(code), minutes)
}

func resetPasswordEmail(name, link string) string {
	safeLink := html.EscapeString(link)
	return fmt.Sprintf(`<!doctype html>
<html>
<body style="margin:0;padding:0;background:#f4f6f8;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6f8;padding:24px 0;">
    <tr><td align="center">
      <table role="presentation" width="480" cellpadding="0" cellspacing="0" style="max-width:480px;width:100%%;background:#ffffff;border-radius:14px;overflow:hidden;box-shadow:0 1px 4px rgba(0,0,0,.06);font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;">
        <tr><td style="background:#2f6feb;height:6px;font-size:0;line-height:0;">&nbsp;</td></tr>
        <tr><td style="padding:36px 40px 8px;">
          <h1 style="margin:0 0 10px;font-size:20px;color:#1c2230;">Reset your password</h1>
          <p style="margin:0 0 24px;font-size:14px;line-height:1.6;color:#5c6675;">Hi %s, click the button below to set a new password.</p>
        </td></tr>
        <tr><td align="center" style="padding:0 40px 8px;">
          <a href="%s" style="display:inline-block;background:#2f6feb;color:#ffffff;text-decoration:none;font-size:15px;font-weight:600;padding:12px 28px;border-radius:10px;">Reset password</a>
        </td></tr>
        <tr><td style="padding:22px 40px 36px;">
          <p style="margin:0 0 6px;font-size:12px;color:#9aa5b1;">Or paste this link into your browser:</p>
          <p style="margin:0 0 12px;font-size:12px;color:#2f6feb;word-break:break-all;">%s</p>
          <p style="margin:0;font-size:12px;color:#9aa5b1;">Didn't request this? You can safely ignore this email.</p>
        </td></tr>
      </table>
      <p style="margin:16px 0 0;font-size:11px;color:#9aa5b1;font-family:-apple-system,'Segoe UI',Roboto,Arial,sans-serif;">go-backend-boilerplate</p>
    </td></tr>
  </table>
</body>
</html>`, html.EscapeString(name), safeLink, safeLink)
}
