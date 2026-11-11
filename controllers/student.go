package controllers

import (
	"fmt"
	"io"
	"math/rand/v2"
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

func (student *StudentController) isEmailCorrect(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := strings.ToLower(parts[1])
	return domain == "aluno.cps.sp.gov.br"
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

	sessionId, err := student.readString(ctx)
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

func (student *StudentController) UpdateProfile(ctx *gin.Context) {
	if !student.AlreadyLogged(ctx) {
		ctx.JSON(http.StatusForbidden, craftResponse(
			"NÃO AUTORIZADO",
			"Você precisa estar logado para editar o perfil.",
			"",
		))
		return
	}

	sessionId, _ := ctx.Cookie("id")

	var data dto.StudentEditData

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, craftResponse(
			"REQUISIÇÃO MAL FORMADA",
			"Seus dados foram mal recebidos pelo nosso servidor.",
			"",
		))
		return
	}

	if err := student.service.UpdateProfile(ctx, sessionId, data); err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Não conseguimos salvar suas alterações.",
			"",
		))
		return
	}

	ctx.JSON(http.StatusOK, craftResponse(
		"SUCESSO",
		"Perfil atualizado com sucesso! Clique no botão abaixo para ver seu novo perfil.",
		`
			<button class="btn-primary" type="button" onclick="window.location.href='/usuario/perfil'">
                Ver
            </button>
		`,
	))
}

func (student *StudentController) PasswordStrengthVerify(ctx *gin.Context) {
	plainPassword, err := student.readString(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Algo falhou. Recarregue a página",
			"",
		))
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

func (student *StudentController) declareEmptyData(data *dto.StudentData) {
	data.Description = "Escreva algo aqui!"
	data.Banner = rand.IntN(9) + 1
	data.Avatar = dto.AvatarData{
		Body:      rand.IntN(9) + 1,
		Shirt:     rand.IntN(10),
		Pants:     rand.IntN(10),
		Accessory: rand.IntN(10),
	}
}

func (student *StudentController) InitPendingRegister(ctx *gin.Context) {
	if student.AlreadyLogged(ctx) {
		ctx.JSON(http.StatusForbidden, craftResponse(
			"USUÁRIO LOGADO",
			"O usuário já está devidamente registrado e com acesso à plataforma.",
			"",
		))
		return
	}

	if !student.InRegisterProcess(ctx) {
		var data dto.StudentData

		student.declareEmptyData(&data)

		err := ctx.ShouldBindJSON(&data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, craftResponse(
				"REQUISIÇÃO MAL FORMADA",
				"Seus dados foram mal recebidos pelo nosso servidor.",
				"",
			))
			return
		}

		if !student.service.ClassExists(ctx, data.Class) {
			ctx.JSON(http.StatusBadRequest, craftResponse(
				"CLASSE INEXISTENTE",
				"A classe fornecida não existe.",
				"",
			))
			return
		}

		data.Email = strings.ToLower(data.Email)

		exists, err := student.service.EmailExists(ctx, data.Email)
		if exists || err != nil {
			ctx.JSON(http.StatusForbidden, craftResponse(
				"EMAIL JÁ UTILIZADO",
				"O email inserido já está sendo utilizado.",
				"",
			))
			return
		}

		correct := student.isEmailCorrect(data.Email)
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
		ctx.JSON(http.StatusForbidden, craftResponse(
			"USUÁRIO LOGADO",
			"O usuário já está devidamente registrado e com acesso à plataforma.",
			"",
		))
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
		ctx.JSON(http.StatusForbidden, craftResponse(
			"USUÁRIO LOGADO",
			"O usuário já está devidamente registrado e com acesso à plataforma.",
			"",
		))
		return
	}

	var submission dto.StudentLogin

	err := ctx.ShouldBindJSON(&submission)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, craftResponse(
			"REQUISIÇÃO MAL FORMADA",
			"Seus dados foram mal recebidos pelo nosso servidor.",
			"",
		))
		return
	}

	submission.Email = strings.ToLower(submission.Email)

	emailExists, err := student.service.EmailExists(ctx, submission.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Não conseguimos processar seu pedido de entrada.",
			"",
		))
		return
	}

	if !emailExists {
		ctx.JSON(http.StatusForbidden, craftResponse(
			"NÃO ENCONTRADO",
			"O email inserido não foi encontrado no nosso sistema. Clique no botão abaixo para ser encaminhado à página de cadastro.",
			`
				<button class="btn-primary" type="button" onclick="window.location.href='/usuario/cadastro'">
            	    Cadastro
            	</button>
			`,
		))
		return
	}

	isLoginValid, err := student.service.Validate(ctx, submission)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Não conseguimos processar seu pedido de entrada.",
			"",
		))
		return
	}

	if !isLoginValid {
		ctx.JSON(http.StatusForbidden, craftResponse(
			"INCORRETO",
			"Seu email ou sua senha estão incorretos.",
			"",
		))
		return
	}

	sessionId, err := student.service.InstantiateSession(ctx, submission.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Não conseguimos fornecer acesso à sua conta. Tente novamente mais tarde..",
			"",
		))
		return
	}

	ctx.SetCookie("id", sessionId, 60*60*24*30, "/", "", true, false) // change to true
	ctx.JSON(http.StatusOK, craftResponse(
		"SUCESSO",
		"Seja bem vindo(a) à plataforma MeuGremio! Clique no botão abaixo para ir até seu novíssimo perfil e personalizá-lo!",
		`
			<button class="btn-primary" type="button" onclick="window.location.href='/usuario/perfil'">
                Perfil
            </button>
		`,
	))
}

func (student *StudentController) InitPasswordRecovery(ctx *gin.Context) {
	email, err := student.readString(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, craftResponse(
			"ERRO INTERNO",
			"Encontramos um problema durante o seu pedido de alteração de senha. Recarregue a página e tente novamente.",
			"",
		))
		return
	}

	email = strings.ToLower(email)

	if len(email) == 0 {
		ctx.JSON(http.StatusForbidden, craftResponse(
			"EMAIL VAZIO",
			"Digite um email para realizar a requisição de alteração de senha.",
			"",
		))
		return
	}

	exists, err := student.service.EmailExists(ctx, email)
	if !exists || err != nil {
		ctx.JSON(http.StatusForbidden, craftResponse(
			"EMAIL INEXISTENTE",
			"O email inserido não está cadastrado na plataforma.",
			"",
		))
		return
	}

	err = student.service.InstantiatePasswordRecovery(ctx, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Encontramos um problema durante o seu pedido de alteração de senha. Recarregue a página e tente novamente.",
			"",
		))
		return
	}

	ctx.JSON(http.StatusOK, craftResponse(
		"ENVIADO",
		"Uma mensagem com as etapas de recuperação de senha foi enviado para seu email.",
		"",
	))
}

func (student *StudentController) AlterPassword(ctx *gin.Context) {
	var submission dto.StudentPasswordUpdate

	err := ctx.ShouldBindJSON(&submission)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, craftResponse(
			"REQUISIÇÃO MAL FORMADA",
			"Seus dados foram mal recebidos pelo nosso servidor.",
			"",
		))
		return
	}

	id, err := student.service.RecoverStudentRecoveryId(ctx, submission.Token)
	if err != nil {
		ctx.JSON(http.StatusForbidden, craftResponse(
			"TOKEN INEXISTENTE",
			"Falhamos em encontrar o token associado a sua requisição de mudança de senha. Se o erro persistir, inicie uma nova requisição.",
			"",
		))
		return
	}

	if len(submission.Password) == 0 {
		return
	}

	err = student.service.ChangePasswordById(ctx, id, submission.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Falhamos alterar sua senha. Tente novamente mais tarde.",
			"",
		))
		return
	}

	err = student.service.DeletePasswordRecovery(ctx, submission.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Falhamos em finalizar com o processo de alteração de senha. Tente novamente mais tarde.",
			"",
		))
		return
	}

	ctx.JSON(http.StatusOK, craftResponse(
		"SUCESSO",
		"Sua senha foi alterada com sucesso! Clique no botão abaixo para entrar em sua conta com a nova senha.",
		`
			<button class="btn-primary" type="button" onclick="window.location.href='/usuario/entrar'">
                Entrar
            </button>
		`,
	))
}

func (student *StudentController) RetrieveClasses(ctx *gin.Context) {
	classes, err := student.service.GetAllClasses(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Não foi possível retornar as classes registradas no sistema. Tente novamente mais tarde.",
			"",
		))
	}
	ctx.JSON(http.StatusOK, classes)
}

func (student *StudentController) RetrieveClass(ctx *gin.Context) {
	classCode, err := student.readString(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"Tivemos um erro no reconhecimento do seu perfil. Recarregue a página.",
			"",
		))
		return
	}

	class, err := student.service.GetClass(ctx, classCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, craftResponse(
			"ERRO INTERNO",
			"A sala definida no seu perfil não foi encontrada. Recarregue a página.",
			"",
		))
	}

	ctx.JSON(http.StatusOK, class)
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
