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
	"fmt"
	"net"
	"sort"
	"strconv"

	"github.com/AlekSi/pointer"
	"github.com/go-sql-driver/mysql"
	api "github.com/percona/pmm/api/agent"

	"github.com/percona/pmm-managed/models"
)

func mysqldExporterConfig(service *models.ServiceRow, exporter *models.AgentRow) *api.SetStateRequest_AgentProcess {
	tp := templatePair(
		pointer.GetString(service.Address),
		pointer.GetString(service.UnixSocket),
		pointer.GetString(exporter.Username),
		pointer.GetString(exporter.Password),
		pointer.GetString(exporter.MetricsURL),
	)

	args := []string{
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
		"-web.listen-address=:" + tp.left + " .listen_port " + tp.right,
	}
	if pointer.GetString(exporter.MetricsURL) != "" {
		args = append(args, "-web.telemetry-path="+*exporter.MetricsURL)
	}

	// TODO Make it configurable. Play safe for now.
	// args = append(args, "-collect.auto_increment.columns")
	// args = append(args, "-collect.info_schema.tables")
	// args = append(args, "-collect.info_schema.tablestats")
	// args = append(args, "-collect.perf_schema.indexiowaits")
	// args = append(args, "-collect.perf_schema.tableiowaits")
	// args = append(args, "-collect.perf_schema.tablelocks")

	sort.Strings(args)

	// TODO TLSConfig: "true", https://jira.percona.com/browse/PMM-1727
	// TODO Other parameters?
	cfg := mysql.NewConfig()
	cfg.User = pointer.GetString(exporter.Username)
	cfg.Passwd = pointer.GetString(exporter.Password)
	if socket := pointer.GetString(service.UnixSocket); socket != "" {
		cfg.Net = "unix"
		cfg.Addr = socket
	} else {
		cfg.Net = "tcp"
		host := pointer.GetString(service.Address)
		port := pointer.GetUint16(service.Port)
		cfg.Addr = net.JoinHostPort(host, strconv.Itoa(int(port)))
	}
	cfg.Timeout = sqlDialTimeout
	dsn := cfg.FormatDSN()

	return &api.SetStateRequest_AgentProcess{
		Type:               api.Type_MYSQLD_EXPORTER,
		TemplateLeftDelim:  tp.left,
		TemplateRightDelim: tp.right,
		Args:               args,
		Env: []string{
			fmt.Sprintf("DATA_SOURCE_NAME=%s", dsn),
		},
	}
}
