package repository

import (
	"time"

	"joi-energy-golang/domain"
)

type PricePlans struct {
	pricePlans    []domain.PricePlan
	meterReadings *MeterReadings
}

func NewPricePlans(pricePlans []domain.PricePlan, meterReadings *MeterReadings) PricePlans {
	return PricePlans{
		pricePlans:    pricePlans,
		meterReadings: meterReadings,
	}
}

func (p *PricePlans) ConsumptionCostOfElectricityReadingsForEachPricePlan(smartMeterId string) map[string]float64 {
	electricityReadings := p.meterReadings.GetReadings(smartMeterId)
	costs := map[string]float64{}
	for _, plan := range p.pricePlans {
		costs[plan.PlanName] = calculateCost(electricityReadings, plan)
	}
	return costs
}

func calculateCost(electricityReadings []domain.ElectricityReading, pricePlan domain.PricePlan) float64 {
	average := calculateAverageReading(electricityReadings)
	timeElapsed := calculateTimeElapsed(electricityReadings)
	averagedCost := average / timeElapsed.Hours()
	return averagedCost * pricePlan.UnitRate
}

func calculateAverageReading(electricityReadings []domain.ElectricityReading) float64 {
	sum := 0.0
	for _, r := range electricityReadings {
		sum += r.Reading
	}
	return sum / float64(len(electricityReadings))
}

func calculateTimeElapsed(electricityReadings []domain.ElectricityReading) time.Duration {
	var first, last time.Time
	for _, r := range electricityReadings {
		if r.Time.Before(first) || first.IsZero() {
			first = r.Time
		}
	}
	for _, r := range electricityReadings {
		if r.Time.After(last) || last.IsZero() {
			last = r.Time
		}
	}
	return last.Sub(first)
}

// GetAllPricePlans returns all price plans
func (p *PricePlans) GetAllPricePlans() ([]domain.PricePlan, error) {
	if len(p.pricePlans) == 0 {
		return nil, domain.ErrNotFound
	}
	return p.pricePlans, nil
}

// GetPricePlanById returns a price plan by its ID
func (p *PricePlans) GetPricePlanById(id string) (domain.PricePlan, error) {
	for _, plan := range p.pricePlans {
		if plan.PlanName == id {
			return plan, nil
		}
	}
	return domain.PricePlan{}, domain.ErrNotFound
}

// GetPricePlanByName returns a price plan by its name
func (p *PricePlans) GetPricePlanByName(name string) (domain.PricePlan, error) {
	for _, plan := range p.pricePlans {
		if plan.PlanName == name {
			return plan, nil
		}
	}
	return domain.PricePlan{}, domain.ErrNotFound
}

// GetPricePlanBySupplier returns all price plans for a given supplier
func (p *PricePlans) GetPricePlanBySupplier(supplier string) ([]domain.PricePlan, error) {
	var plans []domain.PricePlan
	for _, plan := range p.pricePlans {
		if plan.EnergySupplier == supplier {
			plans = append(plans, plan)
		}
	}
	if len(plans) == 0 {
		return nil, domain.ErrNotFound
	}
	return plans, nil
}
