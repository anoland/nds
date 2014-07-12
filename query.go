package nds

import (
	"appengine"
	"appengine/datastore"
)

// Query represents a datastore query. Use NewQuery to create a new query.
type Query struct {
	query *datastore.Query
}

// NewQuery creates a new Query for a specific entity kind.
//
// An empty kind means to return all entities, including entities created and
// managed by other App Engine features, and is called a kindless query.
// Kindless queries cannot include filters or sort orders on property values.
func NewQuery(kind string) *Query {
	return &Query{datastore.NewQuery(kind)}
}

// Ancestor returns a derivative query with an ancestor filter.
// The ancestor should not be nil.
func (q *Query) Ancestor(ancestor *datastore.Key) *Query {
	return &Query{q.query.Ancestor(ancestor)}
}

// Filter returns a derivative query with a field-based filter.
// The filterStr argument must be a field name followed by optional space,
// followed by an operator, one of ">", "<", ">=", "<=", or "=".
// Fields are compared against the provided value using the operator.
// Multiple filters are AND'ed together.
func (q *Query) Filter(filter string, val interface{}) *Query {
	return &Query{q.query.Filter(filter, val)}
}

// Order returns a derivative query with a field-based sort order. Orders are
// applied in the order they are added. The default order is ascending; to sort
// in descending order prefix the fieldName with a minus sign (-).
func (q *Query) Order(field string) *Query {
	return &Query{q.query.Order(field)}
}

// Project returns a derivative query that yields only the given fields. It
// cannot be used with KeysOnly.
func (q *Query) Project(fields ...string) *Query {
	return &Query{q.query.Project(fields...)}
}

// Distinct returns a derivative query that yields de-duplicated entities with
// respect to the set of projected fields. It is only used for projection
// queries.
func (q *Query) Distinct() *Query {
	return &Query{q.query.Distinct()}
}

// KeysOnly returns a derivative query that yields only keys, not keys and
// entities. It cannot be used with projection queries.
func (q *Query) KeysOnly() *Query {
	return &Query{q.query.KeysOnly()}
}

// Limit returns a derivative query that has a limit on the number of results
// returned. A negative value means unlimited.
func (q *Query) Limit(limit int) *Query {
	return &Query{q.query.Limit(limit)}
}

// Offset returns a derivative query that has an offset of how many keys to
// skip over before returning results. A negative value is invalid.
func (q *Query) Offset(offset int) *Query {
	return &Query{q.query.Offset(offset)}
}

// EventualConsistency returns a derivative query that returns eventually
// consistent results.
// It only has an effect on ancestor queries.
func (q *Query) EventualConsistency() *Query {
	return &Query{q.query.EventualConsistency()}
}

// Start returns a derivative query with the given start point.
func (q *Query) Start(c Cursor) *Query {
	return &Query{q.query.Start(c.cursor)}
}

// End returns a derivative query with the given end point.
func (q *Query) End(c Cursor) *Query {
	return &Query{q.query.End(c.cursor)}
}

// Count returns the number of results for the query.
func (q *Query) Count(c appengine.Context) (int, error) {
	return q.query.Count(c)
}

// GetAll runs the query in the given context and returns all keys that match
// that query, as well as appending the values to dst.
//
// vals must have type *[]S or *[]*S, for some struct type S.
//
// Note that PropertyLoadSaver is currently not implemented.
//
// The keys returned by GetAll will be in a 1-1 correspondence with the entities
// added to vals.
//
// If q is a ``keys-only'' query, GetAll ignores vals and only returns the keys.
func (q *Query) GetAll(c appengine.Context,
	vals interface{}) ([]*datastore.Key, error) {
	return q.query.GetAll(c, vals)
}

// Iterator is the result of running a query.
type Iterator struct {
	*datastore.Iterator
}

// Done is returned when a query iteration has completed.
var Done = datastore.Done

// Run runs the query in the given context.
func (q *Query) Run(c appengine.Context) *Iterator {
	return &Iterator{q.query.Run(c)}
}

// Next returns the key of the next result. When there are no more results,
// Done is returned as the error.
func (i *Iterator) Next(val interface{}) (*datastore.Key, error) {
	return i.Iterator.Next(val)
}

// Cursor is an iterator's position. It can be converted to and from an opaque
// string. A cursor can be used from different HTTP requests, but only with a
// query with the same kind, ancestor, filter and order constraints.
type Cursor struct {
	cursor datastore.Cursor
}

// String returns a base-64 string representation of a cursor.
func (c Cursor) String() string {
	return c.cursor.String()
}

// Cursor returns a cursor for the iterator's current location.
func (i *Iterator) Cursor() (Cursor, error) {
	cursor, err := i.Iterator.Cursor()
	return Cursor{cursor}, err
}

// Decode decodes a cursor from its base-64 string representation.
func DecodeCursor(s string) (Cursor, error) {
	cursor, err := datastore.DecodeCursor(s)
	return Cursor{cursor}, err
}
