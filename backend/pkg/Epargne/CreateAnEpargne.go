package epargne

import (
	"fmt"
	"net/http"

	"github.com/RazanakotoMandresy/bank-app-aout/backend/pkg/common/models"
	"github.com/RazanakotoMandresy/bank-app-aout/backend/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// epargne otrn manangombola
// les valeur dans le request de create epargne seron:
// Reason ou nom : la raison de l'epargne
// Value raha ngeda nohon'i value ao anatin'i volan'i moneyUser ilay eparnge de tsy mety (logiaque)
// UUid anle manao (izay manao update anle money an'ilay connecter)
// Date maka ny jours ou du mois anaovana anle epargne
type CreateEpargneRequest struct {
	Name  string `json:"name"`
	Value int    `json:"ValueEpargne"`
	// tous les nombre dates du mois l'argent sera Epargner automatiquemen
	Date uint   `json:"DayEpargne"`
	Type string `json:"type"`
}

func (h handler) CreateEpargne(ctx *gin.Context) {
	// action du createEpargne : ilay userConnecter manao requetes hoe isakin'i inona ilay vola ny tonga automatiquement amin'ilay epargne
	body := new(CreateEpargneRequest)
	if err := ctx.Bind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	userConnectedUUID, err := middleware.ExtractTokenUUID(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}
	user, err := h.GetUserSingleUserFunc(userConnectedUUID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"err": err.Error()})
	}
	// logic stuff
	uuidEpargne := uuid.New()
	user.AutoEpargne = append(user.AutoEpargne, uuidEpargne.String())
	if body.Value > user.Moneys {
		err := fmt.Sprintf("vous ne pouvez pas epargner %v car l'argent sur votre compte est %v", body.Value, user.Moneys)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}
	if (body.Type) != "epargneMensuel" || body.Type != "economie" {
		err := fmt.Sprintf("on n'accepte pas les epargne du type : %v", body.Type)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	epargne := models.Eparge{
		ID:           uuidEpargne,
		Name:         body.Name,
		Value:        body.Value,
		DayPerMounth: body.Date,
		Type:         body.Type,
	}
	h.DB.Create(&epargne)
	// succes found
	ctx.JSON(http.StatusFound, gin.H{"epargne": &epargne})
}

func (h handler) GetUserSingleUserFunc(uuidToFind string) (*models.User, error) {
	var user models.User
	result := h.DB.First(&user, "uuid = ?", uuidToFind)
	if result.Error != nil {
		err := fmt.Errorf("user not in our database err : %v", result.Error)
		return nil, err
	}
	// if success user found user w t uuidToFins
	return &user, nil
}