package mysql

import (
	"context"
	"currency-conversion/internal/dto"
	"currency-conversion/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type exchangeRatesRepo struct {
	DB *gorm.DB
}

func NewExchangeRatesRepo(db *gorm.DB) *exchangeRatesRepo {
	return &exchangeRatesRepo{DB: db}
}

func (r *exchangeRatesRepo) GetExchangeRates(ctx context.Context) (*dto.ExchangeRates, error) {
	logrus.Info("GetExchangeRates called")
	var exchangeRates []entity.ExchangeRates
	if err := r.DB.Find(&exchangeRates).Error; err != nil {
		logrus.Errorf("Failed to fetch exchange rates from database: %v", err)
		return nil, err
	}

	ratesMap := make(map[string]entity.ExchangeRates)
	for _, rate := range exchangeRates {
		ratesMap[rate.Code] = rate
	}

	data := &dto.ExchangeRates{Data: ratesMap}
	logrus.WithField("ExchangeRatesData", data).Debug("Successfully fetched exchange rates data")
	return data, nil
}

func (r *exchangeRatesRepo) AddRates(ctx context.Context, data *dto.ExchangeRates) error {
	logrus.Info("AddRates called")
	tx := r.DB.Begin()
	if tx.Error != nil {
		logrus.Errorf("Failed to begin transaction: %v", tx.Error)
		return tx.Error
	}

	for _, rate := range data.Data {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&rate).Error; err != nil {
			tx.Rollback()
			logrus.Errorf("Failed to add or update exchange rate: %v", err)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		logrus.Errorf("Failed to commit transaction: %v", err)
		return err
	}

	logrus.Info("Exchange rates data inserted or updated in database successfully")
	return nil
}
