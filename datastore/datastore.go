package datastore

// Datastore abstracts over the underlying datastore.
type Datastore interface {
	// Insert a new access token for a given client.
	InsertToken(id string, token string, expireAt int64) error
	// Get the client ID from a non-expired token.
	GetClientID(token string) (string, error)
	// Insert a new sync entity.
	InsertSyncEntity(entity *SyncEntity) error
	// Insert a series of sync entities in a write transaction.
	InsertSyncEntities(entities []*SyncEntity) error
	// Update an existing sync entity.
	UpdateSyncEntity(entity *SyncEntity) (int64, error)
	// Get updates for a specific type which are modified after the time of
	// client token for a given client.
	GetUpdatesForType(dataType int32, clientToken int64, fetchFolders bool, clientID string) ([]SyncEntity, error)
	// TODO: Remove this
	CheckVersion(id string, clientVersion int64) (bool, error)
	// TODO: Remove this
	IsServerDefinedUniqueEntitiesReady(tags []string, clientID string) (bool, error)
}
