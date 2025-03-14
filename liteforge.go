package liteforge

import (
	"github.com/pballentine13/liteforge/internal/orm"
)

type Config = orm.Config

var OpenDB = orm.OpenDB
var CreateTable = orm.CreateTable
var Query = orm.Query
var QueryRow = orm.QueryRow
var Exec = orm.Exec
var BeginTx = orm.BeginTx
var SanitizeInput = orm.SanitizeInput
var GetTableName = orm.GetTableName
var GetFieldInfo = orm.GetFieldInfo

//type UserDataStore = orm.UserDataStore
//type SQLiteDataStore = orm.SQLiteDataStore
//type APIDataStore = orm.APIDataStore

