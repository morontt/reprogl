package repositories

import (
	"database/sql"
	"errors"
	"net"

	"xelbot.com/reprogl/models"
)

type GeolocationRepository struct {
	DB *sql.DB
}

func (gr *GeolocationRepository) FindByIP(ip net.IP) (*models.Geolocation, error) {
	query := `
		SELECT
			gl.ip_long
		FROM geo_location AS gl
		WHERE gl.ip_addr = ?`

	location := &models.Geolocation{IpAddr: ip.String()}

	err := gr.DB.QueryRow(query, ip.String()).Scan(
		&location.ID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.RecordNotFound
		} else {
			return nil, err
		}
	}

	return location, nil
}
