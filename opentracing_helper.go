package gormopentracing

import (
	"context"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/openzipkin/zipkin-go"
	"gorm.io/gorm"
)

const (
	_prefix      = "gorm.opentracing"
	_errorTagKey = "error"
)

var (
	// span.Tag keys
	_tableTagKey = keyWithPrefix("table")
	// span.Log keys
	//_errorLogKey        = keyWithPrefix("error")
	//_resultLogKey       = keyWithPrefix("result")
	_sqlLogKey          = keyWithPrefix("sql")
	_rowsAffectedLogKey = keyWithPrefix("rowsAffected")
)

func keyWithPrefix(key string) string {
	return _prefix + "." + key
}

var (
	opentracingSpanKey = "opentracing:span"
	json               = jsoniter.ConfigCompatibleWithStandardLibrary
)

func (p opentracingPlugin) injectBefore(db *gorm.DB, op operationName) {
	// make sure context could be used
	if db == nil {
		return
	}

	if db.Statement == nil || db.Statement.Context == nil {
		db.Logger.Error(context.TODO(), "could not inject sp from nil Statement.Context or nil Statement")
		return
	}

	sp, _ := p.opt.tracer.StartSpanFromContext(db.Statement.Context, op.String())
	db.InstanceSet(opentracingSpanKey, sp)
}

func (p opentracingPlugin) extractAfter(db *gorm.DB) {
	// make sure context could be used
	if db == nil {
		return
	}
	if db.Statement == nil || db.Statement.Context == nil {
		db.Logger.Error(context.TODO(), "could not extract sp from nil Statement.Context or nil Statement")
		return
	}

	// extract sp from db context
	//sp := opentracing.SpanFromContext(db.Statement.Context)
	v, ok := db.InstanceGet(opentracingSpanKey)
	if !ok || v == nil {
		return
	}

	sp, ok := v.(zipkin.Span)
	if !ok || sp == nil {
		return
	}
	defer sp.Finish()

	// tag and log fields we want.
	tag(sp, db)
}

// tag called after operation
func tag(sp zipkin.Span, db *gorm.DB) {
	if err := db.Error; err != nil {
		sp.Tag(_errorTagKey, "true")
	}

	sp.Tag(_tableTagKey, db.Statement.Table)
	sp.Tag(_rowsAffectedLogKey, strconv.FormatInt(db.Statement.RowsAffected, 10))
	//sp.Tag(_sqlLogKey, db.Statement.SQL.String())
}
