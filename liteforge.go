package liteforge

import (
	"github.com/pballentine13/liteforge/internal/orm"
)

type Config = orm.Config

var OpenDB = orm.OpenDB
var CreateTable = orm.CreateTable
var Create = orm.Create
var Get = orm.Get
var Update = orm.Update
var Delete = orm.Delete
var Query = orm.Query
var Exec = orm.Exec
var BeginTx = orm.BeginTx
var SanitizeInput = orm.SanitizeInput

//type UserDataStore = orm.UserDataStore
//type SQLiteDataStore = orm.SQLiteDataStore
//type APIDataStore = orm.APIDataStore
