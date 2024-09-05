package database

import (
	"api-gin/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConectaDB() {
	conexao := "host=localhost user=postgres password=Rascus/10 dbname=alunos port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(conexao))
	if err != nil {
		log.Panic("Erro ao conectar com DB")
	}
	DB.AutoMigrate(&models.Aluno{})
}
