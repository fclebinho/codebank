package usecases

import (
	"time"

	"github.com/fclebinho/codebank/src/domain/dtos"
	"github.com/fclebinho/codebank/src/domain/entities"
	"github.com/fclebinho/codebank/src/domain/services"
)

type ProcessTranscationUsecase struct {
	TransactionService services.TransactionService
	// Producer           kafka.Producer
}

func NewProcessTranscationUsecase(t services.TransactionService) ProcessTranscationUsecase {
	return ProcessTranscationUsecase{
		TransactionService: t,
	}
}

func (u *ProcessTranscationUsecase) Process(dto dtos.Transaction) (entities.Transaction, error) {
	creditCard := u.hydrateCreditCard(dto)
	balanceAndLimit, err := u.TransactionService.GetCreditCard(*creditCard)
	if err != nil {
		return entities.Transaction{}, err
	}

	creditCard.ID = balanceAndLimit.ID
	creditCard.Limit = balanceAndLimit.Limit
	creditCard.Balance = balanceAndLimit.Balance

	transaction := u.newTransaction(dto, *creditCard)
	transaction.ProcessAndValidate(creditCard)

	err = u.TransactionService.SaveTransaction(*transaction, *creditCard)
	if err != nil {
		return entities.Transaction{}, err
	}

	dto.ID = transaction.ID
	dto.CreatedAt = transaction.CreatedAt

	// json, err := json.Marshal(dto)
	// if err != nil {
	// 	return entities.Transaction{}, nil
	// }

	// err = u.Producer.Publish(string(json), "payments")
	// if err != nil {
	// 	return entities.Transaction{}, nil
	// }

	return *transaction, nil
}

func (u *ProcessTranscationUsecase) hydrateCreditCard(dto dtos.Transaction) *entities.CreditCard {
	creditCard := entities.CreditCard{}
	creditCard.Name = dto.Name
	creditCard.Number = dto.Number
	creditCard.ExpirationMonth = dto.ExpirationMonth
	creditCard.ExpirationYear = dto.ExpirationYear
	creditCard.CVV = dto.CVV

	return &creditCard
}

func (u *ProcessTranscationUsecase) newTransaction(dto dtos.Transaction, creditCard entities.CreditCard) *entities.Transaction {
	transaction := entities.NewTransaction()
	transaction.CreditCardID = creditCard.ID
	transaction.Amount = dto.Amount
	transaction.Store = dto.Store
	transaction.Description = dto.Description
	transaction.CreatedAt = time.Now()

	return transaction
}
