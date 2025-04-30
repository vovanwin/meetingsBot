package store

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --target ./gen ./schema
//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --target ./gen --feature sql/upsert --feature sql/lock --feature sql/execquery --template ./templates ./schema
