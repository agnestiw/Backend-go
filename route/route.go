package route

import (
	"latihan2/app/service"
	"latihan2/app/service/mongo"
	"latihan2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutesPostgres(app *fiber.App) {
	api := app.Group("/api/pg")

	api.Post("/login", service.Login)
	protected := api.Group("", middleware.AuthRequired())
	protected.Get("/profile", service.GetProfile)

	// dengan Pagination, Sorting, & Search
	api.Get("/users", service.GetUsersService)
	api.Get("/users/:id", service.GetUserByIDService)
	// api.Delete("/users/:id", service.SoftDeleteUserService)
	// api.Get("/alumnus", service.GetAlumniService)
	// api.Get("/semua-pekerjaan", service.GetPekerjaanService)

	alumni := protected.Group("/alumni")
	alumni.Get("/", service.GetAlumniService)
	// alumni.Get("/", service.GetAllAlumni)
	alumni.Get("/:id", service.GetAlumniByID)
	alumni.Get("/tahun/:tahun", service.GetAlumniByTahunLulus)
	alumni.Post("/", middleware.AdminOnly(), service.CreateAlumni)
	alumni.Put("/:id", middleware.AdminOnly(), service.UpdateAlumni)
	alumni.Delete("/:id", middleware.AdminOnly(), service.DeleteAlumni)
	// alumni.Delete("/:id", middleware.AdminOnly(), service.SoftDeleteAlumniService)

	pekerjaan := protected.Group("/pekerjaan")
	pekerjaan.Get("/", service.GetPekerjaanService)
	// pekerjaan.Get("/", service.GetAllPekerjaan)
	pekerjaan.Get("/:id", middleware.JWTMiddleware(), service.GetPekerjaanByID)
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), service.GetPekerjaanByAlumniID)
	pekerjaan.Post("/", middleware.AdminOnly(), service.CreatePekerjaan)
	pekerjaan.Put("/:id", middleware.AdminOnly(), service.UpdatePekerjaan)
	// pekerjaan.Delete("/:id", middleware.AdminOnly(), service.DeletePekerjaan)
	pekerjaan.Delete("/soft-delete/:id", service.SoftDeletePekerjaan)
	pekerjaan.Post("/restore/:id", service.RestorePekerjaanService)
	pekerjaan.Get("/trash/:id", service.GetTrashPekerjaanByIDService)
	pekerjaan.Delete("/hard-delete/:id", service.HardDeletePekerjaanService)

}

func SetupRoutesMongo(app *fiber.App) {
	api := app.Group("/api/mg")

	api.Post("/login", mongo.LoginMongo)
	protectedm := api.Group("", middleware.AuthRequiredMongo())

	usersm := protectedm.Group("/users")
	usersm.Get("/", mongo.GetAllUsers)
	usersm.Get("/:id/", mongo.GetUsersByID)

	files := protectedm.Group("/files")
	files.Post("/upload", mongo.UploadFile)
	files.Get("/", mongo.GetAllFiles)
	files.Get("/:id", mongo.GetFileByID)
	files.Get("/open/:id", mongo.GetContentByID)
	files.Delete("/:id", middleware.FileOwnerOrAdmin(), mongo.DeleteFile)

	alumnim := protectedm.Group("/alumni")
	alumnim.Get("/", mongo.GetAllAlumni)
	alumnim.Get("/:id/", mongo.GetAlumniByID)
	alumnim.Post("/", mongo.CreateAlumni)
	alumnim.Put("/:id", mongo.UpdateAlumni)
	alumnim.Delete("/soft-delete/:id", mongo.SoftDeleteAlumni)

	pekerjaanm := protectedm.Group("/pekerjaan")
	pekerjaanm.Get("/", mongo.GetAllPekerjaan)
	pekerjaanm.Get("/alumni/:alumni_id/", mongo.GetPekerjaanByAlumniID)
	pekerjaanm.Get("/:id/", mongo.GetPekerjaanByID)
	pekerjaanm.Post("/", mongo.CreatePekerjaan)
	pekerjaanm.Put("/:id", mongo.UpdatePekerjaan)
	pekerjaanm.Delete("/soft-delete/:id", mongo.SoftDeletePekerjaan)
	pekerjaanm.Get("/trash/:id", mongo.GetTrashPekerjaan)
	pekerjaanm.Post("/restore/:id", mongo.RestorePekerjaan)
	pekerjaanm.Delete("/hard-delete/:id", mongo.HardDeletePekerjaan)
}

