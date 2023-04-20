package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	middlewares "sodality/handlers"
	"sodality/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostContent -> Create a creator content
var PostContent = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)

	var content models.Content
	err := json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	var existingUser models.User
	userID, _ := primitive.ObjectIDFromHex(props["user_id"].(string))

	userCollection := client.Database("sodality").Collection("users")
	err = userCollection.FindOne(r.Context(), bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if existingUser.ID != userID || err == mongo.ErrNoDocuments {
		middlewares.ErrorResponse("user does not exists", rw)
		return
	}

	content.UserID = userID.Hex()
	content.CreatedAt = time.Now().UTC()

	contentCollection := client.Database("sodality").Collection("content")
	result, err := contentCollection.InsertOne(context.TODO(), content)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessResponse(`inserted at `+strings.Replace(string(res), `"`, ``, 2), rw)
})

// GetContentByID -> Get content of user by content id
var GetContentByID = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var content models.Content

	contentID, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("sodality").Collection("content")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: contentID}}).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content id does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	middlewares.SuccessArrRespond(content, rw)
})

var GetOwnContent = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	var allContent []*models.Content

	// opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "fund", Value: -1}})

	collection := client.Database("sodality").Collection("content")
	cursor, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "user_id", Value: props["user_id"].(string)}}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var content models.Content
		err := cursor.Decode(&content)
		if err != nil {
			log.Fatal(err)
		}

		allContent = append(allContent, &content)
	}

	middlewares.SuccessArrRespond(allContent, rw)
})

// GetAllContentOfUser -> Get content of specific user
// var GetAllContentOfUser = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	var challenges []*models.Challenge

// 	collection := client.Database("challenge").Collection("challenges")
// 	cursor, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "coordinator", Value: params["username"]}})
// 	if err != nil {
// 		middlewares.ServerErrResponse(err.Error(), rw)
// 		return
// 	}

// 	for cursor.Next(context.TODO()) {
// 		var challenge models.Challenge
// 		err := cursor.Decode(&challenge)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		challenges = append(challenges, &challenge)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		middlewares.ServerErrResponse(err.Error(), rw)
// 		return
// 	}

// 	middlewares.SuccessChallengeArrRespond(challenges, rw)
// })

// // ListChallenge -> List all the challenges
// var ListChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 	var challenges []*models.Challenge
// 	collection := client.Database("challenge").Collection("challenges")
// 	cursor, err := collection.Find(context.TODO(), bson.D{})
// 	if err != nil {
// 		middlewares.ServerErrResponse(err.Error(), rw)
// 		return
// 	}

// 	for cursor.Next(context.TODO()) {
// 		var challenge models.Challenge
// 		err := cursor.Decode(&challenge)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		challenges = append(challenges, &challenge)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		middlewares.ServerErrResponse(err.Error(), rw)
// 		return
// 	}

// 	middlewares.SuccessChallengeArrRespond(challenges, rw)
// })

// // GetChallengs -> Get challenges for specific user
// var GetChallenges = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	var challenges []*models.Challenge

// 	collection := client.Database("challenge").Collection("challenges")
// 	cursor, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "coordinator", Value: params["username"]}})
// 	if err != nil {
// 		middlewares.ServerErrResponse(err.Error(), rw)
// 		return
// 	}

// 	for cursor.Next(context.TODO()) {
// 		var challenge models.Challenge
// 		err := cursor.Decode(&challenge)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		challenges = append(challenges, &challenge)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		middlewares.ServerErrResponse(err.Error(), rw)
// 		return
// 	}

// 	middlewares.SuccessChallengeArrRespond(challenges, rw)
// })

// CreateChallenge -> Create a challenge
var CreateChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var challenge models.Challenge
	err := json.NewDecoder(r.Body).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	challenge.CreatedAt = time.Now()
	challenge.UpdatedAt = time.Now()
	collection := client.Database("challenge").Collection("challenges")
	result, err := collection.InsertOne(context.TODO(), challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessResponse(`Inserted at `+strings.Replace(string(res), `"`, ``, 2), rw)
})

// GetChallenge -> Get a challenge by id
var GetChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var challenge models.Challenge

	collection := client.Database("challenge").Collection("challenges")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	middlewares.SuccessRespond(challenge, rw)
})

var UpdateChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var challenge models.Challenge

	err := json.NewDecoder(r.Body).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	collection := client.Database("challenge").Collection("challenges")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: challenge}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if res.MatchedCount == 0 {
		middlewares.ErrorResponse("Challenge does not exist", rw)
		return
	}

	middlewares.SuccessResponse("Updated", rw)
})

var DeleteChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

})

var JoinChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var challenge models.Challenge

	collection := client.Database("challenge").Collection("challenges")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	props, _ := r.Context().Value("props").(jwt.MapClaims)

	if challenge.Status == "private" {
		if challenge.Coordinator == props["username"] || challenge.RecipientAddress == props["identity"] {

			challenge.Participants = append(challenge.Participants, challenge.Identity)

			res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "Participants", Value: challenge.Participants}}}})
			if err != nil {
				middlewares.ServerErrResponse(err.Error(), rw)
				return
			}

			if res.MatchedCount == 0 {
				middlewares.ErrorResponse("challenge does not exist", rw)
				return
			}

			middlewares.SuccessRespond(params["id"], rw)
			return
		}
		middlewares.ForbiddenResponse("you have no access for this challenge", rw)
		return
	}

	challenge.Participants = append(challenge.Participants, challenge.Identity)

	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "Participants", Value: challenge.Participants}}}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if res.MatchedCount == 0 {
		middlewares.ErrorResponse("challenge does not exist", rw)
		return
	}

	middlewares.SuccessRespond(params["id"], rw)
})

var UnJoinChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var challenge models.Challenge

	collection := client.Database("challenge").Collection("challenges")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	props, _ := r.Context().Value("props").(jwt.MapClaims)

	identity := props["identity"]
	for i, v := range challenge.Participants {
		if v == identity {
			challenge.Participants = append(challenge.Participants[:i], challenge.Participants[i+1:]...)
			break
		}
	}

	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "Participants", Value: challenge.Participants}}}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if res.MatchedCount == 0 {
		middlewares.ErrorResponse("challenge does not exist", rw)
		return
	}

	middlewares.SuccessResponse("unjoin challenge successfully", rw)
})
