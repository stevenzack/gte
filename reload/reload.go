package reload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/urfave/cli"
)

func ApiCommand(c *cli.Context) error {
	port := c.Int("p")
	if port <= 0 {
		return fmt.Errorf("the port option '-p' is not set")
	}

	r, e := http.NewRequest(http.MethodOptions, "http://localhost:"+strconv.Itoa(port)+"/reload", nil)
	if e != nil {
		log.Println(e)
		return e
	}
	res, e := http.DefaultClient.Do(r)
	if e != nil {
		log.Println(e)
		return e
	}
	defer res.Body.Close()
	b, e := io.ReadAll(res.Body)
	if e != nil {
		log.Println(e)
		return e
	}

	fmt.Println(res.StatusCode, "\t", string(b))
	return nil
}
