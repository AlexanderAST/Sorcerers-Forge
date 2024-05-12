package configs

import (
	"Diploma/internal/model"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
	"time"
)

func GenToken(Email string) string {
	var dateNow = time.Now().Truncate(30 * time.Minute).String()

	hashes := sha256.New()
	hashes.Write([]byte(Email + dateNow))
	hash := hashes.Sum(nil)

	hashString := hex.EncodeToString(hash)

	return hashString[len(hashString)-7 : len(hashString)-1]
}

func RunSmtp(Email string) (string, error) {

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env %v", err)
	}

	emCode := GenToken(Email)

	auth := smtp.PlainAuth("", "forgesorcerers@gmail.com", os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "forgesorcerers@gmail.com", []string{Email}, []byte(
		"Subject: Recovery Code\r\n\r\n"+
			emCode))
	if err != nil {
		return "can't send code to email", err
	}
	return "ok", err
}

func RunSmtpRegister(Email string) (string, error) {

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env %v", err)
	}

	emCode := GenToken(Email)

	auth := smtp.PlainAuth("", "forgesorcerers@gmail.com", os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "forgesorcerers@gmail.com", []string{Email}, []byte(
		"Subject: Registration Code\r\n\r\n"+
			emCode))
	if err != nil {
		return "can't send code to email", err
	}
	return "ok", err
}

func RunSmtpHelp(c *model.Message) (string, error) {

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env %v", err)
	}

	auth := smtp.PlainAuth("", "forgesorcerers@gmail.com", os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "forgesorcerers@gmail.com", []string{"sorcerer.forgehelp@mail.ru"}, []byte(
		"Subject: Обращение в тех. поддержку\r\n\r\n"+
			"Обращение от пользователя: "+c.Name+" "+c.UserEmail+"\r\n\r\n"+
			"Сообщениие: "+c.Message))
	if err != nil {
		return "can't send code to email", err
	}

	auth1 := smtp.PlainAuth("", "forgesorcerers@gmail.com", os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")
	err1 := smtp.SendMail("smtp.gmail.com:587", auth1, "forgesorcerers@gmail.com", []string{c.UserEmail}, []byte(
		"Subject: Обращение в тех. поддержку\r\n\r\n"+
			"Сообщениие: "+"Ваша заявка успешно зарегистрированна в системе, ожидайте ответа от менеджера"))

	if err1 != nil {
		return "can't send code to email", err1
	}

	return "ok", err
}

func RunSmtpOrders(email string, cart []*model.Cart) (string, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env %v", err)
	}

	var messageBody bytes.Buffer
	orderNumber := GenToken(email)

	pdfPath, err := generatePDF(cart, orderNumber)
	if err != nil {
		return "", err
	}

	pdfData, err := ioutil.ReadFile(pdfPath)
	if err != nil {
		return "", err
	}

	messageBody.WriteString(fmt.Sprintf("Subject: Ваш заказ #%s\r\n\r\n", strings.ToUpper(orderNumber)))
	messageBody.WriteString("Ваш заказ:\r\n")

	for _, carts := range cart {
		messageBody.WriteString(fmt.Sprintf("%s %d шт. %d руб.\n", carts.ProductName, carts.Count, carts.Price))
	}

	messageBody.WriteString(fmt.Sprintf("Сумма заказа: %d руб.\n", calculateTotalPrice(cart)))

	auth := smtp.PlainAuth("", "forgesorcerers@gmail.com", os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")
	to := []string{email}
	msg := []byte("To: " + email + "\r\n")
	msg = append(msg, []byte("Subject: Ваш заказ #"+strings.ToUpper(orderNumber)+"\r\n")...)
	msg = append(msg, []byte("Content-Type: multipart/mixed; boundary=boundarystring\r\n\r\n")...)
	msg = append(msg, []byte("--boundarystring\r\n")...)
	msg = append(msg, []byte("Content-Type: text/plain; charset=utf-8\r\n\r\n")...)
	msg = append(msg, []byte("Ваш заказ:\r\n")...)

	for _, carts := range cart {
		msg = append(msg, []byte(fmt.Sprintf("%s %d шт. %d руб.\n", carts.ProductName, carts.Count, carts.Price))...)
	}
	fileName := fmt.Sprintf("Заказ #%s.pdf", orderNumber)
	msg = append(msg, []byte(fmt.Sprintf("Сумма заказа: %d руб.\n", calculateTotalPrice(cart)))...)
	msg = append(msg, []byte("--boundarystring\r\n")...)
	msg = append(msg, []byte("Content-Type: application/pdf\r\n")...)
	msg = append(msg, []byte(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", fileName))...)
	msg = append(msg, pdfData...)
	msg = append(msg, []byte("\r\n--boundarystring--\r\n")...)

	err = smtp.SendMail("smtp.gmail.com:587", auth, "forgesorcerers@gmail.com", to, msg)
	if err != nil {
		return "", fmt.Errorf("can't send email: %v", err)
	}

	err = os.Remove(pdfPath)
	if err != nil {
		return "", fmt.Errorf("error removing PDF file: %v", err)
	}

	return "ok", nil
}

func calculateTotalPrice(products []*model.Cart) int {
	totalPrice := 0
	for _, product := range products {
		totalPrice += product.Price
	}
	return totalPrice
}

func generatePDF(products []*model.Cart, orderNumber string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("Arial", "", "/static/font/arial.ttf")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Содержимое вашего заказа:")
	pdf.Ln(10)

	for _, product := range products {
		pdf.Cell(0, 10, fmt.Sprintf("%s - %d руб.", product.ProductName, product.Price))
		pdf.Ln(10)
	}
	pdf.Cell(40, 10, fmt.Sprintf("Сумма заказа: %d руб.\n", calculateTotalPrice(products)))
	pdf.Ln(10)

	fileName := fmt.Sprintf("Заказ #%s.pdf", orderNumber)

	pdfPath := "static/pdf/" + fileName
	err := pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		return "", err
	}

	return pdfPath, nil
}

func RunSmtpLearn(c *model.Message) (string, error) {

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env %v", err)
	}

	auth := smtp.PlainAuth("", "forgesorcerers@gmail.com", os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "forgesorcerers@gmail.com", []string{"sorcerer.forgehelp@mail.ru"}, []byte(
		"Subject: Заявка на прохождение обучения\r\n\r\n"+
			"Обращение от пользователя: "+c.Name+" "+c.UserEmail+"\r\n\r\n"+
			"Сообщениие: "+c.Message))
	if err != nil {
		return "can't send code to email", err
	}

	auth1 := smtp.PlainAuth("", "forgesorcerers@gmail.com", os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")
	err1 := smtp.SendMail("smtp.gmail.com:587", auth1, "forgesorcerers@gmail.com", []string{c.UserEmail}, []byte(
		"Subject: Заявка на прохождение обучения\r\n\r\n"+
			"Сообщениие: "+"Ваша заявка успешно зарегистрированна в системе, ожидайте ответа от менеджера"))

	if err1 != nil {
		return "can't send code to email", err1
	}

	return "ok", err
}
