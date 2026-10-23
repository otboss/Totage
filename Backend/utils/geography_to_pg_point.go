package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

func GeographyToPgPoint(v interface{}) (pgtype.Point, error) {
	raw, ok := v.([]byte)
	if !ok {
		return pgtype.Point{}, errors.New("last_location is not []byte")
	}

	g, err := ewkb.Unmarshal(raw)
	if err != nil {
		return pgtype.Point{}, err
	}

	p, ok := g.(*geom.Point)
	if !ok {
		return pgtype.Point{}, errors.New("geometry is not POINT")
	}

	coords := p.Coords()

	return pgtype.Point{
		P: pgtype.Vec2{
			X: coords[0], // longitude
			Y: coords[1], // latitude
		},
		Valid: true,
	}, nil
}
