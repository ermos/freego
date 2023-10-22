package command

import (
	"fmt"
	"github.com/ermos/freego/internal/pkg/config"
	"github.com/ermos/freego/internal/pkg/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

type List struct{}

func (l *List) Flags(cmd *cobra.Command) {}

func (l *List) Execute(cmd *cobra.Command, args []string) error {
	t := table.NewWriter()
	t.SetStyle(table.Style{
		Box: table.BoxStyle{
			BottomLeft:       " ",
			BottomRight:      " ",
			BottomSeparator:  " ",
			EmptySeparator:   text.RepeatAndTrim(" ", text.RuneWidthWithoutEscSequences(" ")),
			Left:             " ",
			LeftSeparator:    " ",
			MiddleHorizontal: " ",
			MiddleSeparator:  " ",
			MiddleVertical:   " ",
			PaddingLeft:      " ",
			PaddingRight:     " ",
			PageSeparator:    "\n",
			Right:            " ",
			RightSeparator:   " ",
			TopLeft:          " ",
			TopRight:         " ",
			TopSeparator:     " ",
			UnfinishedRow:    " ~",
		},
	})

	t.AppendHeader(table.Row{"DOMAIN ID", "DOMAIN", "HOST", "PORT", "CREATED"})
	for id, item := range config.GetActiveDomains() {
		t.AppendRow(table.Row{
			id,
			item.Domain,
			item.Host,
			item.Port,
			util.FormatXTimeAgo(item.CreatedAt, "never"),
		})
	}

	fmt.Println(t.Render())

	return nil
}
