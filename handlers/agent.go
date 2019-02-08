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

package handlers

import (
	"context"

	api "github.com/percona/pmm/api/agent"

	"github.com/percona/pmm-managed/services/agents"
)

type AgentServer struct {
	Registry *agents.Registry
}

func (s *AgentServer) Register(context.Context, *api.RegisterRequest) (*api.RegisterResponse, error) {
	panic("not implemented yet")
}

func (s *AgentServer) Connect(stream api.Agent_ConnectServer) error {
	return s.Registry.Run(stream)
}

// check interfaces
var (
	_ api.AgentServer = (*AgentServer)(nil)
)
