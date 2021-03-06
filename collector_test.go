package main

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockDB struct{}

// Implement the SQLDB interface
func (mdb *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows := &sql.Rows{}
	return rows, nil
}

// Implement the SQLDB interface
func (mdb *MockDB) Close() error {
	return nil
}

//func Test_ExtractMetrics(t *testing.T) {
//	t.Parallel()
//	a := assert.New(t)
//
//	t.Run("extract metrics", func(t *testing.T) {
//		c := Collector{
//			db: &MockDB{},
//		}
//
//		metricGroup := &MetricGroup{
//			Labels:  []string{},
//			Metrics: make(map[string]*MetricDesc),
//		}
//
//		actualResult, err := c.extractMetrics(BuildRandomString(1), metricGroup, c.extractKeyValue)
//
//		a.Nil(err)
//		a.NotNil(actualResult)
//	})
//}

func Test_BuildMetricGroup(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("build metrics group", func(t *testing.T) {
		descriptor := MetricDescriptor{
			Prefix: BuildRandomString(3),
			Labels: []string{BuildRandomString(3)},
			MetricProps: []MetricProps{{
				Type: prometheus.GaugeValue,
				Name: BuildRandomString(3),
				Help: BuildRandomString(3),
			}},
		}

		actualResult := buildMetricGroup(descriptor)

		a.IsType(&MetricGroup{}, actualResult)
		a.Equal(descriptor.Labels, actualResult.Labels)
		a.Equal(1, len(actualResult.Metrics))
		a.Equal(descriptor.MetricProps[0].Type, actualResult.Metrics[descriptor.MetricProps[0].Name].Type)
	})
}

func Test_BuildGaugeOpts(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("build gauge opts", func(t *testing.T) {
		props := MetricProps{
			Name: BuildRandomString(3),
			Help: BuildRandomString(3),
		}

		actualResult := buildGaugeOpts(props)

		a.IsType(prometheus.GaugeOpts{}, actualResult)
		a.Equal(props.Name, actualResult.Name)
		a.Equal(props.Help, actualResult.Help)
		a.Equal("pgbouncer", actualResult.Namespace)
	})
}

func Test_BuildCounterOpts(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("build counter opts", func(t *testing.T) {
		props := MetricProps{
			Name: BuildRandomString(3),
			Help: BuildRandomString(3),
		}

		actualResult := buildCounterOpts(props)

		a.IsType(prometheus.CounterOpts{}, actualResult)
		a.Equal(props.Name, actualResult.Name)
		a.Equal(props.Help, actualResult.Help)
		a.Equal("pgbouncer", actualResult.Namespace)
	})
}
