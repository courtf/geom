package wkb

import (
	"encoding/binary"
	"github.com/foobaz/geom"
	"io"
)

func multiPolygonReader(r io.Reader, byteOrder binary.ByteOrder, dimension int) (geom.T, error) {
	var numPolygons uint32
	if err := binary.Read(r, byteOrder, &numPolygons); err != nil {
		return nil, err
	}
	polygons := make([]geom.Polygon, numPolygons)
	for i := range polygons {
		if g, err := Read(r); err == nil {
			var ok bool
			polygons[i], ok = g.(geom.Polygon)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.MultiPolygon(polygons), nil
}

func writeMultiPolygon(w io.Writer, byteOrder binary.ByteOrder, axes uint32, multiPolygon geom.MultiPolygon) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPolygon))); err != nil {
		return err
	}
	for _, polygon := range multiPolygon {
		if err := Write(w, byteOrder, axes, polygon); err != nil {
			return err
		}
	}
	return nil
}
