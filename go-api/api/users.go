package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/gin-gonic/gin"
	// "google.golang.org/appengine/log"
)

func newUserResponse(user db.User) (userResponse, error) {
	var stock []map[string]interface{}
	if user.Stock != nil {
		if unerr := json.Unmarshal(user.Stock, &stock); unerr != nil {
			return userResponse{}, unerr
		}
	}
	return userResponse{
		UserID:      user.UserID,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
		Stock:       stock,
	}, nil
}

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Role        string `json:"role" binding:"required"`
}

type userResponse struct {
	UserID      int64                    `json:"id" binding:"required"`
	Username    string                   `json:"username" binding:"required"`
	Email       string                   `json:"email" binding:"required,email"`
	PhoneNumber string                   `json:"phone_number" binding:"required"`
	Address     string                   `json:"address" binding:"required"`
	Stock       []map[string]interface{} `json:"stock"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var arg db.CreateUserParams
	lenDB, _ := server.store.ListUser(ctx)
	if len(lenDB) == 0 {
		hashPassword, err := utils.GeneratePasswordHash(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		arg = db.CreateUserParams{
			Username:    req.Username,
			Password:    hashPassword,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			Role:        "admin",
			Address:     req.Address,
			Stock:       nil,
		}

	} else {
		pass, _ := utils.GeneratePasswordHash("beforeUpdate")
		arg = db.CreateUserParams{
			Username:    req.Username,
			Password:    pass,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			Role:        req.Role,
			Address:     req.Address,
			Stock:       nil,
		}
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp, _ := newUserResponse(user)

	ctx.JSON(http.StatusOK, resp)
	return
}

type userLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userLoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req userLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var ppl userResponse
	if err := utils.CheckPassword("beforeUpdate", user.Password); err == nil {

		hashPassword, hashErr := utils.GeneratePasswordHash(req.Password)
		if hashErr != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		updatedUser, err := server.store.UpdateUserPasswordFisrtLogin(ctx, db.UpdateUserPasswordFisrtLoginParams{
			UserID:   user.UserID,
			Password: hashPassword,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ppl, _ = newUserResponse(updatedUser)
	} else {

		err = utils.CheckPassword(req.Password, user.Password)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ppl, _ = newUserResponse(user)
	}

	accesToken, err := server.tokenMaker.CreateToken(ppl.Username, server.config.PYTHON_APP_TOKEN_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := userLoginResponse{
		AccessToken: accesToken,
		User:        ppl,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type GetUSerRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req GetUSerRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp, _ := newUserResponse(user)

	ctx.JSON(http.StatusOK, resp)
}

type DeleteUSerRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req DeleteUSerRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteUser(ctx, int64(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Message": "User Deleted Succesfully"})
	return
}

func (server *Server) listUsers(ctx *gin.Context) {
	list_user, err := server.store.ListUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var updateUser []userResponse
	for _, user := range list_user {
		us, _ := newUserResponse(user)
		updateUser = append(updateUser, us)
	}
	ctx.JSON(http.StatusOK, updateUser)
	return
}

type UserToEditUri struct {
	ID int64 `uri:"id" binding:"required"`
}

type UserToEditRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
	Role        string `json:"role" binding:"required,oneof=admin client"`
}

func (server *Server) editUser(ctx *gin.Context) {
	var uri UserToEditUri
	var req UserToEditRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(req.OldPassword, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	hashPassword, hasherr := utils.GeneratePasswordHash(req.NewPassword)
	if hasherr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(hasherr))
		return
	}

	_, err = server.store.EditUserTx(ctx, db.EditUserParams{
		UserID:      uri.ID,
		Password:    hashPassword,
		Role:        req.Role,
		Email:       "",
		Address:     "",
		Username:    "",
		PhoneNumber: "",
	})

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "user update succeful"})
	return
}

type UserToManageUri struct {
	ID int64 `uri:"id" binding:"required"`
}

type UserToManageRequest struct {
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Role        string `json:"role" binding:"required,oneof=client admin"`
}

func (server *Server) manageUser(ctx *gin.Context) {
	var uri UserToManageUri
	var req UserToManageRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.EditUserTx(ctx, db.EditUserParams{
		UserID:      uri.ID,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Username:    req.Username,
		Role:        req.Role,
		Password:    "",
	})

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "user update succeful"})
	return
}

type addAdminStockURIRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type addAdminStockJSONRequest struct {
	UserID   int64 `json:"user_id" binding:"required,min=1"`
	Quantity int64 `json:"quantity" binding:"required,min=1"`
}

func (server *Server) addAdminStock(ctx *gin.Context) {
	var uri addAdminStockURIRequest
	var req addAdminStockJSONRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := server.store.GetProduct(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	admin, err := server.store.GetUserForUpdate(ctx, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updatedAdmin, err := server.store.AddAdminStockTx(ctx, db.AddAdminStockParams{
		Admin:       admin,
		ProducToAdd: product,
		Amount:      req.Quantity,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedAdmin)
	return
}

type searchUser struct {
	SearchWord string `json:"search_word" binding:"required"`
}

func (server *Server) searchUsers(ctx *gin.Context) {
	var req searchUser

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rst, err := server.store.SearchILikeUsers(ctx, req.SearchWord)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// var updateUser []userResponse
	// for _, user := range rst {
	// 	us, _ := newUserResponse(user)
	// 	updateUser = append(updateUser, us)
	// }
	ctx.JSON(http.StatusOK, rst)
	return

}

type addClientStockURIRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type addClientStockJSONRequest struct {
	ProductsID []int64 `json:"products_id" binding:"required"`
	Quantities []int64 `json:"quantities" binding:"required"`
}

func (server *Server) addClientStock(ctx *gin.Context) {
	var uri addClientStockURIRequest
	var req addClientStockJSONRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	admin, err := server.store.GetUserForUpdate(ctx, 1) // add manually the admins_id
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserForUpdate(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var newProducts []db.Product
	for _, id := range req.ProductsID {
		addProduct, err := server.store.GetProduct(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		newProducts = append(newProducts, addProduct)
	}
	updatedData, err := server.store.AddClientStockTx(ctx, db.AddClientStockParams{
		FromAdmin:   admin,
		ToClient:    user,
		ProducToAdd: newProducts,
		Amount:      req.Quantities,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedData)
	return
}

type reduceClientStockURIRequest struct {
	ID int64 `uri:"id" biding:"required"`
}

type reduceClientStockJSONRequest struct {
	ProductsID []int8 `json:"products_id" biding:"required"`
	Quantities []int8 `json:"quantities" biding:"required"`
}

func (server *Server) reduceClientStock(ctx *gin.Context) {
	var uri reduceClientStockURIRequest
	var req reduceClientStockJSONRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var amount int
	for idx, id := range req.ProductsID {
		removeProduct, err := server.store.GetProduct(ctx, int64(id))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		amount += int(removeProduct.UnitPrice) * int(req.Quantities[idx])
	}

	user, err := server.store.GetUser(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	log.Println("Sending STK push")
	trasactionID, err := utils.SendSTK(strconv.Itoa(amount), user.UserID, user.PhoneNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	jsonSoldProduct, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.CreateTransaction(ctx, db.CreateTransactionParams{
		TransactionID: trasactionID,
		Amount:        int32(amount),
		DataSold:      jsonSoldProduct,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"successful": "STK push success"})
	return
}

type mpesaCallbackRequest struct {
	TransactionID string `uri:"id" binding:"required"`
}

func (server *Server) mpesaCallback(ctx *gin.Context) {
	log.Println("Callback Successful")
	var req mpesaCallbackRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := req.TransactionID[len(req.TransactionID)-3:]
	transactionId := req.TransactionID[:len(req.TransactionID)-3]

	intUserID, _ := strconv.Atoi(userId)
	user, err := server.store.GetUserForUpdate(ctx, int64(intUserID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	transaction, err := server.store.GetTransaction(ctx, transactionId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	processMpesaCallbackData(ctx, server, user, transaction)

	ctx.JSON(http.StatusOK, gin.H{"Successful": "Reached"})
	return
}

func processMpesaCallbackData(ctx *gin.Context, server *Server, user db.User, transaction db.Transaction) {
	log.Println("Processing Mpesa Callback")

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		redirectToPythonApp(user, transaction, err)
		return
	}

	var callbackBody map[string]interface{}
	err = json.Unmarshal(body, &callbackBody)
	if err != nil {
		redirectToPythonApp(user, transaction, err)
		return
	}

	log.Println(callbackBody)

	bodyValue := callbackBody["Body"].(map[string]interface{})
	stkCallbackValue := bodyValue["stkCallback"].(map[string]interface{})
	metaData := stkCallbackValue["CallbackMetadata"].(map[string]interface{})
	items := metaData["Item"].([]map[string]interface{})
	resultCode := stkCallbackValue["ResultCode"].(int)
	phoneNumber := items[3]["Value"].(string)
	mpesaReceiptNumber := items[1]["Value"].(string)
	_ = items[0]["Value"].(float64)
	resultDesc := stkCallbackValue["ResultDesc"].(string)

	_, err = server.store.UpdateTransaction(ctx, db.UpdateTransactionParams{
		TransactionID:      transaction.TransactionID,
		MpesaReceiptNumber: mpesaReceiptNumber,
		PhoneNumber:        phoneNumber,
	})
	if err != nil {
		redirectToPythonApp(user, transaction, err)
		return
	}

	if resultCode != 0 {
		redirectToPythonApp(user, transaction, fmt.Errorf(resultDesc))
		return
	}

	var data map[string][]int8
	if unerr := json.Unmarshal(transaction.DataSold, &data); unerr != nil {
		redirectToPythonApp(user, transaction, err)
		return
	}

	var newProducts []db.Product
	for _, id := range data["products_id"] {
		addProduct, err := server.store.GetProduct(ctx, int64(id))
		if err != nil {
			if err == sql.ErrNoRows {
				redirectToPythonApp(user, transaction, err)
				return
			}
			redirectToPythonApp(user, transaction, err)
			return
		}

		newProducts = append(newProducts, addProduct)
	}
	_, err = server.store.ReduceClientStockTx(ctx, db.ReduceClientStockParams{
		Client:         user,
		ProducToReduce: newProducts,
		Amount:         data["quantities"],
		Transaction:    transaction,
	})
	if err != nil {
		redirectToPythonApp(user, transaction, err)
		return
	}

	server.store.ChangeStatus(ctx, db.ChangeStatusParams{
		TransactionID: transaction.TransactionID,
		Status:        true,
	})
}

func redirectToPythonApp(user db.User, transaction db.Transaction, passErr error) {
	pythonEndpoint := "http://0.0.0.0:3030/notify"

	data := gin.H{
		"status":        passErr,
		"user_id":       user.UserID,
		"transactionID": transaction.TransactionID,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	_, err = http.Post(pythonEndpoint, "application/json", bytes.NewBuffer(jsonData))
	fmt.Println(err)
	return
}

type resetPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

func (server *Server) resetPassword(ctx *gin.Context) {
	var req resetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, (10 * time.Minute))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resetPasswordLink := fmt.Sprintf("%v/resetit?token=%v", server.config.PUBLIC_URL, accessToken) // URL + TOKEN for passwordreset
	emailBody := fmt.Sprintf(`
	<h1>Hello %s</h1>
	<p>We received a request to reset your password. Click the link below to reset it:</p>
	<a href="%s" style="display:inline-block; padding:10px 20px; background-color:#007BFF; color:#fff; text-decoration:none; border-radius:5px;">Reset Password</a>
	<h5>The link is valid for 10 Minutes</h5>
`, user.Username, resetPasswordLink)

	err = server.emailSender.SendMail("Reset Password", emailBody, []string{user.Email}, nil, nil, "", nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "accesstoken granted and email send"})
	return
}

type resetItQueryRequest struct {
	Token string `form:"token" binding:"required"`
}

type resetItJSONRequest struct {
	Password string `json:"password" binding:"required"`
}

func (server *Server) resetIt(ctx *gin.Context) {
	var req resetItQueryRequest
	var pass resetItJSONRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&pass); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(req.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	hashPassword, hasherr := utils.GeneratePasswordHash(pass.Password)
	if hasherr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(hasherr))
		return
	}

	user, err := server.store.GetUserByUsename(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user.Password = hashPassword
	_, err = server.store.UpdateUserCredentials(ctx, db.UpdateUserCredentialsParams{
		UserID:      user.UserID,
		Email:       user.Email,
		Password:    hashPassword,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		Username:    user.Username,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"successful": "User password changed successful"})
}
