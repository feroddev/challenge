package main

import (
	"fmt"
	"github/feroddev/challengeV3/internal/configs"
	"github/feroddev/challengeV3/internal/repositories"
	"github/feroddev/challengeV3/internal/services"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: check_password <username> <password>")
		os.Exit(1)
	}

	username := os.Args[1]
	password := os.Args[2]

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	config := configs.LoadConfig()
	db, err := repositories.InitDB(config.Database, logger)
	if err != nil {
		logger.Fatal("Erro ao conectar ao banco de dados", zap.Error(err))
	}

	userRepo := repositories.NewUserRepository(db)
	user, err := userRepo.FindByUsername(username)
	if err != nil {
		logger.Fatal("Usuario nao encontrado", zap.Error(err))
	}

	fmt.Printf("Usuario encontrado: %s\n", user.Username)
	fmt.Printf("Tamanho da senha armazenada: %d\n", len(user.Password))

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Printf("Erro ao verificar senha: %v\n", err)
		
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		fmt.Printf("Nova senha hash: %s\n", string(hashedPassword))
		
		user.Password = string(hashedPassword)
		err = userRepo.Update(user)
		if err != nil {
			fmt.Printf("Erro ao atualizar senha: %v\n", err)
		} else {
			fmt.Println("Senha atualizada com sucesso")
		}
	} else {
		fmt.Println("Senha correta!")
	}
}
