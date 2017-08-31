package jsonHelper

//oauth
type IdTokenRespBody struct{
    Aud     string
    Sub     string
    Name	string
}

type AccessTokenRespBody struct{
    Access_token    string
    Expires_in      float64
    Token_type      string
}

//response after user grants permissions
type OauthRespBody struct{
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