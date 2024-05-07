package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

// goal of this file:
// 1. doesnt import anything
// 2. is accessable from all files that need them when it is called (global, but better, better because it's only global when it is called because it is needed)

// this protects us from an import cycle problem which can hurt the compiler and even stop it from compiling

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template 
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}
