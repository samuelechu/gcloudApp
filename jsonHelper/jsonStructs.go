package jsonHelper

//oauth
type idTokenRespBody struct{
    Aud     string
    Sub     string
    Name	string
}

type accessTokenRespBody struct{
    Access_token    string
    Expires_in      float64
    Token_type      string
}

//response after user grants permissions
type oauthRespBody struct{
	Access_token    string
    Expires_in      float64
    Token_type      string
    Refresh_token 	string
    Id_token 		string
}

//cloudSQL
type User struct{
    Uid     string
    Name    string
}