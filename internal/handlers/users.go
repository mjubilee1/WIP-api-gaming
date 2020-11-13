package handlers

import (
	"net/http"
	"fmt"
	"context"
	"strconv"
	"github.com/gin-gonic/gin"
	"api-gaming/internal/models"
	"github.com/jackc/pgx/v4"
	"api-gaming/internal/util"
	"api-gaming/internal/pkg/email"
	"github.com/dgrijalva/jwt-go"
	"github.com/shaj13/go-guardian/auth"
)

var (
	user = models.User{}
	refreshSecret = []byte(util.ViperEnvVariable("REFRESH_TOKEN_SECRET"))
	getVerifyCode = util.EncodeToString(6)
)

// Todo - Todo Struct
type Todo struct {
	UserID int64 `json:"userId"`
	Title string `json:"title"`
}

func verifyToken(ctx context.Context, r *http.Request, tokenString string) (auth.Info, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := auth.NewDefaultUser(claims["sub"].(string), "", nil, nil)
		return user, nil
	}

	return nil, fmt.Errorf("Invalid token")
}

func register(c *gin.Context) {
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	db, _ := c.Get("db")
	conn := db.(*pgx.Conn)
	err = user.Register(conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Store verification code in database
	storeEmailCode := user.StoreEmailCode(conn, getVerifyCode)
	if storeEmailCode != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	router.LoadHTMLGlob("web/templates/*")
	subject := "Verify Sleepless Gamer's Email"
	receiver := user.Email

	verifyEmail := email.NewRequest([]string{receiver}, subject)
	verifyEmail.Send(map[string]string{"username": user.Username, "verifyCode": getVerifyCode })
}

func login(c *gin.Context) {
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Invalid json provided": err.Error()})
		return
	}
	db, _ := c.Get("db")
	conn := db.(*pgx.Conn)
	validateUser := user.ValidateUser(conn)

	if validateUser != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	ts, err := user.CreateToken(user.ID)
	saveErr := user.SaveJWTMeta(user.ID, ts)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	tokens := map[string]string {
		"accessToken": ts.AccessToken,
		"refreshToken": ts.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}

func newToken(c *gin.Context) {
	tokens, err := user.CreateToken(user.ID)
	if err == nil {
		c.JSON(http.StatusOK, tokens)
		return
	}
	return
}

func refreshToken(c *gin.Context) {
	mapToken := map[string]string{}
  if err := c.ShouldBindJSON(&mapToken); err != nil {
     c.JSON(http.StatusUnprocessableEntity, err.Error())
     return
  }
  refreshToken := mapToken["refreshToken"]
  //verify the token
  token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
     //Make sure that the token method conform to "SigningMethodHMAC"
     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
     }
     return refreshSecret, nil
  })
	//if there is an error, the token must have expired
  if err != nil {
     c.JSON(http.StatusUnauthorized, "Refresh token expired")
     return
  }
  //is token valid?
  if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
     c.JSON(http.StatusUnauthorized, err)
     return
  }
  //Since token is valid, get the uuid:
  claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
  if ok && token.Valid {
     refreshUUID, ok := claims["refreshUUID"].(string) //convert the interface to string
     if !ok {
        c.JSON(http.StatusUnprocessableEntity, err)
        return
     }
		 userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["userID"]), 10, 64)
     if err != nil {
        c.JSON(http.StatusUnprocessableEntity, "Error occurred")
        return
     }
     //Delete the previous Refresh Token
		 deleted, delErr := user.DeleteAuth(refreshUUID)
     if delErr != nil &&  deleted == 0 { //if any goes wrong
        c.JSON(http.StatusUnauthorized, "unauthorized")
        return
     }
    //Create new pairs of refresh and access tokens
     ts, createErr := user.CreateToken(userID)
     if  createErr != nil {
        c.JSON(http.StatusForbidden, createErr.Error())
        return
     }
		//save the tokens metadata to redis
		saveErr := user.SaveJWTMeta(userID, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}
		tokens := map[string]string{
      "access_token":  ts.AccessToken,
  		"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
		} else {
			c.JSON(http.StatusUnauthorized, "refresh expired")
		}
	}

func createTodo(c *gin.Context) {
  var td *Todo
  if err := c.ShouldBindJSON(&td); err != nil {
     c.JSON(http.StatusUnprocessableEntity, "invalid json")
     return
	}

  tokenAuth, err := user.ExtractTokenMetadata(c.Request)
  if err != nil {
     c.JSON(http.StatusUnauthorized, "unauthorized")
     return
  }

	userID, err := user.FetchAuth(tokenAuth)

	if err != nil {
     c.JSON(http.StatusUnauthorized, "unauthorized")
     return
  }

	td.UserID = userID

	//you can proceed to save the Todo to a database
	//but we will just return it to the caller here:
  c.JSON(http.StatusCreated, td)
}

func logout(c *gin.Context) {
	au, err := user.ExtractTokenMetadata(c.Request)
  if err != nil {
     c.JSON(http.StatusUnauthorized, "unauthorized")
     return
  }
	deleted, delErr := user.DeleteAuth(au.AccessUUID)

	if delErr != nil || deleted == 0 { //if any goes wrong
     c.JSON(http.StatusUnauthorized, "unauthorized")
     return
  }
  c.JSON(http.StatusOK, "Successfully logged out")
}

func tokenAuthMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
     err := user.TokenValid(c.Request)
     if err != nil {
        c.JSON(http.StatusUnauthorized, err.Error())
        c.Abort()
        return
     }
     c.Next()
  }
}

func getUserMeta(c *gin.Context) {
	c.JSON(http.StatusOK, "Successfully get user meta")
}

func checkVerifiedEmail(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(*pgx.Conn)
	emailCode := user.IsEmailVerified(conn)

	c.JSON(http.StatusOK, gin.H{
		"verificationCode": emailCode,
	})
}

func addUserHandlers() {
	router.GET("/token", newToken)
	router.GET("/user-meta", tokenAuthMiddleware(), getUserMeta)
	router.POST("/token/refresh", refreshToken)
	router.POST("/register", register)
	router.GET("/verify-email", checkVerifiedEmail)
	router.POST("/login", login)
	router.POST("/token", login)
	router.POST("/todo", tokenAuthMiddleware(), createTodo)
	router.POST("/logout", tokenAuthMiddleware(), logout)
}