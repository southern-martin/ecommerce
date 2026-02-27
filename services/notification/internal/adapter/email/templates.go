package email

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/*.html
var templateFS embed.FS

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFS(templateFS, "templates/*.html"))
}

// OrderData holds the data used to render order-related email templates.
type OrderData struct {
	OrderNumber string
	OrderID     string
	TotalAmount string
	ItemCount   int
	BuyerName   string
}

// RenderOrderConfirmation renders the order confirmation email template.
func RenderOrderConfirmation(data OrderData) (string, error) {
	return render("order_confirmation.html", data)
}

// RenderOrderShipped renders the order shipped email template.
func RenderOrderShipped(data OrderData) (string, error) {
	return render("order_shipped.html", data)
}

// RenderOrderDelivered renders the order delivered email template.
func RenderOrderDelivered(data OrderData) (string, error) {
	return render("order_delivered.html", data)
}

// RenderPaymentConfirmation renders the payment confirmation email template.
func RenderPaymentConfirmation(data OrderData) (string, error) {
	return render("payment_confirmation.html", data)
}

func render(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
