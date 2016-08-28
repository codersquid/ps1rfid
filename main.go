/*
 * Copyright 2015 Derek Bever
 *
 * This file is part of ps1rfid.
 *
 * ps1rfid is free software: you can redistribute it and/or modify it under
 * the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option) any
 * later version.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
 * for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/loansindi/ps1rfid/cfg"
	"github.com/loansindi/ps1rfid/ps1rfid"
)

var cacheDB *bolt.DB

func checkCacheDBForTag(tag string) bool {
	val := ""
	cacheDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("RFIDBucket"))
		val = string(b.Get([]byte(tag)))
		return nil
	})

	if val != "" {
		return true
	}

	return false
}

func addTagToCacheDB(tag string) {
	cacheDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("RFIDBucket"))
		err := b.Put([]byte(tag), []byte(tag))
		return err
	})
}

func main() {

	var settingsFile string
	flag.StringVar(&settingsFile, "config", "./config.toml", "Path to the config file")
	var test = flag.Bool("test", false, "Run server in test mode")
	flag.Parse()
	config, err := cfg.ReadConfig(settingsFile)
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}
	fmt.Printf("Config: %v", config)

	var code string
	robot, err := ps1rfid.NewRobotter(config, *test)
	if err != nil {
		log.Fatal(err)
	}

	// the anonymous function here allows us to call openDoor with splate remaining in scope
	go http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Okay"))
		robot.OpenDoor()
	})
	go http.ListenAndServe(":8080", nil)
	buf := make([]byte, 16)
	for {
		n, err := io.ReadFull(u, buf)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		// We need to strip the stop and start bytes from the tag, so we only assign a certain range of the slice
		code = string(buf[1 : n-3])

		// Now open the cache db to check if it's already here
		cacheDB, err = bolt.Open("rfid-tags.db", 0600, nil)
		if err != nil {
			fmt.Println(err)
		}

		cacheDB.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("RFIDBucket"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})

		// Before checking the site for the code, let's check our cache
		if checkCacheDBForTag(code) == false {
			var request bytes.Buffer
			request.WriteString("https://members.pumpingstationone.org/rfid/check/FrontDoor/")
			request.WriteString(code)
			resp, err := http.Get(request.String())
			if err != nil {
				fmt.Printf("Whoops!")
				os.Exit(1)
			}
			if resp.StatusCode == 200 {

				// We got 200 back, so we're good to add this
				// tag to the cache
				addTagToCacheDB(code)

				fmt.Println("Success!")
				code = ""
				openDoor(*splate)
			} else if resp.StatusCode == 403 {
				fmt.Println("Membership status: Expired")
			} else {
				fmt.Println("Code not found")
			}
		} else {
			// If we're here, we found the tag in the cache, so
			// let's just go and open the door for 'em
			fmt.Println("Success!")
			code = ""
			openDoor(*splate)
		}

		cacheDB.Close()
	}
}
