package main

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid" // WKB encoding/decoding
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	apis_layer "layer_experiment/apis/layer"
)

// SRID=4326;POINT(6.993415176375031 50.61467514050287 231.1882221477385)
// Point3D represents a 3D point with X, Y, Z coordinates and default SRID 4326
type Point3D struct {
	SRID int     `json:"srid"` // Default SRID is 4326
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
}

// NewPoint3D creates a new Point3D with default SRID 4326
func NewPoint3D(x, y, z float64) Point3D {
	return Point3D{
		SRID: 4326,
		X:    x,
		Y:    y,
		Z:    z,
	}
}

// Scan implements sql.Scanner for GORM to read from PostGIS
func (p *Point3D) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var wkb []byte
	switch v := value.(type) {
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
		return fmt.Errorf("invalid WKB format: expected []byte or string, got %T", value)
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
		} else {
			p.SRID = int(srid)
		}
	}

	// Read X, Y, Z coordinates
	if err := binary.Read(r, order, &p.X); err != nil {
		return fmt.Errorf("failed to read X: %w", err)
	}
	if err := binary.Read(r, order, &p.Y); err != nil {
		return fmt.Errorf("failed to read Y: %w", err)
	}
	if err := binary.Read(r, order, &p.Z); err != nil {
		return fmt.Errorf("failed to read Z: %w", err)
	}

	return nil
}

// Value implements driver.Valuer for saving into the database
func (p Point3D) Value() (driver.Value, error) {
	// Ensure SRID is set to default if not already
	if p.SRID == 0 {
		p.SRID = 4326
	}
	return fmt.Sprintf("SRID=%d;POINT(%f %f %f)", p.SRID, p.X, p.Y, p.Z), nil
}

// ToProto converts Point3D to google.protobuf.ListValue (ignoring SRID)
func (p *Point3D) ToProto() (*structpb.ListValue, error) {
	return structpb.NewList([]interface{}{p.X, p.Y, p.Z})
}

// FromProto initializes Point3D from google.protobuf.ListValue (ignoring SRID)
func (p *Point3D) FromProto(lv *structpb.ListValue) error {
	if len(lv.Values) != 3 {
		return fmt.Errorf("expected 3 values in ListValue (X, Y, Z), got %d", len(lv.Values))
	}

	p.X = lv.Values[0].GetNumberValue()
	p.Y = lv.Values[1].GetNumberValue()
	p.Z = lv.Values[2].GetNumberValue()
	p.SRID = 4326 // Default to 4326 always
	return nil
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
	Point Point3D `gorm:"type:geometry(PointZ,4326)"`
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
	arr, _ := structpb.NewList([]interface{}{1.2, 21.2, 78.1})
	g := apis_layer.Geometry{
		Type:        "Point",
		Coordinates: arr,
	}
	// Marshal the Geometry message to json
	bt, err := json.Marshal(g)
	if err != nil {
		fmt.Println("Error marshalling the Geometry message to json: ", err)
	}
	// {"type":"Point","coordinates":[1.2,21.2,78.1]}
	fmt.Println("=> Marshalled Geometry message to json: ", string(bt))

	f := apis_layer.Feature{
		Type: "Feature",
		Properties: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"name": structpb.NewStringValue("test"),
			},
		},
		Geometry: &g,
	}
	// Marshal the Feature message to json
	bt, err = json.Marshal(f)
	if err != nil {
		fmt.Println("Error marshalling the Feature message to json: ", err)
	}
	// {"type":"Feature","geometry":{"type":"point","coordinates":[1.2,21.2,78.1]},"properties":{"name":"test"}}
	fmt.Println("=> Marshalled Feature message to json: ", string(bt))

	db := connectToDB()
	point := Point3D{}
	point.FromProto(arr)
	layer1 := &LayerData{
		Point: point,
	}
	err = db.Debug().Save(layer1).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=> layer1 created: ", layer1.Id)

	// // Get the layer1 from the database
	var layer1Get LayerData
	err = db.Debug().Where("id = ?", layer1.Id).First(&layer1Get).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=> layer1Get: ", layer1.Id, layer1Get.Point.X, layer1Get.Point.Y, layer1Get.Point.Z)

}
