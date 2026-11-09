package routes

import (
	"net/http"
	"strings"
	"tcc/controllers"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(student controllers.StudentController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.FullPath()

		if !student.AlreadyLogged(ctx) && strings.Contains(path, "/usuario/perfil") {
			ctx.Redirect(http.StatusFound, "./cadastro")
			ctx.Abort()
			return
		}

		if student.AlreadyLogged(ctx) && !strings.Contains(path, "/usuario/perfil") {
			ctx.Redirect(http.StatusFound, "./perfil")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func Init(router *gin.Engine, controller controllers.StudentController) {
	router.Static("/static", "./static")

	router.StaticFile("/", "./static/feed.html")

	user_range := router.Group("/usuario")
	user_range.Use(AuthMiddleware(controller))
	{
		user_range.StaticFile("/perfil", "./static/perfil.html")
		user_range.StaticFile("/perfil/editar", "./static/perfil_editar.html")
		user_range.StaticFile("/cadastro", "./static/cadastro.html")
		user_range.StaticFile("/entrar", "./static/login.html")
		user_range.StaticFile("/esqueci", "./static/esqueci_senha.html")
	}

	router.POST("/api/usuario", controller.RetrieveProfile)
	router.POST("/api/usuario/cadastrar", controller.InitPendingRegister)
	router.POST("/api/usuario/confirmar", controller.ConfirmPendingRegister)
	router.POST("/api/usuario/entrar", controller.Login)
	router.POST("/api/usuario/editar", controller.UpdateProfile)

	router.GET("/api/usuario/classes", controller.RetrieveClasses)
	router.POST("/api/usuario/classes", controller.RetrieveClass)

	router.POST("/api/teste/senha", controller.PasswordStrengthVerify)

	router.GET("/students", controller.Retrieve) // teste
}
