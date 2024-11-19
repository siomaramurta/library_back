package utils

import (
	"regexp"
)

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	// Verificar se a senha contém pelo menos uma letra minúscula
	var temMinuscula = regexp.MustCompile(`[a-z]`)
	// Verificar se a senha contém pelo menos uma letra maiúscula
	var temMaiuscula = regexp.MustCompile(`[A-Z]`)
	// Verificar se a senha contém pelo menos um dígito
	var temDigito = regexp.MustCompile(`\d`)
	// Verificar se a senha contém pelo menos um caractere especial
	var temCaractereEspecial = regexp.MustCompile(`[@$!%*#?&]`)
	// Verificar se a senha tem pelo menos 8 caracteres
	var temTamanhoMinimo = regexp.MustCompile(`.{8,}`)

	if !temMinuscula.MatchString(password) {
		return false
	}
	if !temMaiuscula.MatchString(password) {
		return false
	}
	if !temDigito.MatchString(password) {
		return false
	}
	if !temCaractereEspecial.MatchString(password) {
		return false
	}
	if !temTamanhoMinimo.MatchString(password) {
		return false
	}

	// Verificar se não há sequência seguida de 3 números ou 3 letras iguais
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i] == password[i+2] {
			return false
		}
	}
	return true
}
