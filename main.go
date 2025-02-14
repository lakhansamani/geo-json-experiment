package main

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	apis_layer "layer_experiment/apis/layer"
)

// SRID=4326;POINT(6.993415176375031 50.61467514050287 231.1882221477385)

type PointZ4326 apis_layer.Point

// Scan implements sql.Scanner to read from PostGIS (Binary WKB)
func (g *PointZ4326) Scan(val interface{}) error {
	if val == nil {
		return nil
	}

	var wkb []byte
	switch v := val.(type) {
	case []byte:
		wkb = v
	case string:
		// Convert hex string to []byte
		var err error
		wkb, err = hex.DecodeString(v)
		if err != nil {
			return fmt.Errorf("failed to decode hex WKB: %w", err)
		}
	default:
		return fmt.Errorf("invalid WKB format: expected []byte or string, got %T", val)
	}

	r := bytes.NewReader(wkb)

	// Read byte order (1 = Little Endian, 0 = Big Endian)
	var byteOrder byte
	if err := binary.Read(r, binary.LittleEndian, &byteOrder); err != nil {
		return fmt.Errorf("failed to read byte order: %w", err)
	}

	// Determine byte order
	var order binary.ByteOrder
	if byteOrder == 0 {
		order = binary.BigEndian
	} else if byteOrder == 1 {
		order = binary.LittleEndian
	} else {
		return fmt.Errorf("invalid byte order: %d", byteOrder)
	}

	// Read geometry type
	var geomType uint32
	if err := binary.Read(r, order, &geomType); err != nil {
		return fmt.Errorf("failed to read geometry type: %w", err)
	}

	// Handle possible SRID presence (EWKB format)
	const wkbPointZ = 0x80000001 // WKB PointZ type (without SRID flag)
	const ewkbSRIDFlag = 0x20000000

	hasSRID := (geomType & ewkbSRIDFlag) != 0
	geomType = geomType &^ ewkbSRIDFlag // Remove SRID flag if present

	if geomType != wkbPointZ {
		return fmt.Errorf("unexpected WKB type: got %d, expected %d", geomType, wkbPointZ)
	}

	// Read SRID if present
	if hasSRID {
		var srid uint32
		if err := binary.Read(r, order, &srid); err != nil {
			return fmt.Errorf("failed to read SRID: %w", err)
		}
	}

	// Read X, Y, Z coordinates
	if err := binary.Read(r, order, &g.X); err != nil {
		return fmt.Errorf("failed to read X: %w", err)
	}
	if err := binary.Read(r, order, &g.Y); err != nil {
		return fmt.Errorf("failed to read Y: %w", err)
	}
	if err := binary.Read(r, order, &g.Z); err != nil {
		return fmt.Errorf("failed to read Z: %w", err)
	}

	return nil
}

// Value implements the driver Valuer interface for GeometryPoint
func (g PointZ4326) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%f %f %f)", g.X, g.Y, g.Z), nil
}

// struct for the data
type LayerData struct {
	// Id is the unique identifier of the Point
	Id string `gorm:"type:uuid;primaryKey"`
	// CreatedAt is the timestamp when the Point was created
	CreatedAt time.Time `gorm:"autoCreateTime"`
	// LastModifiedAt is the timestamp when the Point was last modified
	LastModifiedAt time.Time `gorm:"autoUpdateTime"`
	// Point is the geometry of the Point
	Point PointZ4326 `gorm:"type:geometry(PointZ,4326)"`
}

// BeforeCreate hook is used to initialize a new ID for the Point entity
func (t *LayerData) BeforeCreate(tx *gorm.DB) error {
	t.Id = uuid.NewString()
	return nil
}

// connect to db and create table
func connectToDB() *gorm.DB {
	// connect to db
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// create table
	db.AutoMigrate(&LayerData{})
	return db
}

func main() {
	db := connectToDB()

	layer1 := &LayerData{
		Point: PointZ4326{X: 6.993415176375031, Y: 50.61467514050287, Z: 231.1882221477385},
	}
	err := db.Debug().Save(layer1).Error
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=> layer1 created: ", layer1.Id)

	// Get the layer1 from the database
	var layer1Get LayerData
	err = db.Debug().First(&layer1Get).Where("id = ?", layer1.Id).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=> layer1Get: ", layer1.Id, layer1Get.Point.X, layer1Get.Point.Y, layer1Get.Point.Z)

}
