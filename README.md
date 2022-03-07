# Cashbag
This is golang sdk for implement promotion schema
### How to use
```go
package main

import "github.com/hualoqueros/cashbag"

func main() {
    promotion := Promo{
		Name:      "TEST",
		StartAt:   time.Now(),
		ExpiredAt: time.Now(),
		Schemas: []Schema{
			Schema{
				AmountType:     AMOUNT_TYPE_SUBTOTAL,
				ConditionType:  CONDITION_TYPE_MIN_PRICE,
				ConditionValue: "5000",
				RewardType:     REWARD_TYPE_DISCOUNT_AMOUNT,
				RewardValue:    "2500",
			},
		},
		AdditionalInfo: "123",
	}

	shoppingCart := ShoppingCart{
		Carts: []Cart{
			Cart{
				Price: 500000,
				Qty:   1,
			},
		},
		Subtotal:   500000,
		GrandTotal: 500000,
	}
	getPromo := NewPromo(promotion)
	rewards, newGrandTotal, totalDeduction, err := getPromo.Calculate(shoppingCart)
}
```

You can adding custom callback function in calculate using `CalculateWithCallback`
```go
...
    promotion := Promo{
		Name:      "TEST",
		StartAt:   time.Now(),
		ExpiredAt: time.Now(),
		Schemas: []Schema{
			Schema{
				AmountType:     AMOUNT_TYPE_SUBTOTAL,
				ConditionType:  CONDITION_TYPE_MIN_PRICE,
				ConditionValue: "5000",
				RewardType:     REWARD_TYPE_DISCOUNT_AMOUNT,
				RewardValue:    "2500",
			},
		},
		AdditionalInfo: "123",
	}

	shoppingCart := ShoppingCart{
		Carts: []Cart{
			Cart{
				AdditionalID: "KFC-123",
				Price:        500000,
				Qty:          1,
			},
		},
		Subtotal:   500000,
		GrandTotal: 500000,
	}
	getPromo := NewPromo(promotion)
	// here we can create our own validation,
	// in this case, we validate the qty minimum is larger than 10
	checkingSKUAvailibilty := func() (err error) {
		for _, cart := range shoppingCart.Carts {
			skuIsNotAvailable := true
			if cart.Qty < 10 {
			    skuIsNotAvailable = false
			}
			if !skuIsNotAvailable {
				return errors.New(fmt.Sprintf("Please add more quantity for SKU %+s", cart.AdditionalID))
			}
		}
		return nil
	}

	_, _, _, err := getPromo.CalculateWithCallback(shoppingCart, checkingSKUAvailibilty)
...
```