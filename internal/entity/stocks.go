package entity

import "time"

type Stocks struct {
  ID int64
  Quantity int
  CreatedAt time.Time
  UpdatedAt time.Time

  Barcode Barcode
  Warehouse Warehouse
  Card Card
  Seller Seller
}
