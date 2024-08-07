package repositories

import (
	"context"
	"errors"

	"github.com/peter6866/foodie/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MenuItemRepository struct {
	collection *mongo.Collection
}

func NewMenuItemRepository(client *mongo.Client) *MenuItemRepository {
	collection := client.Database("foodie").Collection("menuItems")
	return &MenuItemRepository{collection: collection}
}

// create a new menu item
func (r *MenuItemRepository) Create(ctx context.Context, menuItem *models.MenuItem) error {
	result, err := r.collection.InsertOne(ctx, menuItem)
	if err == nil {
		menuItem.ID = result.InsertedID.(primitive.ObjectID)
	}

	return err
}

// get a menu item by ID and createdBy the userID
func (r *MenuItemRepository) GetByID(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*models.MenuItem, error) {
	var item models.MenuItem
	err := r.collection.FindOne(ctx, bson.M{"_id": id, "createdBy": userID}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// get all menu items where createdBy is the userID
func (r *MenuItemRepository) GetAll(ctx context.Context, userID primitive.ObjectID) ([]*models.MenuItem, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"createdBy": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*models.MenuItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// update a menu item
func (r *MenuItemRepository) Update(ctx context.Context, menuItem *models.MenuItem) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": menuItem.ID},
		bson.M{"$set": menuItem},
	)
	return err
}

// delete a menu item
func (r *MenuItemRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
