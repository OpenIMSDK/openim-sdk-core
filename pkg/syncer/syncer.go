package syncer

import (
	"context"
	"reflect"

	"github.com/google/go-cmp/cmp"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/network"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/page"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/tools/log"
)

// New creates a new Syncer instance with the provided synchronization functions.
// This function requires several callback functions to handle different aspects of data synchronization:
// - insert: A function to insert new data.
// - delete: A function to delete existing data.
// - update: A function to update existing data.
// - uuid: A function to generate a unique identifier for each data item.
// - equal: A function to check if two data items are equal.
// - notice: A function to handle notifications during the sync process.
// Panics if insert, delete, update, or uuid functions are not provided.
func New[T, RESP any, V comparable](
	insert func(ctx context.Context, value T) error,
	delete func(ctx context.Context, value T) error,
	update func(ctx context.Context, server T, local T) error,
	uuid func(value T) V,
	equal func(a T, b T) bool,
	notice func(ctx context.Context, state int, server, local T) error,
) *Syncer[T, RESP, V] {
	if insert == nil || update == nil || delete == nil || uuid == nil {
		panic("invalid params")
	}

	// Determine the type of T and remove pointer indirection if necessary.
	var t T
	tof := reflect.TypeOf(&t)
	for tof.Kind() == reflect.Ptr {
		tof = tof.Elem()
	}

	// Return a new Syncer instance with the provided functions and the type as a string.
	return &Syncer[T, RESP, V]{
		insert: insert,
		update: update,
		delete: delete,
		uuid:   uuid,
		equal:  equal,
		notice: notice,
		ts:     tof.String(),
	}
}

// Syncer is a struct that holds functions for synchronizing data.
// It includes functions for inserting, updating, deleting, and notifying,
// as well as functions for generating unique IDs and checking equality of data items.
type Syncer[T, RESP any, V comparable] struct {
	insert                   func(ctx context.Context, server T) error
	update                   func(ctx context.Context, server T, local T) error
	delete                   func(ctx context.Context, local T) error
	batchInsert              func(ctx context.Context, servers []T) error
	deleteAll                func(ctx context.Context, entityID string) error
	notice                   func(ctx context.Context, state int, server, local T) error
	equal                    func(server T, local T) bool
	uuid                     func(value T) V
	batchPageReq             func(entityID string) page.PageReq
	batchPageRespConvertFunc func(resp *RESP) []T
	reqApiRouter             string
	ts                       string // Represents the type of T as a string.
	fullSyncLimit            int64
}

type NoResp struct{}

func New2[T, RESP any, V comparable](opts ...Option[T, RESP, V]) *Syncer[T, RESP, V] {
	// Create a new Syncer instance.
	s := &Syncer[T, RESP, V]{}

	// Apply the options to the Syncer.
	for _, opt := range opts {
		opt(s)
	}

	// Check required functions.
	if s.insert == nil || s.update == nil || s.delete == nil || s.uuid == nil {
		panic("invalid params")
	}

	// Determine the type of T and remove pointer indirection if necessary.
	var t T
	tof := reflect.TypeOf(&t)
	for tof.Kind() == reflect.Ptr {
		tof = tof.Elem()
	}

	// Set the type string.
	s.ts = tof.String()

	return s
}

// Option is a function that configures a Syncer.
type Option[T, RESP any, V comparable] func(*Syncer[T, RESP, V])

// WithInsert sets the insert function for the Syncer.
func WithInsert[T, RESP any, V comparable](f func(ctx context.Context, value T) error) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.insert = f
	}
}

// WithUpdate sets the update function for the Syncer.
func WithUpdate[T, RESP any, V comparable](f func(ctx context.Context, server T, local T) error) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.update = f
	}
}

// WithDelete sets the delete function for the Syncer.
func WithDelete[T, RESP any, V comparable](f func(ctx context.Context, local T) error) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.delete = f
	}
}

// WithBatchInsert sets the batchInsert function for the Syncer.
func WithBatchInsert[T, RESP any, V comparable](f func(ctx context.Context, values []T) error) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.batchInsert = f
	}
}

// WithDeleteAll sets the deleteAll function for the Syncer.
func WithDeleteAll[T, RESP any, V comparable](f func(ctx context.Context, entityID string) error) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.deleteAll = f
	}
}

// WithUUID sets the uuid function for the Syncer.
func WithUUID[T, RESP any, V comparable](f func(value T) V) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.uuid = f
	}
}

// WithEqual sets the equal function for the Syncer.
func WithEqual[T, RESP any, V comparable](f func(a T, b T) bool) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.equal = f
	}
}

// WithNotice sets the notice function for the Syncer.
func WithNotice[T, RESP any, V comparable](f func(ctx context.Context, state int, server, local T) error) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.notice = f
	}
}

// WithBatchPageReq sets the batchPageReq for the Syncer.
func WithBatchPageReq[T, RESP any, V comparable](f func(entityID string) page.PageReq) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.batchPageReq = f
	}
}

// WithBatchPageRespConvertFunc sets the batchPageRespConvertFunc function for the Syncer.
func WithBatchPageRespConvertFunc[T, RESP any, V comparable](f func(resp *RESP) []T) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.batchPageRespConvertFunc = f
	}
}

// WithReqApiRouter sets the reqApiRouter for the Syncer.
func WithReqApiRouter[T, RESP any, V comparable](router string) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.reqApiRouter = router
	}
}

// WithFullSyncLimit sets the fullSyncLimit for the Syncer.
func WithFullSyncLimit[T, RESP any, V comparable](limit int64) Option[T, RESP, V] {
	return func(s *Syncer[T, RESP, V]) {
		s.fullSyncLimit = limit
	}
}

// NewSyncer creates a new Syncer with the provided options.
func NewSyncer[T, RESP any, V comparable](opts ...Option[T, RESP, V]) *Syncer[T, RESP, V] {
	syncer := &Syncer[T, RESP, V]{}
	for _, opt := range opts {
		opt(syncer)
	}
	return syncer
}

// eq is a helper function to check equality of two data items.
// It uses the equal function if provided; otherwise, it falls back to the cmp.Equal function.
func (s *Syncer[T, RESP, V]) eq(server T, local T) bool {
	if s.equal != nil {
		return s.equal(server, local)
	}
	return cmp.Equal(server, local)
}

// onNotice is a helper function to handle notifications.
// It calls the Syncer's notice function and the provided notice function in sequence if they are not nil.
func (s *Syncer[T, RESP, V]) onNotice(ctx context.Context, state int, server, local T, fn func(ctx context.Context, state int, server, local T) error) error {
	if s.notice != nil {
		if err := s.notice(ctx, state, server, local); err != nil {
			return err
		}
	}
	if fn != nil {
		if err := fn(ctx, state, server, local); err != nil {
			return err
		}
	}
	return nil
}

// Sync synchronizes server data with local data.
// Sync synchronizes the data between the server and local storage.
// It takes a context, two slices of data (serverData and localData),
// a notice function to handle notifications during the sync process,
// and a variadic parameter skipDeletion to control deletion behavior.
func (s *Syncer[T, RESP, V]) Sync(ctx context.Context, serverData []T, localData []T, notice func(ctx context.Context, state int, server, local T) error, skipDeletionAndNotice ...bool) (err error) {
	defer func() {
		// Log the outcome of the synchronization process.
		if err == nil {
			log.ZDebug(ctx, "sync success", "type", s.ts)
		} else {
			log.ZError(ctx, "sync failed", err, "type", s.ts)
		}
	}()
	skipDeletion := false
	skipNotice := false
	if len(skipDeletionAndNotice) > 0 {
		skipDeletion = skipDeletionAndNotice[0]
	}
	if len(skipDeletionAndNotice) > 1 {
		skipNotice = skipDeletionAndNotice[1]
	}

	// If both server and local data are empty, log and return.
	if len(serverData) == 0 && len(localData) == 0 {
		log.ZDebug(ctx, "sync both the server and client are empty", "type", s.ts)
		return nil
	}

	// Convert local data into a map for easier lookup.
	localMap := datautil.SliceToMap(localData, func(item T) V {
		return s.uuid(item)
	})

	// Iterate through server data to sync with local data.
	for i := range serverData {
		server := serverData[i]
		id := s.uuid(server)
		local, ok := localMap[id]

		// If the item doesn't exist locally, insert it.
		if !ok {
			log.ZDebug(ctx, "sync insert", "type", s.ts, "server", server)
			if err := s.insert(ctx, server); err != nil {
				log.ZError(ctx, "sync insert failed", err, "type", s.ts, "server", server, "local", local)
				return err
			}
			if !skipNotice {
				if err := s.onNotice(ctx, Insert, server, local, notice); err != nil {
					log.ZError(ctx, "sync notice insert failed", err, "type", s.ts, "server", server, "local", local)
					return err
				}
			}
			continue
		}

		// Remove the item from the local map as it's found.
		delete(localMap, id)

		// If the local and server items are equal, notify and continue.
		if s.eq(server, local) {
			if !skipNotice {
				if err := s.onNotice(ctx, Unchanged, local, server, notice); err != nil {
					log.ZError(ctx, "sync notice unchanged failed", err, "type", s.ts, "server", server, "local", local)
					return err
				}
			}
			continue
		}

		log.ZDebug(ctx, "sync update", "type", s.ts, "server", server, "local", local)
		// Update the local item with server data.
		if err := s.update(ctx, server, local); err != nil {
			log.ZError(ctx, "sync update failed", err, "type", s.ts, "server", server, "local", local)
			return err
		}
		if !skipNotice {
			if err := s.onNotice(ctx, Update, server, local, notice); err != nil {
				log.ZError(ctx, "sync notice update failed", err, "type", s.ts, "server", server, "local", local)
				return err
			}
		}
	}

	// Check the skipDeletion flag; if set, skip deletion.
	if skipDeletion {
		return nil
	}
	log.ZDebug(ctx, "sync delete", "type", s.ts, "localMap", localMap)
	// Delete any local items not present in server data.
	for id := range localMap {
		local := localMap[id]
		if err := s.delete(ctx, local); err != nil {
			log.ZError(ctx, "sync delete failed", err, "type", s.ts, "local", local)
			return err
		}
		var server T
		if !skipNotice {
			if err := s.onNotice(ctx, Delete, server, local, notice); err != nil {
				log.ZError(ctx, "sync notice delete failed", err, "type", s.ts, "local", local)
				return err
			}
		}
	}
	return nil
}

func (s *Syncer[T, RESP, V]) FullSync(ctx context.Context, entityID string) (err error) {
	defer func() {
		// Log the outcome of the synchronization process.
		if err == nil {
			log.ZDebug(ctx, "full sync success", "type", s.ts)
		} else {
			log.ZError(ctx, "full sync failed", err, "type", s.ts)
		}
	}()

	//// If server data is empty, log and return
	//if len(serverData) == 0 {
	//	log.ZDebug(ctx, "full sync server data is empty", "type", s.ts)
	//	return nil
	//}

	// Clear local table data
	if err = s.deleteAll(ctx, entityID); err != nil {
		return errs.New("full sync delete all failed", "err", err.Error(), "type", s.ts)
	}

	// Get batch req
	batchReq := s.batchPageReq(entityID)

	// Batch page pull data and insert server data
	if err = network.FetchAndInsertPagedData(ctx, s.reqApiRouter, batchReq, s.batchPageRespConvertFunc,
		s.batchInsert, s.insert, s.fullSyncLimit); err != nil {
		return errs.New("full sync batch insert failed", "err", err.Error(), "type", s.ts)
	}

	return nil

}
