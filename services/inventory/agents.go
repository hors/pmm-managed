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

package inventory

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	api "github.com/percona/pmm/api/inventory"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/agents"
)

// AgentsService works with inventory API Agents.
type AgentsService struct {
	q *reform.Querier
	r *agents.Registry
}

func NewAgentsService(q *reform.Querier, r *agents.Registry) *AgentsService {
	return &AgentsService{
		q: q,
		r: r,
	}
}

// makeAgent converts database row to Inventory API Agent.
func (as *AgentsService) makeAgent(ctx context.Context, row *models.AgentRow) (api.Agent, error) {
	switch row.AgentType {
	case models.PMMAgentType:
		return &api.PMMAgent{
			AgentId: row.AgentID,
			NodeId:  row.RunsOnNodeID,
		}, nil

	case models.NodeExporterType:
		return &api.NodeExporter{
			AgentId: row.AgentID,
			NodeId:  row.RunsOnNodeID,
		}, nil

	case models.MySQLdExporterType:
		services, err := models.ServicesForAgent(as.q, row.AgentID)
		if err != nil {
			return nil, err
		}
		if len(services) != 1 {
			return nil, errors.Errorf("expected exactly one Services, got %d", len(services))
		}

		return &api.MySQLdExporter{
			AgentId:      row.AgentID,
			RunsOnNodeId: row.RunsOnNodeID,
			ServiceId:    services[0].ServiceID,
			Username:     pointer.GetString(row.Username),
		}, nil

	default:
		panic(fmt.Errorf("unhandled AgentRow type %s", row.AgentType))
	}
}

func (as *AgentsService) get(ctx context.Context, id string) (*models.AgentRow, *models.AgentNode, error) {
	if id == "" {
		return nil, nil, status.Error(codes.InvalidArgument, "Empty Agent ID.")
	}

	row := &models.AgentRow{AgentID: id}
	switch err := as.q.Reload(row); err {
	case nil:
		return row, nil, nil
	case reform.ErrNoRows:
		return nil, nil, status.Errorf(codes.NotFound, "Agent with ID %q not found.", id)
	default:
		return nil, nil, errors.WithStack(err)
	}
}

func (as *AgentsService) checkUniqueID(ctx context.Context, id string) error {
	if id == "" {
		panic("empty Agent ID")
	}

	row := &models.AgentRow{AgentID: id}
	switch err := as.q.Reload(row); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Agent with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

// List selects all Agents in a stable order for a given service.
func (as *AgentsService) List(ctx context.Context, filters AgentFilters) ([]api.Agent, error) {
	agentRows, err := as.agentsByFilters(filters)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO That loop makes len(agentRows) SELECTs, that can be slow. Optimize when needed.
	res := make([]api.Agent, len(agentRows))
	for i, row := range agentRows {
		agent, err := as.makeAgent(ctx, row)
		if err != nil {
			return nil, err
		}
		res[i] = agent
	}
	return res, nil
}

// Get selects a single Agent by ID.
func (as *AgentsService) Get(ctx context.Context, id string) (api.Agent, error) {
	row, _, err := as.get(ctx, id)
	if err != nil {
		return nil, err
	}
	return as.makeAgent(ctx, row)
}

// AddPMMAgent inserts pmm-agent Agent with given parameters.
func (as *AgentsService) AddPMMAgent(ctx context.Context, nodeID string) (api.Agent, error) {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// TODO Check runs-on Node: it must be BM, VM, DC (i.e. not remote, AWS RDS, etc.)

	id := "/agent_id/" + uuid.New().String()
	if err := as.checkUniqueID(ctx, id); err != nil {
		return nil, err
	}

	ns := NewNodesService(as.q)
	if _, err := ns.get(ctx, nodeID); err != nil {
		return nil, err
	}

	row := &models.AgentRow{
		AgentID:      id,
		AgentType:    models.PMMAgentType,
		RunsOnNodeID: nodeID,
	}
	if err := as.q.Insert(row); err != nil {
		return nil, errors.WithStack(err)
	}

	err := as.q.Insert(&models.AgentNode{
		AgentID: row.AgentID,
		NodeID:  nodeID,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	agent, err := as.makeAgent(ctx, row)
	return agent, err
}

// AddNodeExporter inserts node_exporter Agent with given parameters.
func (as *AgentsService) AddNodeExporter(ctx context.Context, nodeID string, disabled bool) (api.Agent, error) {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// TODO Check runs-on Node: it must be BM, VM, DC (i.e. not remote, AWS RDS, etc.)

	id := "/agent_id/" + uuid.New().String()
	if err := as.checkUniqueID(ctx, id); err != nil {
		return nil, err
	}

	ns := NewNodesService(as.q)
	if _, err := ns.get(ctx, nodeID); err != nil {
		return nil, err
	}

	row := &models.AgentRow{
		AgentID:      id,
		AgentType:    models.NodeExporterType,
		RunsOnNodeID: nodeID,
	}
	if err := as.q.Insert(row); err != nil {
		return nil, errors.WithStack(err)
	}

	err := as.q.Insert(&models.AgentNode{
		AgentID: row.AgentID,
		NodeID:  nodeID,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// send new state to pmm-agents
	agents, err := as.agentsByFilters(AgentFilters{RunsOnNodeID: nodeID})
	if err != nil {
		return nil, err
	}
	for _, agent := range agents {
		if agent.AgentType != models.PMMAgentType {
			continue
		}
		as.r.SendSetStateRequest(ctx, agent.AgentID)
	}

	return as.makeAgent(ctx, row)
}

// AddMySQLdExporter inserts mysqld_exporter Agent with given parameters.
func (as *AgentsService) AddMySQLdExporter(ctx context.Context, nodeID string, disabled bool, serviceID string, username, password *string) (api.Agent, error) {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// TODO Check runs-on Node: it must be BM, VM, DC (i.e. not remote, AWS RDS, etc.)

	id := "/agent_id/" + uuid.New().String()
	if err := as.checkUniqueID(ctx, id); err != nil {
		return nil, err
	}

	ns := NewNodesService(as.q)
	if _, err := ns.get(ctx, nodeID); err != nil {
		return nil, err
	}

	ss := NewServicesService(as.q)
	if _, err := ss.get(ctx, serviceID); err != nil {
		return nil, err
	}

	row := &models.AgentRow{
		AgentID:      id,
		AgentType:    models.MySQLdExporterType,
		RunsOnNodeID: nodeID,
		Username:     username,
		Password:     password,
	}
	if err := as.q.Insert(row); err != nil {
		return nil, errors.WithStack(err)
	}

	err := as.q.Insert(&models.AgentNode{
		AgentID: row.AgentID,
		NodeID:  nodeID,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = as.q.Insert(&models.AgentService{
		AgentID:   row.AgentID,
		ServiceID: serviceID,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// send new state to pmm-agents
	agents, err := as.agentsByFilters(AgentFilters{RunsOnNodeID: nodeID})
	if err != nil {
		return nil, err
	}
	for _, agent := range agents {
		if agent.AgentType != models.PMMAgentType {
			continue
		}
		as.r.SendSetStateRequest(ctx, agent.AgentID)
	}

	return as.makeAgent(ctx, row)
}

/*
// SetDisabled enables or disables Agent by ID.
func (as *AgentsService) SetDisabled(ctx context.Context, id string, disabled bool) error {
	row, _, err := as.get(ctx, id)
	if err != nil {
		return err
	}

	row.Disabled = disabled
	err = as.q.Update(row)
	return errors.WithStack(err)
}
*/

// Remove deletes Agent by ID.
func (as *AgentsService) Remove(ctx context.Context, id string) error {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// ID is not 0.

	row, _, err := as.get(ctx, id)
	if err != nil {
		return err
	}

	if _, err = as.q.DeleteFrom(models.AgentServiceView, "WHERE agent_id = "+as.q.Placeholder(1), id); err != nil {
		return errors.WithStack(err)
	}
	if _, err = as.q.DeleteFrom(models.AgentNodeView, "WHERE agent_id = "+as.q.Placeholder(1), id); err != nil {
		return errors.WithStack(err)
	}

	if err = as.q.Delete(row); err != nil {
		return errors.WithStack(err)
	}

	as.r.SendSetStateRequest(ctx, id)

	if row.AgentType == models.PMMAgentType {
		as.r.Kick(ctx, id)
	}

	return nil
}

// AgentFilters represents filters for agents list.
type AgentFilters struct {
	// Return only Agents running on that Node.
	RunsOnNodeID string
	// Return only Agents that provide insights for that Node.
	NodeID string
	// Return only Agents that provide insights for that Service.
	ServiceID string
}

// TODO Move to models package to avoid import cycle between inventory and agents packages.
func (as *AgentsService) agentsByFilters(filters AgentFilters) ([]*models.AgentRow, error) {
	var tail string
	var args []interface{}
	switch {
	case filters.RunsOnNodeID != "":
		tail = fmt.Sprintf("WHERE runs_on_node_id = %s", as.q.Placeholder(1))
		args = []interface{}{filters.RunsOnNodeID}
	case filters.NodeID != "":
		agentNodes, err := as.q.SelectAllFrom(models.AgentNodeView, "WHERE node_id = ?", filters.NodeID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		for _, str := range agentNodes {
			args = append(args, str.(*models.AgentNode).AgentID)
		}
		if len(args) == 0 {
			return []*models.AgentRow{}, nil
		}
		tail = fmt.Sprintf("WHERE agent_id IN (%s)", strings.Join(as.q.Placeholders(1, len(args)), ", "))
	case filters.ServiceID != "":
		agentServices, err := as.q.SelectAllFrom(models.AgentServiceView, "WHERE service_id = ?", filters.ServiceID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		for _, str := range agentServices {
			args = append(args, str.(*models.AgentService).AgentID)
		}
		if len(args) == 0 {
			return []*models.AgentRow{}, nil
		}

		tail = fmt.Sprintf("WHERE agent_id IN (%s)", strings.Join(as.q.Placeholders(1, len(args)), ", "))
	}

	tail += " ORDER BY agent_id"

	structs, err := as.q.SelectAllFrom(models.AgentRowTable, tail, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	agentRows := make([]*models.AgentRow, len(structs))
	for i, str := range structs {
		agentRows[i] = str.(*models.AgentRow)
	}
	return agentRows, nil
}
