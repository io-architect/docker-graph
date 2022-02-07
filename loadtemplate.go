package main

import (
        "embed"
        "html/template"
        "io/fs"

        "github.com/gin-gonic/gin"
)

func loadTemplate(engine *gin.Engine, assets embed.FS) error {
        fs.WalkDir(assets, "templates", func(path string, d fs.DirEntry, err error) error {
                if err != nil {
                        panic(err)
                }
                if d.IsDir() {
                        return nil
                }

                tmplByte, err := fs.ReadFile(assets, path)
                if err != nil {
                        return err
                }
                tmpl, err := template.New(path).Funcs(engine.FuncMap).Parse(string(tmplByte[:]))
                if err != nil {
                        return err
                }
                engine.SetHTMLTemplate(tmpl)

                return nil
        })
        return nil
}
