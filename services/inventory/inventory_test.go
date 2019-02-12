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
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	api "github.com/percona/pmm/api/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mysql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/tests"
)

func TestNodes(t *testing.T) {
	sqlDB := tests.OpenTestDB(t)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	ctx := logger.Set(context.Background(), t.Name())

	setup := func(t *testing.T) (ns *NodesService, teardown func(t *testing.T)) {
		uuid.SetRand(new(tests.IDReader))

		db := reform.NewDB(sqlDB, mysql.Dialect, reform.NewPrintfLogger(t.Logf))
		tx, err := db.Begin()
		require.NoError(t, err)

		teardown = func(t *testing.T) {
			require.NoError(t, tx.Rollback())
		}
		ns = NewNodesService(tx.Querier)
		return
	}

	t.Run("Basic", func(t *testing.T) {
		ns, teardown := setup(t)
		defer teardown(t)

		actualNodes, err := ns.List(ctx)
		require.NoError(t, err)
		require.Len(t, actualNodes, 1) // PMMServerNodeType

		actualNode, err := ns.Add(ctx, models.GenericNodeType, "test-bm", nil, nil)
		require.NoError(t, err)
		expectedNode := &api.GenericNode{
			NodeId:   "/node_id/00000000-0000-4000-8000-000000000001",
			NodeName: "test-bm",
		}
		assert.Equal(t, expectedNode, actualNode)

		actualNode, err = ns.Get(ctx, "/node_id/00000000-0000-4000-8000-000000000001")
		require.NoError(t, err)
		assert.Equal(t, expectedNode, actualNode)

		actualNode, err = ns.Change(ctx, "/node_id/00000000-0000-4000-8000-000000000001", "test-bm-new")
		require.NoError(t, err)
		expectedNode = &api.GenericNode{
			NodeId:   "/node_id/00000000-0000-4000-8000-000000000001",
			NodeName: "test-bm-new",
		}
		assert.Equal(t, expectedNode, actualNode)

		actualNodes, err = ns.List(ctx)
		require.NoError(t, err)
		require.Len(t, actualNodes, 2)
		assert.Equal(t, expectedNode, actualNodes[0])

		err = ns.Remove(ctx, "/node_id/00000000-0000-4000-8000-000000000001")
		require.NoError(t, err)
		actualNode, err = ns.Get(ctx, "/node_id/00000000-0000-4000-8000-000000000001")
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Node with ID "/node_id/00000000-0000-4000-8000-000000000001" not found.`), err)
		assert.Nil(t, actualNode)
	})

	t.Run("GetEmptyID", func(t *testing.T) {
		ns, teardown := setup(t)
		defer teardown(t)

		actualNode, err := ns.Get(ctx, "")
		tests.AssertGRPCError(t, status.New(codes.InvalidArgument, `Empty Node ID.`), err)
		assert.Nil(t, actualNode)
	})

	t.Run("AddNameEmpty", func(t *testing.T) {
		ns, teardown := setup(t)
		defer teardown(t)

		_, err := ns.Add(ctx, models.GenericNodeType, "", nil, nil)
		tests.AssertGRPCError(t, status.New(codes.InvalidArgument, `Empty Node name.`), err)
	})

	t.Run("AddNameNotUnique", func(t *testing.T) {
		ns, teardown := setup(t)
		defer teardown(t)

		_, err := ns.Add(ctx, models.GenericNodeType, "test", pointer.ToString("test"), nil)
		require.NoError(t, err)

		_, err = ns.Add(ctx, models.RemoteNodeType, "test", nil, nil)
		tests.AssertGRPCError(t, status.New(codes.AlreadyExists, `Node with name "test" already exists.`), err)
	})

	t.Run("AddHostnameNotUnique", func(t *testing.T) {
		ns, teardown := setup(t)
		defer teardown(t)

		_, err := ns.Add(ctx, models.GenericNodeType, "test1", pointer.ToString("test"), nil)
		require.NoError(t, err)

		_, err = ns.Add(ctx, models.GenericNodeType, "test2", pointer.ToString("test"), nil)
		require.NoError(t, err)
	})

	t.Run("AddInstanceRegionNotUnique", func(t *testing.T) {
		ns, teardown := setup(t)
		defer teardown(t)

		_, err := ns.Add(ctx, models.RemoteAmazonRDSNodeType, "test1", pointer.ToString("test-instance"), pointer.ToString("test-region"))
		require.NoError(t, err)

		_, err = ns.Add(ctx, models.RemoteAmazonRDSNodeType, "test2", pointer.ToString("test-instance"), pointer.ToString("test-region"))
		expected := status.New(codes.AlreadyExists, `Node with instance "test-instance" and region "test-region" already exists.`)
		tests.AssertGRPCError(t, expected, err)
	})

	t.Run("ChangeNotFound", func(t *testing.T) {
		ns, teardown := setup(t)
		defer teardown(t)

		_, err := ns.Change(ctx, "no-such-id", "test-bm-new")
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Node with ID "no-such-id" not found.`), err)
	})

	t.Run("ChangeNameNotUnique", func(t *testing.T) {
		ns, teardown := setup(t)
		defer teardown(t)

		_, err := ns.Add(ctx, models.RemoteNodeType, "test-remote", nil, nil)
		require.NoError(t, err)

		rdsNode, err := ns.Add(ctx, models.RemoteAmazonRDSNodeType, "test-rds", nil, nil)
		require.NoError(t, err)

		_, err = ns.Change(ctx, rdsNode.ID(), "test-remote")
		tests.AssertGRPCError(t, status.New(codes.AlreadyExists, `Node with name "test-remote" already exists.`), err)
	})

	t.Run("RemoveNotFound", func(t *testing.T) {
		ns, teardown := setup(t)
		defer teardown(t)

		err := ns.Remove(ctx, "no-such-id")
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Node with ID "no-such-id" not found.`), err)
	})
}

func TestServices(t *testing.T) {
	sqlDB := tests.OpenTestDB(t)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	ctx := logger.Set(context.Background(), t.Name())

	setup := func(t *testing.T) (ss *ServicesService, teardown func(t *testing.T)) {
		uuid.SetRand(new(tests.IDReader))

		db := reform.NewDB(sqlDB, mysql.Dialect, reform.NewPrintfLogger(t.Logf))
		tx, err := db.Begin()
		require.NoError(t, err)

		teardown = func(t *testing.T) {
			require.NoError(t, tx.Rollback())
		}
		ss = NewServicesService(tx.Querier)
		return
	}

	t.Run("Basic", func(t *testing.T) {
		ss, teardown := setup(t)
		defer teardown(t)

		actualServices, err := ss.List(ctx)
		require.NoError(t, err)
		require.Len(t, actualServices, 0)

		actualService, err := ss.AddMySQL(ctx, "test-mysql", models.PMMServerNodeID, pointer.ToString("127.0.0.1"), pointer.ToUint16(3306), nil)
		require.NoError(t, err)
		expectedService := &api.MySQLService{
			ServiceId:   "/service_id/00000000-0000-4000-8000-000000000001",
			ServiceName: "test-mysql",
			NodeId:      models.PMMServerNodeID,
			Address:     "127.0.0.1",
			Port:        3306,
		}
		assert.Equal(t, expectedService, actualService)

		actualService, err = ss.Get(ctx, "/service_id/00000000-0000-4000-8000-000000000001")
		require.NoError(t, err)
		assert.Equal(t, expectedService, actualService)

		actualService, err = ss.Change(ctx, "/service_id/00000000-0000-4000-8000-000000000001", "test-mysql-new")
		require.NoError(t, err)
		expectedService = &api.MySQLService{
			ServiceId:   "/service_id/00000000-0000-4000-8000-000000000001",
			ServiceName: "test-mysql-new",
			NodeId:      models.PMMServerNodeID,
			Address:     "127.0.0.1",
			Port:        3306,
		}
		assert.Equal(t, expectedService, actualService)

		actualServices, err = ss.List(ctx)
		require.NoError(t, err)
		require.Len(t, actualServices, 1)
		assert.Equal(t, expectedService, actualServices[0])

		err = ss.Remove(ctx, "/service_id/00000000-0000-4000-8000-000000000001")
		require.NoError(t, err)
		actualService, err = ss.Get(ctx, "/service_id/00000000-0000-4000-8000-000000000001")
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Service with ID "/service_id/00000000-0000-4000-8000-000000000001" not found.`), err)
		assert.Nil(t, actualService)
	})

	t.Run("GetEmptyID", func(t *testing.T) {
		ss, teardown := setup(t)
		defer teardown(t)

		actualNode, err := ss.Get(ctx, "")
		tests.AssertGRPCError(t, status.New(codes.InvalidArgument, `Empty Service ID.`), err)
		assert.Nil(t, actualNode)
	})

	t.Run("AddNameNotUnique", func(t *testing.T) {
		ss, teardown := setup(t)
		defer teardown(t)

		_, err := ss.AddMySQL(ctx, "test-mysql", models.PMMServerNodeID, pointer.ToString("127.0.0.1"), pointer.ToUint16(3306), nil)
		require.NoError(t, err)

		_, err = ss.AddMySQL(ctx, "test-mysql", models.PMMServerNodeID, pointer.ToString("127.0.0.1"), pointer.ToUint16(3306), nil)
		tests.AssertGRPCError(t, status.New(codes.AlreadyExists, `Service with name "test-mysql" already exists.`), err)
	})

	t.Run("AddNodeNotFound", func(t *testing.T) {
		ss, teardown := setup(t)
		defer teardown(t)

		_, err := ss.AddMySQL(ctx, "test-mysql", "no-such-id", pointer.ToString("127.0.0.1"), pointer.ToUint16(3306), nil)
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Node with ID "no-such-id" not found.`), err)
	})

	t.Run("ChangeNotFound", func(t *testing.T) {
		ss, teardown := setup(t)
		defer teardown(t)

		_, err := ss.Change(ctx, "no-such-id", "test-mysql-new")
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Service with ID "no-such-id" not found.`), err)
	})

	t.Run("ChangeNameNotUnique", func(t *testing.T) {
		ss, teardown := setup(t)
		defer teardown(t)

		_, err := ss.AddMySQL(ctx, "test-mysql", models.PMMServerNodeID, pointer.ToString("127.0.0.1"), pointer.ToUint16(3306), nil)
		require.NoError(t, err)

		s, err := ss.AddMySQL(ctx, "test-mysql-2", models.PMMServerNodeID, pointer.ToString("127.0.0.2"), pointer.ToUint16(3306), nil)
		require.NoError(t, err)

		_, err = ss.Change(ctx, s.ID(), "test-mysql")
		tests.AssertGRPCError(t, status.New(codes.AlreadyExists, `Service with name "test-mysql" already exists.`), err)
	})

	t.Run("RemoveNotFound", func(t *testing.T) {
		ss, teardown := setup(t)
		defer teardown(t)

		err := ss.Remove(ctx, "no-such-id")
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Service with ID "no-such-id" not found.`), err)
	})
}

func TestAgents(t *testing.T) {
	sqlDB := tests.OpenTestDB(t)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	ctx := logger.Set(context.Background(), t.Name())

	setup := func(t *testing.T) (ns *NodesService, ss *ServicesService, as *AgentsService, teardown func(t *testing.T)) {
		uuid.SetRand(new(tests.IDReader))

		db := reform.NewDB(sqlDB, mysql.Dialect, reform.NewPrintfLogger(t.Logf))
		tx, err := db.Begin()
		require.NoError(t, err)

		r := new(mockRegistry)
		teardown = func(t *testing.T) {
			require.NoError(t, tx.Rollback())
			r.AssertExpectations(t)
		}
		ns = NewNodesService(tx.Querier)
		ss = NewServicesService(tx.Querier)
		as = NewAgentsService(tx.Querier, r)
		return
	}

	t.Run("Basic", func(t *testing.T) {
		ns, ss, as, teardown := setup(t)
		defer teardown(t)

		actualAgents, err := as.List(ctx, AgentFilters{})
		require.NoError(t, err)
		require.Len(t, actualAgents, 0)

		actualAgent, err := as.AddNodeExporter(ctx, models.PMMServerNodeID, true)
		require.NoError(t, err)
		expectedNodeExporterAgent := &api.NodeExporter{
			AgentId: "/agent_id/00000000-0000-4000-8000-000000000001",
			NodeId:  models.PMMServerNodeID,
		}
		assert.Equal(t, expectedNodeExporterAgent, actualAgent)

		actualAgent, err = as.Get(ctx, "/agent_id/00000000-0000-4000-8000-000000000001")
		require.NoError(t, err)
		assert.Equal(t, expectedNodeExporterAgent, actualAgent)

		s, err := ss.AddMySQL(ctx, "test-mysql", models.PMMServerNodeID, pointer.ToString("127.0.0.1"), pointer.ToUint16(3306), nil)
		require.NoError(t, err)

		n, err := ns.Add(ctx, models.GenericNodeType, "new node name", nil, nil)
		require.NoError(t, err)

		actualAgent, err = as.AddMySQLdExporter(ctx, n.ID(), false, s.ID(), pointer.ToString("username"), nil)
		require.NoError(t, err)
		expectedMySQLdExporterAgent := &api.MySQLdExporter{
			AgentId:      "/agent_id/00000000-0000-4000-8000-000000000004",
			RunsOnNodeId: n.ID(),
			ServiceId:    s.ID(),
			Username:     "username",
		}
		assert.Equal(t, expectedMySQLdExporterAgent, actualAgent)

		actualAgent, err = as.Get(ctx, "/agent_id/00000000-0000-4000-8000-000000000004")
		require.NoError(t, err)
		assert.Equal(t, expectedMySQLdExporterAgent, actualAgent)

		// err = as.SetDisabled(ctx, "/agent_id/00000000-0000-4000-8000-000000000001", true)
		// require.NoError(t, err)
		// expectedMySQLdExporterAgent.Disabled = true
		// actualAgent, err = as.Get(ctx, "/agent_id/00000000-0000-4000-8000-000000000001")
		// require.NoError(t, err)
		// assert.Equal(t, expectedMySQLdExporterAgent, actualAgent)

		actualAgents, err = as.List(ctx, AgentFilters{})
		require.NoError(t, err)
		require.Len(t, actualAgents, 2)
		assert.Equal(t, expectedNodeExporterAgent, actualAgents[0])
		assert.Equal(t, expectedMySQLdExporterAgent, actualAgents[1])

		actualAgents, err = as.List(ctx, AgentFilters{ServiceID: s.ID()})
		require.NoError(t, err)
		require.Len(t, actualAgents, 1)
		assert.Equal(t, expectedMySQLdExporterAgent, actualAgents[0])

		actualAgents, err = as.List(ctx, AgentFilters{RunsOnNodeID: n.ID()})
		require.NoError(t, err)
		require.Len(t, actualAgents, 1)
		assert.Equal(t, expectedMySQLdExporterAgent, actualAgents[0])

		actualAgents, err = as.List(ctx, AgentFilters{NodeID: models.PMMServerNodeID})
		require.NoError(t, err)
		require.Len(t, actualAgents, 1)
		assert.Equal(t, expectedNodeExporterAgent, actualAgents[0])

		err = as.Remove(ctx, "/agent_id/00000000-0000-4000-8000-000000000001")
		require.NoError(t, err)
		actualAgent, err = as.Get(ctx, "/agent_id/00000000-0000-4000-8000-000000000001")
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Agent with ID "/agent_id/00000000-0000-4000-8000-000000000001" not found.`), err)
		assert.Nil(t, actualAgent)

		err = as.Remove(ctx, "/agent_id/00000000-0000-4000-8000-000000000004")
		require.NoError(t, err)
		actualAgent, err = as.Get(ctx, "/agent_id/00000000-0000-4000-8000-000000000004")
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Agent with ID "/agent_id/00000000-0000-4000-8000-000000000004" not found.`), err)
		assert.Nil(t, actualAgent)

		actualAgents, err = as.List(ctx, AgentFilters{})
		require.NoError(t, err)
		require.Len(t, actualAgents, 0)
	})

	t.Run("GetEmptyID", func(t *testing.T) {
		_, _, as, teardown := setup(t)
		defer teardown(t)

		actualNode, err := as.Get(ctx, "")
		tests.AssertGRPCError(t, status.New(codes.InvalidArgument, `Empty Agent ID.`), err)
		assert.Nil(t, actualNode)
	})

	t.Run("AddPMMAgent", func(t *testing.T) {
		_, _, as, teardown := setup(t)
		defer teardown(t)

		actualAgent, err := as.AddPMMAgent(ctx, models.PMMServerNodeID)
		require.NoError(t, err)
		expectedPMMAgent := &api.PMMAgent{
			AgentId: "/agent_id/00000000-0000-4000-8000-000000000001",
			NodeId:  models.PMMServerNodeID,
		}
		assert.Equal(t, expectedPMMAgent, actualAgent)
	})

	t.Run("AddNodeNotFound", func(t *testing.T) {
		_, _, as, teardown := setup(t)
		defer teardown(t)

		_, err := as.AddNodeExporter(ctx, "no-such-id", true)
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Node with ID "no-such-id" not found.`), err)
	})

	t.Run("AddServiceNotFound", func(t *testing.T) {
		_, _, as, teardown := setup(t)
		defer teardown(t)

		_, err := as.AddMySQLdExporter(ctx, models.PMMServerNodeID, false, "no-such-id", pointer.ToString("username"), nil)
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Service with ID "no-such-id" not found.`), err)
	})

	// t.Run("DisableNotFound", func(t *testing.T) {
	// 	_, _, as, teardown := setup(t)
	// 	defer teardown(t)

	// 	err := as.SetDisabled(ctx, "no-such-id", true)
	// 	tests.AssertGRPCError(t, status.New(codes.NotFound, `Agent with ID "no-such-id" not found.`), err)
	// })

	t.Run("RemoveNotFound", func(t *testing.T) {
		_, _, as, teardown := setup(t)
		defer teardown(t)

		err := as.Remove(ctx, "no-such-id")
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Agent with ID "no-such-id" not found.`), err)
	})
}
