package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
)

func CreateTrain(o *entities.Player, name string) (*entities.Train, error) {
	t := Model.NewTrain(o, name)
	AddOpLog("CreateTrain", o, t)
	return t, nil
}

func DeployTrain(o *entities.Player, t *entities.Train, l *entities.RailLine) error {
	if err := CheckAuth(o, t); err != nil {
		return err
	}
	if err := CheckAuth(o, l); err != nil {
		return err
	}
	if !l.IsRing() {
		return fmt.Errorf("try to deploy unringed RailLine: %v", l)
	}
	var start *entities.LineTask
	for _, lt := range l.Tasks {
		start = lt
		break
	}
	t.SetTask(start)
	AddOpLog("DeployTrain", o, t, start)
	return nil
}

func UnDeployTrain(o *entities.Player, t *entities.Train) error {
	if err := CheckAuth(o, t); err != nil {
		return err
	}
	t.UnLoad()
	t.SetTask(nil)
	return nil
}

func RemoveTrain(o *entities.Player, id uint) error {
	if t, err := Model.DeleteIf(o, entities.TRAIN, id); err != nil {
		return err
	} else {
		AddOpLog("RemoveTrain", o, t)
		return nil
	}
}
