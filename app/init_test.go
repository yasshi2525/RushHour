package app

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	. "github.com/yasshi2525/RushHour/app/models/entities"
)

var (
	TestPoint0      = Point{X: 0, Y: 0}
	TestPointNear   = Point{X: 1, Y: 1}
	TestPointCenter = Point{X: 2, Y: 2}
	TestPointFar    = Point{X: 3, Y: 3}
)

func TestInitGame(t *testing.T) {
	testDB, err := gorm.Open("mysql", "rushhourgo:rushhourgo@/rushhourgo?parseTime=True&loc=Asia%2FTokyo")
	defer testDB.Close()

	if err != nil {
		panic("failed to connect database")
	}

	testDB.AutoMigrate(
		&Company{},
		&Residence{},
		&Human{},
		&Player{},
		&RailNode{},
		&RailEdge{},
		&Platform{},
		&Gate{},
		&Station{},
	)

	t.Run("create Human", func(t *testing.T) {
		var (
			from   Residence
			to     Company
			target Human
		)

		testDB.FirstOrCreate(&from, Residence{
			Point: TestPointNear,
		})

		testDB.FirstOrCreate(&to, Company{
			Point: TestPointFar,
		})

		testDB.FirstOrCreate(&target, Human{
			FromRefer: from.ID,
			ToRefer:   to.ID,
			Point:     TestPointCenter,
		})

		testDB.Delete(&from)
		testDB.Delete(&to)
		testDB.Delete(&target)
	})

	t.Run("create Station", func(t *testing.T) {
		var (
			testPlayer Player
			node       RailNode
			platform   Platform
			gate       Gate
		)

		testDB.FirstOrCreate(
			&testPlayer,
			Player{
				DisplayName: "testPlayer",
				Password:    "testPassword",
			})

		testDB.FirstOrCreate(
			&node,
			RailNode{
				Owner: Owner{OwnerRefer: testPlayer.ID},
				Point: TestPointCenter,
			})

		testDB.FirstOrCreate(
			&platform,
			Platform{
				Owner: Owner{OwnerRefer: testPlayer.ID},
				On:    node,
			})

		testDB.FirstOrCreate(
			&gate,
			Gate{
				Owner: Owner{OwnerRefer: testPlayer.ID},
			})

		testDB.Delete(&gate)
		testDB.Delete(&platform)
		testDB.Delete(&node)
		testDB.Delete(&testPlayer)
	})

}
