package service

import "context"

type InvoiceService struct{}

func NewInvoiceService() InvoiceService {
	return InvoiceService{}
}

func (s InvoiceService) RebuildInvoice(invoiceID string, ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	_ = invoiceID
	return nil
}
