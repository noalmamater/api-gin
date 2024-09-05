package models

import (
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

type Aluno struct {
	gorm.Model
	Nome string `json:"nome" validate:"nonzero,regexp=^[a-zA-Z]*$"`
	CPF  string `json:"cpf" validate:"nonzero,len=11,regexp=^[0-9]*$"`
	RG   string `json:"rg" validate:"nonzero,min=4,max=9,regexp=^[0-9]*[xX]?$"`
}

func ValidaDados(aluno *Aluno) error {
	if err := validator.Validate(aluno); err != nil {
		return err
	}
	return nil
}
