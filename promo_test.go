package cashbag

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPromoTypeDiscountConditionMinPriceWithSubTotalShouldSuccess(t *testing.T) {
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
	assert.Nil(t, err)
	assert.NotNil(t, rewards)
	assert.NotNil(t, totalDeduction)
	assert.Less(t, newGrandTotal, shoppingCart.GrandTotal)
}

func TestPromoTypeDiscountConditionPriceRangeWithSubTotalShouldSuccess(t *testing.T) {
	promotion := Promo{
		Name:      "TEST",
		StartAt:   time.Now(),
		ExpiredAt: time.Now(),
		Schemas: []Schema{
			Schema{
				AmountType:     AMOUNT_TYPE_SUBTOTAL,
				ConditionType:  CONDITION_TYPE_RANGE_PRICE,
				ConditionValue: "50000|100000",
				RewardType:     REWARD_TYPE_DISCOUNT_AMOUNT,
				RewardValue:    "2500",
			},
		},
		AdditionalInfo: "123",
	}

	shoppingCart := ShoppingCart{
		Carts: []Cart{
			Cart{
				Price: 50000,
				Qty:   1,
			},
			Cart{
				Price: 10000,
				Qty:   1,
			},
		},
		Subtotal:   60000,
		GrandTotal: 60000,
	}
	getPromo := NewPromo(promotion)
	rewards, newGrandTotal, totalDeduction, err := getPromo.Calculate(shoppingCart)
	assert.Nil(t, err)
	assert.NotNil(t, rewards)
	assert.NotNil(t, totalDeduction)
	assert.Less(t, newGrandTotal, shoppingCart.GrandTotal)
	assert.Equal(t, totalDeduction, float32(2500))
}

func TestPromoTypeDiscountConditionPriceRangeWithSubTotalShouldFailed(t *testing.T) {
	promotion := Promo{
		Name:      "TEST",
		StartAt:   time.Now(),
		ExpiredAt: time.Now(),
		Schemas: []Schema{
			Schema{
				AmountType:     AMOUNT_TYPE_SUBTOTAL,
				ConditionType:  CONDITION_TYPE_RANGE_PRICE,
				ConditionValue: "50000|100000",
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
			Cart{
				Price: 10000,
				Qty:   1,
			},
		},
		Subtotal:   510000,
		GrandTotal: 510000,
	}
	getPromo := NewPromo(promotion)
	rewards, newGrandTotal, totalDeduction, err := getPromo.Calculate(shoppingCart)
	assert.Nil(t, err)
	assert.NotNil(t, rewards)
	assert.NotNil(t, totalDeduction)
	assert.Less(t, newGrandTotal, shoppingCart.GrandTotal)
	assert.NotEqual(t, totalDeduction, float32(100000))
}

func TestPromoTypeProductShouldSuccess(t *testing.T) {
	// REWARD_TYPE_PRODUCT
	promotion := Promo{
		Name:      "TEST",
		StartAt:   time.Now(),
		ExpiredAt: time.Now(),
		Schemas: []Schema{
			Schema{
				AmountType:     AMOUNT_TYPE_SUBTOTAL,
				ConditionType:  CONDITION_TYPE_MIN_PRICE,
				ConditionValue: "5000",
				RewardType:     REWARD_TYPE_PRODUCT,
				RewardValue:    `{"products":[{"id":"123","name": "Beef Burger","image":"http"}],"qty": 1}`,
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
	assert.Nil(t, err)
	assert.NotNil(t, rewards)
	assert.NotNil(t, totalDeduction)
	assert.Equal(t, shoppingCart.GrandTotal, newGrandTotal)
}

func TestPromoTypeDiscountShouldFailed(t *testing.T) {
	promotion := Promo{
		Name:      "TEST",
		StartAt:   time.Now(),
		ExpiredAt: time.Now(),
		Schemas: []Schema{
			Schema{
				AmountType:     AMOUNT_TYPE_SUBTOTAL,
				ConditionType:  CONDITION_TYPE_MIN_PRICE,
				ConditionValue: "Value is Not Number",
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
		Subtotal:   0,
		GrandTotal: 500000,
	}
	getPromo := NewPromo(promotion)
	_, _, _, err := getPromo.Calculate(shoppingCart)
	assert.NotNil(t, err)

}
