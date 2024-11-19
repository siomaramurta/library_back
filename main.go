package main

import (
	"database/sql"
	"log"

	_ "backend/docs" // Importando o pacote docs
	"backend/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"     // Importa os arquivos do Swagger
	ginSwagger "github.com/swaggo/gin-swagger" // Importa o manipulador Swagger para o Gin
)

func main() {
	// Conexão com o banco de dados PostgreSQL
	db, err := sql.Open("postgres", "postgres://postgres:14041989@localhost:5432/biblioteca?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	// Adicione o endpoint para a documentação Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Adicione as rotas da sua aplicação
	routes.UserRoutes(router, db)
	routes.BookRoutes(router, db)
	routes.EmprestimoRoutes(router, db)

	// Inicie o servidor
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
