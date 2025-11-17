package spam

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/DavidHoenisch/oxidation9/internal/models"
)

func Help() {

}

func Run(params models.FuncParams) error {
	url := params.Args["url"].(string)
	count := params.Args["count"].(string)

	cleanedCount, err := strconv.Atoi(count)
	if err != nil {
		return err
	}

	for x := range cleanedCount {
		fmt.Printf("================= REQUEST [%d] ================= \n", x)
		data, err := http.Get(url)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		for name, values := range data.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)

			}
		}
	}

	return nil
}
