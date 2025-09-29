package logger

import (
	"time"

	"go.uber.org/zap"
)

// Field is a type alias for zap.Field, which is used to represent structured logging fields.
type Field = zap.Field

func String(key string, val string) Field          { return zap.String(key, val) }
func Bool(key string, val bool) Field              { return zap.Bool(key, val) }
func Int(key string, val int) Field                { return zap.Int(key, val) }
func Int64(key string, val int64) Field            { return zap.Int64(key, val) }
func Uint(key string, val uint) Field              { return zap.Uint(key, val) }
func Uint64(key string, val uint64) Field          { return zap.Uint64(key, val) }
func Float64(key string, val float64) Field        { return zap.Float64(key, val) }
func Time(key string, val time.Time) Field         { return zap.Time(key, val) }
func Duration(key string, val time.Duration) Field { return zap.Duration(key, val) }
func Err(err error) Field                          { return zap.Error(err) }
func Any(key string, val any) Field                { return zap.Any(key, val) }
func Strings(key string, val []string) Field       { return zap.Strings(key, val) }
func Ints(key string, val []int) Field             { return zap.Ints(key, val) }
func Int64s(key string, val []int64) Field         { return zap.Int64s(key, val) }
func Float64s(key string, val []float64) Field     { return zap.Float64s(key, val) }
func Bools(key string, val []bool) Field           { return zap.Bools(key, val) }
func Reflect(key string, val any) Field            { return zap.Reflect(key, val) }
func Namespace(key string) Field                   { return zap.Namespace(key) }
func Skip() Field                                  { return zap.Skip() }
