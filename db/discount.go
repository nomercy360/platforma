package db

import (
	"time"
)

type Discount struct {
	ID         int64      `db:"id" json:"id"`
	Code       string     `db:"code" json:"code"`
	IsActive   bool       `db:"is_active" json:"is_active"`
	Type       string     `db:"type" json:"type"`
	UsageLimit int        `db:"usage_limit" json:"usage_limit"`
	UsageCount int        `db:"usage_count" json:"usage_count"`
	StartsAt   time.Time  `db:"starts_at" json:"starts_at"`
	EndsAt     *time.Time `db:"ends_at" json:"ends_at" extensions:"x-nullable"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`
	Value      int        `db:"value" json:"value"`
}

func (d Discount) IsValid() bool {
	now := time.Now()

	// should be active, not expired, and not reached usage limit - if limit is 0, it means unlimited
	return d.IsActive &&
		d.StartsAt.Before(now) &&
		(d.EndsAt == nil || d.EndsAt.After(now)) &&
		(d.UsageLimit == 0 || d.UsageCount < d.UsageLimit)
}

type DiscountQuery struct {
	Code string
	ID   int64
}

func (s Storage) GetDiscount(query DiscountQuery) (*Discount, error) {
	var discount Discount

	var args interface{}

	q := `
		SELECT id, code, is_active, type, usage_limit, usage_count, starts_at, ends_at, created_at, updated_at, deleted_at, value
		FROM discounts
	`

	if query.Code != "" {
		q += " WHERE code = ?"
		args = query.Code
	} else {
		q += " WHERE id = ?"
		args = query.ID
	}

	err := s.db.QueryRow(q, args).Scan(
		&discount.ID,
		&discount.Code,
		&discount.IsActive,
		&discount.Type,
		&discount.UsageLimit,
		&discount.UsageCount,
		&discount.StartsAt,
		&discount.EndsAt,
		&discount.CreatedAt,
		&discount.UpdatedAt,
		&discount.DeletedAt,
		&discount.Value,
	)

	if IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &discount, nil
}

func (s Storage) UpdateDiscountUsageCount(id int64) error {
	query := `
		UPDATE discounts
		SET usage_count = usage_count + 1
		WHERE id = ?;
	`

	_, err := s.db.Exec(query, id)

	return err
}
