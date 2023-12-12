package main

import (
	"errors"
	"final-project/cmd/data"
	"fmt"
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "home.page.gohtml", nil)
}
func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "login.page.gohtml", nil)
}

func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())

	// parse form post
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// get email and password from form post
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credential")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// check password
	validPassword, err := user.PasswordMatches(password)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credential")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !validPassword {
		msg := Message{
			To:      email,
			Subject: "Failed login attempt!",
			Data:    "Invalid Login attempt!"}

		app.sendEmail(msg)

		app.Session.Put(r.Context(), "error", "Invalid credential")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// okay, so log in user
	app.Session.Put(r.Context(), "UserID", user.ID)
	app.Session.Put(r.Context(), "User", user)

	app.Session.Put(r.Context(), "flash", "successful login!")
	// redirect the user
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "register.page.gohtml", nil)

}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	// create a user

	u := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		Active:    0,
	}

	_, err = u.Insert(u)
	if err != nil {
		app.Session.Put(r.Context(), "error", "unable to create user")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// send an activated email
	url := fmt.Sprintf("http://localhost:8080/activate?email=%s", u.Email)
	signedURL := GenerateTokenFromString(url)
	app.InfoLog.Println(signedURL)

	msg := Message{
		To:       u.Email,
		Subject:  "Activate you account",
		Template: "confirmation-email",
		Data:     template.HTML(signedURL),
	}

	app.sendEmail(msg)

	app.Session.Put(r.Context(), "flash", "confirmation email sent. Check you email.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	// subscribe the user to an account
}

func (app *Config) TestEmail(w http.ResponseWriter, r *http.Request) {
	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encryption:  "none",
		FromAddress: "info@company.com",
		FromName:    "info",
	}

	msg := Message{
		To:      "me@here.com",
		Subject: "Test Email",
		Data:    "Hello world",
	}

	m.sendMail(msg, make(chan error))
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate url
	url := r.RequestURI
	testUrl := fmt.Sprintf("http://localhost:%s%s", webPort, url)
	okay := VerifyToken(testUrl)

	if !okay {
		app.Session.Put(r.Context(), "error", "invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// activate account
	u, err := app.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.Session.Put(r.Context(), "error", "user not found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	u.Active = 1
	err = u.Update()
	if err != nil {
		app.Session.Put(r.Context(), "error", "unable to update user")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "flash", "account activated. you can now login")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	// generate an invoice

	// send an email with attachments

	// send an email with the invoice attached

	// subscribe the user to an account
}

func (app *Config) chooseSubscription(w http.ResponseWriter, r *http.Request) {
	if !app.Session.Exists(r.Context(), "UserID") {
		app.Session.Put(r.Context(), "warning", "you must login to see this page")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	plan, err := app.Models.Plan.GetAll()
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	dataMap := make(map[string]any)
	dataMap["plans"] = plan

	app.Render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}

func (app *Config) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	// get the id of the plan that is choosen
	id := r.URL.Query().Get("id")

	planID, err := strconv.Atoi(id)
	if err != nil {
		app.ErrorLog.Println("Error getting plan: ", err)
	}

	// get the plan from the database
	plan, err := app.Models.Plan.GetOne(planID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "unable to find plan")
		http.Redirect(w, r, "/members/plan", http.StatusSeeOther)
		return
	}

	// get user from the session
	user, ok := app.Session.Get(r.Context(), "User").(data.User)
	if !ok {
		app.Session.Put(r.Context(), "error", "log in first!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//NOTE: generate invoice and email it
	app.Wait.Add(1)
	func() {
		defer app.Wait.Done()

		invoice, err := app.getInvoice(user, plan)
		if err != nil {
			app.ErrorChan <- err
		}

		msg := Message{
			To:       user.Email,
			Subject:  "Your email",
			Data:     invoice,
			Template: "invoice",
		}

		app.sendEmail(msg)
	}()

	// Send an email with the manual attached
	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()

		pdf := app.generateManual(user, plan)
		err := pdf.OutputFileAndClose(fmt.Sprintf("./tmp/%d_manual.pdf", user.ID))
		if err != nil {
			app.ErrorChan <- err
			return
		}

		msg := Message{
			To:      user.Email,
			Subject: "Your Manual",
			Data:    "Your user manual is attached",
			AttachmentMap: map[string]string{
				"manual.pdf": fmt.Sprintf("./tmp/%d_manual.pdf", user.ID),
			},
		}

		app.sendEmail(msg)
		// test app error chan
		app.ErrorChan <- errors.New("Some custom error")
	}()

	//TODO Subscribe the user to an account
	err = app.Models.Plan.SubscribeUserToPlan(user, *plan)
	if err != nil {
		app.Session.Put(r.Context(), "error", "error subscribing to plan!")
		http.Redirect(w, r, "/members/plan", http.StatusSeeOther)
		return
	}

	u, err := app.Models.User.GetOne(user.ID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "error getting user from database")
		http.Redirect(w, r, "/members/plan", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "User", u)
	// redirect
	app.Session.Put(r.Context(), "flash", "subscribed!")
	http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
}

func (app *Config) generateManual(u data.User, plan *data.Plan) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)

	importer := gofpdi.NewImporter()

	time.Sleep(5 * time.Second)

	t := importer.ImportPage(pdf, "./pdf/manual.pdf", 1, "/MediaBox")
	pdf.AddPage()

	importer.UseImportedTemplate(pdf, t, 0, 0, 215.9, 0)

	pdf.SetX(75)
	pdf.SetY(150)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", u.FirstName, u.LastName), "", "C", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", plan.PlanName), "", "C", false)

	return pdf
}
func (app *Config) getInvoice(u data.User, plan *data.Plan) (string, error) {
	return plan.PlanAmountFormatted, nil
}
