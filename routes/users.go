package routes

import (
	"backend/models"
	"backend/utils"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, db *sql.DB) {
	router.POST("/users", func(c *gin.Context) {
		createUser(c, db)
	})
	router.GET("/users", func(c *gin.Context) {
		listUsers(c, db)
	})
	router.GET("/users/:id", func(c *gin.Context) {
		getOneUser(c, db)
	})
	router.PUT("/users/:id", func(c *gin.Context) {
		updateUser(c, db)
	})
	router.DELETE("/users/:id", func(c *gin.Context) {
		deleteUser(c, db)
	})
}

// @Summary Cria um novo usuário
// @Description Cria um novo usuário com as informações fornecidas
// @Tags users
// @Accept json
// @Produce json
// @Param usuario body models.User true "Informações do Usuário"
// @Success 201 {object} models.User
// @Failure 400 {object} utils.ErrorResponse "E-mail inválido ou senha inválida"
// @Failure 500 {object} utils.ErrorResponse "Erro interno do servidor"
// @Router /users [post]
func createUser(c *gin.Context, db *sql.DB) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
		return
	}

	if !utils.IsValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Endereço de email inválido"})
		return
	}

	if !utils.IsValidPassword(user.PasswordHash) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Senha inválida"})
		return
	}

	exists, err := emailExists(db, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Erro ao verificar o e-mail"})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "E-mail já cadastrado"})
		return
	}

	_, err = db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", user.Name, user.Email, user.PasswordHash)
	if err != nil {
		log.Printf("Erro ao criar usuário: %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Erro ao criar um usuário"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func emailExists(db *sql.DB, email string) (bool, error) {
	var existingEmail string
	err := db.QueryRow("SELECT email FROM users WHERE email = $1", email).Scan(&existingEmail)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return existingEmail != "", nil
}

// @Summary Lista todos os usuarios
// @Description Lista todos os usuarios disponíveis
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} utils.ErrorResponse
// @Router /users [get]
func listUsers(c *gin.Context, db *sql.DB) {
	var users []models.User

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Printf("Erro ao executar a consulta: %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Erro ao listar usuarios"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
			log.Printf("Erro ao escanear usuário: %v", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Erro ao listar os usuarios"})
			return
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Erro ao listar os resultados: %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Erro ao listar os usuarios"})
		return
	}

	log.Printf("Número de usuários encontrados: %d", len(users))
	c.JSON(http.StatusOK, users)
}

// @Summary Obtém um usuario pelo ID
// @Description Obtém detalhes de um usuario com base no ID fornecido
// @Tags users
// @Produce json
// @Param id path string true "ID do Usuario"
// @Success 200 {object} models.User
// @Failure 404 {object} utils.ErrorResponse
// @Router /users/{id} [get]
func getOneUser(c *gin.Context, db *sql.DB) {
	var user models.User

	id := c.Param("id")
	err := db.QueryRow("SELECT * FROM users  WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		log.Printf("Erro ao executar a consulta: %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Erro ao buscar usuário"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Atualiza um usuário
// @Description Atualiza as informações de um usuário existente
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID do Usuário"
// @Param usuario body models.User true "Informações atualizadas do Usuário"
// @Success 200 {object} models.User
// @Failure 400 {object} utils.ErrorResponse "Erro na requisição"
// @Failure 404 {object} utils.ErrorResponse "Usuário não encontrado"
// @Failure 500 {object} utils.ErrorResponse "Erro interno do servidor"
// @Router /users/{id} [put]
func updateUser(c *gin.Context, db *sql.DB) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "ID inválido"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
		return
	}

	if user.IsDeleted {
		softDeleteUser(c, db, id)
		return
	}

	currentTimestamp := time.Now().Unix()

	result, err := db.Exec("UPDATE users SET name = $1, email = $2, password_hash = $3, updated_at = $4 WHERE id = $5",
		user.Name, user.Email, user.PasswordHash, currentTimestamp, id)
	if err != nil {
		log.Printf("Erro ao atualizar usuário: %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Erro ao atualizar usuário"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Erro ao verificar atualização"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorResponse{Error: "Usuário não encontrado"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func softDeleteUser(c *gin.Context, db *sql.DB, userID int) {
	var user models.User
	user.UpdateTimestamps(true)

	_, err := db.Exec("UPDATE users SET deleted_at = $1 WHERE id = $2",
		user.DeletedAt, userID)
	if err != nil {
		log.Printf("Erro ao atualizar usuário como deletado: %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Erro ao deletar usuário"})
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{Message: "Usuário deletado com sucesso"})
}

// @Summary Deleta um usuário
// @Description Remove um usuário da base de dados com base no ID fornecido
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID do Usuário a ser deletado"
// @Success 204 {object} utils.SuccessResponse "Usuário deletado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro na requisição"
// @Failure 404 {object} utils.ErrorResponse "Usuário não encontrado"
// @Failure 500 {object} utils.ErrorResponse "Erro interno do servidor"
// @Router /users/{id} [delete]
func deleteUser(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Printf("Erro ao deletar usuário: %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Erro ao excluir usuario"})
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{Message: "Usuario deletado com sucesso"})
}
