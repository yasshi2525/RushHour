package app

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	. "github.com/yasshi2525/RushHour/app/models/entities"
)

var (
	in              = []*Step{}
	out             = []*Step{}
	TestPoint0      = Junction{Point: Point{X: 0, Y: 0}, In: in, Out: out}
	TestPointNear   = Junction{Point: Point{X: 1, Y: 1}, In: in, Out: out}
	TestPointCenter = Point{X: 2, Y: 2}
	TestPointFar    = Junction{Point: Point{X: 3, Y: 3}, In: in, Out: out}
)

func TestInitGame(t *testing.T) {
	testDB, err := gorm.Open("mysql", "rushhourgo:rushhourgo@/rushhourgo?parseTime=True&loc=Asia%2FTokyo")
	defer testDB.Close()

	if err != nil {
		panic("failed to connect database")
	}

	//testDB.LogMode(true)

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

	testDB.Model(&Human{}).AddForeignKey("from_id", "residences(id)", "RESTRICT", "RESTRICT")
	testDB.Model(&Human{}).AddForeignKey("to_id", "companies(id)", "RESTRICT", "RESTRICT")

	t.Run("create Human", func(t *testing.T) {
		var (
			from       Residence
			to         Company
			target     Human
			slave      Human
			fetch      Human
			slavefetch Human
		)

		testDB.FirstOrCreate(&from, Residence{
			Junction: TestPointNear,
		})

		testDB.FirstOrCreate(&to, Company{
			Junction: TestPointFar,
		})

		testDB.FirstOrCreate(&target, Human{
			FromID: from.ID,
			ToID:   to.ID,
			Point:  TestPointCenter,
			On:     OnGround,
		})

		testDB.FirstOrCreate(&slave, Human{
			FromID: from.ID,
			ToID:   to.ID,
			Point:  TestPointCenter,
			On:     OnTrain,
		})

		testDB.Preload("From").Preload("To").Find(&fetch, target.ID)
		testDB.Preload("From").Preload("To").Find(&slavefetch, slave.ID)

		if &fetch == nil {
			t.Error("Human is nil")
		}

		if fetch.From.ID != from.ID {
			t.Error("Human.From is nil")
		}
		if fetch.To.ID != to.ID {
			t.Error("Human.To is nil")
		}
		if fetch.X == 0 || fetch.Y == 0 {
			t.Error("Human.X/Y is nil")
		}

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
			station    Station
			fetch      Station
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
				Ownable: Ownable{OwnerID: testPlayer.ID},
				Point:   TestPointCenter,
			})

		testDB.FirstOrCreate(
			&platform,
			Platform{
				Ownable: Ownable{OwnerID: testPlayer.ID},
				On:      &node,
			})

		testDB.FirstOrCreate(
			&gate,
			Gate{
				Ownable: Ownable{OwnerID: testPlayer.ID},
			})

		testDB.FirstOrCreate(
			&station,
			Station{
				Ownable:    Ownable{OwnerID: testPlayer.ID},
				GateID:     gate.ID,
				PlatformID: platform.ID,
			})

		testDB.Find(&fetch, station.ID)

		if &fetch.Owner == nil {
			t.Error("Station.Owner is nil")
		}
		if &fetch.Gate == nil {
			t.Error("Station.Gate is nil")
		}
		if &fetch.Platform == nil {
			t.Error("Station.Platform is nil")
		}

		testDB.Delete(&station)
		testDB.Delete(&gate)
		testDB.Delete(&platform)
		testDB.Delete(&node)
		testDB.Delete(&testPlayer)
	})

}
