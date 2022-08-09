package basket

import (
	"log"
	"time"

	entities "github.com/mataliksamil/Go_Bootcamp_Final/entities"
)

// checks every discount options updates TotalCost and TotalVATCost of the basket
func ApplyDiscount(user_id string) error {

	totalCost, totalVAT, err := calculateTotalCost(user_id)
	if err != nil {
		return err
	}
	// resultTotal[][] for final result
	// FOR THE SAKE OF FUNCTIONALITY
	// I SWEAR I HAD TO DO THIS IM SO SORRY  :(
	var resultTotal [4][2]float64

	resultTotal[0][0], resultTotal[0][1] = totalCost, totalVAT
	if returnIfAboveMonthly(user_id) {
		resultTotal[1][0], resultTotal[1][1] = totalCost*0.9, totalVAT*0.9
	} else {
		resultTotal[1][0], resultTotal[1][1] = totalCost, totalVAT
	}

	if returnIfFourthOrder(user_id) {
		resultTotal[2][0], resultTotal[2][1], err = Discount4thOrder(user_id)
		if err != nil {
			return err
		}
	} else {
		resultTotal[2][0], resultTotal[2][1] = totalCost, totalVAT
	}

	resultTotal[3][0], resultTotal[3][1], err = Discount4thItem(user_id)
	if err != nil {
		return err
	}

	for i := 1; i < 4; i++ {
		if resultTotal[i][1] < resultTotal[0][1] {
			resultTotal[0][0], resultTotal[0][1] = resultTotal[i][0], resultTotal[i][1]
		}
	}

	updateTotals(resultTotal[0][0], resultTotal[0][1], user_id)
	return nil
}

func returnIfAboveMonthly(user_id string) bool {
	monthAgo := time.Now().AddDate(0, -1, 0)
	monthlyCostVAT := 0.0
	baskets := &[]entities.Basket{}
	err := dbConnect.Model(baskets).
		Where("user_id=?", user_id).
		Where("updated_at>?", monthAgo).
		Select()
	//basket.BasketProducts[].Product
	if err != nil {
		log.Printf("No previous orders, Reason: %v\n", err)
		return false
	}
	for _, b := range *baskets {
		if b.UpdatedAt.After(monthAgo) {
			monthlyCostVAT = monthlyCostVAT + b.TotalVAT
		}
	}
	return monthlyCostVAT >= MOUNTHLYTOTAL
}

func returnIfFourthOrder(user_id string) bool {
	baskets := &[]entities.Basket{}
	// JOIN for whole structure under a user
	err := dbConnect.Model(baskets).
		Where("user_id=?", user_id).Order("updated_at DESC").
		Select()
	//basket.BasketProducts[].Product
	if err != nil {
		log.Printf("No previous orders, Reason: %v\n", err)
		return false
	}

	if len(*baskets) == 0 {
		return false
	}
	return len(*baskets)/4 == 0
}
