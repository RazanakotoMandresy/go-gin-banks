package money

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/RazanakotoMandresy/go-gin-banks/pkg/common/models"
	"github.com/RazanakotoMandresy/go-gin-banks/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (h handler) SendMoney(ctx *gin.Context) {
	// extracttion de l'uuid depuis le bearer
	uuidConnectedStr, err := middleware.ExtractTokenUUID(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	userToSend := ctx.Param("uuid")
	body := new(sendMoneyRequest)

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	value := body.Value
	if value == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "value cannot be nul"})
		return
	}
	userConnected, err := h.GetUserByuuid(uuidConnectedStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return

	}
	// maka anle userRecepteur tokony par uuid na par AppUserName
	userRecepteur, err := h.GetUserByuuid(userToSend)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}
	_, found := slices.BinarySearch(userConnected.BlockedAcc, userRecepteur.AppUserName)
	if found {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("vous avez deja bloquer %v", userRecepteur.AppUserName)})
		return
	}
	// check si l'envoyeur essayent d'envoyer plus d'argent que ce qu'il en a
	if value > userConnected.Moneys {
		err := fmt.Errorf("impossible d'envoyer votre argent %v l'argent que vous voulez envoyer est %v", userConnected.Moneys, value)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}
	// check si l'userConnecter est la meme que celui qui il essaye d'envoyer de l'argent
	if uuidConnectedStr == userRecepteur.UUID {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "vous ne pouvez pas envoyer de l'argent a vous meme "})
		return
	}
	// 	message si tous se passe bien
	// message := fmt.Sprintf("%v a envoye un argent d'un montant de %v a %v", userConnected.AppUserName, value, userRecepteur.AppUserName)
	userConnected.Moneys = userConnected.Moneys - value
	userRecepteur.Moneys = (userRecepteur.Moneys + value)
	// save les money ao anaty user satria ireo tables fatsy uuid simplement ril
	h.DB.Save(userRecepteur)
	h.DB.Save(userConnected)
	// la models ho creena
	moneyTransaction, err := h.dbManipulationSendMoney(userConnected, userRecepteur, body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}
	ctx.JSON(http.StatusOK, &moneyTransaction)
}

func (h handler) dbManipulationSendMoney(userConnected, userRecepteur *models.User, body *sendMoneyRequest) (*models.Money, error) {
	var moneyTransaction models.Money
	res := h.DB.Where("send_by = ? AND sent_to = ?", userConnected.UUID, userRecepteur.UUID).First(&moneyTransaction)
	resume := fmt.Sprintf("%v a envoyer la somme de %v a %v a l'instant%v ", userConnected.AppUserName, body.Value, userRecepteur.AppUserName, time.Now())
	if res.Error != nil {
		fmt.Printf("no transaction betwen %s and %s creating a new one...", userConnected.UUID, userRecepteur.UUID)
		moneyTransaction.ID = uuid.New()
		moneyTransaction.SendBy = userConnected.UUID
		moneyTransaction.SentTo = userRecepteur.UUID
		moneyTransaction.Totals = body.Value
		moneyTransaction.Resume = resume
		moneyTransaction.SendByImg = userConnected.Image
		moneyTransaction.SendToImg = userRecepteur.Image
		moneyTransaction.SentToName = userRecepteur.AppUserName
		moneyTransaction.MoneyTransite = append(moneyTransaction.MoneyTransite, body.Value)

		result := h.DB.Create(moneyTransaction)
		if result.Error != nil {
			return nil, fmt.Errorf("creationraw %v", result.Error)
		}
		return &moneyTransaction, nil
	}
	fmt.Println("transaction already exist")
	moneyTransaction.MoneyTransite = append(moneyTransaction.MoneyTransite, body.Value)
	moneyTransaction.Totals = getTotals(moneyTransaction.MoneyTransite)
	h.DB.Save(&moneyTransaction)
	return &moneyTransaction, nil
}
func getTotals(money pq.Int32Array) int32 {
	var totals int32
	for i := 0; i < len(money); i++ {
		totals += money[i]
	}
	return totals
}

// user req que se soit uuid na appUserName
func (h handler) GetUserByuuid(userReq string) (*models.User, error) {
	var users models.User
	res := h.DB.Where("uuid = ? OR app_user_name = ?", userReq, userReq).First(&users)
	if res.Error != nil {
		return nil, fmt.Errorf(" %v pas dans uuid et AppUserName", userReq)
	}
	return &users, nil
}
