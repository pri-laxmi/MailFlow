package main

import (
	"log"

	"github.com/pri-laxmi/MailFlow/configs"
	"github.com/pri-laxmi/MailFlow/internal/apis"
	"github.com/pri-laxmi/MailFlow/internal/database"
	"github.com/pri-laxmi/MailFlow/internal/email"
	"github.com/pri-laxmi/MailFlow/internal/models"
	"github.com/pri-laxmi/MailFlow/internal/queue"
	"github.com/pri-laxmi/MailFlow/internal/repository"
	"github.com/pri-laxmi/MailFlow/internal/utils"
	"github.com/pri-laxmi/MailFlow/internal/worker"
)

func main() {
	//load config
	cfg := configs.LoadConfig()
	jwtManager := utils.NewJWTManager(cfg.JWTSecret)
	//connect to database
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Database connected successfully")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Contact{},
		&models.Template{},
		&models.Campaign{},
		&models.Job{},
		&models.JobLog{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conrelid = 'contacts'::regclass
				  AND contype = 'p'
			) THEN
				ALTER TABLE "contacts" ADD CONSTRAINT "contacts_pkey" PRIMARY KEY ("id");
			END IF;
		END $$;
	`).Error; err != nil {
		log.Fatalf("failed to ensure contacts id is a primary key: %v", err)
	}
	if err := db.AutoMigrate(
		&models.Campaign{},
		&models.Job{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	jobQueue := queue.NewQueue(1000)
	jobRepo := repository.NewJobRepository(db)
	jobLogRepo := repository.NewJobLogRepository(db)
	emailSender := email.NewMockSender()
	workerPool := worker.NewWorkerPool(
		jobQueue,
		jobRepo,
		jobLogRepo,
		emailSender,
		4, // number of workers
	)
	workerPool.Start()
	pendingJobs, err := jobRepo.FindPending()
	if err != nil {
		log.Fatalf("failed to load pending jobs: %v", err)
	}
	for i := range pendingJobs {
		jobQueue.Push(&pendingJobs[i])
	}
	log.Printf("queued %d pending jobs for processing", len(pendingJobs))
	//setup router
	router := apis.SetRoutes(cfg, db, jwtManager, jobQueue)
	port := getPort(cfg)
	log.Printf("server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}

func getPort(cfg *configs.Config) string {
	if cfg.ServerPort == "" {
		return "8080"
	}
	return cfg.ServerPort
}
