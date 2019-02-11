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

package agents

import (
	"testing"

	"github.com/AlekSi/pointer"
	api "github.com/percona/pmm/api/agent"
	"github.com/stretchr/testify/assert"

	"github.com/percona/pmm-managed/models"
)

func TestMySQLdExporterConfig(t *testing.T) {
	mysql := &models.ServiceRow{
		Address: pointer.ToString("1.2.3.4"),
		Port:    pointer.ToUint16(3306),
	}
	exporter := &models.AgentRow{
		Username: pointer.ToString("username"),
		Password: pointer.ToString("password"),
	}
	actual := mysqldExporterConfig(mysql, exporter)
	expected := &api.SetStateRequest_AgentProcess{
		Type:               api.Type_MYSQLD_EXPORTER,
		TemplateLeftDelim:  "{{",
		TemplateRightDelim: "}}",
		Args: []string{
			"-collect.binlog_size",
			"-collect.global_status",
			"-collect.global_variables",
			"-collect.info_schema.innodb_metrics",
			"-collect.info_schema.processlist",
			"-collect.info_schema.query_response_time",
			"-collect.info_schema.userstats",
			"-collect.perf_schema.eventswaits",
			"-collect.perf_schema.file_events",
			"-collect.slave_status",
			"-web.listen-address=:{{ .listen_port }}",
		},
		Env: []string{
			"DATA_SOURCE_NAME=username:password@tcp(1.2.3.4:3306)/?timeout=5s",
		},
	}
	assert.Equal(t, expected, actual)
}
