/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2019.09
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package gmongodb

import (
    "fmt"
    "time"
    "strings"
    "os"
    "context"

    "go4api/cmd"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

var MongoDBClients map[string]*mongo.Client

type MongoDBExec struct {
    TargetMongoDB string
    Database string
    Collection string
    CmdSlice []string
    CmdAffectedCount int
    CmdResults interface{}
}

func InitMongoDBConnection () {
    MongoDBClients = make(map[string]*mongo.Client)

    mongs := cmd.GetMongoDBConfig()

    // master only
    for k, v := range mongs {
        ip := v.Ip
        port := v.Port
        user := v.UserName
    
        pw := ""
        pwV := v.Password
        pwV = strings.Replace(pwV, "${", "", -1)
        pwV = strings.Replace(pwV, "}", "", -1)
        if len(pwV) > 0 {
            pw = os.Getenv(pwV)
        }
        
        client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + user + ":" + pw + "@" + ip + ":" + fmt.Sprint(port)))
        if err != nil {
            panic(err)
        }

        ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
        err = client.Connect(ctx)

        if err != nil {
            panic(err)
        }

        err = client.Ping(ctx, nil)
        if err != nil {
            panic(err)
        }

        key := strings.ToLower(k)
        MongoDBClients[key] = client
    }
} 


func Run (cmdStr string) (int, interface{}, string) {
    var err error
    mongoExecStatus := ""
    
    tDb := "master"
    sl := strings.Split(cmdStr, ", ")
    cmdSlice := sl[2:]

    mongoDBExec := &MongoDBExec{tDb, sl[0], sl[1], cmdSlice, 0, ""}
    err = mongoDBExec.Do()

    if err == nil {
        mongoExecStatus = "cmdSuccess"
    } else {
        mongoExecStatus = "cmdFailed"
    }

    return mongoDBExec.CmdAffectedCount, mongoDBExec.CmdResults, mongoExecStatus
}


func (mongoDBExec *MongoDBExec) Do () error {
    client := MongoDBClients[mongoDBExec.TargetMongoDB]

    var err error
    var res interface{}

    fmt.Println(">>>>>>>>: ", mongoDBExec.CmdSlice[0])

    //
    switch strings.ToUpper(mongoDBExec.CmdSlice[0]) {
        case "GET":
            collection := client.Database(mongoDBExec.Database).Collection(mongoDBExec.Collection)

            ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

            // get data
            objID, _ := primitive.ObjectIDFromHex(mongoDBExec.CmdSlice[2])
            filter := bson.M{"_id": objID}
            res := collection.FindOne(ctx, filter)

            rr, err := res.DecodeBytes()

            if err != nil {
                panic(err)
            }

            if err == nil {
                mongoDBExec.CmdAffectedCount = 1
                mongoDBExec.CmdResults = rr
            }

            fmt.Println(">>>>>>>>: ", string(rr))

        case "UPDATE":
            collection := client.Database(mongoDBExec.Database).Collection(mongoDBExec.Collection)

            ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

            // update data
            objID, _ := primitive.ObjectIDFromHex(mongoDBExec.CmdSlice[2])
            res, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{mongoDBExec.CmdSlice[3]: mongoDBExec.CmdSlice[4]}})

            if err != nil {
                panic(err)
            }

            if err == nil {
                mongoDBExec.CmdAffectedCount = 1
                mongoDBExec.CmdResults = res
            }

            fmt.Println(">>>>>>>>: ", res)

        default:
            mongoDBExec.CmdAffectedCount = -1
            fmt.Println("!! Warning, Command ", mongoDBExec.CmdSlice, " is not supported currently, will enhance it later")
    }

    return err
}
