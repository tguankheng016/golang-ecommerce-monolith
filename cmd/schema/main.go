package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"ariga.io/atlas-provider-gorm/gormschema"
	identityModel "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
)

func main() {
	sb := &strings.Builder{}
	loadModels(sb)

	io.WriteString(os.Stdout, sb.String())
}

func loadModels(sb *strings.Builder) {
	models := []interface{}{
		&identityModel.User{},
		&identityModel.Role{},
		&identityModel.UserRolePermission{},
		&identityModel.UserToken{},
	}
	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	sb.WriteString(stmts)
	sb.WriteString(";\n")
}
