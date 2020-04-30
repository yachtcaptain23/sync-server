package datastore

// Datastore abstracts over the underlying datastore.
type Datastore interface {
	// Insert a new access token for a given client.
	InsertClientToken(id string, token string, expireAt int64) error
	// Get the client ID from a non-expired token.
	GetClientID(token string) (string, error)
	// Insert a new sync entity.
	InsertSyncEntity(entity *SyncEntity) error
	// Insert a series of sync entities in a write transaction.
	InsertSyncEntities(entities []*SyncEntity) error
	// Update an existing sync entity.
	UpdateSyncEntity(entity *SyncEntity) (bool, error)
	// Get updates for a specific type which are modified after the time of
	// client token for a given client.
	GetUpdatesForType(dataType int, clientToken int64, fetchFolders bool, clientID string) ([]SyncEntity, error)
	// Check if a server-defined unique tag is in the datastore.
	HasServerDefinedUniqueTag(clientID string, tag string) (bool, error)
}
