package route

import (
	"latihan2/app/service"
	"latihan2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/login", service.Login)
	protected := api.Group("", middleware.AuthRequired())
	protected.Get("/profile", service.GetProfile)

	// dengan Pagination, Sorting, & Search
	api.Get("/users", service.GetUsersService)
	// api.Delete("/users/:id", service.SoftDeleteUserService)
	api.Get("/users/:id", service.GetUserByIDService)
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
	pekerjaan.Get("/:id",  middleware.JWTMiddleware() ,service.GetPekerjaanByID)
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), service.GetPekerjaanByAlumniID) 
	pekerjaan.Post("/", middleware.AdminOnly(), service.CreatePekerjaan)             
	pekerjaan.Put("/:id", middleware.AdminOnly(), service.UpdatePekerjaan)           
	// pekerjaan.Delete("/:id", middleware.AdminOnly(), service.DeletePekerjaan)  
	pekerjaan.Delete("/soft-delete/:id", service.SoftDeletePekerjaan)
	pekerjaan.Post("/restore/:id", service.RestorePekerjaanService)
	pekerjaan.Get("/trash/:id", service.GetTrashPekerjaanByIDService)
	pekerjaan.Delete("/hard-delete/:id", service.HardDeletePekerjaanService)
}
