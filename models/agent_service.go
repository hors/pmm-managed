// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package models

import (
	"time"

	"gopkg.in/reform.v1"
)

//go:generate reform

// AgentService implements many-to-many relationship between Agents and Services.
//reform:agent_services
type AgentService struct {
	AgentID   string    `reform:"agent_id"`
	ServiceID string    `reform:"service_id"`
	CreatedAt time.Time `reform:"created_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
//nolint:unparam
func (as *AgentService) BeforeInsert() error {
	now := time.Now().Truncate(time.Microsecond).UTC()
	as.CreatedAt = now
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
//nolint:unparam
func (as *AgentService) BeforeUpdate() error {
	panic("AgentService should not be updated")
}

// AfterFind implements reform.AfterFinder interface.
//nolint:unparam
func (as *AgentService) AfterFind() error {
	as.CreatedAt = as.CreatedAt.UTC()
	return nil
}

// check interfaces
var (
	_ reform.BeforeInserter = (*AgentService)(nil)
	_ reform.BeforeUpdater  = (*AgentService)(nil)
	_ reform.AfterFinder    = (*AgentService)(nil)
)
