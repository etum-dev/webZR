package scan

import(
	"context"

	"github.com/joho/godotenv"
	"github.com/shadowscatcher/shodan"
	"github.com/shadowscatcher/shodan/search"

	"log"
	"net/http"
	"os"
)

func AuthShodan() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("where teh FUCK is the env file")
	}
	shodankey := os.Getenv("SHODAN_API_KEY")
	return shodankey
}

func SearchShodan(){
	shodan_env := AuthShodan()
	//TODO: exclude honeypot tags
	// we have to filter this after results unless enterprise user.
	shodanQuery := search.Params{
		Page:1,
		Query: search.Query{
			All: "101",
		},
	}
	

	client, _ := shodan.GetClient(shodan_env, http.DefaultClient, true)
	ctx := context.Background()
	result, err := client.Search(ctx, shodanQuery)
	
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)


}