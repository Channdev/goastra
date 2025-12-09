/*
 * GoAstra CLI - Ent Schema Template
 *
 * Generates Ent ORM schema definitions.
 * Provides User entity and base schema patterns.
 */
package ent

// UserSchemaGo returns the User entity schema template.
func UserSchemaGo() string {
	return `package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

/*
 * User holds the schema definition for the User entity.
 * Represents user accounts in the system.
 */
type User struct {
	ent.Schema
}

/*
 * Fields of the User entity.
 */
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").
			NotEmpty().
			Unique().
			Comment("User's email address"),

		field.String("password").
			Sensitive().
			Comment("Bcrypt hashed password"),

		field.String("name").
			NotEmpty().
			Comment("User's display name"),

		field.String("role").
			Default("user").
			Comment("User role: admin, user, etc."),

		field.Bool("active").
			Default(true).
			Comment("Whether the user account is active"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Account creation timestamp"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Last update timestamp"),
	}
}

/*
 * Edges of the User entity.
 * Define relationships to other entities here.
 */
func (User) Edges() []ent.Edge {
	return nil
}

/*
 * Indexes of the User entity.
 */
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email"),
		index.Fields("role"),
		index.Fields("created_at"),
	}
}
`
}

// BaseMixinGo returns a mixin template for common fields.
func BaseMixinGo() string {
	return `package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

/*
 * TimestampsMixin adds created_at and updated_at fields to schemas.
 * Embed this mixin in your schemas for automatic timestamp management.
 */
type TimestampsMixin struct {
	mixin.Schema
}

/*
 * Fields of the TimestampsMixin.
 */
func (TimestampsMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

/*
 * SoftDeleteMixin adds soft delete functionality.
 * Embed this mixin to add deleted_at field instead of hard deletes.
 */
type SoftDeleteMixin struct {
	mixin.Schema
}

/*
 * Fields of the SoftDeleteMixin.
 */
func (SoftDeleteMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}
`
}
