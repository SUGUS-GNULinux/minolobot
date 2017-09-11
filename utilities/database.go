// Copyright 2017-2018 SUGUS GNU/Linux <sugus@us.es>
//
// This file is part of Minolobot.
//
//     Minolobot is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     Minolobot is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with Minolobot.  If not, see <http://www.gnu.org/licenses/>.

// Package config holds everything related to configuration
package utilities

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

var (
	// BoltDB wrapper
	BoltDB *bolt.DB
)

const BoltDBPath = "minolobot.db"

// Connect to the database
func init() {
	var err error
	// Connect to Bolt
	if BoltDB, err = bolt.Open(BoltDBPath, 0600, &bolt.Options{Timeout: 1 * time.Second}); err != nil {
		log.Fatal("Bolt Driver Error", err)
	}
	// defer BoltDB.Close()  // Not close the connection because we need it permanetly open
}

// Update makes a modification to Bolt
func Update(bucketName string, key string, dataStruct interface{}) error {
	err := BoltDB.Update(func(tx *bolt.Tx) error {
		// Create the bucket
		bucket, e := tx.CreateBucketIfNotExists([]byte(bucketName))
		if e != nil {
			return e
		}

		// Encode the record
		encodedRecord, e := json.Marshal(dataStruct)
		if e != nil {
			return e
		}

		// Store the record
		if e = bucket.Put([]byte(key), encodedRecord); e != nil {
			return e
		}
		return nil
	})
	return err
}

// View retrieves a record in Bolt
func View(bucketName string, key string, dataStruct interface{}) error {
	err := BoltDB.View(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		// Retrieve the record
		v := b.Get([]byte(key))
		if len(v) < 1 {
			return bolt.ErrInvalid
		}

		// Decode the record
		e := json.Unmarshal(v, &dataStruct)
		if e != nil {
			return e
		}

		return nil
	})

	return err
}

// ViewAll retrieves all records in Bolt
func ViewAll(bucketName string, dataStruct []interface{}) error {
	err := BoltDB.View(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		// Retrieve the records
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// Decode the record
			var value interface{}
			err := json.Unmarshal(v, &value)
			if err != nil {
				log.Println(err)
			}
			dataStruct = append(dataStruct, value)
		}
		return nil
	})

	return err
}

// Delete removes a record from Bolt
func Delete(bucketName string, key string) error {
	err := BoltDB.Update(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		return b.Delete([]byte(key))
	})
	return err
}
