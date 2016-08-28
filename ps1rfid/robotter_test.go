// -*- Mode: Go; indent-tabs-mode: t -*-
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

package ps1rfid

import (
	"testing"

	"github.com/hybridgroup/gobot/platforms/gpio"
)

func TestHappyTestRobotter(t *testing.T) {
	tbot := NewRobotter(TESTMODE)
	var pin gpio.DirectPinDriver
	err := tbot.openDoor(pin)
	if err != nil {
		t.Error("Expected no error for TESTMODE")
	}
}
func TestSadTestRobotter(t *testing.T) {
	tbot := NewRobotter(ERRMODE)
	var pin gpio.DirectPinDriver
	err := tbot.openDoor(pin)
	if err == nil {
		t.Error("Expected error for TESTMODE")
	}
}
