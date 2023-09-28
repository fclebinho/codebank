package services

import "github.com/fclebinho/codebank/src/domain/entities"

type TransactionService interface {
	SaveTransaction(transaction entities.Transaction, creditCard entities.CreditCard) error
	GetCreditCard(creditCard entities.CreditCard) (entities.CreditCard, error)
	CreateCreditCard(creditCard entities.CreditCard) error
}
