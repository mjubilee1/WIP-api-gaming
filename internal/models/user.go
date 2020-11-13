package models

import (
	"api-gaming/internal/config"
	"api-gaming/internal/util"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	accessSecret = []byte(util.ViperEnvVariable("ACCESS_TOKEN_SECRET"))
 	refreshSecret = []byte(util.ViperEnvVariable("REFRESH_TOKEN_SECRET"))
)

// User is a representation of a user
type User struct {
	ID int64 `json:"userID"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	GamerTag string `json:"gamerTag"`
	Country string `json:"country"`
	DOB string `json:"dob"` // date of birth
	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
	IsVerified bool `json:"isVerified"`
	CreatedAt time.Time `json:"createdAt"`
}

// TokenDetails - House token definitions, their expiration periods and UUIDs.
type TokenDetails struct {
  AccessToken  string `json:"accessToken"`
  RefreshToken string `json:"refreshToken"`
  AccessUUID   string `json:"accessID"`
  RefreshUUID  string `json:"refreshID"`
  AtExpires    int64 `json:"atExpires"`
  RtExpires    int64 `json:"rtExpires"`
}

// AccessDetails - Contains the metadata that we will need to make a lookup in redis
type AccessDetails struct {
    AccessUUID string
    UserID   uint64
}

// EmailCode - Contains the email verification code that we will need to make a lookup in database.
type EmailCode struct {
	Code string
}

// Register user account
func (u *User) Register(conn *pgx.Conn) error {
	userLookup := User{}

	if len(u.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	u.Email = strings.ToLower(u.Email)
	getEmail := conn.QueryRow(context.Background(), "SELECT user_id from users WHERE email = $1", u.Email)
	err := getEmail.Scan(&userLookup.Email)
	if err != pgx.ErrNoRows {
		return fmt.Errorf("A user with the email already exists")
	}

	getUsername := conn.QueryRow(context.Background(), "SELECT user_id from users WHERE username = $1", u.Username)
	err = getUsername.Scan(&userLookup.Username)
	if err != pgx.ErrNoRows {
		return fmt.Errorf("A user with the username already exists")
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("there was an error creating the account")
	}

	u.Password = string(pwdHash)
	now := time.Now()
	_, err = conn.Exec(context.Background(), "INSERT INTO users (first_name,last_name,gamer_tag,country,is_verified,dob, email, username, password, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
	u.FirstName, u.LastName, u.GamerTag, u.Country, u.IsVerified, u.DOB, u.Email, u.Username, u.Password, now)

	return err
}

// ValidateUser - Checks to make sure password is correct and user is active
func (u *User) ValidateUser(conn *pgx.Conn) error {
	userLookup := User{}
	getUsername := conn.QueryRow(context.Background(), "SELECT username from users WHERE username = $1", u.Username)
	err := getUsername.Scan(&userLookup.Username)

	if err != nil {
		log.Println("Invalid username")
	}

	getPassword := conn.QueryRow(context.Background(), "SELECT password from users WHERE username = $1", u.Username)
	err = getPassword.Scan(&userLookup.Password)

	// Checks if database will return us a row with the request username provided.
	if err != nil {
		log.Println("A password was found.", err)
	}

	// Compares the request password with the database password and returns err.
	err = comparePasswords(userLookup.Password, []byte(u.Password))

	if err != nil {
		log.Println("Passwords does not match whats in database.", err)
	}

	return err
}

func comparePasswords(hashedPwd string, plainPwd []byte) error {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte (hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err
}

// StoreEmailCode - store email verification code in database
func (u *User) StoreEmailCode(conn *pgx.Conn, verifyCode string) error {
	_,err := conn.Exec(context.Background(), "INSERT INTO email_codes (username, code) VALUES ($1, $2)", u.Username, verifyCode)
	return err
}

// IsEmailVerified - check if email has been verified
func(u *User) IsEmailVerified(conn *pgx.Conn) string {
	emailCode := EmailCode{}
	getEmailCode := conn.QueryRow(context.Background(), "SELECT code from email_codes WHERE username = $1", u.Username)
	err := getEmailCode.Scan(&emailCode.Code)

	if err != nil {
			return "error"
	}

	_, err = conn.Exec(context.Background(), "UPDATE users SET is_verified = $1 WHERE username = $2", true, u.Username)

	if err != nil {
		log.Println("Error was thrown", err)
	}

	return emailCode.Code
}

// CreateToken - returns the jwt token to be used
func (u *User) CreateToken(userID int64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix() // Access Token expires after 15 minutes
	td.AccessUUID = uuid.NewV4().String()

  td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() // Refresh Token expires after 7 days
	td.RefreshUUID = td.AccessUUID + "++" + strconv.Itoa(int(userID))

	var err error
	// Creating Access Token
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["userID"] = userID
	accessTokenClaims["accessUUID"] = td.AccessUUID
	accessTokenClaims["exp"] = td.AtExpires
	atc := jwt.NewWithClaims(jwt.SigningMethodHS256, &accessTokenClaims)
	td.AccessToken, err = atc.SignedString(accessSecret)

	if err != nil {
		return nil, err
	}


	// Creating Refresh Token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["userID"] = userID
	refreshTokenClaims["refreshUUID"] = td.RefreshUUID
	refreshTokenClaims["exp"] = td.RtExpires
	rtc := jwt.NewWithClaims(jwt.SigningMethodHS256, &refreshTokenClaims)
	td.RefreshToken, err = rtc.SignedString(refreshSecret)

	if err != nil {
		return nil, err
	}

	return td, nil
}

// SaveJWTMeta - Used to save the JWTs metadata
func (u *User) SaveJWTMeta(userID int64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) // Converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	var redisClient = config.RedisConn()

	// Importantly, use defer and the connection's Close() method to
	// ensure that the connection is always returned to the pool before
	// Method() exits.
	errAccessToken := redisClient.Set(context.Background(),td.AccessUUID, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if errAccessToken != nil {
		return errAccessToken
	}
	errRefreshToken := redisClient.Set(context.Background(),td.RefreshUUID, strconv.Itoa(int(userID)), rt.Sub(now)).Err()
	if errRefreshToken != nil {
		return errRefreshToken
	}

	return nil
}

// ExtractToken - Extract the token from the request header
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")

  //normally Authorization the_token_xxx
  strArr := strings.Split(bearToken, " ")
  if len(strArr) == 2 {
     return strArr[1]
  }
  return ""
}

// ExtractTokenMetadata - Extract the token metadata that will lookup in our redis store,
// to extract the token
func(u *User) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
  token, err := VerifyToken(r)
  if err != nil {
     return nil, err
  }
  claims, ok := token.Claims.(jwt.MapClaims)
  if ok && token.Valid {
     accessUUID, ok := claims["accessUUID"].(string)
     if !ok {
        return nil, err
     }
     userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["userID"]), 10, 64)
     if err != nil {
        return nil, err
		 }
     return &AccessDetails{
        AccessUUID: accessUUID,
        UserID:   userID,
     }, nil
	}
  return nil, err
}

// VerifyToken - Get the token string, then proceed to check the signing method.
func VerifyToken(r *http.Request) (*jwt.Token, error) {
  tokenString := ExtractToken(r)
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
     //Make sure that the token method conform to "SigningMethodHMAC"
     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
     }
     return accessSecret, nil
  })
  if err != nil {
     return nil, err
  }
  return token, nil
}

// TokenValid - Then we check the validity of this token, whether it is still useful or it has expired.
func(u *User) TokenValid(r *http.Request) error {
  token, err := VerifyToken(r)
  if err != nil {
     return err
  }
  if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
     return err
  }
  return nil
}

// FetchAuth - Accepts the AccessDetails, then looks it up in redis.
// If the record is not found, it may mean the token has expired,
// hence an error is thrown.
func(u *User) FetchAuth(authD *AccessDetails) (int64, error) {
	var redisClient = config.RedisConn()

  userid, err := redisClient.Get(context.Background(), authD.AccessUUID).Result()
  if err != nil {
     return 0, err
  }
  userID, _ := strconv.ParseInt(userid, 10, 64)
  return userID, nil
}

// DeleteAuth - Enables us delete a JWT metadata from redis
func(u *User) DeleteAuth(givenUUID string) (int64,error) {
	var redisClient = config.RedisConn()

  deleted, err := redisClient.Del(context.Background(),givenUUID).Result()

	if err != nil {
     return 0, err
  }
  return deleted, nil
}
