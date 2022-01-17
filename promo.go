package cashbag

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	AMOUNT_TYPE_SUBTOTAL   = "SUBTOTAL"
	AMOUNT_TYPE_GRANDTOTAL = "GRANDTOTAL"

	CONDITION_TYPE_MIN_PRICE         = "MIN_PRICE"
	CONDITION_TYPE_RANGE_PRICE       = "RANGE_PRICE"
	CONDITION_TYPE_SPESIFIC_SKU      = "SPESIFIC_SKU"
	CONDITION_TYPE_SPESIFIC_CATEGORY = "SPESIFIC_CATEGORY"

	REWARD_TYPE_DISCOUNT_AMOUNT     = "DISCOUNT_AMOUNT"
	REWARD_TYPE_DISCOUNT_PERCENTAGE = "DISCOUNT_PERCENTAGE"
	REWARD_TYPE_PRODUCT             = "PRODUCT"

	SCHEMA_STATUS_ACTIVE   = "ACTIVE"
	SCHEMA_STATUS_INACTIVE = "INACTIVE"
)

/**
 * 1.define initial struct
 * 2.checking promotion with schema
 * 3.give result
 * */

type Promo struct {
	Name           string
	StartAt        time.Time
	ExpiredAt      time.Time
	Schemas        []Schema
	AdditionalInfo interface{}
}

type Schema struct {
	AmountType     string
	ConditionType  string
	ConditionValue string
	RewardType     string
	RewardValue    string
	Status         string
}

type ShoppingCart struct {
	Carts      []Cart
	Subtotal   float32
	GrandTotal float32
}

type Cart struct {
	Price float32
	Qty   int
}

type Reward struct {
	RewardType   string
	RewardValue  string
	RewardResult string
}

type RewardProduct struct {
	Products []ProductReward `json:"products"`
	Qty      int             `json:"qty"`
}

type ProductReward struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type PromoRewarded struct {
	Carts             []Cart
	Reward            []Reward
	TotalAmountReward float32
}

func NewPromo(promo Promo) Promo {
	return promo
}

func (p *Promo) Calculate(shoppingCart ShoppingCart) (rewards []Reward, grandTotal float32, amountOfDeduction float32, err error) {
	for _, schema := range p.Schemas {
		switch schema.ConditionType {
		case CONDITION_TYPE_MIN_PRICE:
			valid, amountWillDeduct, err := checkConditionMinRange(shoppingCart, schema)
			if err != nil {
				return []Reward{}, 0, 0, err
			}
			if valid {
				rewardAmount, reward, err := getReward(amountWillDeduct, schema)
				if err != nil {
					return []Reward{}, 0, 0, err
				}
				rewards = append(rewards, reward)
				amountOfDeduction += rewardAmount
			}
			break
		}
	}
	grandTotal = shoppingCart.GrandTotal - amountOfDeduction
	return
}

func checkConditionMinRange(shoppingCart ShoppingCart, schema Schema) (valid bool, amountWillDeduct float32, err error) {
	switch schema.AmountType {
	case AMOUNT_TYPE_SUBTOTAL:
		limit, err := strconv.ParseFloat(schema.ConditionValue, 32)
		if err != nil {
			return false, 0, err
		}

		if shoppingCart.Subtotal < float32(limit) {
			return false, 0, nil
		}
		valid = true
		amountWillDeduct = shoppingCart.Subtotal
		break
	}
	// case AMOUNT_TYPE_SUBTOTAL:
	// 	log.Println("MASUK")
	// 	rangePrice := strings.Split(schema.ConditionValue, "|")
	// 	lowRange, err := strconv.ParseFloat(rangePrice[0], 32)
	// 	if err != nil {
	// 		return false, 0, err
	// 	}
	// 	highRange, err := strconv.ParseFloat(rangePrice[1], 32)
	// 	if err != nil {
	// 		return false, 0, err
	// 	}
	// 	if float32(lowRange) > shoppingCart.Subtotal && shoppingCart.Subtotal > float32(highRange) {
	// 		return false, 0, nil
	// 	}
	// 	valid = false
	// 	amountWillDeduct = shoppingCart.GrandTotal
	// 	break
	// }
	return
}

func getReward(amountWillDeduct float32, schema Schema) (rewardAmount float32, reward Reward, err error) {
	switch schema.RewardType {
	case REWARD_TYPE_PRODUCT:
		var rewardProduct RewardProduct
		json.Unmarshal([]byte(schema.RewardValue), &rewardProduct)
		reward = Reward{
			RewardType:  schema.RewardType,
			RewardValue: schema.RewardValue,
		}
	case REWARD_TYPE_DISCOUNT_AMOUNT:
		rewardValue, err := strconv.ParseFloat(schema.RewardValue, 32)
		if err != nil {
			return 0, Reward{}, err
		}
		if amountWillDeduct > 0 {
			rewardAmount += float32(rewardValue)
			reward = Reward{
				RewardType:  schema.RewardType,
				RewardValue: schema.RewardValue,
			}
		}
	}

	return
}
