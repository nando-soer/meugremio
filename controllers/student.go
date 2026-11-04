package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"tcc/dto"
	"tcc/services"

	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

type StudentController struct {
	service services.StudentService
}

type GenericResponse struct {
	Title   string `json:"title"`
	Details string `json:"details,omitempty"`
	Html    string `json:"html"`
}

func CreateStudentController(service services.StudentService) StudentController {
	return StudentController{service: service}
}

func craftResponse(title string, details string, additional string) GenericResponse {
	return GenericResponse{
		Title:   title,
		Details: details,
		Html: fmt.Sprintf(`
			<div class="modal-header">
        	  <span>%s</span>
        	  <button type="button" class="modal-close" id="modalClose" aria-label="Fechar modal">
        	    <img src="/static/icones/v branco/19.png" alt="">
        	  </button>
      		</div>

      		<div class="modal-body">
      		  <p class="confirmation-text">
      		    %s
      		  </p>
      		  %s
      		</div>
		`, title, details, additional),
	}
}

func (student *StudentController) readString(ctx *gin.Context) (string, error) {
	requestBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return "", err
	}
	return string(requestBytes), err
}

func (student *StudentController) AlreadyLogged(ctx *gin.Context) bool {
	id, err := ctx.Cookie("id")
	if err != nil {
		return false
	}

	idExists, err := student.service.ValidateSessionId(ctx, id)
	if err != nil {
		return false
	}

	return idExists
}

func (student StudentController) InRegisterProcess(ctx *gin.Context) bool {
	id, err := ctx.Cookie("pending_id")
	if err != nil {
		return false
	}

	pendingRegisterExists, err := student.service.DoesPendingRegisterExists(ctx, id)
	if err != nil {
		return false
	}

	return pendingRegisterExists
}

func (student *StudentController) RetrieveProfile(ctx *gin.Context) {
	if !student.AlreadyLogged(ctx) {
		// TODO: better error handling
		ctx.String(http.StatusForbidden, "User is not logged")
		return
	}

	sessionId, err := ctx.Cookie("id")
	if err != nil {
		// TODO: better error handling
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	data, err := student.service.RecoverSession(ctx, sessionId)
	if err != nil {
		// TODO: better error handling
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (student *StudentController) PasswordStrengthVerify(ctx *gin.Context) {
	plainPassword, err := student.readString(ctx)
	if err != nil {
		ctx.String(http.StatusForbidden, "Failed to read password")
		return
	}

	entropy := passwordvalidator.GetEntropy(plainPassword)
	strength := ""

	if entropy > 60 {
		strength = "SENHA FORTE"
	} else if entropy > 34 && entropy < 60 {
		strength = "SENHA MÉDIA"
	} else if entropy < 34 {
		strength = "SENHA FRACA"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"strength": strength,
		"entropy":  entropy,
	})
}

func (student *StudentController) InitPendingRegister(ctx *gin.Context) {
	if student.AlreadyLogged(ctx) {
		ctx.JSON(http.StatusForbidden, GenericResponse{
			Title:   "USUÁRIO LOGADO",
			Details: "O usuário já está devidamente registrado e com acesso à plataforma.",
			Html:    "",
		})
		return
	}

	if !student.InRegisterProcess(ctx) {
		var data dto.StudentData

		err := ctx.ShouldBindJSON(&data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, craftResponse(
				"REQUISIÇÃO MAL FORMADA",
				"Seus dados foram mal recebidos pelo nosso servidor.",
				"",
			))
			return
		}

		data.Email = strings.ToLower(data.Email)

		exists, err := student.service.IsEmailAlreadyUsed(ctx, data.Email)
		if exists || err != nil {
			ctx.JSON(http.StatusForbidden, craftResponse(
				"EMAIL JÁ UTILIZADO",
				"O email inserido já está sendo utilizado.",
				"",
			))
			return
		}

		correct := student.service.IsEmailCorrect(ctx, data.Email)
		if !correct {
			ctx.JSON(http.StatusForbidden, craftResponse(
				"EMAIL PROIBIDO",
				"O email inserido não pode ser utilizado. Use sua conta institucional.",
				"",
			))
			return
		}

		pendingId, err := student.service.InstantiatePendingRegister(ctx, data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, craftResponse(
				"ERRO INTERNO",
				"Não conseguimos receber seu registro por conta de um erro interno. Tente novamente mais tarde.",
				"",
			))
			return
		}

		ctx.SetCookie("pending_id", pendingId, 60*10, "/", "", true, false) // change to true
	}

	ctx.JSON(http.StatusOK, craftResponse(
		"CONFIRMAÇÃO",
		"Digite o código de confirmação enviado no seu E-mail Institucional.",
		`
			<input type="text" class="text-input code-input" id="confirmCode" maxlength="8">
			<button class="btn-primary" id="confirmRegister">Finalizar</button>
		`,
	))
}

func (student *StudentController) ConfirmPendingRegister(ctx *gin.Context) {
	if student.AlreadyLogged(ctx) {
		ctx.JSON(http.StatusForbidden, GenericResponse{
			Title:   "USUÁRIO LOGADO",
			Details: "O usuário já está devidamente registrado e com acesso à plataforma.",
			Html:    "",
		})
		return
	}

	pendingId, err := ctx.Cookie("pending_id")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Não conseguimos reconhecer o seu registro.",
			"",
		))
		return
	}

	plainCode, err := student.readString(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Não conseguimos ler o código de confirmação.",
			"",
		))
		return
	}

	isCodeValid, err := student.service.ValidatePendingRegister(ctx, pendingId, plainCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Não conseguimos validar seu pedido de registro. Tente novamente mais tarde.",
			"",
		))
		return
	}

	if !isCodeValid {
		ctx.JSON(http.StatusUnauthorized, craftResponse(
			"CÓDIGO INCORRETO",
			"O código que você inseriu está incorreto. Tente novamente.",
			`
				<input type="text" class="text-input code-input" id="confirmCode" maxlength="8">
				<button class="btn-primary" id="confirmRegister">Finalizar</button>
			`,
		))
		return
	}

	err = student.service.FulfillPendingRegister(ctx, pendingId)
	if err != nil {
		// TODO: better error handling
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Não conseguimos concluir seu registro.",
			"",
		))
		return
	}

	ctx.SetCookie("pending_id", "", -1, "/", "", true, false)
	ctx.JSON(http.StatusOK, craftResponse(
		"CÓDIGO CORRETO",
		"Você foi cadastrado em nossa plataforma com sucesso!<br>Para acessar agora, basta clicar no botão abaixo e prosseguir para a página de entrada.",
		`
			<button class="btn-primary" type="button" onclick="window.location.href='/usuario/entrar'">
                Entrar
            </button>
		`,
	))
}

func (student *StudentController) Login(ctx *gin.Context) {
	if student.AlreadyLogged(ctx) {
		// TODO: better error handling
		ctx.String(http.StatusForbidden, "User already logged")
		return
	}

	var submission dto.StudentLogin

	err := ctx.ShouldBindJSON(&submission)
	if err != nil {
		// TODO: better error handling
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	isLoginValid, err := student.service.Validate(ctx, submission)
	if err != nil {
		// TODO: better error handling
		ctx.String(http.StatusInternalServerError, err.Error())
	}

	if !isLoginValid {
		// TODO: better error handling
		ctx.String(http.StatusUnauthorized, "Email or password is incorrect")
		return
	}

	sessionId, err := student.service.InstantiateSession(ctx, submission.Email)
	if err != nil {
		// TODO: better error handling
		ctx.String(http.StatusInternalServerError, err.Error())
	}

	ctx.SetCookie("id", sessionId, 60*60*24, "/", "", true, false) // change to true
	ctx.String(http.StatusOK, sessionId)
}

// ------------------------------------------- Test ---------------------------------------------

func (student *StudentController) Retrieve(ctx *gin.Context) {
	email := ctx.Query("email")

	if email == "" {
		students, err := student.service.QueryAll(ctx)

		if err != nil {
			// TODO: better error handling
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, students)
		return
	}

	data, err := student.service.QueryNonSensibleDataByEmail(ctx, email)

	if err != nil {
		// TODO: better error handling
		ctx.String(http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, data)
}
