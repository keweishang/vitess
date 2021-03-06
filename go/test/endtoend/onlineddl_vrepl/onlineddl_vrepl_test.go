/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package onlineddl

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"vitess.io/vitess/go/mysql"
	"vitess.io/vitess/go/sqltypes"
	"vitess.io/vitess/go/vt/schema"
	throttlebase "vitess.io/vitess/go/vt/vttablet/tabletserver/throttle/base"

	"vitess.io/vitess/go/test/endtoend/cluster"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	clusterInstance  *cluster.LocalProcessCluster
	vtParams         mysql.ConnParams
	httpClient       = throttlebase.SetupHTTPClient(time.Second)
	throttlerAppName = "vreplication"

	hostname              = "localhost"
	keyspaceName          = "ks"
	cell                  = "zone1"
	schemaChangeDirectory = ""
	totalTableCount       = 4
	createTable           = `
		CREATE TABLE %s (
			id bigint(20) NOT NULL,
			test_val bigint unsigned NOT NULL DEFAULT 0,
			msg varchar(64),
			PRIMARY KEY (id)
		) ENGINE=InnoDB;`
	// To verify non online-DDL behavior
	alterTableNormalStatement = `
		ALTER TABLE %s
			ADD COLUMN non_online int UNSIGNED NOT NULL DEFAULT 0`
	// A trivial statement which must succeed and does not change the schema
	alterTableTrivialStatement = `
		ALTER TABLE %s
			ENGINE=InnoDB`
	// The following statement is valid
	alterTableSuccessfulStatement = `
		ALTER TABLE %s
			MODIFY id bigint UNSIGNED NOT NULL,
			ADD COLUMN vrepl_col int NOT NULL DEFAULT 0,
			ADD INDEX idx_msg(msg)`
	// The following statement will fail because vreplication requires shared PRIMARY KEY columns
	alterTableFailedStatement = `
		ALTER TABLE %s
			DROP PRIMARY KEY,
			DROP COLUMN vrepl_col`
	// We will run this query while throttling vreplication
	alterTableThrottlingStatement = `
		ALTER TABLE %s
			DROP COLUMN vrepl_col`
	onlineDDLCreateTableStatement = `
		CREATE TABLE %s (
			id bigint NOT NULL,
			test_val bigint unsigned NOT NULL DEFAULT 0,
			online_ddl_create_col INT NOT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB;`
	onlineDDLDropTableStatement = `
		DROP TABLE %s`
	onlineDDLDropTableIfExistsStatement = `
		DROP TABLE IF EXISTS %s`
	insertRowStatement = `
		INSERT INTO %s (id, test_val) VALUES (%d, 1)
	`
	selectCountRowsStatement = `
		SELECT COUNT(*) AS c FROM %s
	`
	countInserts int64
	insertMutex  sync.Mutex
)

func fullWordUUIDRegexp(uuid, searchWord string) *regexp.Regexp {
	return regexp.MustCompile(uuid + `.*?\b` + searchWord + `\b`)
}
func fullWordRegexp(searchWord string) *regexp.Regexp {
	return regexp.MustCompile(`.*?\b` + searchWord + `\b`)
}

func TestMain(m *testing.M) {
	defer cluster.PanicHandler(nil)
	flag.Parse()

	exitcode, err := func() (int, error) {
		clusterInstance = cluster.NewCluster(cell, hostname)
		schemaChangeDirectory = path.Join("/tmp", fmt.Sprintf("schema_change_dir_%d", clusterInstance.GetAndReserveTabletUID()))
		defer os.RemoveAll(schemaChangeDirectory)
		defer clusterInstance.Teardown()

		if _, err := os.Stat(schemaChangeDirectory); os.IsNotExist(err) {
			_ = os.Mkdir(schemaChangeDirectory, 0700)
		}

		clusterInstance.VtctldExtraArgs = []string{
			"-schema_change_dir", schemaChangeDirectory,
			"-schema_change_controller", "local",
			"-schema_change_check_interval", "1"}

		clusterInstance.VtTabletExtraArgs = []string{
			"-enable-lag-throttler",
			"-throttle_threshold", "1s",
			"-heartbeat_enable",
			"-heartbeat_interval", "250ms",
			"-migration_check_interval", "5s",
		}
		clusterInstance.VtGateExtraArgs = []string{
			"-ddl_strategy", "online",
		}

		if err := clusterInstance.StartTopo(); err != nil {
			return 1, err
		}

		// Start keyspace
		keyspace := &cluster.Keyspace{
			Name: keyspaceName,
		}

		if err := clusterInstance.StartUnshardedKeyspace(*keyspace, 2, true); err != nil {
			return 1, err
		}
		if err := clusterInstance.StartKeyspace(*keyspace, []string{"1"}, 1, false); err != nil {
			return 1, err
		}

		vtgateInstance := clusterInstance.NewVtgateInstance()
		// set the gateway we want to use
		vtgateInstance.GatewayImplementation = "tabletgateway"
		// Start vtgate
		if err := vtgateInstance.Setup(); err != nil {
			return 1, err
		}
		// ensure it is torn down during cluster TearDown
		clusterInstance.VtgateProcess = *vtgateInstance
		vtParams = mysql.ConnParams{
			Host: clusterInstance.Hostname,
			Port: clusterInstance.VtgateMySQLPort,
		}

		return m.Run(), nil
	}()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	} else {
		os.Exit(exitcode)
	}

}

func throttleResponse(tablet *cluster.Vttablet, path string) (resp *http.Response, respBody string, err error) {
	apiURL := fmt.Sprintf("http://%s:%d/%s", tablet.VttabletProcess.TabletHostname, tablet.HTTPPort, path)
	resp, err = httpClient.Get(apiURL)
	if err != nil {
		return resp, respBody, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	respBody = string(b)
	return resp, respBody, err
}

func throttleApp(tablet *cluster.Vttablet, app string) (*http.Response, string, error) {
	return throttleResponse(tablet, fmt.Sprintf("throttler/throttle-app?app=%s&duration=1h", app))
}

func unthrottleApp(tablet *cluster.Vttablet, app string) (*http.Response, string, error) {
	return throttleResponse(tablet, fmt.Sprintf("throttler/unthrottle-app?app=%s", app))
}

func TestSchemaChange(t *testing.T) {
	defer cluster.PanicHandler(t)
	assert.Equal(t, 2, len(clusterInstance.Keyspaces[0].Shards))
	testWithInitialSchema(t)
	t.Run("alter non_online", func(t *testing.T) {
		_ = testOnlineDDLStatement(t, alterTableNormalStatement, string(schema.DDLStrategyDirect), "vtctl", "non_online")
		insertRows(t, 2)
		testRows(t)
	})
	t.Run("successful online alter, vtgate", func(t *testing.T) {
		insertRows(t, 2)
		uuid := testOnlineDDLStatement(t, alterTableSuccessfulStatement, "online", "vtgate", "vrepl_col")
		checkRecentMigrations(t, uuid, schema.OnlineDDLStatusComplete)
		testRows(t)
		checkCancelMigration(t, uuid, false)
		checkRetryMigration(t, uuid, false)
	})
	t.Run("successful online alter, vtctl", func(t *testing.T) {
		insertRows(t, 2)
		uuid := testOnlineDDLStatement(t, alterTableTrivialStatement, "online", "vtctl", "vrepl_col")
		checkRecentMigrations(t, uuid, schema.OnlineDDLStatusComplete)
		testRows(t)
		checkCancelMigration(t, uuid, false)
		checkRetryMigration(t, uuid, false)
	})
	t.Run("throttled migration", func(t *testing.T) {
		insertRows(t, 2)
		for i := range clusterInstance.Keyspaces[0].Shards {
			throttleApp(clusterInstance.Keyspaces[0].Shards[i].Vttablets[0], throttlerAppName)
			defer unthrottleApp(clusterInstance.Keyspaces[0].Shards[i].Vttablets[0], throttlerAppName)
		}
		uuid := testOnlineDDLStatement(t, alterTableThrottlingStatement, "online", "vtgate", "vrepl_col")
		checkRecentMigrations(t, uuid, schema.OnlineDDLStatusRunning)
		testRows(t)
		checkCancelMigration(t, uuid, true)
		time.Sleep(2 * time.Second)
		checkRecentMigrations(t, uuid, schema.OnlineDDLStatusFailed)
	})
	t.Run("failed migration", func(t *testing.T) {
		insertRows(t, 2)
		uuid := testOnlineDDLStatement(t, alterTableFailedStatement, "online", "vtgate", "vrepl_col")
		checkRecentMigrations(t, uuid, schema.OnlineDDLStatusFailed)
		testRows(t)
		checkCancelMigration(t, uuid, false)
		checkRetryMigration(t, uuid, true)
		// migration will fail again
	})
	t.Run("cancel all migrations: nothing to cancel", func(t *testing.T) {
		// no migrations pending at this time
		time.Sleep(10 * time.Second)
		checkCancelAllMigrations(t, 0)
	})
	t.Run("cancel all migrations: some migrations to cancel", func(t *testing.T) {
		for i := range clusterInstance.Keyspaces[0].Shards {
			throttleApp(clusterInstance.Keyspaces[0].Shards[i].Vttablets[0], throttlerAppName)
			defer unthrottleApp(clusterInstance.Keyspaces[0].Shards[i].Vttablets[0], throttlerAppName)
		}
		// spawn n migrations; cancel them via cancel-all
		var wg sync.WaitGroup
		count := 4
		for i := 0; i < count; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = testOnlineDDLStatement(t, alterTableThrottlingStatement, "online", "vtgate", "vrepl_col")
			}()
		}
		wg.Wait()
		checkCancelAllMigrations(t, count)
	})
	t.Run("Online DROP, vtctl", func(t *testing.T) {
		uuid := testOnlineDDLStatement(t, onlineDDLDropTableStatement, "online", "vtctl", "")
		checkRecentMigrations(t, uuid, schema.OnlineDDLStatusComplete)
		checkCancelMigration(t, uuid, false)
		checkRetryMigration(t, uuid, false)
	})
	t.Run("Online CREATE, vtctl", func(t *testing.T) {
		uuid := testOnlineDDLStatement(t, onlineDDLCreateTableStatement, "online", "vtctl", "online_ddl_create_col")
		checkRecentMigrations(t, uuid, schema.OnlineDDLStatusComplete)
		checkCancelMigration(t, uuid, false)
		checkRetryMigration(t, uuid, false)
	})
	t.Run("Online DROP TABLE IF EXISTS, vtgate", func(t *testing.T) {
		uuid := testOnlineDDLStatement(t, onlineDDLDropTableIfExistsStatement, "online", "vtgate", "")
		checkRecentMigrations(t, uuid, schema.OnlineDDLStatusComplete)
		checkCancelMigration(t, uuid, false)
		checkRetryMigration(t, uuid, false)
		// this table existed
		checkTables(t, schema.OnlineDDLToGCUUID(uuid), 1)
	})
	t.Run("Online DROP TABLE IF EXISTS for nonexistent table, vtgate", func(t *testing.T) {
		uuid := testOnlineDDLStatement(t, onlineDDLDropTableIfExistsStatement, "online", "vtgate", "")
		checkRecentMigrations(t, uuid, schema.OnlineDDLStatusComplete)
		checkCancelMigration(t, uuid, false)
		checkRetryMigration(t, uuid, false)
		// this table did not exist
		checkTables(t, schema.OnlineDDLToGCUUID(uuid), 0)
	})
	t.Run("Online DROP TABLE for nonexistent table, expect error, vtgate", func(t *testing.T) {
		uuid := testOnlineDDLStatement(t, onlineDDLDropTableStatement, "online", "vtgate", "")
		checkRecentMigrations(t, uuid, schema.OnlineDDLStatusFailed)
		checkCancelMigration(t, uuid, false)
		checkRetryMigration(t, uuid, true)
	})
}

func insertRow(t *testing.T) {
	insertMutex.Lock()
	defer insertMutex.Unlock()

	tableName := fmt.Sprintf("vt_onlineddl_test_%02d", 3)
	sqlQuery := fmt.Sprintf(insertRowStatement, tableName, countInserts)
	r := vtgateExecQuery(t, sqlQuery, "")
	require.NotNil(t, r)
	countInserts++
}

func insertRows(t *testing.T, count int) {
	for i := 0; i < count; i++ {
		insertRow(t)
	}
}

func testRows(t *testing.T) {
	insertMutex.Lock()
	defer insertMutex.Unlock()

	tableName := fmt.Sprintf("vt_onlineddl_test_%02d", 3)
	sqlQuery := fmt.Sprintf(selectCountRowsStatement, tableName)
	r := vtgateExecQuery(t, sqlQuery, "")
	require.NotNil(t, r)
	row := r.Named().Row()
	require.NotNil(t, row)
	require.Equal(t, countInserts, row.AsInt64("c", 0))
}

func testWithInitialSchema(t *testing.T) {
	// Create 4 tables
	var sqlQuery = "" //nolint
	for i := 0; i < totalTableCount; i++ {
		sqlQuery = fmt.Sprintf(createTable, fmt.Sprintf("vt_onlineddl_test_%02d", i))
		err := clusterInstance.VtctlclientProcess.ApplySchema(keyspaceName, sqlQuery)
		require.Nil(t, err)
	}

	// Check if 4 tables are created
	checkTables(t, "", totalTableCount)
}

// testOnlineDDLStatement runs an online DDL, ALTER statement
func testOnlineDDLStatement(t *testing.T, alterStatement string, ddlStrategy string, executeStrategy string, expectColumn string) (uuid string) {
	tableName := fmt.Sprintf("vt_onlineddl_test_%02d", 3)
	sqlQuery := fmt.Sprintf(alterStatement, tableName)
	if executeStrategy == "vtgate" {
		row := vtgateExec(t, ddlStrategy, sqlQuery, "").Named().Row()
		if row != nil {
			uuid = row.AsString("uuid", "")
		}
	} else {
		var err error
		uuid, err = clusterInstance.VtctlclientProcess.ApplySchemaWithOutput(keyspaceName, sqlQuery, ddlStrategy)
		assert.NoError(t, err)
	}
	uuid = strings.TrimSpace(uuid)
	fmt.Println("# Generated UUID (for debug purposes):")
	fmt.Printf("<%s>\n", uuid)

	strategy, _, err := schema.ParseDDLStrategy(ddlStrategy)
	assert.NoError(t, err)

	if !strategy.IsDirect() {
		time.Sleep(time.Second * 20)
	}

	if expectColumn != "" {
		checkMigratedTable(t, tableName, expectColumn)
	}
	return uuid
}

// checkTables checks the number of tables in the first two shards.
func checkTables(t *testing.T, showTableName string, expectCount int) {
	for i := range clusterInstance.Keyspaces[0].Shards {
		checkTablesCount(t, clusterInstance.Keyspaces[0].Shards[i].Vttablets[0], showTableName, expectCount)
	}
}

// checkTablesCount checks the number of tables in the given tablet
func checkTablesCount(t *testing.T, tablet *cluster.Vttablet, showTableName string, expectCount int) {
	query := fmt.Sprintf(`show tables like '%%%s%%';`, showTableName)
	queryResult, err := tablet.VttabletProcess.QueryTablet(query, keyspaceName, true)
	require.Nil(t, err)
	assert.Equal(t, expectCount, len(queryResult.Rows))
}

// checkRecentMigrations checks 'OnlineDDL <keyspace> show recent' output. Example to such output:
// +------------------+-------+--------------+----------------------+--------------------------------------+----------+---------------------+---------------------+------------------+
// |      Tablet      | shard | mysql_schema |     mysql_table      |            migration_uuid            | strategy |  started_timestamp  | completed_timestamp | migration_status |
// +------------------+-------+--------------+----------------------+--------------------------------------+----------+---------------------+---------------------+------------------+
// | zone1-0000003880 |     0 | vt_ks        | vt_onlineddl_test_03 | a0638f6b_ec7b_11ea_9bf8_000d3a9b8a9a | online   | 2020-09-01 17:50:40 | 2020-09-01 17:50:41 | complete         |
// | zone1-0000003884 |     1 | vt_ks        | vt_onlineddl_test_03 | a0638f6b_ec7b_11ea_9bf8_000d3a9b8a9a | online   | 2020-09-01 17:50:40 | 2020-09-01 17:50:41 | complete         |
// +------------------+-------+--------------+----------------------+--------------------------------------+----------+---------------------+---------------------+------------------+

func checkRecentMigrations(t *testing.T, uuid string, expectStatus schema.OnlineDDLStatus) {
	result, err := clusterInstance.VtctlclientProcess.OnlineDDLShowRecent(keyspaceName)
	assert.NoError(t, err)
	fmt.Println("# 'vtctlclient OnlineDDL show recent' output (for debug purposes):")
	fmt.Println(result)
	assert.Equal(t, len(clusterInstance.Keyspaces[0].Shards), strings.Count(result, uuid))
	// We ensure "full word" regexp becuase some column names may conflict
	expectStatusRegexp := fullWordUUIDRegexp(uuid, string(expectStatus))
	m := expectStatusRegexp.FindAllString(result, -1)
	assert.Equal(t, len(clusterInstance.Keyspaces[0].Shards), len(m))
}

// checkCancelMigration attempts to cancel a migration, and expects rejection
func checkCancelMigration(t *testing.T, uuid string, expectCancelPossible bool) {
	result, err := clusterInstance.VtctlclientProcess.OnlineDDLCancelMigration(keyspaceName, uuid)
	fmt.Println("# 'vtctlclient OnlineDDL cancel <uuid>' output (for debug purposes):")
	fmt.Println(result)
	assert.NoError(t, err)

	var r *regexp.Regexp
	if expectCancelPossible {
		r = fullWordRegexp("1")
	} else {
		r = fullWordRegexp("0")
	}
	m := r.FindAllString(result, -1)
	assert.Equal(t, len(clusterInstance.Keyspaces[0].Shards), len(m))
}

// checkCancelAllMigrations all pending migrations
func checkCancelAllMigrations(t *testing.T, expectCount int) {
	result, err := clusterInstance.VtctlclientProcess.OnlineDDLCancelAllMigrations(keyspaceName)
	fmt.Println("# 'vtctlclient OnlineDDL cancel-all' output (for debug purposes):")
	fmt.Println(result)
	assert.NoError(t, err)

	r := fullWordRegexp(fmt.Sprintf("%d", expectCount))
	m := r.FindAllString(result, -1)
	assert.Equal(t, len(clusterInstance.Keyspaces[0].Shards), len(m))
}

// checkRetryMigration attempts to retry a migration, and expects rejection
func checkRetryMigration(t *testing.T, uuid string, expectRetryPossible bool) {
	result, err := clusterInstance.VtctlclientProcess.OnlineDDLRetryMigration(keyspaceName, uuid)
	fmt.Println("# 'vtctlclient OnlineDDL retry <uuid>' output (for debug purposes):")
	fmt.Println(result)
	assert.NoError(t, err)

	var r *regexp.Regexp
	if expectRetryPossible {
		r = fullWordRegexp("1")
	} else {
		r = fullWordRegexp("0")
	}
	m := r.FindAllString(result, -1)
	assert.Equal(t, len(clusterInstance.Keyspaces[0].Shards), len(m))
}

// checkMigratedTables checks the CREATE STATEMENT of a table after migration
func checkMigratedTable(t *testing.T, tableName, expectColumn string) {
	for i := range clusterInstance.Keyspaces[0].Shards {
		createStatement := getCreateTableStatement(t, clusterInstance.Keyspaces[0].Shards[i].Vttablets[0], tableName)
		assert.Contains(t, createStatement, expectColumn)
	}
}

// getCreateTableStatement returns the CREATE TABLE statement for a given table
func getCreateTableStatement(t *testing.T, tablet *cluster.Vttablet, tableName string) (statement string) {
	queryResult, err := tablet.VttabletProcess.QueryTablet(fmt.Sprintf("show create table %s;", tableName), keyspaceName, true)
	require.Nil(t, err)

	assert.Equal(t, len(queryResult.Rows), 1)
	assert.Equal(t, len(queryResult.Rows[0]), 2) // table name, create statement
	statement = queryResult.Rows[0][1].ToString()
	return statement
}

func vtgateExecQuery(t *testing.T, query string, expectError string) *sqltypes.Result {
	t.Helper()

	ctx := context.Background()
	conn, err := mysql.Connect(ctx, &vtParams)
	require.Nil(t, err)
	defer conn.Close()

	qr, err := conn.ExecuteFetch(query, 1000, true)
	if expectError == "" {
		require.NoError(t, err)
	} else {
		require.Error(t, err, "error should not be nil")
		assert.Contains(t, err.Error(), expectError, "Unexpected error")
	}
	return qr
}

func vtgateExec(t *testing.T, ddlStrategy string, query string, expectError string) *sqltypes.Result {
	t.Helper()

	ctx := context.Background()
	conn, err := mysql.Connect(ctx, &vtParams)
	require.Nil(t, err)
	defer conn.Close()

	setSession := fmt.Sprintf("set @@ddl_strategy='%s'", ddlStrategy)
	_, err = conn.ExecuteFetch(setSession, 1000, true)
	assert.NoError(t, err)

	qr, err := conn.ExecuteFetch(query, 1000, true)
	if expectError == "" {
		require.NoError(t, err)
	} else {
		require.Error(t, err, "error should not be nil")
		assert.Contains(t, err.Error(), expectError, "Unexpected error")
	}
	return qr
}
