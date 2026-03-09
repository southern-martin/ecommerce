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

// WelcomeData holds data for the welcome email template.
type WelcomeData struct {
	Name    string
	Email   string
	ShopURL string
}

// PasswordResetData holds data for the password reset email template.
type PasswordResetData struct {
	Name      string
	ResetURL  string
	ExpiresIn string
}

// SellerApprovedData holds data for the seller approved email template.
type SellerApprovedData struct {
	Name         string
	StoreName    string
	DashboardURL string
}

// ReturnUpdateData holds data for the return update email template.
type ReturnUpdateData struct {
	BuyerName    string
	ReturnNumber string
	Status       string
	Reason       string
}

// PromotionData holds data for the promotion email template.
type PromotionData struct {
	Name           string
	PromotionTitle string
	Description    string
	CouponCode     string
	DiscountText   string
	ExpiresAt      string
	ShopURL        string
}

// RenderWelcome renders the welcome email template.
func RenderWelcome(data WelcomeData) (string, error) {
	return render("welcome.html", data)
}

// RenderPasswordReset renders the password reset email template.
func RenderPasswordReset(data PasswordResetData) (string, error) {
	return render("password_reset.html", data)
}

// RenderSellerApproved renders the seller approved email template.
func RenderSellerApproved(data SellerApprovedData) (string, error) {
	return render("seller_approved.html", data)
}

// RenderReturnUpdate renders the return update email template.
func RenderReturnUpdate(data ReturnUpdateData) (string, error) {
	return render("return_update.html", data)
}

// RenderPromotion renders the promotion email template.
func RenderPromotion(data PromotionData) (string, error) {
	return render("promotion.html", data)
}

func render(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
