package midtrans

import (
	"errors"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/yogarn/filkompedia-be/entity"
)

type IMidtrans interface {
	NewTransactionToken(orderId string, amount int64, user *entity.User) (*snap.Response, error)
}

type Midtrans struct {
	Client snap.Client
}

func NewMidtrans() IMidtrans {
	client := snap.Client{}
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	return &Midtrans{
		Client: client,
	}
}

func (m *Midtrans) NewTransactionToken(orderId string, amount int64, user *entity.User) (*snap.Response, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderId,
			GrossAmt: amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Username,
			Email: user.Email,
		},
		//todo add more details
	}

	snapResp, err := m.Client.CreateTransaction(req)
	var midtransErr *midtrans.Error
	if errors.As(err, &midtransErr) && midtransErr == nil {
		return snapResp, nil
	}
	return snapResp, err
}
