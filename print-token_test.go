package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_tryParse(t *testing.T) {
	const testStr = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMDkxNTY2OTA3MTgzMjI5NDA1ODY5NDU0MTQyMDMwIiwibGlkIjoiUUFSVUkiLCJpc3MiOiJodHRwOi8vc3NvIiwiYXV0aG9yaXRpZXMiOlsiUk9MRV9UUlVTVEVEX0NMSUVOVCIsIlJPTEVfQ0xJRU5UIl0sImNsaWVudF9pZCI6ImNsb3VkOC1hY2NlcHRhbmNlIiwiYXVkIjoiY2xvdWQ4LWFjY2VwdGFuY2UiLCJzY2QiOiJhbnkiLCJzY29wZSI6WyJhbnkiXSwiZXhwIjoxNjEwNDcyOTQ1LCJpYXQiOjE2MTA0NzIzNDUsImp0aSI6Ijc3NzRlMDgwLWRkZmQtNGJiNS1iNzk1LTFjNjQ0ZGExMjcxZCIsImJybyI6ZmFsc2UsImNpZCI6IjEwOTUwMjY4Mzk4MTgyMDMyOTYwMTk4NTQxNDIwMzQifQ.dNyCW4qAC5g5CppRFE4OaUZrUWBMwfPD47Qxlu6cKslg0e7PTZF2MVz9O0NuqU8Pd7AoQb1XNcdxVVww4r4ByCsZFgF3Qi9DTTBC9izlO2kwiTo9vkGXVB-aug_O3_p0OqtvK4rhHrkslg7WySdmZAH_XYGOeZtN1BWxUi0kaayRr0fOeOU-lNdD7HbJNRXBC0P3uVUZXIuZ9CXiTJk6RwFPCpLgr8KwqwgDbnrIbJjQF0Vs1n0yBWFssZyfTGIOfRxKQZbRUPgdZZVUJRvpR3PWppDcI7JoFCdNln9PBuu1sOn0E-GDz7O89rpQ40DEVn3CQiucvwhsa5ZUFYYoVg`
	name, token, err := tryParse(testStr)
	require.NoError(t, err)
	assert.True(t, token.Valid)
	assert.EqualValues(t, "docker-sso", name)
	assert.Error(t, token.Claims.Valid()) // Token is expired.
	mapClaims := token.Claims.(jwt.MapClaims)
	assert.EqualValues(t, "any", mapClaims["scd"])
}
