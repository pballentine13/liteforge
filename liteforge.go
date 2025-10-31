package liteforge

import (
	"github.com/pballentine13/liteforge/internal/orm"
	"github.com/pballentine13/liteforge/pkg/model"
)

type Config = orm.Config
type Datastore = orm.Datastore

// Repository is the high-level, model-centric interface for CRUD operations.
type Repository = model.Repository

// NewRepository creates a new model-centric repository.
var NewRepository = model.NewORMRepository

var OpenDB = orm.OpenDB
var CreateTable = orm.CreateTable
var Query = orm.Query
var QueryRow = orm.QueryRow
var Exec = orm.Exec
var BeginTx = orm.BeginTx
var Sanitize = orm.Sanitize
var GetTableName = orm.GetTableName
var GetFieldInfo = orm.GetFieldInfo

//type UserDataStore = orm.UserDataStore
//type SQLiteDataStore = orm.SQLiteDataStore
//type APIDataStore = orm.APIDataStore
