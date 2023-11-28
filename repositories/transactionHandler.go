package repositories

import (
	"context"
	"database/sql"
)

func BeginTransaction(jobsRepository *JobsRepository, weatherRepository *WeatherRepository) error {
	ctx := context.Background()
	transaction, err := weatherRepository.dbHandler.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	jobsRepository.transaction = transaction
	weatherRepository.transaction = transaction

	return nil
}

func RollbackTransaction(jobsRepository *JobsRepository, weatherRepository *WeatherRepository) error {
	transaction := jobsRepository.transaction

	jobsRepository.transaction = nil
	weatherRepository.transaction = nil

	return transaction.Rollback()
}

func CommitTransaction(jobsRepository *JobsRepository, weatherRepository *WeatherRepository) error {
	transaction := jobsRepository.transaction

	jobsRepository.transaction = nil
	weatherRepository.transaction = nil

	return transaction.Commit()
}
