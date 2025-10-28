// package route

// import (
// 	"latihan2/app/service"
// 	"latihan2/app/service/mongo"
// 	"latihan2/middleware"

// 	"github.com/gofiber/fiber/v2"
// )

// func SetupRoutes(app *fiber.App) {
// 	api := app.Group("/api")

// 	api.Post("/login", service.Login)
// 	protected := api.Group("", middleware.AuthRequired())
// 	protected.Get("/profile", service.GetProfile)

// 	// dengan Pagination, Sorting, & Search
// 	api.Get("/users", service.GetUsersService)
// 	api.Get("/users/:id", service.GetUserByIDService)
// 	// api.Delete("/users/:id", service.SoftDeleteUserService)
// 	// api.Get("/alumnus", service.GetAlumniService)
// 	// api.Get("/semua-pekerjaan", service.GetPekerjaanService)

// 	alumni := protected.Group("/alumni")
// 	alumni.Get("/", service.GetAlumniService)
// 	// alumni.Get("/", service.GetAllAlumni)
// 	alumni.Get("/:id", service.GetAlumniByID)
// 	alumni.Get("/tahun/:tahun", service.GetAlumniByTahunLulus)
// 	alumni.Post("/", middleware.AdminOnly(), service.CreateAlumni)
// 	alumni.Put("/:id", middleware.AdminOnly(), service.UpdateAlumni)
// 	alumni.Delete("/:id", middleware.AdminOnly(), service.DeleteAlumni)
// 	// alumni.Delete("/:id", middleware.AdminOnly(), service.SoftDeleteAlumniService)

// 	pekerjaan := protected.Group("/pekerjaan")
// 	pekerjaan.Get("/", service.GetPekerjaanService)
// 	// pekerjaan.Get("/", service.GetAllPekerjaan)
// 	pekerjaan.Get("/:id", middleware.JWTMiddleware(), service.GetPekerjaanByID)
// 	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), service.GetPekerjaanByAlumniID)
// 	pekerjaan.Post("/", middleware.AdminOnly(), service.CreatePekerjaan)
// 	pekerjaan.Put("/:id", middleware.AdminOnly(), service.UpdatePekerjaan)
// 	// pekerjaan.Delete("/:id", middleware.AdminOnly(), service.DeletePekerjaan)
// 	pekerjaan.Delete("/soft-delete/:id", service.SoftDeletePekerjaan)
// 	pekerjaan.Post("/restore/:id", service.RestorePekerjaanService)
// 	pekerjaan.Get("/trash/:id", service.GetTrashPekerjaanByIDService)
// 	pekerjaan.Delete("/hard-delete/:id", service.HardDeletePekerjaanService)

// }

// func SetupRoutesMongo(app *fiber.App) {
// 	api := app.Group("/api")

// 	// -----------------------------
// 	// LOGIN — tanpa middleware
// 	// -----------------------------
// 	api.Post("/mongo/login", mongo.LoginMongo)

// 	// -----------------------------
// 	// Semua route SELAIN login perlu token
// 	// -----------------------------
// 	protectedm := api.Group("/mongo", middleware.AuthRequired())

// 	usersm := protectedm.Group("/users-m")
// 	usersm.Get("/mongo/", mongo.GetAllUsers)
// 	usersm.Get("/mongo/:id/", mongo.GetUsersByID)

// 	alumnim := protectedm.Group("/alumni-m")
// 	alumnim.Get("/mongo/", mongo.GetAllAlumni)
// 	alumnim.Get("/mongo/:id/", middleware.JWTMiddleware(), mongo.GetAlumniByID)
// 	alumnim.Post("/mongo/", middleware.AdminOnly(), mongo.CreateAlumni)
// 	alumnim.Put("/mongo/:id", middleware.AdminOnly(), mongo.UpdateAlumni)
// 	alumnim.Delete("/mongo/soft-delete/:id", middleware.AdminOnly(), mongo.SoftDeleteAlumni)

// 	// Grup '/files' baru, dilindungi oleh 'AuthRequired'
// 	files := protectedm.Group("/files")
// 	files.Post("/upload", mongo.UploadFile) // Handler dari fileService.go
// 	files.Get("/", mongo.GetAllFiles)
// 	files.Get("/:id", mongo.GetFileByID)
// 	files.Delete("/:id", mongo.DeleteFile)

// 	pekerjaanm := protectedm.Group("/pekerjaan-m")
// 	pekerjaanm.Get("/mongo/", mongo.GetAllPekerjaan)
// 	pekerjaanm.Get("/mongo/:id/", middleware.JWTMiddleware(), mongo.GetPekerjaanByID)
// 	pekerjaanm.Get("/mongo/alumni/:alumni_id/", middleware.AdminOnly(), mongo.GetPekerjaanByAlumniID)
// 	pekerjaanm.Post("/mongo/", middleware.AdminOnly(), mongo.CreatePekerjaan)
// 	pekerjaanm.Put("/mongo/:id", middleware.AdminOnly(), mongo.UpdatePekerjaan)
// 	pekerjaanm.Delete("/mongo/soft-delete/:id", middleware.AdminOnly(), mongo.SoftDeletePekerjaan)
// 	pekerjaanm.Post("/mongo/restore/:id", middleware.AdminOnly(), mongo.RestorePekerjaan)
// 	pekerjaanm.Delete("/mongo/hard-delete/:id", middleware.AdminOnly(), mongo.HardDeletePekerjaan)

// }

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
	files.Get("/open/:id", mongo.GetFileContentByID)
	files.Delete("/:id", middleware.FileOwnerOrAdmin(), mongo.DeleteFile)

	alumnim := protectedm.Group("/alumni")
	alumnim.Get("/", mongo.GetAllAlumni)
	alumnim.Get("/:id/", middleware.JWTMiddleware(), mongo.GetAlumniByID)
	alumnim.Post("/", middleware.AdminOnly(), mongo.CreateAlumni)
	alumnim.Put("/:id", middleware.AdminOnly(), mongo.UpdateAlumni)
	alumnim.Delete("/soft-delete/:id", middleware.AdminOnly(), mongo.SoftDeleteAlumni)

	pekerjaanm := protectedm.Group("/pekerjaan")
	pekerjaanm.Get("/", mongo.GetAllPekerjaan)
	pekerjaanm.Get("/:id/", middleware.JWTMiddleware(), mongo.GetPekerjaanByID)
	pekerjaanm.Get("/alumni/:alumni_id/", middleware.AdminOnly(), mongo.GetPekerjaanByAlumniID)
	pekerjaanm.Post("/", middleware.AdminOnly(), mongo.CreatePekerjaan)
	pekerjaanm.Put("/:id", middleware.AdminOnly(), mongo.UpdatePekerjaan)
	pekerjaanm.Delete("/soft-delete/:id", middleware.AdminOnly(), mongo.SoftDeletePekerjaan)
	pekerjaanm.Post("/restore/:id", middleware.AdminOnly(), mongo.RestorePekerjaan)
	pekerjaanm.Delete("/hard-delete/:id", middleware.AdminOnly(), mongo.HardDeletePekerjaan)
}

