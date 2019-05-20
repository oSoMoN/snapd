// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018-2019 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package snapstate_test

import (
	. "gopkg.in/check.v1"

	"github.com/snapcore/snapd/overlord/snapstate"
	"github.com/snapcore/snapd/overlord/snapstate/snapstatetest"
	"github.com/snapcore/snapd/overlord/state"
)

type deviceCtxSuite struct {
	st *state.State
}

var _ = Suite(&deviceCtxSuite{})

func (s *deviceCtxSuite) SetUpTest(c *C) {
	s.st = state.New(nil)
}

func (s *deviceCtxSuite) TestDevicePastSeedingTooEarly(c *C) {
	s.st.Lock()
	defer s.st.Unlock()

	r := snapstatetest.MockDeviceModel(nil)
	defer r()

	expectedErr := &snapstate.ChangeConflictError{
		Message: "too early for operation, device not yet seeded or" +
			" device model not acknowledged",
		ChangeKind: "seed",
	}

	// not seeded, no model assertion
	_, err := snapstate.DevicePastSeeding(s.st, nil)
	c.Assert(err, DeepEquals, expectedErr)

	// seeded, no model assertion
	s.st.Set("seeded", true)
	_, err = snapstate.DevicePastSeeding(s.st, nil)
	c.Assert(err, DeepEquals, expectedErr)
}

func (s *deviceCtxSuite) TestDevicePastSeedingProvided(c *C) {
	s.st.Lock()
	defer s.st.Unlock()

	r := snapstatetest.MockDeviceContext(nil)
	defer r()

	expectedErr := &snapstate.ChangeConflictError{
		Message: "too early for operation, device not yet seeded or" +
			" device model not acknowledged",
		ChangeKind: "seed",
	}

	deviceCtx1 := &snapstatetest.TrivialDeviceContext{DeviceModel: MakeModel(nil)}

	// not seeded
	_, err := snapstate.DevicePastSeeding(s.st, deviceCtx1)
	c.Assert(err, DeepEquals, expectedErr)

	// seeded
	s.st.Set("seeded", true)
	deviceCtx, err := snapstate.DevicePastSeeding(s.st, deviceCtx1)
	c.Assert(err, IsNil)
	c.Assert(deviceCtx, Equals, deviceCtx1)
}

func (s *deviceCtxSuite) TestDevicePastSeedingReady(c *C) {
	s.st.Lock()
	defer s.st.Unlock()

	// seeded and model assertion
	s.st.Set("seeded", true)

	r := snapstatetest.MockDeviceModel(DefaultModel())
	defer r()

	deviceCtx, err := snapstate.DevicePastSeeding(s.st, nil)
	c.Assert(err, IsNil)
	c.Check(deviceCtx.Model().Model(), Equals, "baz-3000")
}

func (s *deviceCtxSuite) TestDeviceCtxFromStateReady(c *C) {
	s.st.Lock()
	defer s.st.Unlock()

	// model assertion but not seeded yet
	r := snapstatetest.MockDeviceModel(DefaultModel())
	defer r()

	deviceCtx, err := snapstate.DeviceCtxFromState(s.st, nil)
	c.Assert(err, IsNil)
	c.Check(deviceCtx.Model().Model(), Equals, "baz-3000")
}

func (s *deviceCtxSuite) TestDeviceCtxFromStateProvided(c *C) {
	s.st.Lock()
	defer s.st.Unlock()

	r := snapstatetest.MockDeviceContext(nil)
	defer r()

	deviceCtx1 := &snapstatetest.TrivialDeviceContext{DeviceModel: MakeModel(nil)}

	// not seeded
	deviceCtx, err := snapstate.DeviceCtxFromState(s.st, deviceCtx1)
	c.Assert(err, IsNil)
	c.Assert(deviceCtx, Equals, deviceCtx1)

	// seeded
	s.st.Set("seeded", true)
	deviceCtx, err = snapstate.DeviceCtxFromState(s.st, deviceCtx1)
	c.Assert(err, IsNil)
	c.Assert(deviceCtx, Equals, deviceCtx1)
}

func (s *deviceCtxSuite) TestDeviceCtxFromStateTooEarly(c *C) {
	s.st.Lock()
	defer s.st.Unlock()

	r := snapstatetest.MockDeviceModel(nil)
	defer r()

	expectedErr := &snapstate.ChangeConflictError{
		Message: "too early for operation, device model " +
			"not yet acknowledged",
		ChangeKind: "seed",
	}

	// not seeded, no model assertion
	_, err := snapstate.DeviceCtxFromState(s.st, nil)
	c.Assert(err, DeepEquals, expectedErr)

	// seeded, no model assertion
	s.st.Set("seeded", true)
	_, err = snapstate.DeviceCtxFromState(s.st, nil)
	c.Assert(err, DeepEquals, expectedErr)
}