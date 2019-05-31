package catsql

import (
  "log"

  "github.com/go-pg/pg"
)

type postModelInitHook func(db *pg.DB) error

var postModelInitHooks = make([]postModelInitHook, 0, 8)

func RegisterPostModelInitHook(hook postModelInitHook) {
  postModelInitHooks = append(postModelInitHooks, hook)
}

type initLogger struct { }

func (d initLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d initLogger) AfterQuery(q *pg.QueryEvent) {
  // TODO: for test and prod, would be nice to compare what gets generated here
  // with what's in the local git clone for integrity.
  if env.IsDev() {
    if file == nil {
      initialize it
    }
  	file.Println(q.FormattedQuery())
  }
}

func InitializeDB(modelDefs ...interface{}) {
  RawConnect()
  db.AddQueryHook(initLogger{})

  createOptions := pg.CreateTableOptions{
    FKConstraints : true,
  }

  for _, modelDef := range modelDefs {
    if err := db.CreateTable(modeDef, &createOptions); err != nil {
      // TODO: can we get this to print the struct def? If not, just name using
      // reflect?
      log.Panicf("Could not create table for %v", modelDef)
    }
  }

  for _, hook := range postModelInitHooks {
    if err := hook(db); err != nil {
      log.Panicf(`Could not execute post-model init hook: %+v; %s`, hook, err)
    }
  }
}
