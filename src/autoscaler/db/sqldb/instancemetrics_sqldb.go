package sqldb

import (
	"context"
	"database/sql"
	"time"

	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/db"
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/models"
	"code.cloudfoundry.org/lager/v3"
	"github.com/jmoiron/sqlx"
)

type InstanceMetricsSQLDB struct {
	logger   lager.Logger
	dbConfig db.DatabaseConfig
	sqldb    *sqlx.DB
}

func NewInstanceMetricsSQLDB(dbConfig db.DatabaseConfig, logger lager.Logger) (*InstanceMetricsSQLDB, error) {
	database, err := db.GetConnection(dbConfig.URL)
	if err != nil {
		return nil, err
	}

	sqldb, err := sqlx.Open(database.DriverName, database.DSN)
	if err != nil {
		logger.Error("failed-open-instancemetrics-db", err, lager.Data{"dbConfig": dbConfig})
		return nil, err
	}

	err = sqldb.Ping()
	if err != nil {
		sqldb.Close()
		logger.Error("failed-ping-instancemetrics-db", err, lager.Data{"dbConfig": dbConfig})
		return nil, err
	}

	sqldb.SetConnMaxLifetime(dbConfig.ConnectionMaxLifetime)
	sqldb.SetMaxIdleConns(int(dbConfig.MaxIdleConnections))
	sqldb.SetMaxOpenConns(int(dbConfig.MaxOpenConnections))
	sqldb.SetConnMaxIdleTime(dbConfig.ConnectionMaxIdleTime)

	return &InstanceMetricsSQLDB{
		sqldb:    sqldb,
		logger:   logger,
		dbConfig: dbConfig,
	}, nil
}

func (idb *InstanceMetricsSQLDB) Close() error {
	err := idb.sqldb.Close()
	if err != nil {
		idb.logger.Error("failed-close-instancemetrics-db", err, lager.Data{"dbConfig": idb.dbConfig})
		return err
	}
	return nil
}

func (idb *InstanceMetricsSQLDB) SaveMetric(metric *models.AppInstanceMetric) error {
	query := idb.sqldb.Rebind("INSERT INTO appinstancemetrics(appid, instanceindex, collectedat, name, unit, value, timestamp) values(?, ?, ?, ?, ?, ?, ?)")
	_, err := idb.sqldb.Exec(query, metric.AppId, metric.InstanceIndex, metric.CollectedAt, metric.Name, metric.Unit, metric.Value, metric.Timestamp)

	if err != nil {
		idb.logger.Error("failed-insert-instancemetric-into-appinstancemetrics-table", err, lager.Data{"query": query, "metric": metric})
	}
	return err
}

func (idb *InstanceMetricsSQLDB) SaveMetricsInBulk(metrics []*models.AppInstanceMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	txn, err := idb.sqldb.BeginTxx(ctx, nil)
	if err != nil {
		idb.logger.Error("failed-to-start-transaction", err)
		return err
	}

	sqlStr := "INSERT INTO appinstancemetrics(appid, instanceindex, collectedat, name, unit, value, timestamp) VALUES (:app_id, :instance_index, :collected_at, :name, :unit, :value, :timestamp)"

	_, err = txn.NamedExec(sqlStr, metrics)
	if err != nil {
		idb.logger.Error("failed-to-execute-statement", err)
		_ = txn.Rollback()
		return err
	}

	err = txn.Commit()
	if err != nil {
		idb.logger.Error("failed-to-commit-transaction", err)
		_ = txn.Rollback()
		return err
	}

	return nil
}

func (idb *InstanceMetricsSQLDB) RetrieveInstanceMetrics(appid string, instanceIndex int, name string, start int64, end int64, orderType db.OrderType) ([]*models.AppInstanceMetric, error) {
	var orderStr string
	if orderType == db.ASC {
		orderStr = db.ASCSTR
	} else {
		orderStr = db.DESCSTR
	}
	query := idb.sqldb.Rebind("SELECT instanceindex, collectedat, unit, value, timestamp FROM appinstancemetrics WHERE " +
		" appid = ? " +
		" AND name = ? " +
		" AND timestamp >= ?" +
		" AND timestamp <= ?" +
		" ORDER BY timestamp " + orderStr + ", instanceindex")

	queryByInstanceIndex := idb.sqldb.Rebind("SELECT instanceindex, collectedat, unit, value, timestamp FROM appinstancemetrics WHERE " +
		" appid = ? " +
		" AND instanceindex = ?" +
		" AND name = ? " +
		" AND timestamp >= ?" +
		" AND timestamp <= ?" +
		" ORDER BY timestamp " + orderStr)

	if end < 0 {
		end = time.Now().UnixNano()
	}
	var rows *sql.Rows
	var err error
	if instanceIndex >= 0 {
		rows, err = idb.sqldb.Query(queryByInstanceIndex, appid, instanceIndex, name, start, end)
		if err != nil {
			idb.logger.Error("failed-retrieve-instancemetrics-from-appinstancemetrics-table", err,
				lager.Data{"query": query, "appid": appid, "instanceindex": instanceIndex, "metricName": name, "start": start, "end": end, "orderType": orderType})
			return nil, err
		}
	} else {
		rows, err = idb.sqldb.Query(query, appid, name, start, end)
		if err != nil {
			idb.logger.Error("failed-retrieve-instancemetrics-from-appinstancemetrics-table", err,
				lager.Data{"query": query, "appid": appid, "metricName": name, "start": start, "end": end, "orderType": orderType})
			return nil, err
		}
	}

	defer func() { _ = rows.Close() }()

	mtrcs := []*models.AppInstanceMetric{}
	var index uint32
	var collectedAt, timestamp int64
	var unit, value string

	for rows.Next() {
		if err := rows.Scan(&index, &collectedAt, &unit, &value, &timestamp); err != nil {
			idb.logger.Error("failed-scan-instancemetric-from-search-result", err)
			return nil, err
		}

		length := len(mtrcs)
		if (length > 0) && (timestamp == mtrcs[length-1].Timestamp) && (index == mtrcs[length-1].InstanceIndex) {
			continue
		}

		metric := models.AppInstanceMetric{
			AppId:         appid,
			InstanceIndex: index,
			CollectedAt:   collectedAt,
			Name:          name,
			Unit:          unit,
			Value:         value,
			Timestamp:     timestamp,
		}
		mtrcs = append(mtrcs, &metric)
	}
	return mtrcs, rows.Err()
}
func (idb *InstanceMetricsSQLDB) PruneInstanceMetrics(before int64) error {
	query := idb.sqldb.Rebind("DELETE FROM appinstancemetrics WHERE timestamp <= ?")
	_, err := idb.sqldb.Exec(query, before)
	if err != nil {
		idb.logger.Error("failed-prune-instancemetric-from-appinstancemetrics-table", err, lager.Data{"query": query, "before": before})
	}

	return err
}
func (idb *InstanceMetricsSQLDB) GetDBStatus() sql.DBStats {
	return idb.sqldb.Stats()
}
