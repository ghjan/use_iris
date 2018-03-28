package main

//Go实战--也许最快的Go语言Web框架kataras/iris初识三(Redis、leveldb、BoltDB)
//https://blog.csdn.net/wangshubo1989/article/details/78344183

import (
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"

	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/boltdb"
)

func main() {
	db, _ := boltdb.New("./sessions.db", 0666, "users")
	// use different go routines to sync the database
	db.Async(true)

	// close and unlock the database when control+C/cmd+C pressed
	iris.RegisterOnInterrupt(func() {
		db.Close()
	})

	sess := sessions.New(sessions.Config{
		Cookie:  "sessionscookieid",
		Expires: 45 * time.Minute, // <=0 means unlimited life
	})

	//
	// IMPORTANT:
	//
	sess.UseDatabase(db)

	// the rest of the code stays the same.
	app := iris.New()

	app.Get("/", func(ctx context.Context) {
		ctx.Writef("You should navigate to the /set, /get, /delete, /clear,/destroy instead")
	})
	app.Get("/set", func(ctx context.Context) {
		s := sess.Start(ctx)
		//set session values
		s.Set("name", "iris")

		//test if setted here
		ctx.Writef("All ok session setted to: %s", s.GetString("name"))
	})

	app.Get("/set/{key}/{value}", func(ctx context.Context) {
		key, value := ctx.Params().Get("key"), ctx.Params().Get("value")
		s := sess.Start(ctx)
		// set session values
		s.Set(key, value)

		// test if setted here
		ctx.Writef("All ok session setted to %s: %s", key, s.GetString(key))
	})

	app.Get("/get", func(ctx context.Context) {
		// get a specific key, as string, if no found returns just an empty string
		name := sess.Start(ctx).GetString("name")

		ctx.Writef("The name on the /get was : %s", name)
	})

	app.Get("/get/{key}", func(ctx context.Context) {
		// get a specific key, as string, if no found returns just an empty string
		name := sess.Start(ctx).GetString(ctx.Params().Get("key"))
		ctx.Writef("The key/value on the /get was %s : %s", ctx.Params().Get("key"), name)
	})

	app.Get("/delete", func(ctx context.Context) {
		// delete a specific key
		sess.Start(ctx).Delete("name")
	})

	app.Get("/delete/{key}", func(ctx context.Context) {
		// delete a specific key
		key := ctx.Params().Get("key")
		sess.Start(ctx).Delete(key)
		ctx.Writef("The key on the /delete was %s", key)

	})

	app.Get("/clear", func(ctx context.Context) {
		// removes all entries
		sess.Start(ctx).Clear()
	})

	app.Get("/destroy", func(ctx context.Context) {
		//destroy, removes the entire session data and cookie
		sess.Destroy(ctx)
	})

	app.Get("/update", func(ctx context.Context) {
		// updates expire date with a new date
		sess.ShiftExpiration(ctx)
	})

	app.Run(iris.Addr(":8080"))
}
