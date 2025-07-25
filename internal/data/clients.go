package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Client struct {
	ID                     int64     `json:"id"`
	CreatedAt              time.Time `json:"created_at"`
	FirstName              string    `json:"first_name"`
	MiddleName             string    `json:"middle_name,omitempty"`
	LastName               string    `json:"last_name"`
	Suffix                 string    `json:"suffix,omitempty"`
	Email                  string    `json:"email"`
	Mobile                 string    `json:"mobile,omitempty"`
	Username               string    `json:"username"`
	RimanUserId            int64     `json:"riman_user_id"`
	Status                 bool      `json:"status"`
	OrganizationType       string    `json:"organization_type,omitempty"`
	SignupDate             string    `json:"signup_date"`
	AnniversaryDate        string    `json:"anniversary_date"`
	AccountType            string    `json:"account_type"`
	SponsorUsername        string    `json:"sponsor_username,omitempty"`
	MemberId               string    `json:"member_id"`
	Rank                   string    `json:"rank,omitempty"`
	EnrollmentDate         string    `json:"enrollment_date,omitempty"`
	PersonalOrdersVolume   float64   `json:"personal_orders_volume,omitempty"`
	PersonalClientsVolume  float64   `json:"personal_clients_volume,omitempty"`
	TotalPersonalVolume    float64   `json:"total_personal_volume,omitempty"`
	CurrentMonthSp         float64   `json:"current_month_sp,omitempty"`
	CurrentMonthBp         float64   `json:"current_month_bp,omitempty"`
	LastOrderDate          string    `json:"last_order_date,omitempty"`
	LastOrderId            int64     `json:"last_order_id,omitempty"`
	LastOrderSp            float64   `json:"last_order_sp,omitempty"`
	LastOrderBp            float64   `json:"last_order_bp,omitempty"`
	LifetimeSpend          float64   `json:"lifetime_spend,omitempty"`
	MostRecent12MonthSpend float64   `json:"most_recent_12_month_spend,omitempty"`
	Data                   any       `json:"data,omitempty"`
	PasswordHash           string    `json:"password_hash"`
}

type ClientModel struct {
	DB *sql.DB
}

func (m ClientModel) GetAll() ([]*Client, Metadata, error) {

	query := fmt.Sprintf(`
	SELECT count(*) OVER(),id, created_at, first_name, middle_name, last_name, suffix, email, mobile, username, riman_user_id, status, organization_type, signup_date, anniversary_date, account_type, sponsor_username, member_id, rank, enrollment_date, personal_orders_volume, personal_clients_volume, total_personal_volume, current_month_sp, current_month_bp, last_order_date, last_order_id, last_order_sp, last_order_bp, lifetime_spend, most_recent_12_month_spend, data
	FROM clients
	`)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	totalRecords := 0
	clients := []*Client{}

	for rows.Next() {
		var client Client

		err := rows.Scan(
			&totalRecords,
			&client.ID,
			&client.CreatedAt,
			&client.FirstName,
			&client.MiddleName,
			&client.LastName,
			&client.Suffix,
			&client.Email,
			&client.Mobile,
			&client.Username,
			&client.RimanUserId,
			&client.Status,
			&client.OrganizationType,
			&client.SignupDate,
			&client.AnniversaryDate,
			&client.AccountType,
			&client.SponsorUsername,
			&client.MemberId,
			&client.Rank,
			&client.EnrollmentDate,
			&client.PersonalOrdersVolume,
			&client.PersonalClientsVolume,
			&client.TotalPersonalVolume,
			&client.CurrentMonthSp,
			&client.CurrentMonthBp,
			&client.LastOrderDate,
			&client.LastOrderId,
			&client.LastOrderSp,
			&client.LastOrderBp,
			&client.LifetimeSpend,
			&client.MostRecent12MonthSpend,
			&client.Data,
		)
		if err != nil {
			return nil, Metadata{}, err

		}

		clients = append(clients, &client)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, 1, 20)

	return clients, metadata, nil
}
