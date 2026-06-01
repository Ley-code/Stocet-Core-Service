package email

import (
	"fmt"
	"strings"
	"time"
)

// AlertEmailData holds all the values substituted into the alert email template.
type AlertEmailData struct {
	Ticker        string
	AlertType     string // "price" | "percent_change" | "new_news"
	ConditionText string // human-readable: "Price went above ETB 1,200"
	Price         float64
	PercentChange float64
	TriggeredAt   time.Time
	DeepLink      string // e.g. "https://stocet.com/markets?symbol=CBE"
}

// BuildAlertEmail returns (subject, htmlBody) for a triggered alert email.
func BuildAlertEmail(d AlertEmailData) (string, string) {
	subject := fmt.Sprintf("[STOCET] Alert Triggered: %s", d.Ticker)

	changeColor := "#22c55e" // green
	changePrefix := "+"
	if d.PercentChange < 0 {
		changeColor = "#ef4444" // red
		changePrefix = ""
	}

	pctStr := fmt.Sprintf("%s%.2f%%", changePrefix, d.PercentChange)
	priceStr := fmt.Sprintf("%.2f", d.Price)
	timeStr := d.TriggeredAt.UTC().Format("02 Jan 2006 · 15:04 UTC")

	typeLabel := "Price Alert"
	switch d.AlertType {
	case "percent_change":
		typeLabel = "% Change Alert"
	case "new_news":
		typeLabel = "News Alert"
	}

	deepLink := d.DeepLink
	if deepLink == "" {
		deepLink = fmt.Sprintf("https://stocet.com/markets?symbol=%s", d.Ticker)
	}

	var b strings.Builder
	b.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Stocet Alert</title>
</head>
<body style="margin:0;padding:0;background-color:#080b0f;font-family:'Courier New',Courier,monospace;">
<table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:#080b0f;">
<tr><td align="center" style="padding:32px 16px;">
<table width="560" cellpadding="0" cellspacing="0" border="0" style="max-width:560px;width:100%;">

  <!-- HEADER -->
  <tr>
    <td style="background-color:#0d1117;border:1px solid #1e2a38;border-radius:8px 8px 0 0;padding:18px 24px;">
      <table width="100%" cellpadding="0" cellspacing="0" border="0">
        <tr>
          <td>
            <span style="font-family:'Courier New',monospace;font-size:15px;font-weight:bold;color:#3b82f6;letter-spacing:5px;">STOCET</span>
            <span style="font-family:'Courier New',monospace;font-size:9px;color:#374151;letter-spacing:3px;margin-left:8px;">TERMINAL</span>
          </td>
          <td align="right">
            <span style="font-family:'Courier New',monospace;font-size:9px;color:#4b5563;letter-spacing:2px;text-transform:uppercase;">`)
	b.WriteString(typeLabel)
	b.WriteString(`</span>
          </td>
        </tr>
      </table>
    </td>
  </tr>

  <!-- ALERT BODY -->
  <tr>
    <td style="background-color:#0a0d11;border-left:1px solid #1e2a38;border-right:1px solid #1e2a38;padding:28px 24px;">

      <!-- Ticker + price row -->
      <table width="100%" cellpadding="0" cellspacing="0" border="0" style="border-bottom:1px solid #1e2a38;padding-bottom:20px;margin-bottom:20px;">
        <tr>
          <td>
            <div style="font-family:'Courier New',monospace;font-size:32px;font-weight:bold;color:#f8fafc;letter-spacing:2px;">`)
	b.WriteString(d.Ticker)
	b.WriteString(`</div>
            <div style="font-family:'Courier New',monospace;font-size:11px;color:#6b7280;margin-top:4px;letter-spacing:1px;">ETHIOPIAN SECURITIES EXCHANGE</div>
          </td>
          <td align="right" style="vertical-align:top;">
            <div style="font-family:'Courier New',monospace;font-size:24px;font-weight:bold;color:#f8fafc;">ETB `)
	b.WriteString(priceStr)
	b.WriteString(`</div>
            <div style="font-family:'Courier New',monospace;font-size:13px;color:`)
	b.WriteString(changeColor)
	b.WriteString(`;text-align:right;margin-top:4px;">`)
	b.WriteString(pctStr)
	b.WriteString(` (24h)</div>
          </td>
        </tr>
      </table>

      <!-- Condition + time rows -->
      <table width="100%" cellpadding="0" cellspacing="0" border="0">
        <tr>
          <td style="padding:10px 0;border-bottom:1px solid #1a2332;">
            <div style="font-family:'Courier New',monospace;font-size:9px;color:#4b5563;letter-spacing:2px;text-transform:uppercase;margin-bottom:4px;">Alert Condition</div>
            <div style="font-family:'Courier New',monospace;font-size:13px;color:#d1d5db;">`)
	b.WriteString(d.ConditionText)
	b.WriteString(`</div>
          </td>
        </tr>
        <tr>
          <td style="padding:10px 0;">
            <div style="font-family:'Courier New',monospace;font-size:9px;color:#4b5563;letter-spacing:2px;text-transform:uppercase;margin-bottom:4px;">Triggered At</div>
            <div style="font-family:'Courier New',monospace;font-size:13px;color:#d1d5db;">`)
	b.WriteString(timeStr)
	b.WriteString(`</div>
          </td>
        </tr>
      </table>

    </td>
  </tr>

  <!-- CTA -->
  <tr>
    <td style="background-color:#0d1117;border:1px solid #1e2a38;border-top:none;border-radius:0 0 8px 8px;padding:24px;text-align:center;">
      <a href="`)
	b.WriteString(deepLink)
	b.WriteString(`"
         style="display:inline-block;background-color:#3b82f6;color:#ffffff;font-family:'Courier New',monospace;font-size:11px;font-weight:bold;letter-spacing:3px;text-transform:uppercase;text-decoration:none;padding:13px 36px;border-radius:4px;">
        VIEW ON STOCET &rarr;
      </a>
      <div style="margin-top:8px;">
        <a href="https://stocet.com/alerts"
           style="font-family:'Courier New',monospace;font-size:9px;color:#4b5563;text-decoration:none;letter-spacing:1px;">
          Manage your alerts
        </a>
      </div>
    </td>
  </tr>

  <!-- FOOTER -->
  <tr>
    <td style="padding:20px 0;text-align:center;">
      <div style="font-family:'Courier New',monospace;font-size:9px;color:#374151;line-height:1.6;">
        You are receiving this because you set up an alert on Stocet.<br>
        Ethiopian Capital Markets Terminal &mdash; stocet.com
      </div>
    </td>
  </tr>

</table>
</td></tr>
</table>
</body>
</html>`)

	return subject, b.String()
}
