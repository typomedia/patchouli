package boltdb

import (
	"encoding/json"
	"github.com/typomedia/patchouli/app/helper"
	"sort"
	"time"

	"github.com/typomedia/patchouli/app/structs"
	"go.etcd.io/bbolt"
)

var Config = bbolt.Options{
	Timeout: 1 * time.Second,
}

func New() *Bolt {
	db, err := bbolt.Open("patchouli.boltdb", 0600, &Config)
	if err != nil {
		panic(err)
	}
	return &Bolt{db}
}

type Bolt struct {
	db *bbolt.DB
}

func (bolt *Bolt) GetMachine(key string) (interface{}, error) {
	var data []byte
	var result structs.Machine
	err := bolt.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("machine"))
		data = bucket.Get([]byte(key))
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}

func (bolt *Bolt) GetSystem(key string) (interface{}, error) {
	var data []byte
	var result structs.System
	err := bolt.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("systems"))
		data = bucket.Get([]byte(key))
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}

func (bolt *Bolt) Get(key string, result interface{}, bucket string) error {
	var data []byte
	err := bolt.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		data = bucket.Get([]byte(key))
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (bolt *Bolt) GetConfig() (structs.Config, error) {
	var data []byte
	var result structs.Config
	err := bolt.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("config"))
		data = bucket.Get([]byte("main"))
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}

func (bolt *Bolt) GetLastByName(id string, bucket string) ([]byte, error) {
	var result []byte
	err := bolt.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket != nil { // bucket may not exist
			cursor := bucket.Cursor()
			for k, v := cursor.Last(); k != nil; k, v = cursor.Prev() {
				var update structs.Update
				err := json.Unmarshal(v, &update)
				if err != nil {
					return nil
				}
				if update.Machine == id {
					result = v
					break
				}
			}
		}
		return nil
	})
	return result, err
}

func (bolt *Bolt) Set(key string, value interface{}, bucket string) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = bolt.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		bucket.Put([]byte(key), data)
		return nil
	})
	return err
}

func (bolt *Bolt) Delete(key string, bucket string) error {
	err := bolt.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		err := bucket.Delete([]byte(key))
		return err
	})
	return err
}

func (bolt *Bolt) GetAll(bucket string) ([][]byte, error) {
	var result [][]byte
	err := bolt.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			result = append(result, v)
		}
		return nil
	})
	return result, err
}

func (bolt *Bolt) GetAllUpdatesByMachineId(id string) ([][]byte, error) {
	var result [][]byte
	err := bolt.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("history"))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var update structs.Update
			err := json.Unmarshal(v, &update)
			if err != nil {
				return err
			}
			if update.Machine == id {
				result = append(result, v)
			}
		}
		return nil
	})
	return result, err
}

func (bolt *Bolt) GetAllMachinesByOperatorId(id string) (structs.Machines, error) {
	var result structs.Machines
	machines, err := bolt.GetAllMachines(false)
	if err != nil {
		return nil, err
	}
	for _, m := range machines {
		if m.Operator.Id != id {
			continue
		}
		result = append(result, m)
	}
	return result, nil
}

func (bolt *Bolt) GetAllMachinesBySystemId(systemId string) (structs.Machines, error) {
	var result structs.Machines
	machines, err := bolt.GetAllMachines(false)

	for _, m := range machines {
		if m.System.Id != systemId {
			continue
		}
		result = append(result, m)
	}
	return result, err
}

func (bolt *Bolt) GetAllMachines(onlyActive bool) (structs.Machines, error) {
	var result structs.Machines
	config, err := bolt.GetConfig()
	machines, err := bolt.GetAll("machine")
	if err != nil {
		return nil, err
	}
	for _, v := range machines {
		m := structs.Machine{}
		err = json.Unmarshal(v, &m)
		if err != nil {
			return nil, err
		}
		if m.Inactive && onlyActive {
			continue
		}
		if m.Interval == 0 {
			m.Interval = config.General.Interval
		}
		result = append(result, m)
	}

	return result, nil
}

func (bolt *Bolt) GetActiveMachines() (structs.Machines, error) {
	Machines, err := bolt.GetAllMachines(true)
	if err != nil {
		return nil, err
	}

	config, err := bolt.GetConfig()
	if err != nil {
		return nil, err
	}

	for i := range Machines {
		currentDate := time.Now()

		lastUpdate, _ := bolt.GetLastByName(Machines[i].Id, "history")

		if lastUpdate != nil {
			update := structs.Update{}
			err = json.Unmarshal(lastUpdate, &update)
			if err != nil {
				return nil, err
			}
			Machines[i].Update = update
		}

		if Machines[i].Update.Date == "" {
			Machines[i].Update.Date = "0000-00-00"
			Machines[i].Status = "danger"
			continue
		}

		Machines[i].Update.Date = helper.UnixToDateString(Machines[i].Update.Date)

		date, err := time.Parse(time.DateOnly, Machines[i].Update.Date)
		if err != nil {
			return nil, err
		}

		var interval int
		// check if machine has custom interval
		if Machines[i].Interval != 0 {
			interval = Machines[i].Interval
		} else {
			interval = config.General.Interval
		}

		Machines[i].Days = interval - int(currentDate.Sub(date).Hours()/24)

		if Machines[i].Days <= interval/3 && Machines[i].Days > 0 {
			Machines[i].Status = "warning"
		} else if Machines[i].Days <= interval && Machines[i].Days > 0 {
			Machines[i].Status = "success"
		} else {
			Machines[i].Status = "danger"
		}
	}

	// sort machines by oldest update first
	sort.Slice(Machines, func(i, j int) bool {
		return Machines[i].Days < Machines[j].Days
	})

	return Machines, nil
}

func (bolt *Bolt) GetAllOperators() (structs.Operators, error) {
	operators, _ := bolt.GetAll("operator")

	Operators := structs.Operators{}
	InactiveOperators := structs.Operators{}
	for _, v := range operators {
		operator := structs.Operator{}
		err := json.Unmarshal(v, &operator)
		if err != nil {
			return nil, err
		}

		if operator.Inactive {
			InactiveOperators = append(InactiveOperators, operator)
		} else {
			Operators = append(Operators, operator)
		}
	}
	return append(Operators, InactiveOperators...), nil
}

func (bolt *Bolt) GetOperatorById(id string) (structs.Operator, error) {
	var operator structs.Operator
	err := bolt.Get(id, &operator, "operator")
	if err != nil {
		return operator, err
	}
	return operator, nil
}

func (bolt *Bolt) GetMachinesBySystem(systemId string) (structs.Machines, error) {
	machines, err := bolt.GetAll("machine")
	if err != nil {
		return nil, err
	}
	machinesOfSystem := structs.Machines{}
	for _, v := range machines {
		machine := structs.Machine{}
		err := json.Unmarshal(v, &machine)
		if err != nil {
			return nil, err
		}
		if machine.System.Id == systemId {
			machinesOfSystem = append(machinesOfSystem, machine)
		}
	}

	return machinesOfSystem, nil
}

func (bolt *Bolt) GetMachinesByOperator(operatorId string) (structs.Machines, error) {
	machines, err := bolt.GetAll("machine")
	if err != nil {
		return nil, err
	}
	machinesOfOperator := structs.Machines{}
	for _, v := range machines {
		machine := structs.Machine{}
		err := json.Unmarshal(v, &machine)
		if err != nil {
			return nil, err
		}
		if machine.Operator.Id == operatorId {
			machinesOfOperator = append(machinesOfOperator, machine)
		}
	}

	return machinesOfOperator, nil
}

func (bolt *Bolt) Close() error {
	return bolt.db.Close()
}

func (bolt *Bolt) SetBucket(name string) error {
	return bolt.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		return nil
	})
}
