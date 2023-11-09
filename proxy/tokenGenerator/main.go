package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lautarojayat/backoffice/proxy"
)

var help string = `
Usage: genToken KEY_FOLDER [update | -v]

Examples:
- genToken $(pwd)/proxy/.tmp
	This will try to read jwt-sample.json in the root of the project and turn into a signed jwt

- genToken $(pwd)/proxy/.tmp update
	This will do the same as the previous command, but will update exp and iat fields so it is a valid token
	
- genToken $(pwd)/proxy/.tmp -v
	This will produce a more verbose output with instructions on what to do with the output
`

var verboseOutputTemplate string = `
Append the following token into the 'Auth' header before making a request:
		
%s
		
If you want to just get the token do something like:
	export MY_TOKEN_ENV=$(genToken $(pwd)/.tmp/proxy update)
`

func main() {
	var keyPath string

	if len(os.Args) > 1 {
		keyPath = os.Args[1]
	}

	if keyPath == "" {
		fmt.Printf("%s", help)
		return
	}

	overrideExp := false
	verbose := false

	if len(os.Args[2:]) > 0 {
		for _, arg := range os.Args[2:] {
			switch arg {
			case "update":
				if overrideExp {
					fmt.Printf("%s", help)
					return
				}
				overrideExp = true
			case "-v":
				if verbose {
					fmt.Printf("%s", help)
					return
				}
				verbose = true
			default:
				fmt.Printf("%s", help)
				return
			}
		}
	}

	privateKey, _ := proxy.ReadKeys(keyPath)
	path := "./jwt-sample.json"
	tokenBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("could not read json sample. error=%q", err)
	}

	claims := proxy.JWTStructure{}
	err = json.Unmarshal(tokenBytes, &claims)
	if err != nil {
		fmt.Printf("could not unmarshall token into struct. error=%q", err)
	}
	if overrideExp {
		claims.IssuedAt = jwt.NewNumericDate(time.Now())
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(100 * time.Hour))
	}

	token := proxy.CreateRSASignedToken(claims, privateKey)

	if !verbose {
		fmt.Printf("%s", token)
		return
	}

	fmt.Printf(verboseOutputTemplate, token)

}
