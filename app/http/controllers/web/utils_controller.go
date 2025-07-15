package web

import (
	"fmt"
	"net/http"
	"os"
	"semita/app/helpers"
	"semita/config"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
	"gopkg.in/gomail.v2"
)

func IndexPDF(c *gin.Context) {
	authData := helpers.AuthSessionService(c.Writer, c.Request, "PDF", nil)
	tmpl := template.Must(template.ParseFiles("resources/utils/pdf.html", config.MainLayoutFilePath))
	err := tmpl.Execute(c.Writer, authData)
	if err != nil {
		fmt.Println("Error al ejecutar la plantilla:", err)
		return
	}
}

func GenerateNewPDF(c *gin.Context) {
	var pdf = gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4}) // Configuración de la página A4
	pdf.AddPage()                                        // Agregar una nueva página

	var errorLoadFont = pdf.AddTTFFont("roboto", "public/fonts/Roboto/static/Roboto-Regular.ttf")
	if errorLoadFont != nil {
		fmt.Println("Error al cargar la fuente:", errorLoadFont)
		return
	}

	var errorSetFont = pdf.SetFont("roboto", "", 14)
	if errorSetFont != nil {
		fmt.Println("Error al establecer la fuente:", errorSetFont)
		return
	}

	pdf.SetXY(50, 50) // Establecer la posición inicial del texto
	var errorCell = pdf.Cell(nil, "Hola Mundo desde Golang")
	if errorCell != nil {
		fmt.Println("Error al agregar la celda:", errorCell)
		return
	}

	// Creando directorio de salida
	var errorCreatePath = os.MkdirAll("storage/files", 0775)
	if errorCreatePath != nil {
		fmt.Println("Error al crear el directorio:", errorCreatePath)
		return
	}

	var errorWriteFile = pdf.WritePdf("storage/files/ejemplo.pdf")
	if errorWriteFile != nil {
		fmt.Println("Error al escribir el archivo PDF:", errorWriteFile)
		return
	}

	c.Redirect(http.StatusSeeOther, "/pdf")
	c.Abort()
}

func IndexExcel(c *gin.Context) {
	authData := helpers.AuthSessionService(c.Writer, c.Request, "Excel", nil)
	tmpl := template.Must(template.ParseFiles("resources/utils/excel.html", config.MainLayoutFilePath))
	err := tmpl.Execute(c.Writer, authData)
	if err != nil {
		fmt.Println("Error al ejecutar la plantilla:", err)
		return
	}
}

func GenerateNewExcel(c *gin.Context) {
	var newExcel = excelize.NewFile()

	defer func() {
		if errorNewExcelClose := newExcel.Close(); errorNewExcelClose != nil {
			fmt.Println("Error al cerrar el archivo Excel:", errorNewExcelClose)
		}
	}()

	var shheetName = "Sheet1"

	var sheet, errorSheet = newExcel.NewSheet(shheetName)
	if errorSheet != nil {
		fmt.Println("Error al crear la hoja:", errorSheet)
		return
	}

	newExcel.SetCellValue(shheetName, "A1", "Id")
	newExcel.SetCellValue(shheetName, "B1", "Nombre")
	newExcel.SetCellValue(shheetName, "C1", "Apellido")
	newExcel.SetActiveSheet(sheet)

	var userData = [][]string{
		{"1", "Juan", "Pérez"},
		{"2", "María", "Gómez"},
		{"3", "Pedro", "López"},
	}

	for index, user := range userData {
		newExcel.SetCellValue(shheetName, fmt.Sprintf("A%d", index+2), user[0])
		newExcel.SetCellValue(shheetName, fmt.Sprintf("B%d", index+2), user[1])
		newExcel.SetCellValue(shheetName, fmt.Sprintf("C%d", index+2), user[2])
	}

	var errorCreatePath = os.MkdirAll("storage/files", 0775)
	if errorCreatePath != nil {
		fmt.Println("Error al crear el directorio:", errorCreatePath)
		return
	}

	var errorWriteFile = newExcel.SaveAs("storage/files/ejemplo.xlsx")
	if errorWriteFile != nil {
		fmt.Println("Error al escribir el archivo Excel:", errorWriteFile)
		return
	}

	c.Redirect(http.StatusSeeOther, "/excel")
	c.Abort()
}

func IndexQR(c *gin.Context) {
	authData := helpers.AuthSessionService(c.Writer, c.Request, "QR", nil)
	tmpl := template.Must(template.ParseFiles("resources/utils/qr.html", config.MainLayoutFilePath))
	err := tmpl.Execute(c.Writer, authData)
	if err != nil {
		fmt.Println("Error al ejecutar la plantilla:", err)
		return
	}
}

func GenerateNewQR(c *gin.Context) {
	var qrText string = "http://localhost:8080"

	var qrCode, errorGenerateQR = qrcode.Encode(qrText, qrcode.Medium, 256)
	if errorGenerateQR != nil {
		fmt.Println("Error al generar el código QR:", errorGenerateQR)
		return
	}

	var errorCreatePath = os.MkdirAll("storage/files", 0775)
	if errorCreatePath != nil {
		fmt.Println("Error al crear el directorio:", errorCreatePath)
		return
	}

	var errorWriteFile = os.WriteFile("storage/files/ejemplo.png", qrCode, 0775)
	if errorWriteFile != nil {
		fmt.Println("Error al escribir el archivo QR:", errorWriteFile)
		return
	}

	c.Redirect(http.StatusSeeOther, "/qr")
	c.Abort()
}

func IndexSendEmail(c *gin.Context) {
	authData := helpers.AuthSessionService(c.Writer, c.Request, "Email", nil)
	tmpl := template.Must(template.ParseFiles("resources/utils/email.html", config.MainLayoutFilePath))
	err := tmpl.Execute(c.Writer, authData)
	if err != nil {
		fmt.Println("Error al ejecutar la plantilla:", err)
		return
	}
}

func GenerateNewEmail(c *gin.Context) {
	// Aquí iría la lógica para enviar el correo electrónico
	// Puedes usar un paquete como "net/smtp" o "github.com/go-gomail/gomail"
	// para enviar correos electrónicos en Go.

	var message = gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", "user01@wisus.dev")
	message.SetHeader("To", "wisusdev@gmail.com")
	message.SetHeader("Subject", "Hola desde Golang")

	// Set email body
	message.SetBody("text/html", "<h1>Hola Mundo</h1><p>Este es un correo electrónico enviado desde Golang.</p>")

	var dialer = gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, "35bb66fbafb3f8", "8e8b7e8f026938")
	var errorDialer = dialer.DialAndSend(message)
	if errorDialer != nil {
		fmt.Println("Error al enviar el correo electrónico:", errorDialer)
		return
	} else {
		fmt.Println("Correo electrónico enviado con éxito.")
	}

	// Por ahora, solo redirigimos a la página de inicio.
	c.Redirect(http.StatusSeeOther, "/email")
	c.Abort()
}
