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

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	api "github.com/percona/pmm/api/inventory"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// ServicesService works with inventory API Services.
type ServicesService struct {
	q *reform.Querier
	// r *agents.Registry
}

func NewServicesService(q *reform.Querier) *ServicesService {
	return &ServicesService{
		q: q,
		// r: r,
	}
}

// makeService converts database row to Inventory API Service.
func makeService(row *models.ServiceRow) api.Service {
	switch row.ServiceType {
	case models.MySQLServiceType:
		return &api.MySQLService{
			ServiceId:   row.ServiceID,
			ServiceName: row.ServiceName,
			NodeId:      row.NodeID,
			Address:     pointer.GetString(row.Address),
			Port:        uint32(pointer.GetUint16(row.Port)),
			UnixSocket:  pointer.GetString(row.UnixSocket),
		}

	default:
		panic(fmt.Errorf("unhandled ServiceRow type %s", row.ServiceType))
	}
}

func (ss *ServicesService) get(ctx context.Context, id string) (*models.ServiceRow, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Service ID.")
	}

	row := &models.ServiceRow{ServiceID: id}
	switch err := ss.q.Reload(row); err {
	case nil:
		return row, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "Service with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

func (ss *ServicesService) checkUniqueID(ctx context.Context, id string) error {
	if id == "" {
		panic("empty Service ID")
	}

	row := &models.ServiceRow{ServiceID: id}
	switch err := ss.q.Reload(row); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Service with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

func (ss *ServicesService) checkUniqueName(ctx context.Context, name string) error {
	_, err := ss.q.FindOneFrom(models.ServiceRowTable, "service_name", name)
	switch err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Service with name %q already exists.", name)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

// List selects all Services in a stable order.
func (ss *ServicesService) List(ctx context.Context) ([]api.Service, error) {
	structs, err := ss.q.SelectAllFrom(models.ServiceRowTable, "ORDER BY service_id")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]api.Service, len(structs))
	for i, str := range structs {
		row := str.(*models.ServiceRow)
		res[i] = makeService(row)
	}
	return res, nil
}

// Get selects a single Service by ID.
func (ss *ServicesService) Get(ctx context.Context, id string) (api.Service, error) {
	row, err := ss.get(ctx, id)
	if err != nil {
		return nil, err
	}
	return makeService(row), nil
}

// AddMySQL inserts MySQL Service with given parameters.
func (ss *ServicesService) AddMySQL(ctx context.Context, name string, nodeID string, address *string, port *uint16, unixSocket *string) (api.Service, error) {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// Both address and socket can't be empty, etc.

	id := "/service_id/" + uuid.New().String()
	if err := ss.checkUniqueID(ctx, id); err != nil {
		return nil, err
	}
	if err := ss.checkUniqueName(ctx, name); err != nil {
		return nil, err
	}

	ns := NewNodesService(ss.q)
	if _, err := ns.get(ctx, nodeID); err != nil {
		return nil, err
	}

	row := &models.ServiceRow{
		ServiceID:   id,
		ServiceType: models.MySQLServiceType,
		ServiceName: name,
		NodeID:      nodeID,
		Address:     address,
		Port:        port,
		UnixSocket:  unixSocket,
	}
	if err := ss.q.Insert(row); err != nil {
		return nil, errors.WithStack(err)
	}
	return makeService(row), nil
}

// Change updates Service by ID.
func (ss *ServicesService) Change(ctx context.Context, id string, name string) (api.Service, error) {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// ID is not 0, name is not empty and valid.

	if err := ss.checkUniqueName(ctx, name); err != nil {
		return nil, err
	}

	row, err := ss.get(ctx, id)
	if err != nil {
		return nil, err
	}

	row.ServiceName = name
	if err = ss.q.Update(row); err != nil {
		return nil, errors.WithStack(err)
	}
	return makeService(row), nil
}

// Remove deletes Service by ID.
func (ss *ServicesService) Remove(ctx context.Context, id string) error {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// ID is not 0.

	// TODO check absence of Agents

	err := ss.q.Delete(&models.ServiceRow{ServiceID: id})
	if err == reform.ErrNoRows {
		return status.Errorf(codes.NotFound, "Service with ID %q not found.", id)
	}
	return errors.WithStack(err)
}
