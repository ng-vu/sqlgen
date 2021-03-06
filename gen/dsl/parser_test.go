package dsl_test

import (
	"testing"

	. "github.com/ng-vu/sqlgen/gen/dsl"
	. "github.com/ng-vu/sqlgen/mock"
)

func TestDSL(t *testing.T) {
	t.Run("Full simple declaration", func(t *testing.T) {
		src := `generate Account (plural Accounts) from "account";`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), src+"\n")
		AssertEqual(t, len(file.Declarations), 1)
	})

	t.Run("Error: Syntax error", func(t *testing.T) {
		src := `generate Account (plural Accounts from "account";`
		_, err := ParseString("test", src)
		AssertErrorEqual(t, err, "Error at test:1:35: syntax error")
	})

	t.Run("Spacing", func(t *testing.T) {
		src := `generate Account()from"account"`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate Account from "account";`+"\n")
		AssertEqual(t, len(file.Declarations), 1)
	})

	t.Run("Simplified 1", func(t *testing.T) {
		src := `generate`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate {} from "{}";`+"\n")
	})

	t.Run("Simplified 2", func(t *testing.T) {
		src := `generate Account`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate Account from "{}";`+"\n")
	})

	t.Run("Simplified 3", func(t *testing.T) {
		src := `generate from account`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate {} from "account";`+"\n")
	})

	t.Run("Simplified 4", func(t *testing.T) {
		src := `generate from "account"`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate {} from "account";`+"\n")
	})

	t.Run("Simplified with options", func(t *testing.T) {
		src := `generate (plural Accounts)`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate {} (plural Accounts) from "{}";`+"\n")
	})

	t.Run("Commonly use", func(t *testing.T) {
		src := `generate Account from account`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate Account from "account";`+"\n")
	})

	t.Run("Empty option", func(t *testing.T) {
		src := `generate Account () from "account"`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate Account from "account";`+"\n")
	})

	t.Run("Table name with escape characters", func(t *testing.T) {
		src := `generate Account () from "schema\"account"`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate Account from "schema\"account";`+"\n")
		AssertEqual(t, file.Declarations[0].TableName, `schema"account`)
	})

	t.Run("Table name with quotation", func(t *testing.T) {
		src := `generate Account () from "schema.account"`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate Account from "schema.account";`+"\n")
		AssertEqual(t, file.Declarations[0].TableName, `schema.account`)
	})

	t.Run("Table name with schema", func(t *testing.T) {
		src := `generate Account () from schema.account`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate Account from "schema"."account";`+"\n")
		AssertEqual(t, file.Declarations[0].SchemaName, `schema`)
		AssertEqual(t, file.Declarations[0].TableName, `account`)
	})

	t.Run("Table name with schema and quotation", func(t *testing.T) {
		src := `generate Account () from "schema"."account"`
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), `generate Account from "schema"."account";`+"\n")
		AssertEqual(t, file.Declarations[0].SchemaName, `schema`)
		AssertEqual(t, file.Declarations[0].TableName, `account`)
	})

	t.Run("Multiple declarations", func(t *testing.T) {
		src := `
generate Account from account;
generate User (plural Users) from "user"
`
		expected := `
generate Account from "account";
generate User (plural Users) from "user";
`[1:]
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), expected)
		AssertEqual(t, len(file.Declarations), 2)
	})

	t.Run("Auto semicolon insertion", func(t *testing.T) {
		src := `generate generate Account generate from account`
		expected := `
generate {} from "{}";
generate Account from "{}";
generate {} from "account";
`[1:]
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), expected)
		AssertEqual(t, len(file.Declarations), 3)
	})
}

func TestJoin(t *testing.T) {
	t.Run("Full syntax", func(t *testing.T) {
		src := `
generate UserJoinAccount
	from "user"    (User)    as u
	join "account" (Account) as a on u.id = a.user_id
`
		expected := `
generate UserJoinAccount
    from "user" (User) as u
    join "account" (Account) as a on u.id = a.user_id;
`[1:]
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), expected)
		AssertEqual(t, len(file.Declarations[0].Joins), 2)
	})

	t.Run("Full syntax with 3 joins", func(t *testing.T) {
		src := `
generate UserJoinAccount
	from "user"         (User)        as u
	join "account_user" (AccountUser) as au on u.id = au.user_id
	full join "account"      (Account)     as a  on a.id = au.account_id;
`
		expected := `
generate UserJoinAccount
    from "user" (User) as u
    join "account_user" (AccountUser) as au on u.id = au.user_id
    full join "account" (Account) as a on a.id = au.account_id;
`[1:]
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), expected)
		AssertEqual(t, len(file.Declarations[0].Joins), 3)
	})

	t.Run("Simplified, keep spacing and double quotes", func(t *testing.T) {
		src := `
generate UserJoinAccount
	from user
	join account_user on "user".id  = account_user.user_id
	full join account      on account.id = account_user.account_id
`
		expected := `
generate UserJoinAccount
    from "user"
    join "account_user" on "user".id  = account_user.user_id
    full join "account" on account.id = account_user.account_id;
`[1:]
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), expected)
		AssertEqual(t, len(file.Declarations[0].Joins), 3)
	})

	t.Run("Simplified, use ` for on condition", func(t *testing.T) {
		src := `
generate UserJoinAccount
	from user
	join account_user on ` + "`" + `"user".id  = account_user.user_id` + "`" + `
	full join account      on ` + "`" + `account.id = account_user.account_id` + "`" + `
`
		expected := `
generate UserJoinAccount
    from "user"
    join "account_user" on "user".id  = account_user.user_id
    full join "account" on account.id = account_user.account_id;
`[1:]
		file, err := ParseString("test", src)
		AssertNoError(t, err)
		AssertEqual(t, file.String(), expected)
		AssertEqual(t, len(file.Declarations[0].Joins), 3)
	})
}
