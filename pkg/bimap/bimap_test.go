package bimap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBimap_GetMissingReturnsFalse(t *testing.T) {
	m := New()
	val, ok := m.Get("missing")
	assert.Empty(t, val)
	assert.False(t, ok)
}

func TestBimap_GetByValMissingReturnsFalse(t *testing.T) {
	m := New()
	id, ok := m.GetByVal("missing")
	assert.Empty(t, id)
	assert.False(t, ok)
}

func TestBimap_GetExistingReturnsTrue(t *testing.T) {
	m := New()
	m.Put("id1", "val1")
	val, ok := m.Get("id1")
	assert.Equal(t, "val1", val)
	assert.True(t, ok)
}

func TestBimap_GetByValExistingReturnsTrue(t *testing.T) {
	m := New()
	m.Put("id1", "val1")
	id, ok := m.GetByVal("val1")
	assert.Equal(t, "id1", id)
	assert.True(t, ok)
}

func TestBimap_GetEmptyKeyReturnsTrue(t *testing.T) {
	m := New()
	m.Put("", "val1")
	val, ok := m.Get("")
	assert.Equal(t, "val1", val)
	assert.True(t, ok)
}

func TestBimap_GetByValEmptyVal(t *testing.T) {
	m := New()
	m.Put("id1", "")
	id, ok := m.GetByVal("")
	assert.Equal(t, "id1", id)
	assert.True(t, ok)
}

func TestBimap_DeleteMissingReturnsFalse(t *testing.T) {
	m := New()
	val, ok := m.Delete("missing")
	assert.Empty(t, val)
	assert.False(t, ok)
}

func TestBimap_DeleteByValMissingReturnsFalse(t *testing.T) {
	m := New()
	id, ok := m.DeleteByVal("missing")
	assert.Empty(t, id)
	assert.False(t, ok)
}

func TestBimap_DeleteExistingGetReturnsFalse(t *testing.T) {
	m := New()
	m.Put("id1", "val1")
	val, ok := m.Delete("id1")
	assert.Equal(t, "val1", val)
	assert.True(t, ok)
	val, ok = m.Get("id1")
	assert.Empty(t, val)
	assert.False(t, ok)
}

func TestBimap_DeleteByValExistingReturnsTrue(t *testing.T) {
	m := New()
	m.Put("id1", "val1")
	id, ok := m.DeleteByVal("val1")
	assert.Equal(t, "id1", id)
	assert.True(t, ok)
	id, ok = m.GetByVal("val1")
	assert.Empty(t, id)
	assert.False(t, ok)
}

func TestBimap_ReplaceByIdGetReturnsTrue(t *testing.T) {
	m := New()
	m.Put("id1", "val")
	m.Put("id1", "val1")
	val, ok := m.Get("id1")
	assert.Equal(t, "val1", val)
	assert.True(t, ok)
}

func TestBimap_ReplaceByValGetByValReturnsTrue(t *testing.T) {
	m := New()
	m.Put("id", "val1")
	m.Put("id1", "val1")
	id, ok := m.GetByVal("val1")
	assert.Equal(t, "id1", id)
	assert.True(t, ok)
}

func TestBimap_ReplaceByIdAndValGetReturnsTrue(t *testing.T) {
	m := New()
	m.Put("id1", "val1")
	m.Put("id1", "val1")
	val, ok := m.Get("id1")
	assert.Equal(t, "val1", val)
	assert.True(t, ok)
}

func TestBimap_ReplaceByIdAndValGetByValReturnsTrue(t *testing.T) {
	m := New()
	m.Put("id1", "val1")
	m.Put("id1", "val1")
	id, ok := m.GetByVal("val1")
	assert.Equal(t, "id1", id)
	assert.True(t, ok)
}

func TestBimap_IterateEmptyMap(t *testing.T) {
	m := New()
	m.Iterate(
		func(key string, val string) {
			t.FailNow()
		})
}

func TestBimap_IterateMap(t *testing.T) {
	m := New()
	m.Put("id1", "val1")
	m.Put("id2", "val2")
	vals := map[string]string{}
	m.Iterate(
		func(key string, val string) {
			vals[key] = val
		})
	assert.Equal(t, map[string]string{"id1": "val1", "id2": "val2"}, vals)
}
