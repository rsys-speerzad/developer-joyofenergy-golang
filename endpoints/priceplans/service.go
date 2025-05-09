package priceplans

import (
	"sort"

	"joi-energy-golang/domain"
	"joi-energy-golang/repository"
)

type Service interface {
	CompareAllPricePlans(smartMeterId string) (domain.PricePlanComparisons, error)
	RecommendPricePlans(smartMeterId string, limit uint64) (domain.PricePlanRecommendation, error)
	GetAllPricePlans() ([]domain.PricePlan, error)
}

type service struct {
	pricePlans *repository.PricePlans
	accounts   *repository.Accounts
}

func NewService(
	pricePlans *repository.PricePlans,
	accounts *repository.Accounts,
) Service {
	return &service{
		pricePlans: pricePlans,
		accounts:   accounts,
	}
}

func (s *service) CompareAllPricePlans(smartMeterId string) (domain.PricePlanComparisons, error) {
	pricePlanId := s.accounts.PricePlanIdForSmartMeterId(smartMeterId)
	consumptionsForPricePlans := s.pricePlans.ConsumptionCostOfElectricityReadingsForEachPricePlan(smartMeterId)
	if len(consumptionsForPricePlans) == 0 {
		return domain.PricePlanComparisons{}, domain.ErrNotFound
	}
	return domain.PricePlanComparisons{
		PricePlanId:          pricePlanId,
		PricePlanComparisons: consumptionsForPricePlans,
	}, nil
}

func (s *service) RecommendPricePlans(smartMeterId string, limit uint64) (domain.PricePlanRecommendation, error) {
	consumptionsForPricePlans := s.pricePlans.ConsumptionCostOfElectricityReadingsForEachPricePlan(smartMeterId)
	if len(consumptionsForPricePlans) == 0 {
		return domain.PricePlanRecommendation{}, domain.ErrNotFound
	}
	var recommendations []domain.SingleRecommendation
	for k, v := range consumptionsForPricePlans {
		recommendations = append(recommendations, domain.SingleRecommendation{
			Key:   k,
			Value: v,
		})
	}
	sort.Slice(recommendations, func(i, j int) bool { return recommendations[i].Value < recommendations[j].Value })

	if limit > 0 && limit < uint64(len(recommendations)) {
		recommendations = recommendations[:limit]
	}

	return domain.PricePlanRecommendation{Recommendations: recommendations}, nil
}

func (s *service) GetAllPricePlans() ([]domain.PricePlan, error) {
	pricePlans, err := s.pricePlans.GetAllPricePlans()
	if err != nil {
		return nil, err
	}
	if len(pricePlans) == 0 {
		return nil, domain.ErrNotFound
	}
	return pricePlans, nil
}
