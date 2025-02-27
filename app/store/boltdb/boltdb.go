package boltdb

import (
	"encoding/json"
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

func (bolt *Bolt) Delete(key string) error {
	err := bolt.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("machine"))
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

func (bolt *Bolt) GetAllByName(id string, bucket string) ([][]byte, error) {
	var result [][]byte
	err := bolt.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
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

func (bolt *Bolt) GetAllByOperatorId(id string, bucket string) ([][]byte, error) {
	var result [][]byte
	err := bolt.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var machine structs.Machine
			err := json.Unmarshal(v, &machine)
			if err != nil {
				return err
			}
			if machine.Operator.Id == id {
				result = append(result, v)
			}
		}
		return nil
	})
	return result, err
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
