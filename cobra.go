package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"reflect"

	_ "context"
	_ "github.com/uptrace/bun"
	_ "strings"
)

func cobraErrorHandler(err error) {
	panic(err)
}

func NewListCommand(pt interface{}) *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   fmt.Sprintf("List %s", reflect.TypeOf(pt).Name()),
		Run: func(cmd *cobra.Command, args []string) {

			elemType := reflect.TypeOf(pt)
			elemSlice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
			elemSlicePtr := reflect.New(elemSlice.Type())

			err := DB.NewSelect().Model(elemSlicePtr.Interface()).Scan(cmd.Context())
			if err != nil {
				cobraErrorHandler(err)
				return
			}

			cols := []interface{}{}
			fields := []string{}
			format := []func(interface{}) interface{}{}

			meta := ParseStructMeta(elemType)

			for name, meta := range meta {

				cols = append(cols, meta.DisplayName)
				fields = append(fields, name)

				var format1 func(interface{}) interface{}

				/*
					for _, tagPart := range tagParts[1:] {
						if strings.HasPrefix(tagPart, "foreign=") {

							foreign := strings.TrimPrefix(tagPart, "foreign=")
							foreignParts := strings.Split(foreign, ":")
							foreignModel := foreignParts[0]
							foreignField := foreignParts[1]

							format1 = func(v interface{}) interface{} {

								var r string
								bun.NewRawQuery(DB,
									"SELECT " + foreignField + " FROM " + foreignModel + " WHERE id = ?",
									v,
								).Scan(context.Background(), &r)

								return r
							}
						}
					}
				*/

				format = append(format, format1)

			}

			t := table.New(cols...)
			t.WithHeaderFormatter(color.New(color.Bold).SprintfFunc())

			for i := 0; i < elemSlicePtr.Elem().Len(); i++ {
				row := []interface{}{}
				for j, field := range fields {
					c := elemSlicePtr.Elem().Index(i).FieldByName(field).Interface()

					if format[j] != nil {
						c = format[j](c)
					}

					row = append(row, c)
				}

				t.AddRow(row...)
			}

			t.Print()
		},
	}
}

func NewCreateCommand(pt interface{}) *cobra.Command {

	t := reflect.TypeOf(pt)
	v := reflect.New(t)

	var mapForeignToPk = map[string]func(k interface{}) string{}

	cc := &cobra.Command{
		Use:     "new",
		Aliases: []string{"create", "add"},
		Short:   fmt.Sprintf("Create %s", reflect.TypeOf(pt).Name()),
		Run: func(cmd *cobra.Command, args []string) {

			for i := 0; i < t.NumField(); i++ {
				m := mapForeignToPk[t.Field(i).Name]
				if m != nil {
					v.Elem().Field(i).SetString(m(v.Elem().Field(i).Interface()))
				}
			}

			_, err := DB.NewInsert().Model(v.Interface()).Exec(cmd.Context())
			if err != nil {
				cobraErrorHandler(err)
				return
			}
		},
	}

	meta := ParseStructMeta(t)

	for name, meta := range meta {

		field, _ := t.FieldByName(name)

		iif := v.Elem().Field(meta.structFieldIndex).Addr().Interface()

		if field.Type.Kind() == reflect.String {
			cc.Flags().StringVarP(iif.(*string), meta.DisplayName, meta.Short, "", meta.Description)
		} else if field.Type.Kind() == reflect.Int {
			cc.Flags().IntVarP(iif.(*int), meta.DisplayName, meta.Short, 0, meta.Description)
		} else if field.Type.Kind() == reflect.Bool {
			cc.Flags().BoolVarP(iif.(*bool), meta.DisplayName, meta.Short, false, meta.Description)
		} else if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.String {
			cc.Flags().StringSliceVarP(iif.(*[]string), meta.DisplayName, meta.Short, []string{}, meta.Description)
		} else {
			log.Warn(fmt.Sprintf("field %s has unsupported cli type %s.%s",
				field.Name, field.Type.PkgPath(), field.Type.Name()))
			continue
		}

		if meta.Required {
			cc.MarkFlagRequired(meta.DisplayName)
		}

		/*
			if meta.Foreign != nil {
				cc.RegisterFlagCompletionFunc(meta.DisplayName,
					func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
						var r []string
						DB.NewSelect().Model(meta.Foreign).Scan(cmd.Context(), &r)
						return r, cobra.ShellCompDirectiveNoFileComp
					},
				)
			}
		*/
	}

	return cc

}
