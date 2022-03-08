package cashbag

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	Price        float32
	Qty          int
	AdditionalID string
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

type CallbackFunction func() (err error)

func NewPromo(promo Promo) Promo {
	return promo
}

func (p *Promo) Calculate(shoppingCart ShoppingCart) (rewards []Reward, grandTotal float32, amountOfDeduction float32, err error) {
	for _, schema := range p.Schemas {
		switch schema.ConditionType {
		case CONDITION_TYPE_RANGE_PRICE:
			valid, amountWillDeduct, err := checkConditionPriceRange(shoppingCart, schema)
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
				grandTotal = shoppingCart.GrandTotal - amountOfDeduction
				return rewards, grandTotal, amountOfDeduction, err
			}
		case CONDITION_TYPE_MIN_PRICE:
			valid, amountWillDeduct, err := checkConditionMinPrice(shoppingCart, schema)
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
				grandTotal = shoppingCart.GrandTotal - amountOfDeduction
				return rewards, grandTotal, amountOfDeduction, err
			}
			break
		}
	}

	return

}

func (p *Promo) CalculateWithCallback(shoppingCart ShoppingCart, callback CallbackFunction) (rewards []Reward, grandTotal float32, amountOfDeduction float32, err error) {
	for _, schema := range p.Schemas {
		switch schema.ConditionType {
		case CONDITION_TYPE_RANGE_PRICE:
			valid, amountWillDeduct, err := checkConditionPriceRange(shoppingCart, schema)
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
				grandTotal = shoppingCart.GrandTotal - amountOfDeduction
				return rewards, grandTotal, amountOfDeduction, err
			}
		case CONDITION_TYPE_MIN_PRICE:
			valid, amountWillDeduct, err := checkConditionMinPrice(shoppingCart, schema)
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
				grandTotal = shoppingCart.GrandTotal - amountOfDeduction
				return rewards, grandTotal, amountOfDeduction, err
			}
			break
		}
	}
	err = callback()
	grandTotal = shoppingCart.GrandTotal - amountOfDeduction
	if err == nil {
		object := p.Schemas[0].AmountType
		condition := p.Schemas[0].ConditionType
		desc := ""
		priceWording := ""
		if condition == CONDITION_TYPE_RANGE_PRICE {
			desc = strings.ReplaceAll(p.Schemas[0].ConditionValue, "|", " - ")
			desc = fmt.Sprintf("%s", desc)
			priceWording = "Range Price is"
		}
		if condition == CONDITION_TYPE_MIN_PRICE {
			desc = strings.ReplaceAll(p.Schemas[0].ConditionValue, "|", " - ")
			desc = fmt.Sprintf("%s", desc)
			priceWording = "Minimum Price is"
		}

		err = errors.New(strings.ToLower(fmt.Sprintf("%s %s %s", object, priceWording, desc)))
	}
	return
}

func checkConditionMinPrice(shoppingCart ShoppingCart, schema Schema) (valid bool, amountWillDeduct float32, err error) {
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
	return
}

func checkConditionPriceRange(shoppingCart ShoppingCart, schema Schema) (valid bool, amountWillDeduct float32, err error) {
	switch schema.AmountType {
	case AMOUNT_TYPE_SUBTOTAL:
		priceWillCompare := shoppingCart.Subtotal
		rangePrice := strings.Split(schema.ConditionValue, "|")
		lowRange, err := strconv.ParseFloat(rangePrice[0], 32)
		if err != nil {
			return false, 0, err
		}
		highRange, err := strconv.ParseFloat(rangePrice[1], 32)
		if err != nil {
			return false, 0, err
		}

		if float32(lowRange) > priceWillCompare {
			return false, 0, nil
		}
		if priceWillCompare > float32(highRange) {
			return false, 0, nil
		}
		valid = true
		amountWillDeduct = shoppingCart.GrandTotal
		break
	}
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
