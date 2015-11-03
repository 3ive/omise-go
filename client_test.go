package omise_test

import (
	. "github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	a "github.com/stretchr/testify/assert"
	"testing"
)

// TODO: Include token retrieval and customer creation into the test flow.
const (
	TestCustomer = "cust_test_51w32b7kmr04n6bkcp8"
	TestCard     = "card_test_51w32k09s2eu8azi12r"
)

func TestClient_Charge(t *testing.T) {
	client, op := NewClient(testKeys()), &operations.CreateCharge{
		Amount:   204842,
		Currency: "thb",
		Customer: TestCustomer,
		Card:     TestCard,
	}

	charge, e := client.CreateCharge(op)
	if !(a.NoError(t, e) && a.NotNil(t, charge)) {
		return
	}

	a.Equal(t, op.Amount, charge.Amount)
	a.Equal(t, op.Currency, charge.Currency)
}

func TestClient_InvalidCharge(t *testing.T) {
	client := NewClient(testKeys())

	_, e := client.CreateCharge(&operations.CreateCharge{
		Amount:   12345,
		Currency: "oms", // OMISE DOLLAR, why not?
		Customer: TestCustomer,
		Card:     TestCard,
	})
	a.EqualError(t, e, "(400/invalid_charge) currency is currently not supported")
}