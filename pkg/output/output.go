package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kumarsaurabh/killport/pkg/discover"
)

func Text(w io.Writer, listeners []discover.Listener) {
	for _, l := range listeners {
		_, _ = fmt.Fprintf(w, "PID=%d USER=%s CMD=%s ADDR=%s\n", l.PID, l.User, l.Cmd, l.Address)
	}
}

func JSON(w io.Writer, listeners []discover.Listener) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(listeners)
}

func CSV(w io.Writer, listeners []discover.Listener) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	if err := writer.Write([]string{"pid", "user", "cmd", "address"}); err != nil {
		return err
	}
	for _, l := range listeners {
		if err := writer.Write([]string{fmt.Sprintf("%d", l.PID), l.User, l.Cmd, l.Address}); err != nil {
			return err
		}
	}
	return nil
}

func Format(w io.Writer, format string, listeners []discover.Listener) error {
	s := strings.ToLower(format)
	switch s {
	case "text":
		Text(w, listeners)
		return nil
	case "json":
		return JSON(w, listeners)
	case "csv":
		return CSV(w, listeners)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func PrintDefault(listeners []discover.Listener) {
	_ = Format(os.Stdout, "text", listeners)
}
