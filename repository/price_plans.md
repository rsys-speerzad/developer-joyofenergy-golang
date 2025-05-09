# price_plans.go Explanation

The `price_plans.go` file defines a repository for managing price plans and calculating the cost of electricity consumption for each price plan. It provides methods to calculate costs based on electricity readings and price plan details. Below is a detailed explanation of its components:

---

## `PricePlans` Struct
```go
type PricePlans struct {
    pricePlans    []domain.PricePlan
    meterReadings *MeterReadings
}
```
The `PricePlans` struct contains two fields:
- **`pricePlans`**: A slice of `domain.PricePlan` objects, representing the available price plans.
- **`meterReadings`**: A pointer to a `MeterReadings` repository, used to fetch electricity readings for smart meters.

---

## Constructor: `NewPricePlans`
```go
func NewPricePlans(pricePlans []domain.PricePlan, meterReadings *MeterReadings) PricePlans {
    return PricePlans{
        pricePlans:    pricePlans,
        meterReadings: meterReadings,
    }
}
```
This function initializes and returns a new `PricePlans` instance:
1. **Input**:
   - A slice of `domain.PricePlan` objects.
   - A `MeterReadings` repository.
2. **Output**: Returns a `PricePlans` struct with the provided price plans and meter readings repository.

---

## Method: `ConsumptionCostOfElectricityReadingsForEachPricePlan`
```go
func (p *PricePlans) ConsumptionCostOfElectricityReadingsForEachPricePlan(smartMeterId string) map[string]float64 {
    electricityReadings := p.meterReadings.GetReadings(smartMeterId)
    costs := map[string]float64{}
    for _, plan := range p.pricePlans {
        costs[plan.PlanName] = calculateCost(electricityReadings, plan)
    }
    return costs
}
```
This method calculates the cost of electricity consumption for each price plan:
1. **Input**: Takes a `smartMeterId` of type `string`.
2. **Operation**:
   - Fetches electricity readings for the given `smartMeterId` using the `MeterReadings` repository.
   - Iterates over the available price plans and calculates the cost for each plan using the `calculateCost` function.
3. **Output**: Returns a map where the keys are price plan names (`string`) and the values are the calculated costs (`float64`).

---

## Helper Function: `calculateCost`
```go
func calculateCost(electricityReadings []domain.ElectricityReading, pricePlan domain.PricePlan) float64 {
    average := calculateAverageReading(electricityReadings)
    timeElapsed := calculateTimeElapsed(electricityReadings)
    averagedCost := average / timeElapsed.Hours()
    return averagedCost * pricePlan.UnitRate
}
```
This function calculates the cost of electricity consumption for a specific price plan:
1. **Input**:
   - A slice of `domain.ElectricityReading` objects.
   - A `domain.PricePlan` object.
2. **Operation**:
   - Calculates the average electricity reading using `calculateAverageReading`.
   - Calculates the time elapsed between the first and last readings using `calculateTimeElapsed`.
   - Computes the cost by multiplying the average consumption rate by the price plan's unit rate.
3. **Output**: Returns the calculated cost as a `float64`.

---

## Helper Function: `calculateAverageReading`
```go
func calculateAverageReading(electricityReadings []domain.ElectricityReading) float64 {
    sum := 0.0
    for _, r := range electricityReadings {
        sum += r.Reading
    }
    return sum / float64(len(electricityReadings))
}
```
This function calculates the average electricity reading:
1. **Input**: A slice of `domain.ElectricityReading` objects.
2. **Operation**: Sums up all the readings and divides by the number of readings.
3. **Output**: Returns the average reading as a `float64`.

---

## Helper Function: `calculateTimeElapsed`
```go
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
```
This function calculates the time elapsed between the first and last electricity readings:
1. **Input**: A slice of `domain.ElectricityReading` objects.
2. **Operation**:
   - Finds the earliest and latest timestamps in the readings.
   - Calculates the duration between these two timestamps.
3. **Output**: Returns the time elapsed as a `time.Duration`.

---

## Key Points
- **Integration with MeterReadings**: The `PricePlans` repository relies on the `MeterReadings` repository to fetch electricity readings.
- **Cost Calculation**: The repository calculates costs based on average consumption rates and price plan unit rates.
- **Helper Functions**: The `calculateAverageReading` and `calculateTimeElapsed` functions encapsulate specific calculations, improving code readability and reusability.
- **Extensibility**: Additional methods can be added to support more complex pricing models or additional features.

---

## Usage
This repository is used in the service layer (e.g., `priceplans/service.go`) to calculate electricity consumption costs for different price plans. It acts as a data access and computation layer.

---

## Example
Here’s an example of how the `PricePlans` repository can be used:
```go
pricePlans := []domain.PricePlan{
    {PlanName: "Plan A", UnitRate: 0.15},
    {PlanName: "Plan B", UnitRate: 0.20},
}

meterReadings := NewMeterReadings(map[string][]domain.ElectricityReading{
    "smart-meter-1": {
        {Time: time.Now().Add(-1 * time.Hour), Reading: 10.0},
        {Time: time.Now(), Reading: 20.0},
    },
})

pricePlansRepo := NewPricePlans(pricePlans, &meterReadings)
costs := pricePlansRepo.ConsumptionCostOfElectricityReadingsForEachPricePlan("smart-meter-1")
fmt.Println(costs)
// Output: map[Plan A:calculated_cost Plan B:calculated_cost]
```

This file is part of the repository layer of the application, responsible for managing price plans and calculating electricity consumption costs.