package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/EmilioCliff/inventory-system/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
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
	UserID          int64                    `json:"id" binding:"required"`
	Username        string                   `json:"username" binding:"required"`
	Email           string                   `json:"email" binding:"required,email"`
	PhoneNumber     string                   `json:"phone_number" binding:"required"`
	Address         string                   `json:"address" binding:"required"`
	Stock           []map[string]interface{} `json:"stock"`
	StockValue      int64                    `json:"stock_value,omitempty"`
	AdminStockValue float64                  `json:"admin_stock_value,omitempty"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var arg db.CreateUserParams
	lenDB, _ := server.store.ListUser(ctx, db.ListUserParams{
		Limit:  5,
		Offset: 0,
	})
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

	// user, err := server.store.CreateUser(ctx, arg)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }
	createUserResult, err := server.store.CreateUserTx(ctx, db.CreateUserTxParams{
		CreateUserParams: arg,
		AfterCreate: func(user db.User) error {
			sendPayload := &worker.SendEmailVerifyPayload{
				Username: user.Username,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(5 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, *sendPayload, opts...)
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := server.redis.Del(ctx, ListUsers+fmt.Sprintf(":1")).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp, _ := newUserResponse(createUserResult.User)

	ctx.JSON(http.StatusOK, resp)
	return
}

type calculatePriceRequest struct {
	ProductIDs []int64 `json:"product_ids" binding:"required"`
	Quantities []int64 `json:"quantities" binding:"required"`
}

type calculatePriceResponse struct {
	TotalAmount int64 `json:"total_amount" binding:"required"`
}

func (server *Server) calculatePrice(ctx *gin.Context) {
	var req calculatePriceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var amount int64
	for idx, productID := range req.ProductIDs {
		unitPrice, err := server.store.GetProductPrice(ctx, productID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		amount += int64(unitPrice) * req.Quantities[idx]
	}

	ctx.JSON(http.StatusOK, calculatePriceResponse{
		TotalAmount: amount,
	})
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

	var value db.StockValue
	if req.ID == 1 {
		count, err := server.store.TotalStockValue(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		resp.StockValue = count
		// also add the admins total amount of money
		totalAdminStock := 0.0
		for _, d := range resp.Stock {
			productID := d["productID"].(float64)
			quantity := d["productQuantity"].(float64)

			product, err := server.store.GetProduct(ctx, int64(productID))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			totalAdminStock += float64(product.UnitPrice) * quantity
		}

		resp.AdminStockValue = totalAdminStock

	} else {
		value, err = server.store.GetUserStockValue(ctx, int32(req.ID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		resp.StockValue = value.Value
	}

	err = server.setCache(ctx, GetUser+fmt.Sprintf("%v", req.ID), resp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, resp)
	return
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

	err := server.store.DeleteUserTx(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// err := server.store.DeleteUser(ctx, int64(req.ID))
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		ctx.JSON(http.StatusNotFound, errorResponse(err))
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	if err := server.redis.Del(ctx, ListUsers+fmt.Sprintf(":1")).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Message": "User Deleted Succesfully"})
	return
}

type listUserRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
}

type listUserResponse struct {
	Data     []userResponse     `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("%v:%v", ctx.Request.URL.Path, req.PageID)
	cacheData, err := server.redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		log.Info().Msgf("cached hit for: %v", cacheKey)
		ctx.Data(http.StatusOK, "application/json", cacheData)
		return
	}

	list_user, err := server.store.ListUser(ctx, db.ListUserParams{
		Limit:  PageSize,
		Offset: (req.PageID - 1) * PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var updateUser []userResponse
	for _, user := range list_user {
		us, _ := newUserResponse(user)
		updateUser = append(updateUser, us)
	}

	totalUser, err := server.store.CountUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := totalUser / int64(PageSize)
	if totalUser/int64(PageSize) != 0 {
		totalPages++
	}

	rsp := listUserResponse{
		Data: updateUser,
		Metadata: PaginationMetadata{
			CurrentPage: req.PageID,
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalUser),
		},
	}

	err = server.setCache(ctx, cacheKey, rsp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
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

	if err := server.redis.Del(ctx, GetUser+fmt.Sprintf("%v", uri.ID)).Err(); err != nil {
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

	if err := server.redis.Del(ctx, GetUser+fmt.Sprintf("%v", uri.ID)).Err(); err != nil {
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
	Quantity int64 `json:"quantity" binding:"required"`
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

	if err := server.redis.Del(ctx, GetUser+fmt.Sprintf("1")).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

type requstStockUri struct {
	UserID int32 `uri:"id" binding:"required"`
}

type requestStockQuery struct {
	Products  []int32 `json:"products" binding:"required"`
	Quantites []int32 `json:"quantities" binding:"required"`
}

func (server *Server) requestStock(ctx *gin.Context) {
	var req requstStockUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqQuery requestStockQuery
	if err := ctx.ShouldBindJSON(&reqQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	reqPayload := worker.RequestStockPayload{
		UsernameID: req.UserID,
		Products:   reqQuery.Products,
		Quantities: reqQuery.Quantites,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(5 * time.Second),
		asynq.Queue(worker.QueueDefault),
	}

	if err := server.taskDistributor.DistributeSendRequestToAdmin(ctx, reqPayload, opts...); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "send_successfully")
	return
}

type searchUser struct {
	SearchWord string `form:"search_word" binding:"required"`
}

func (server *Server) searchUser(ctx *gin.Context) {
	var req searchUser

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var pgQuery pgtype.Text
	if err := pgQuery.Scan(req.SearchWord); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rst, err := server.store.SearchILikeUsers(ctx, pgQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rst)
	return
}

type addClientStockURIRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type addClientStockJSONRequest struct {
	ProductsID  []int64 `json:"products_id" binding:"required"`
	Quantities  []int64 `json:"quantities" binding:"required"`
	InvoiceDate string  `json:"invoice_date" binding:"required"`
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

	admin, err := server.store.GetUser(ctx, 1) // add manually the admins_id
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var adminProducts []map[string]interface{}
	if unerr := json.Unmarshal(admin.Stock, &adminProducts); unerr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

	var newProducts []db.Product
	for idx, id := range req.ProductsID {
		log.Info().Int("productId", int(id)).Int("quantity", int(req.Quantities[idx])).Msg("add client stock log")
		addProduct, err := server.store.GetProduct(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// check if admin has enough stock to sell
		for _, adminProduct := range adminProducts {
			if idAdmin, ok := adminProduct["productID"].(float64); ok {
				idInt := int64(idAdmin)
				if idInt == id {
					quantityInt := adminProduct["productQuantity"].(float64)
					// quantityInt := quantityFloat
					if int64(quantityInt)-req.Quantities[idx] < 0 {
						ctx.JSON(http.StatusNotAcceptable, gin.H{"error": fmt.Sprintf("Not enough in stock to distribute: %s = %v", addProduct.ProductName, req.Quantities[idx])})
						return
						// return fmt.Errorf("Not enough in inventory %v - %v to sell %v", adminProduct["productName"], adminProduct["productQuantity"], arg.Amount[index])
					}
				}
			}
		}

		newProducts = append(newProducts, addProduct)
	}

	invoiceDate, err := time.Parse("2006-01-02", req.InvoiceDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updatedData, err := server.store.AddClientStockTx(ctx, db.AddClientStockParams{
		FromAdmin:   admin,
		ToClient:    user,
		ProducToAdd: newProducts,
		Amount:      req.Quantities,
		AfterProcess: func(data []map[string]interface{}) error {
			invoiceTaskPayload := &worker.GenerateInvoiceAndSendEmailPayload{
				User: user,
				// Products:    newProducts,
				// Amount:      req.Quantities,
				InvoiceDate: invoiceDate,
				InvoiceData: data,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.Queue(worker.QueueDefault),
			}

			return server.taskDistributor.DistributeGenerateAndSendInvoice(ctx, *invoiceTaskPayload, opts...)

		},
	})
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, errorResponse(err))
		return
	}

	if err := server.redis.Del(ctx, GetUser+fmt.Sprintf("%v", uri.ID)).Err(); err != nil {
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
	// Amount     int64   `json:"amount" biding:"required"`
	ProductsID []int32 `json:"products_id" biding:"required"`
	Quantities []int64 `json:"quantities" biding:"required"`
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

	user, err := server.store.GetUser(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var jsonUserStock []map[string]interface{}
	if err := json.Unmarshal(user.Stock, &jsonUserStock); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var amount int64
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
		amount += int64(removeProduct.UnitPrice) * (req.Quantities[idx])

		for _, data := range jsonUserStock {
			dataID := data["productID"].(float64)
			if dataID == float64(id) {
				// if id == int8(data["productID"].(float64)) {
				if int64(data["productQuantity"].(float64)) < int64(req.Quantities[idx]) {
					log.Info().Int("productId", int(id)).Int("quantity", int(req.Quantities[idx])).Msg("reduce client stock log")
					ctx.JSON(http.StatusNotAcceptable, errorResponse(fmt.Errorf("Not enough in stock to sell: Product: %v InStock: %v", data["productName"], data["productQuantity"])))
					return
				}
				// }
			}
		}
	}

	jsonSoldProduct, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sendSTKPayload := &worker.SendSTKPayload{
		User:            user,
		Amount:          amount,
		TransactionData: jsonSoldProduct,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.Queue(worker.QueueCritical),
	}
	err = server.taskDistributor.DistributeSendSTK(ctx, *sendSTKPayload, opts...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := server.redis.Del(ctx, GetUser+fmt.Sprintf("%v", uri.ID)).Err(); err != nil {
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
	log.Info().Msg("In mpesaCallbackURL from safaricom")
	var req mpesaCallbackRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := req.TransactionID[len(req.TransactionID)-3:]
	transactionId := req.TransactionID[:len(req.TransactionID)-3]

	body, _ := io.ReadAll(ctx.Request.Body)

	var callbackBody map[string]interface{}
	_ = json.Unmarshal(body, &callbackBody)

	// intUserID, _ := strconv.Atoi(userId)
	// user, err := server.store.GetUserForUpdate(ctx, int64(intUserID))
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		ctx.JSON(http.StatusNotFound, errorResponse(err))
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	// transaction, err := server.store.GetTransaction(ctx, transactionId)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		ctx.JSON(http.StatusNotFound, errorResponse(err))
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	processMpesaCallbackPayload := &worker.ProcessMpesaCallbackPayload{
		UserID:        userId,
		TransactionID: transactionId,
		Body:          callbackBody,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.Queue(worker.QueueCritical),
	}

	err := server.taskDistributor.DistributeProcessMpesaCallback(ctx, *processMpesaCallbackPayload, opts...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Add to redis queue
	// go func() {
	// 	processMpesaCallbackData(ctx, server, user, transaction)
	// }()

	ctx.JSON(http.StatusOK, gin.H{"Successful": "Reached"})
	return
}

type reduceClientProductByAdminRequest struct {
	Amount             int64   `json:"amount"`
	ProductsID         []int32 `json:"products_id"`
	Quantities         []int64 `json:"quantities"`
	PhoneNumber        string  `json:"phone_number"`
	MpesaReceiptNumber string  `json:"mpesa_receipt_number"`
	Description        string  `json:"description"`
}

type reduceClientProductByAdminUri struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) reduceClientProductByAdmin(ctx *gin.Context) {
	var uri reduceClientProductByAdminUri
	var req reduceClientProductByAdminRequest

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

	var jsonUserStock []map[string]interface{}
	if err := json.Unmarshal(user.Stock, &jsonUserStock); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var productIdInt64 []int64
	var productToAdd []db.Product
	var amount int64
	for idx, id := range req.ProductsID {

		productIdInt64 = append(productIdInt64, int64(id))

		removeProduct, err := server.store.GetProduct(ctx, int64(id))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		productToAdd = append(productToAdd, removeProduct)
		amount += int64(removeProduct.UnitPrice) * (req.Quantities[idx])

		for _, data := range jsonUserStock {
			dataID := data["productID"].(float64)
			if dataID == float64(id) {
				if int64(data["productQuantity"].(float64)) < int64(req.Quantities[idx]) {
					log.Info().Int("productId", int(id)).Int("quantity", int(req.Quantities[idx])).Msg("reduce client stock log")
					ctx.JSON(http.StatusNotAcceptable, errorResponse(fmt.Errorf("Not enough in stock to sell: Product: %v InStock: %v", data["productName"], data["productQuantity"])))
					return
				}
			}
		}
	}

	if req.Amount != amount {
		ctx.JSON(http.StatusNotAcceptable, errorResponse(fmt.Errorf("Amount does not match. Expected: %v Received: %v", amount, req.Amount)))
		return
	}

	transactionSoldData := map[string][]int64{
		"products_id": productIdInt64,
		"quantities":  req.Quantities,
	}

	jsonSoldProduct, err := json.Marshal(transactionSoldData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	reduceStockByAdminPayload := &worker.ReduceClientStockByAdminPayload{
		Amount:             amount,
		ProducToReduce:     productToAdd,
		Quantities:         req.Quantities,
		PhoneNumber:        req.PhoneNumber,
		MpesaReceiptNumber: req.MpesaReceiptNumber,
		Description:        req.Description,
		UserID:             int32(user.UserID),
		TransactionData:    jsonSoldProduct,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.Queue(worker.QueueDefault),
	}

	err = server.taskDistributor.DistributeSendReduceClientStockAdmin(ctx, *reduceStockByAdminPayload, opts...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := server.redis.Del(ctx, GetUser+fmt.Sprintf("%v", uri.ID)).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "added successfully"})
}

func redirectToPythonApp(user db.User, transaction db.Transaction, passErr error) {
	pythonEndpoint := "https://inventory-system.railway.internal/notify"

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

	sendResetPasswordPayload := &worker.SendResetPasswordEmail{
		Email: req.Email,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.Queue(worker.QueueCritical),
	}

	if err := server.taskDistributor.DistributeSendResetPasswordEmail(ctx, *sendResetPasswordPayload, opts...); err != nil {
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

	// user.Password = hashPassword
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
