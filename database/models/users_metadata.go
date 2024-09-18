// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// UsersMetadatum is an object representing the database table.
type UsersMetadatum struct {
	PubKey          string      `boil:"pub_key" json:"pub_key" toml:"pub_key" yaml:"pub_key"`
	CreatedAt       time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt       time.Time   `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	DeletedAt       null.Time   `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`
	Content         null.String `boil:"content" json:"content,omitempty" toml:"content" yaml:"content,omitempty"`
	FollowListEvent null.String `boil:"follow_list_event" json:"follow_list_event,omitempty" toml:"follow_list_event" yaml:"follow_list_event,omitempty"`

	R *usersMetadatumR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L usersMetadatumL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UsersMetadatumColumns = struct {
	PubKey          string
	CreatedAt       string
	UpdatedAt       string
	DeletedAt       string
	Content         string
	FollowListEvent string
}{
	PubKey:          "pub_key",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
	Content:         "content",
	FollowListEvent: "follow_list_event",
}

var UsersMetadatumTableColumns = struct {
	PubKey          string
	CreatedAt       string
	UpdatedAt       string
	DeletedAt       string
	Content         string
	FollowListEvent string
}{
	PubKey:          "users_metadata.pub_key",
	CreatedAt:       "users_metadata.created_at",
	UpdatedAt:       "users_metadata.updated_at",
	DeletedAt:       "users_metadata.deleted_at",
	Content:         "users_metadata.content",
	FollowListEvent: "users_metadata.follow_list_event",
}

// Generated where

var UsersMetadatumWhere = struct {
	PubKey          whereHelperstring
	CreatedAt       whereHelpertime_Time
	UpdatedAt       whereHelpertime_Time
	DeletedAt       whereHelpernull_Time
	Content         whereHelpernull_String
	FollowListEvent whereHelpernull_String
}{
	PubKey:          whereHelperstring{field: "\"users_metadata\".\"pub_key\""},
	CreatedAt:       whereHelpertime_Time{field: "\"users_metadata\".\"created_at\""},
	UpdatedAt:       whereHelpertime_Time{field: "\"users_metadata\".\"updated_at\""},
	DeletedAt:       whereHelpernull_Time{field: "\"users_metadata\".\"deleted_at\""},
	Content:         whereHelpernull_String{field: "\"users_metadata\".\"content\""},
	FollowListEvent: whereHelpernull_String{field: "\"users_metadata\".\"follow_list_event\""},
}

// UsersMetadatumRels is where relationship names are stored.
var UsersMetadatumRels = struct {
}{}

// usersMetadatumR is where relationships are stored.
type usersMetadatumR struct {
}

// NewStruct creates a new relationship struct
func (*usersMetadatumR) NewStruct() *usersMetadatumR {
	return &usersMetadatumR{}
}

// usersMetadatumL is where Load methods for each relationship are stored.
type usersMetadatumL struct{}

var (
	usersMetadatumAllColumns            = []string{"pub_key", "created_at", "updated_at", "deleted_at", "content", "follow_list_event"}
	usersMetadatumColumnsWithoutDefault = []string{"pub_key"}
	usersMetadatumColumnsWithDefault    = []string{"created_at", "updated_at", "deleted_at", "content", "follow_list_event"}
	usersMetadatumPrimaryKeyColumns     = []string{"pub_key"}
	usersMetadatumGeneratedColumns      = []string{}
)

type (
	// UsersMetadatumSlice is an alias for a slice of pointers to UsersMetadatum.
	// This should almost always be used instead of []UsersMetadatum.
	UsersMetadatumSlice []*UsersMetadatum
	// UsersMetadatumHook is the signature for custom UsersMetadatum hook methods
	UsersMetadatumHook func(context.Context, boil.ContextExecutor, *UsersMetadatum) error

	usersMetadatumQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	usersMetadatumType                 = reflect.TypeOf(&UsersMetadatum{})
	usersMetadatumMapping              = queries.MakeStructMapping(usersMetadatumType)
	usersMetadatumPrimaryKeyMapping, _ = queries.BindMapping(usersMetadatumType, usersMetadatumMapping, usersMetadatumPrimaryKeyColumns)
	usersMetadatumInsertCacheMut       sync.RWMutex
	usersMetadatumInsertCache          = make(map[string]insertCache)
	usersMetadatumUpdateCacheMut       sync.RWMutex
	usersMetadatumUpdateCache          = make(map[string]updateCache)
	usersMetadatumUpsertCacheMut       sync.RWMutex
	usersMetadatumUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var usersMetadatumAfterSelectMu sync.Mutex
var usersMetadatumAfterSelectHooks []UsersMetadatumHook

var usersMetadatumBeforeInsertMu sync.Mutex
var usersMetadatumBeforeInsertHooks []UsersMetadatumHook
var usersMetadatumAfterInsertMu sync.Mutex
var usersMetadatumAfterInsertHooks []UsersMetadatumHook

var usersMetadatumBeforeUpdateMu sync.Mutex
var usersMetadatumBeforeUpdateHooks []UsersMetadatumHook
var usersMetadatumAfterUpdateMu sync.Mutex
var usersMetadatumAfterUpdateHooks []UsersMetadatumHook

var usersMetadatumBeforeDeleteMu sync.Mutex
var usersMetadatumBeforeDeleteHooks []UsersMetadatumHook
var usersMetadatumAfterDeleteMu sync.Mutex
var usersMetadatumAfterDeleteHooks []UsersMetadatumHook

var usersMetadatumBeforeUpsertMu sync.Mutex
var usersMetadatumBeforeUpsertHooks []UsersMetadatumHook
var usersMetadatumAfterUpsertMu sync.Mutex
var usersMetadatumAfterUpsertHooks []UsersMetadatumHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UsersMetadatum) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usersMetadatumAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UsersMetadatum) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usersMetadatumBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UsersMetadatum) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usersMetadatumAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UsersMetadatum) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usersMetadatumBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UsersMetadatum) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usersMetadatumAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UsersMetadatum) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usersMetadatumBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UsersMetadatum) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usersMetadatumAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UsersMetadatum) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usersMetadatumBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UsersMetadatum) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range usersMetadatumAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUsersMetadatumHook registers your hook function for all future operations.
func AddUsersMetadatumHook(hookPoint boil.HookPoint, usersMetadatumHook UsersMetadatumHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		usersMetadatumAfterSelectMu.Lock()
		usersMetadatumAfterSelectHooks = append(usersMetadatumAfterSelectHooks, usersMetadatumHook)
		usersMetadatumAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		usersMetadatumBeforeInsertMu.Lock()
		usersMetadatumBeforeInsertHooks = append(usersMetadatumBeforeInsertHooks, usersMetadatumHook)
		usersMetadatumBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		usersMetadatumAfterInsertMu.Lock()
		usersMetadatumAfterInsertHooks = append(usersMetadatumAfterInsertHooks, usersMetadatumHook)
		usersMetadatumAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		usersMetadatumBeforeUpdateMu.Lock()
		usersMetadatumBeforeUpdateHooks = append(usersMetadatumBeforeUpdateHooks, usersMetadatumHook)
		usersMetadatumBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		usersMetadatumAfterUpdateMu.Lock()
		usersMetadatumAfterUpdateHooks = append(usersMetadatumAfterUpdateHooks, usersMetadatumHook)
		usersMetadatumAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		usersMetadatumBeforeDeleteMu.Lock()
		usersMetadatumBeforeDeleteHooks = append(usersMetadatumBeforeDeleteHooks, usersMetadatumHook)
		usersMetadatumBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		usersMetadatumAfterDeleteMu.Lock()
		usersMetadatumAfterDeleteHooks = append(usersMetadatumAfterDeleteHooks, usersMetadatumHook)
		usersMetadatumAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		usersMetadatumBeforeUpsertMu.Lock()
		usersMetadatumBeforeUpsertHooks = append(usersMetadatumBeforeUpsertHooks, usersMetadatumHook)
		usersMetadatumBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		usersMetadatumAfterUpsertMu.Lock()
		usersMetadatumAfterUpsertHooks = append(usersMetadatumAfterUpsertHooks, usersMetadatumHook)
		usersMetadatumAfterUpsertMu.Unlock()
	}
}

// OneG returns a single usersMetadatum record from the query using the global executor.
func (q usersMetadatumQuery) OneG(ctx context.Context) (*UsersMetadatum, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single usersMetadatum record from the query.
func (q usersMetadatumQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UsersMetadatum, error) {
	o := &UsersMetadatum{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: failed to execute a one query for users_metadata")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all UsersMetadatum records from the query using the global executor.
func (q usersMetadatumQuery) AllG(ctx context.Context) (UsersMetadatumSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all UsersMetadatum records from the query.
func (q usersMetadatumQuery) All(ctx context.Context, exec boil.ContextExecutor) (UsersMetadatumSlice, error) {
	var o []*UsersMetadatum

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "database: failed to assign all query results to UsersMetadatum slice")
	}

	if len(usersMetadatumAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all UsersMetadatum records in the query using the global executor
func (q usersMetadatumQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all UsersMetadatum records in the query.
func (q usersMetadatumQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to count users_metadata rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q usersMetadatumQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q usersMetadatumQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "database: failed to check if users_metadata exists")
	}

	return count > 0, nil
}

// UsersMetadata retrieves all the records using an executor.
func UsersMetadata(mods ...qm.QueryMod) usersMetadatumQuery {
	mods = append(mods, qm.From("\"users_metadata\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"users_metadata\".*"})
	}

	return usersMetadatumQuery{q}
}

// FindUsersMetadatumG retrieves a single record by ID.
func FindUsersMetadatumG(ctx context.Context, pubKey string, selectCols ...string) (*UsersMetadatum, error) {
	return FindUsersMetadatum(ctx, boil.GetContextDB(), pubKey, selectCols...)
}

// FindUsersMetadatum retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUsersMetadatum(ctx context.Context, exec boil.ContextExecutor, pubKey string, selectCols ...string) (*UsersMetadatum, error) {
	usersMetadatumObj := &UsersMetadatum{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"users_metadata\" where \"pub_key\"=$1", sel,
	)

	q := queries.Raw(query, pubKey)

	err := q.Bind(ctx, exec, usersMetadatumObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: unable to select from users_metadata")
	}

	if err = usersMetadatumObj.doAfterSelectHooks(ctx, exec); err != nil {
		return usersMetadatumObj, err
	}

	return usersMetadatumObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *UsersMetadatum) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UsersMetadatum) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("database: no users_metadata provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(usersMetadatumColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	usersMetadatumInsertCacheMut.RLock()
	cache, cached := usersMetadatumInsertCache[key]
	usersMetadatumInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			usersMetadatumAllColumns,
			usersMetadatumColumnsWithDefault,
			usersMetadatumColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(usersMetadatumType, usersMetadatumMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(usersMetadatumType, usersMetadatumMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"users_metadata\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"users_metadata\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "database: unable to insert into users_metadata")
	}

	if !cached {
		usersMetadatumInsertCacheMut.Lock()
		usersMetadatumInsertCache[key] = cache
		usersMetadatumInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single UsersMetadatum record using the global executor.
// See Update for more documentation.
func (o *UsersMetadatum) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the UsersMetadatum.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UsersMetadatum) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	usersMetadatumUpdateCacheMut.RLock()
	cache, cached := usersMetadatumUpdateCache[key]
	usersMetadatumUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			usersMetadatumAllColumns,
			usersMetadatumPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("database: unable to update users_metadata, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"users_metadata\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, usersMetadatumPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(usersMetadatumType, usersMetadatumMapping, append(wl, usersMetadatumPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update users_metadata row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by update for users_metadata")
	}

	if !cached {
		usersMetadatumUpdateCacheMut.Lock()
		usersMetadatumUpdateCache[key] = cache
		usersMetadatumUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q usersMetadatumQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q usersMetadatumQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all for users_metadata")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected for users_metadata")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o UsersMetadatumSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UsersMetadatumSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("database: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usersMetadatumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"users_metadata\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, usersMetadatumPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all in usersMetadatum slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected all in update all usersMetadatum")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *UsersMetadatum) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns, opts...)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UsersMetadatum) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("database: no users_metadata provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(usersMetadatumColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	usersMetadatumUpsertCacheMut.RLock()
	cache, cached := usersMetadatumUpsertCache[key]
	usersMetadatumUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			usersMetadatumAllColumns,
			usersMetadatumColumnsWithDefault,
			usersMetadatumColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			usersMetadatumAllColumns,
			usersMetadatumPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("database: unable to upsert users_metadata, could not build update column list")
		}

		ret := strmangle.SetComplement(usersMetadatumAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(usersMetadatumPrimaryKeyColumns) == 0 {
				return errors.New("database: unable to upsert users_metadata, could not build conflict column list")
			}

			conflict = make([]string, len(usersMetadatumPrimaryKeyColumns))
			copy(conflict, usersMetadatumPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"users_metadata\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(usersMetadatumType, usersMetadatumMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(usersMetadatumType, usersMetadatumMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "database: unable to upsert users_metadata")
	}

	if !cached {
		usersMetadatumUpsertCacheMut.Lock()
		usersMetadatumUpsertCache[key] = cache
		usersMetadatumUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single UsersMetadatum record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *UsersMetadatum) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single UsersMetadatum record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UsersMetadatum) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("database: no UsersMetadatum provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), usersMetadatumPrimaryKeyMapping)
	sql := "DELETE FROM \"users_metadata\" WHERE \"pub_key\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete from users_metadata")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by delete for users_metadata")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q usersMetadatumQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q usersMetadatumQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("database: no usersMetadatumQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from users_metadata")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for users_metadata")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o UsersMetadatumSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UsersMetadatumSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(usersMetadatumBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usersMetadatumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"users_metadata\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, usersMetadatumPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from usersMetadatum slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for users_metadata")
	}

	if len(usersMetadatumAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *UsersMetadatum) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: no UsersMetadatum provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UsersMetadatum) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUsersMetadatum(ctx, exec, o.PubKey)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UsersMetadatumSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("database: empty UsersMetadatumSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UsersMetadatumSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UsersMetadatumSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), usersMetadatumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"users_metadata\".* FROM \"users_metadata\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, usersMetadatumPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "database: unable to reload all in UsersMetadatumSlice")
	}

	*o = slice

	return nil
}

// UsersMetadatumExistsG checks if the UsersMetadatum row exists.
func UsersMetadatumExistsG(ctx context.Context, pubKey string) (bool, error) {
	return UsersMetadatumExists(ctx, boil.GetContextDB(), pubKey)
}

// UsersMetadatumExists checks if the UsersMetadatum row exists.
func UsersMetadatumExists(ctx context.Context, exec boil.ContextExecutor, pubKey string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"users_metadata\" where \"pub_key\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, pubKey)
	}
	row := exec.QueryRowContext(ctx, sql, pubKey)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "database: unable to check if users_metadata exists")
	}

	return exists, nil
}

// Exists checks if the UsersMetadatum row exists.
func (o *UsersMetadatum) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return UsersMetadatumExists(ctx, exec, o.PubKey)
}
