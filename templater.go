package templater

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
)

type Templater interface {
	Tmpl() *template.Template
}

type TemplaterOptions struct {
	Embedded         fs.FS
	FuncMap          template.FuncMap
	Pattern          string
	RootTemplateName string
}

func New(dir string, options TemplaterOptions) (Templater, error) {
	if options.Pattern == "" {
		options.Pattern = "*"
	}
	if options.Embedded != nil {
		return NewEmbeddedTemplater(dir, options.Embedded, options.FuncMap, options.Pattern, options.RootTemplateName)
	}
	return NewDiskTemplater(dir, options.FuncMap, options.Pattern, options.RootTemplateName), nil
}

type DiskTemplater struct {
	fs               fs.FS
	funcMap          template.FuncMap
	pattern          string
	rootTemplateName string
}

func NewDiskTemplater(dir string, funcMap template.FuncMap, pattern, rootTemplateName string) Templater {
	return &DiskTemplater{
		fs:               os.DirFS(dir),
		funcMap:          funcMap,
		pattern:          pattern,
		rootTemplateName: rootTemplateName,
	}
}

func (d *DiskTemplater) Tmpl() *template.Template {
	tmpl := template.New(d.rootTemplateName)
	tmpl.Funcs(d.funcMap)
	tmpl = template.Must(tmpl.ParseFS(d.fs, d.pattern))
	return tmpl
}

type EmbeddedTemplater struct {
	tmpl *template.Template
}

func NewEmbeddedTemplater(dir string, embedded fs.FS, funcMap template.FuncMap, pattern, rootTemplateName string) (Templater, error) {
	f, err := fs.Sub(embedded, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get subdirectory of embedded file system: %w", err)
	}
	tmpl := template.New(rootTemplateName)
	tmpl.Funcs(funcMap)
	tmpl = template.Must(tmpl.ParseFS(f, pattern))
	return EmbeddedTemplater{
		tmpl: tmpl,
	}, nil
}

func (e EmbeddedTemplater) Tmpl() *template.Template {
	return e.tmpl
}
