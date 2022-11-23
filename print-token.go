package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var (
	knownKeys = map[string]string{
		// Note that docker-sso uses same keypair as pilot.
		//"docker-sso": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjZfkt7DNOEkg3b7DmMcBAQQqJmxM+/fZhf4tJUyVWShncsYRiKSPzogtYE4U/U33YZFgLkJphmXI6grUhtIy3bgUt6DyTq0SR08i56Za5trjKsnuA37h6FrH+eqUQMbJdjRD2M6jOQHJJ59jNMjYSdDE9V65Vd1feRbhzoYA2i9v2D+zifBGq7c8LeULUjufuwDgcUC3ajzeLL7TkO4qNIjpaQNOkhe/cJ5jlUH0f2G8FRm8YWbYUydw5mcBSLENoSRnup2TDbW5TgthGIiAEkE5xUyDBgtFrY01WJhRFo9BBcRF3a1sHvryTQi0mZAPNlCTjZzwEtaHce90pdL2UwIDAQAB\n-----END PUBLIC KEY-----",
		"dev":   "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxKRuyp5dmfKR4CdJpYBvM3r9jenMWHShPSSwRgvBYETds5hE/lFv20vDem9Xd+5tDQvR/rQGD/ZWpFPKad9XBIP2gDPKtbdpMzrSc0CZpXz82rN2/U2aj2A6OhDpJ9b2be1uVrTTO32HOqorR7kmCZP26r/ftxPDtR1S0ewyWbvOffDaH6/zmDz5zdId4aWwWNBtbrazhseH7UY2GBcR1xJcESP+mdHfwkkTadzK7ONyF2MVhxWySr6I//WGkRgHubQHd3cQd9Wtv3jnAsZU1up9s5xV+mfeSmZzamV531Y+DV+1W9zqrUdeRqxTpo+90K+uUjFE3NKaKB7jnMGm0QIDAQAB\n-----END PUBLIC KEY-----",
		"pilot": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjZfkt7DNOEkg3b7DmMcBAQQqJmxM+/fZhf4tJUyVWShncsYRiKSPzogtYE4U/U33YZFgLkJphmXI6grUhtIy3bgUt6DyTq0SR08i56Za5trjKsnuA37h6FrH+eqUQMbJdjRD2M6jOQHJJ59jNMjYSdDE9V65Vd1feRbhzoYA2i9v2D+zifBGq7c8LeULUjufuwDgcUC3ajzeLL7TkO4qNIjpaQNOkhe/cJ5jlUH0f2G8FRm8YWbYUydw5mcBSLENoSRnup2TDbW5TgthGIiAEkE5xUyDBgtFrY01WJhRFo9BBcRF3a1sHvryTQi0mZAPNlCTjZzwEtaHce90pdL2UwIDAQAB\n-----END PUBLIC KEY-----",
		"local": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjGDYGz//GZ0xVJ0Ol5hr7xXBfSCBqVxTWeFKmE/nj84ALdlJowTGagflfcTg5KPYqUumFzcyGP5RCN53MUdmMTJ8Ja6HcDq4RNhDgYnRCBHDkDj+ZorOjNUFM9u0h2d/UnN5V18qK+Ckxlz87vmYS1qWwXPD2aPcX6Oxo2HELlk/cMs94HHHurFAGXxz8AUQ1IAN959JDVErpv8757quJfOHOY13PDsCPghvER/nu+0xBy04l5DdHlDTu6bE37acxYpfv81SX2xm+N72i6oyoBFKOuAnI3dj1snlW5x3jCNBlQu3bLzesVLRJh43ZdmwpG2i7pq9Y0sFYelGPLhz6QIDAQAB\n-----END PUBLIC KEY-----",
		"prod":  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArBPEBmeY1bUdOCPR6p6UjK3MbvkTmwRv1BTUIjopvQubVpP547m9UbSkyJBOgzKqeaSlhd/6Jte14fOBNC387AfVpql2h545tdv/dvETGWUB10lhyhQoHCP+aMA/c+BivjLE7CK1Z4J7tbstN5b4RfF6reEw0SCf6IEW8sz/TiKK172TJW9aOedNiM4R1vCbnH2S6j/JSfg4TYiYXSH1+MZPYUyNpGPNp5mm1Y5YCmeWKrQkh76JG6U9abkx5OLOX0Q5cKRmVaxevvCSi87K1Ceb5Zy4XRTNJxbUweC4CM/fLwRMTO8TSa0o1oDmvM2ke+dpTXDQ4V9szk6gnnsEVwIDAQAB\n-----END PUBLIC KEY-----",
		"box8":  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0dyoadRxVEZ/FYPg5nvZFgo0dpZnQnDenjf+KKybIihn++2XB/MLOAXLABn9ddEQCAljpVH+hBKVjGtRlWzwFRlrP5dYD6qMd2bpMITNZbHE/6LzPq8JE+wXQib1aqh/GJ70iAR86Pui6OJq3Z1rcbiWE47P62VFUcw36qvLJrXNnxnVIcG1AB6yrNpnY4sPwmnvDQPFY77jNAd9xjl12lLvmVre8y764amjtODSB85zkOrLtl1SvN0xB/Jsf8B8P1PciH9RiGaSYNmQ0CzrDvzTSYosKy9TVwSkBM9fzbR1PPf/PIsEAyXP62EZT5X5B3vUXLAgxx1ywV9eonMEbwIDAQAB\n-----END PUBLIC KEY-----",
	}
	parser = jwt.Parser{
		SkipClaimsValidation: true,
	}
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		keyName, token, err := tryParse(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse %+q: %s\n", scanner.Text(), err)
			continue
		}
		if keyName != "" {
			fmt.Fprintf(os.Stderr, "Signed by '%s'\n", keyName)
		}
		fmt.Fprintln(os.Stderr, "Valid:", token.Valid)

		b, err := json.MarshalIndent(token.Claims, "", "  ")
		abortIfError(err)
		fmt.Println(string(b))
	}
}

func tryParse(text string) (string, *jwt.Token, error) {
	for name, keyText := range knownKeys {
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyText))
		abortIfError(err)

		token, err := parser.Parse(text, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err == nil {
			return name, token, nil
		}
	}

	token, _, err := parser.ParseUnverified(text, jwt.MapClaims{})
	return "", token, err
}

func decode(b []byte) string {
	r := base64.NewDecoder(base64.RawURLEncoding, bytes.NewReader(b))
	text, err := ioutil.ReadAll(r)
	abortIfError(err)

	var decoded map[string]interface{}
	abortIfError(json.Unmarshal(text, &decoded))
	text, err = json.MarshalIndent(&decoded, "", "  ")
	abortIfError(err)
	return string(text)
}

func abortIfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
