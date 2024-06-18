package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/SerjRamone/vaultme/internal/models"
)

var (
	errMsgStartTx  = "starting transaction error: %w"
	errMsgRollback = "rolling back transaction error"
	errMsgCommit   = "committing transaction error: %w"
)

// GetItem gets user's item by a given ID
func (db *DB) GetItem(ctx context.Context, userID string, itemId string) (*models.Item, error) {
	query := `SELECT i.id, i.user_id, i.name, i.type, r.version, i.created_at, r.updated_at, d.data
	FROM item as i
	LEFT JOIN item_data as d ON i.id = d.item_id
	WHERE i.user_id =  $1 AND i.id = $2`

	var i models.Item

	row := db.pool.QueryRow(ctx, query, userID, itemId)
	if err := row.Scan(&i.ID, &i.UserID, &i.Name, &i.Type, &i.Version, &i.CreatedAt, &i.UpdatedAt, &i.Data); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("item not found: %w", err)
		}
		return nil, fmt.Errorf("getting item error: %w", err)
	}

	return &i, nil
}

// CreateItem creates new item
func (db *DB) CreateItem(ctx context.Context, userID string, item *models.ItemDTO) (*models.Item, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(errMsgStartTx, err)
	}

	defer func(tx pgx.Tx) {
		if err := tx.Rollback(ctx); err != nil {
			if !errors.Is(err, pgx.ErrTxClosed) {
				db.log.Error(errMsgRollback, zap.Error(err))
			}
		}
	}(tx)

	query := `INSERT INTO item (user_id, name, type) VALUES ($1, $2, $3)
	RETURNING id, user_id, name, type, version, created_at, updated_at;`

	var i models.Item
	row := tx.QueryRow(ctx, query, userID, item.Name, item.Type)
	if err := row.Scan(&i.ID, &i.UserID, &i.Name, &i.Type, &i.Version, &i.CreatedAt, &i.UpdatedAt); err != nil {
		return nil, fmt.Errorf("creating item error: %w", err)
	}

	if err := db.addItemData(ctx, tx, i.ID, item.Data); err != nil {
		return nil, fmt.Errorf("add item data error: %w", err)
	}
	i.Data = item.Data

	if err := db.updateItemMeta(ctx, tx, i.ID, item.Meta); err != nil {
		return nil, fmt.Errorf("add item meta error: %w", err)
	}
	i.Meta = item.Meta

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(errMsgCommit, err)
	}

	return &i, nil
}

// UpdateItem updates item
func (db *DB) UpdateItem(ctx context.Context, userID string, item *models.Item) (*models.Item, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(errMsgStartTx, err)
	}

	defer func(tx pgx.Tx) {
		if err := tx.Rollback(ctx); err != nil {
			if !errors.Is(err, pgx.ErrTxClosed) {
				db.log.Error(errMsgRollback, zap.Error(err))
			}
		}
	}(tx)

	var r models.Item

	query := `INSERT INTO item (id, user_id, name, type, version) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (id) DO UPDATE SET name = $2, version = $4, updated_at = CURRENT_TIMESTAMP
	RETURNING id, user_id, name, type, version, created_at, updated_at;`
	row := tx.QueryRow(ctx, query, item.ID, userID, item.Name, item.Type, item.Version)
	if err := row.Scan(
		&r.ID, &r.UserID, &r.Name, &r.Type, &r.Version, &r.CreatedAt, &r.UpdatedAt); err != nil {
		return nil, fmt.Errorf("row scan error: %w", err)
	}

	if err := db.addItemData(ctx, tx, r.ID, item.Data); err != nil {
		return nil, fmt.Errorf("adding item data error: %w", err)
	}

	if err := db.updateItemMeta(ctx, tx, r.ID, item.Meta); err != nil {
		return nil, fmt.Errorf("updating intem meta error: %w", err)
	}

	r.Data = item.Data
	r.Meta = item.Meta

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(errMsgCommit, err)
	}

	return &r, nil
}

// ListItems returns list of items
func (db *DB) ListItems(ctx context.Context, userID string, limit int, offset int) ([]*models.Item, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(errMsgStartTx, err)
	}

	defer func(tx pgx.Tx) {
		if err := tx.Rollback(ctx); err != nil {
			if !errors.Is(err, pgx.ErrTxClosed) {
				db.log.Error(errMsgRollback, zap.Error(err))
			}
		}
	}(tx)

	query := `SELECT i.id, i.user_id, i.name, i.type, i.version, i.created_at, i.updated_at, d.data
	FROM item as i
	LEFT JOIN item_data as d ON i.id = d.item_id
	WHERE i.user_id = $1
	LIMIT $2 
	OFFSET $3`

	rows, err := tx.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get items froms storage error: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	var ids []string
	for rows.Next() {
		var i models.Item
		err := rows.Scan(&i.ID, &i.UserID, &i.Name, &i.Type, &i.Version, &i.CreatedAt, &i.UpdatedAt, &i.Data)
		if err != nil {
			return nil, fmt.Errorf("rows scan error: %w", err)
		}
		items = append(items, &i)
		ids = append(ids, i.ID)
	}

	itemMeta, err := db.getItemMeta(ctx, tx, ids)
	if err != nil {
		return nil, fmt.Errorf("getting meta error: %w", err)
	}

	for _, i := range items {
		if mi, ok := itemMeta[i.ID]; ok {
			i.Meta = mi
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(errMsgCommit, err)
	}

	return items, nil
}

func (db *DB) addItemData(ctx context.Context, tx pgx.Tx, itemID string, itemData []byte) error {
	query := `INSERT INTO item_data (item_id, data) VALUES ($1, $2)
	ON CONFLICT (item_id) DO UPDATE SET data = $2`
	if _, err := tx.Exec(ctx, query, itemID, itemData); err != nil {
		return fmt.Errorf("upsert item data error: %w", err)
	}

	return nil
}

func (db *DB) addMeta(ctx context.Context, tx pgx.Tx, itemID string, m *models.Meta) error {
	query := `INSERT INTO meta (item_id, tag, text) VALUES ($1, $2, $3);`

	if _, err := tx.Exec(ctx, query, itemID, m.Tag, m.Text); err != nil {
		return fmt.Errorf("upsert meta error: %w", err)
	}

	return nil
}

func (db *DB) getItemMeta(ctx context.Context, tx pgx.Tx, itemIDs []string) (map[string][]*models.Meta, error) {
	query := `SELECT item_id, tag, text
	FROM meta
	WHERE item_id = ANY($1);`
	rows, err := tx.Query(ctx, query, itemIDs)
	if err != nil {
		return nil, fmt.Errorf("getting meta error: %w", err)
	}

	defer rows.Close()

	metas := make(map[string][]*models.Meta)

	for rows.Next() {
		itemID := ""
		var m models.Meta

		if err := rows.Scan(&itemID, &m.Tag, &m.Text); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		metas[itemID] = append(metas[itemID], &m)
	}

	return metas, nil
}

func (db *DB) updateItemMeta(ctx context.Context, tx pgx.Tx, itemID string, metas []*models.Meta) error {
	query := `DELETE FROM meta WHERE item_id = $1;`
	if _, err := tx.Exec(ctx, query, itemID); err != nil {
		return fmt.Errorf("an occured error while cleaning up metadata, err: %w", err)
	}
	for _, m := range metas {
		if err := db.addMeta(ctx, tx, itemID, m); err != nil {
			return fmt.Errorf("update meta for item (ID: %s) error: %w", itemID, err)
		}
	}
	return nil
}
